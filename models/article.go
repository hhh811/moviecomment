package models

import (
	"time"
)

type Article struct {
	Id      int
	User    *User     `orm:"rel(fk)"`
	Title   string    `orm:"size(60)"`
	Content string    `orm:"type(text)"`
	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
}
