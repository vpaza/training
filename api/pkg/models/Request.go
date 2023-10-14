package models

import (
	"time"

	adh "github.com/adh-partnership/api/pkg/database/models"
	"github.com/vpaza/training/api/pkg/database"
)

type Request struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;not null"`
	Type         string     `gorm:"type:varchar(255);not null"`
	Creator      *adh.User  `gorm:"foreignKey:CreatorID;references:CID;not null"`
	CreatorID    string     `gorm:"type:varchar(255);not null"`
	Start        *time.Time `gorm:"type:datetime;not null"`
	End          *time.Time `gorm:"type:datetime;not null"`
	Title        string     `gorm:"type:varchar(255);not null"`
	Position     string     `gorm:"type:varchar(255);not null"`
	Notes        string     `gorm:"type:mediumtext;not null"`
	Accepted     bool       `gorm:"type:tinyint(1);not null"`
	AcceptedBy   *adh.User  `gorm:"foreignKey:AcceptedByID;references:CID"`
	AcceptedByID string     `gorm:"type:varchar(255)"`
	AcceptedAt   *time.Time `gorm:"type:datetime"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateRequest(r *Request) error {
	return database.DB.Create(r).Error
}

func FindRequest(id uint) (*Request, error) {
	r := &Request{}
	if err := database.DB.Where("id = ?", id).First(r).Error; err != nil {
		return nil, err
	}

	return r, nil
}

func FindRequests() ([]*Request, error) {
	var rs []*Request
	if err := database.DB.Find(&rs).Error; err != nil {
		return nil, err
	}

	return rs, nil
}

func FindRequestsByCreatorID(id string) ([]*Request, error) {
	var rs []*Request
	if err := database.DB.Where("creator_id = ?", id).Find(&rs).Error; err != nil {
		return nil, err
	}

	return rs, nil
}

func FindRequestsByAcceptedByID(id string) ([]*Request, error) {
	var rs []*Request
	if err := database.DB.Where("accepted_by_id = ?", id).Find(&rs).Error; err != nil {
		return nil, err
	}

	return rs, nil
}

func UpdateRequest(r *Request) error {
	return database.DB.Save(r).Error
}

func DeleteRequest(r *Request) error {
	return database.DB.Delete(r).Error
}
