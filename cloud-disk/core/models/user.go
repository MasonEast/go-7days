package models

import "time"

type User struct {
	Id        int `xorm:"id"`
	Identity  string
	UserName      string `xorm:"username"`
	Password  string `xorm:"password"`
	Email     string `xorm:"email"`
	CreatedTime time.Time `xorm:"create_time"`
}

func (table User) TableName() string {
	return "user"
}
