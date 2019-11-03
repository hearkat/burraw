package main

import "log"

func ERR(s ...interface{}) {
	log.Println("[ERROR]", s)
}

func WARN(s ...interface{}) {
	log.Println("[WARNING]", s)
}

func LOG(s ...interface{}) {
	log.Println("[LOG]", s)
}

func DEBUG(s ...interface{}) {
	log.Println("[DEBUG]", s)
}
