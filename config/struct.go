package config

import "time"

type configuration struct {
	Port   string
	DbHost string
	DbPort string
	DbName string
	DbUser string
	DbPwd  string
	Driver string
	Limit  int

	EmailAdmin    string
	EmailSend     string
	EmailPassword string

	ReadTimeout time.Duration
}
