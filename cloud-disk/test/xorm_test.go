package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"cloud-disk/core/models"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func TestXormTest(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:Ljf941118@tcp(127.0.0.1:3306)/cloud-disk?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}

	data := make([]*models.User, 0)
	err = engine.Find(&data)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dst.String())
}

