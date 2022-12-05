package main

import (
	"fmt"
	"github.com/Sjhzjxc/go_db"
)

func main() {
	db, err := go_db.DefaultDb("msmh", "zxw123456", "www.cosck.com:3306", "msmh", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	var schemas []string
	err = db.DB.Table("books").Select("book_name").Limit(10).Find(&schemas).Error
	if err != nil {
		fmt.Println(err)
	}
}
