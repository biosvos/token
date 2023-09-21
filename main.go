package main

import (
	"log"
	"time"
)

func main() {
	data := `{"a":"b"}`
	creator, err := NewTokenCreator("localhost.dev.key")
	if err != nil {
		log.Panicf("%+v", err)
	}
	token, err := creator.Create([]byte(data), 60*time.Second)
	if err != nil {
		log.Panicf("%+v", err)
	}
	log.Println("token", token)

	parser, err := NewTokeParser("localhost.dev.crt")
	if err != nil {
		log.Panicf("%+v", err)
	}

	parse, err := parser.Parse(token)
	if err != nil {
		log.Panicf("%+v", err)
	}

	log.Println("data", parse)
}
