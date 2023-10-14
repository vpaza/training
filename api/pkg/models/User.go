package models

import (
	adh "github.com/adh-partnership/api/pkg/database/models"
	"github.com/vpaza/training/api/pkg/database"
	"gorm.io/gorm/clause"
)

func FindUser(id string) (*adh.User, error) {
	u := &adh.User{}
	if err := database.DB.Preload(clause.Associations).Where("c_id = ?", id).First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}
