package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type User struct {
	Uid   int
	Name  string
	Phone string
}

func (u User) tableName() string {
	return "User"
}

func init() {
	DB, _ = sql.Open("mysql", "shennong_oss:test1234@tcp(9.135.133.29:3306)/shennong_oss")
}

func GetUser(uid int) (u User, err error) {
	err = DB.QueryRow("select uid,name,phone from `user` where uid=?", uid).Scan(&u.Uid, &u.Name, &u.Phone)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return u, err
	}
	return u, nil
}
