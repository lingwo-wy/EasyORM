package session

import (
	// "testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}
