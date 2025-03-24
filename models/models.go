package models

type Service struct {
	Id          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
}

type Version struct {
	Id        uint   `gorm:"primaryKey" json:"id"`
	ServiceId uint   `gorm:"not null" json:"service_id"`
	Number    string `json:"number"`
}
