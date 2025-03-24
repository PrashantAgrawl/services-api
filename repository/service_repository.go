package repository

import (
	"context"
	"database/sql"
	"services-api/logger"
	"services-api/models"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewServiceRepository(db *gorm.DB, logger *logger.Logger) *ServiceRepository {
	return &ServiceRepository{db: db, logger: logger}
}

func (r *ServiceRepository) GetServices(ctx context.Context, name, sort string, page, pageSize int) ([]models.Service, error) {
	r.logger.Info(ctx, "Querying services", "name_filter", name, "sort", sort, "page", page, "page_size", pageSize)

	query := `
		SELECT s.id AS service_id, 
			s.name AS service_name, 
			s.description AS service_description,
			v.id AS version_id,
			v.service_id AS version_service_id,
			v.number AS version_number
		FROM services s
		LEFT JOIN versions v ON s.id = v.service_id
	`
	var args []interface{}

	/*	if name != "" {
			query += " WHERE s.name LIKE $1"
			args = append(args, "%"+name+"%")
		}

		if sort != "" {
			query += " ORDER BY " + sort
		}
	*/
	offset := (page - 1) * pageSize
	query += " LIMIT $1 OFFSET $2"
	args = append(args, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		r.logger.Error(ctx, "Failed to query services", "error", err)
		return nil, err
	}
	defer rows.Close()

	servicesMap := make(map[uint]*models.Service)
	for rows.Next() {
		var serviceId uint
		var serviceName, serviceDescription string
		var versionID, versionServiceID sql.NullInt64
		var versionNumber sql.NullString

		err := rows.Scan(&serviceId, &serviceName, &serviceDescription, &versionID, &versionServiceID, &versionNumber)
		if err != nil {
			r.logger.Error(ctx, "Failed to scan row", "error", err)
			return nil, err
		}

		if _, exists := servicesMap[serviceId]; !exists {
			servicesMap[serviceId] = &models.Service{
				Id:          serviceId,
				Name:        serviceName,
				Description: serviceDescription,
				Versions:    []models.Version{},
			}
		}

		if versionID.Valid {
			servicesMap[serviceId].Versions = append(servicesMap[serviceId].Versions, models.Version{
				Id:        uint(versionID.Int64),
				ServiceId: uint(versionServiceID.Int64),
				Number:    versionNumber.String,
			})
		}
	}

	var services []models.Service
	for _, service := range servicesMap {
		services = append(services, *service)
	}

	return services, nil
}

func (r *ServiceRepository) GetService(ctx context.Context, id string) (models.Service, error) {
	r.logger.Info(ctx, "Fetching service", "id", id)

	var service models.Service
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, name, description
		FROM services
		WHERE id = $1
	`, id).Scan(&service).Error

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Error(ctx, "Service not found", "id", id)
		} else {
			r.logger.Error(ctx, "Failed to fetch service", "id", id, "error", err)
		}
		return models.Service{}, err
	}

	var versions []models.Version
	err = r.db.WithContext(ctx).Raw(`
		SELECT id, service_id, number
		FROM versions
		WHERE service_id = $1
	`, id).Scan(&versions).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to fetch versions", "id", id, "error", err)
		return models.Service{}, err
	}
	service.Versions = versions

	return service, nil
}

func (r *ServiceRepository) CreateServiceRaw(ctx context.Context, name, description string, versions []string) error {
	r.logger.Info(ctx, "Creating service with raw SQL", "name", name)

	var serviceId uint
	err := r.db.WithContext(ctx).Raw("INSERT INTO services (name, description) VALUES ($1, $2) RETURNING id", name, description).Scan(&serviceId).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to insert service", "name", name, "error", err)
		return err
	}

	for _, version := range versions {
		err := r.db.WithContext(ctx).Exec("INSERT INTO versions (service_id, number) VALUES ($1, $2)", serviceId, version).Error
		if err != nil {
			r.logger.Error(ctx, "Failed to insert version", "service_id", serviceId, "number", version, "error", err)
			return err
		}
	}

	r.logger.Info(ctx, "Service created with raw SQL", "id", serviceId)
	return nil
}
