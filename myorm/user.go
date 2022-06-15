package myorm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `gorm:"unique, not null"`
	Passwd   string  `gorm:"not null"`
	Ip       string  `gorm:"unique, not null"`
	Conn     []*Conn `gorm:"many2many:user_conn;association_jointable_foreignkey:conn_id"`
	IsAdmin  bool    `gorm:"default:false"`
}

func (User) TableName() string {
	return "user"
}

// user 和  conn 多对多
type Conn struct {
	gorm.Model
	Local   string  `gorm:"unique,default:null"`
	Remote  string  `gorm:"unique,not null"`
	Svcname string  `gorm:"unique,not null"`
	User    []*User `gorm:"many2many:user_conn;association_jointable_foreignkey:user_id"`
}

func (Conn) TableName() string {
	return "conn"
}

type Workflow struct {
	gorm.Model
	Username  string
	Localport string
	Svcname   string
	Pass      uint8 `gorm:"default:1"` // 1 审核中  2 通过  3  拒绝
}
