package models

import (
	"time"
)

const (
	User_Role_Common    = iota // 普通用户
	User_Role_Commenter = iota // 影评人
	User_Role_Admin     = iota // 管理员
)

type User struct {
	Id       int
	UserName string    `orm:"size(30);unique"`
	NickName string    `orm:"size(30)"`
	Password string    `orm:"size(128)"`
	Role     int       `orm:"inxex"`
	Created  time.Time `orm:"auto_now_add"`
	Updated  time.Time `orm:"auto_now"`
}
