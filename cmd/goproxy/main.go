package main

import (
	"log"

	"github.com/algrvvv/goproxy/internal"
)

func main() {
	serv, err := internal.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	serv.Run()
	select {}
}
