package main

import (
	"Week02/src/api"
	"fmt"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover error: %s", err)
		}
	}()
	fmt.Println("123")
	err := api.Init()
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
}

type errString struct {
	s string
}

func New(msg string) errString {
	return errString{s: msg}
}

func (e *errString) Error() string {
	return e.s
}
