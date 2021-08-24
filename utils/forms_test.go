package utils_test

import (
	"moviecomment/controllers/base"
	"moviecomment/modules/auth"
	"testing"
)

func TestLoginForm(t *testing.T) {
	form := auth.LoginForm{}
	//fSets := utils.NewFormSets(&form, nil)
	rt := base.BaseRouter{}
	rt.SetFormSets(&form)
	//fmt.Printf("%v", fSets)
	//t.Logf("%v", fSets)
}
