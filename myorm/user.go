package myorm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique, not null"`
	Passwd   string `gorm:"not null"`
	Ip       string `gorm:"unique, not null"`
	Conn     []Conn `gorm:"many2many:user_conn;"`
	// Conns    []*Conn `gorm:"-"`
	IsAdmin bool `gorm:"default:false"`
}

func (User) TableName() string {
	return "user"
}

// user 和  conn 多对多
type Conn struct {
	gorm.Model
	// port
	Local string `gorm:"unique,default:null"`
	// ip:port
	Remote  string `gorm:"unique,not null"`
	Svcname string `gorm:"unique,not null"`
	// Users   []*User `gorm:"-"`
	User []User `gorm:"many2many:user_conn;"`
}

func (Conn) TableName() string {
	return "conn"
}

type Workflow struct {
	gorm.Model
	Username  string
	Localport string
	Svcname   string
	Pass      bool `gorm:"default:false"`
}
