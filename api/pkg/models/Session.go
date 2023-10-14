package models

import (
	"time"

	adh "github.com/adh-partnership/api/pkg/database/models"
	"github.com/vpaza/training/api/pkg/database"
)

type Session struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;not null"`
	Request      *Request   `gorm:"foreignKey:RequestID;references:ID;not null"`
	RequestID    uint       `gorm:"not null"`
	Start        *time.Time `gorm:"type:datetime;not null"`
	End          *time.Time `gorm:"type:datetime;not null"`
	Position     string     `gorm:"type:varchar(255);not null"`
	Notes        string     `gorm:"type:mediumtext;not null"`
	Student      *adh.User  `gorm:"foreignKey:StudentID;references:CID;not null"`
	StudentID    string     `gorm:"type:varchar(255);not null"`
	Instructor   *adh.User  `gorm:"foreignKey:InstructorID;references:CID;not null"`
	InstructorID string     `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateSession(s *Session) error {
	return database.DB.Create(s).Error
}

func FindSession(id uint) (*Session, error) {
	s := &Session{}
	if err := database.DB.Where("id = ?", id).First(s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func FindSessions() ([]*Session, error) {
	s := []*Session{}
	if err := database.DB.Find(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func FindSessionByRequestID(id uint) (*Session, error) {
	s := &Session{}
	if err := database.DB.Where("request_id = ?", id).First(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func FindSessionsByStudentID(id string) ([]*Session, error) {
	s := []*Session{}
	if err := database.DB.Where("student_id = ?", id).Find(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func FindSessionsByInstructorID(id string) ([]*Session, error) {
	s := []*Session{}
	if err := database.DB.Where("instructor_id = ?", id).Find(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func UpdateSession(s *Session) error {
	return database.DB.Save(s).Error
}

func DeleteSession(s *Session) error {
	return database.DB.Delete(s).Error
}
