package db

import (
	"context"
	"services-api/config"
	"services-api/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitDB(ctx context.Context, log *logger.Logger) (*gorm.DB, error) {
	log.Info(ctx, "initializing db")

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	}

	db, err := gorm.Open(postgres.Open(config.DBDSN), gormConfig)
	if err != nil {
		log.Error(ctx, "failed to connect to postgres db", "error", err)
		return nil, err
	}

	log.Info(ctx, "db connection done")
	return db, nil
}

// This can be used while initialising the DB once or inserting the data manually in the DB using the below query
func SeedData(ctx context.Context, db *gorm.DB, log *logger.Logger) {

	// Insert into services
	servicesQuery := `
		INSERT INTO services (name, description) VALUES
		('User Management', 'Handles user authentication and profiles'),
		('Payment Service', 'Processes payments and transactions')
		RETURNING id
	`
	rows, err := db.WithContext(ctx).Raw(servicesQuery).Rows()
	if err != nil {
		log.Error(ctx, "Failed to seed services", "error", err)
		return
	}
	defer rows.Close()

	var serviceIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			log.Error(ctx, "Failed to scan service ID", "error", err)
			return
		}
		serviceIDs = append(serviceIDs, id)
	}
	log.Info(ctx, "Seeded services table", "ids", serviceIDs)

	// Insert into versions using retrieved service IDs
	versionsQuery := `
		INSERT INTO versions (service_id, number) VALUES
		($1, '1.0.0'),
		($1, '1.1.0'),
		($2, '2.0.0')
	`
	if err := db.WithContext(ctx).Exec(versionsQuery, serviceIDs[0], serviceIDs[1]).Error; err != nil {
		log.Error(ctx, "Failed to seed versions", "error", err)
		return
	}
	log.Info(ctx, "Seeded versions table")
}
