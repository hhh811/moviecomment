package auth

import (
	"moviecomment/models"

	"github.com/astaxie/beego/session"
)

func GetUserIdFromSession(sess session.Store) int {
	if id, ok := sess.Get("auth_user_id").(int); ok && id > 0 {
		return id
	}
	return 0
}

// get user if key exists in session
func GetUserFromSession(user *models.User, sess session.Store) bool {
	id := GetUserIdFromSession(sess)
	if id > 0 {
		u := models.User{Id: id}
		if u.Read() == nil {
			*user = u
			return true
		}
	}

	return false
}
