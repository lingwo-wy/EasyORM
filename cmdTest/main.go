package main

import (
	"fmt"
	"EasyORM/easyorm"
	_ "github.com/mattn/go-sqlite3"
	"EasyORM/session"
	"EasyORM/log"
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func tests () {
	
	engine, _ := easyorm.NewEngine("sqlite3", "easy.db")
	defer engine.Close()
	s := engine.NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		fmt.Println("失败")
	}
}
func main() {
	
	tests()

	
}
