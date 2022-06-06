package myorm

import (
	"gorm.io/gorm"
)

type Sshinfo struct {
	gorm.Model
	Username string `gorm:"not null"`
	Passwd   string `gorm:"not null"`
	Port     string `gorm:"not null"`
	Host     string `gorm:"not null"`
}

func (Sshinfo) TableName() string {
	return "sshinfo"
}
