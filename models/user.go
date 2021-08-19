package models

import (
	"fmt"
	"moviecomment/setting"
	"moviecomment/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
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
	Rands    string    `orm:"size(10)"` // 加盐
	Created  time.Time `orm:"auto_now_add"`
	Updated  time.Time `orm:"auto_now"`
}

func (m *User) Insert() error {
	m.Rands = GetUserSalt()
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) String() string {
	return strconv.FormatInt(int64(m.Id), 10)
}

// return url of user page
func (m *User) Link() string {
	return fmt.Sprintf("%suser/%s", setting.AppUrl, m.UserName)
}

// 加盐，返回一个随机字符串
func GetUserSalt() string {
	return utils.GetRandomString(10) // TODO specify random string len
}
