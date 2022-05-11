package models

import "time"

type UserRepository struct {
	Id                 int
	Identity           string
	UserName       string
	ParentId           int64
	RepositoryIdentity string
	Ext                string
	Name               string
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated"`
}

func (table UserRepository) TableName() string {
	return "user_repository"
}
