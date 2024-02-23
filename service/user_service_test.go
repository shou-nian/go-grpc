package service

import (
	"regexp"
	"testing"
)

func TestUserServiceClient_Login(t *testing.T) {
	//re, err := regexp.Compile(`(?m)^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
	//
	//if err != nil {
	//	t.Error(err)
	//}
	//if ok := re.MatchString("a123@qq.com"); !ok {
	//	t.Error("invalid email address")
	//}

	if ok, err := regexp.MatchString(`(?m)^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`, "a12@qq.com"); !ok {
		if err != nil {
			t.Error(err)
		}
		t.Error("invalid email")
	}
}
