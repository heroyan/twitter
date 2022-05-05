package dao

import "testing"

func TestNewRedisDao(t *testing.T) {
	rd := NewRedisDao("localhost:6379", "", 0)
	userName := "yanshuifa"
	user, err := rd.GetUser(userName)
	if err != nil {
		t.Errorf("get user error: %+v", err)
		return
	}
	if user == nil {
		t.Logf("%s not found", userName)
		return
	}
	t.Logf("found user: %+v", user)
}
