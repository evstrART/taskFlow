package main

import (
	"github.com/evstrART/taskFlow"
	"log"
)

func main() {
	srv := new(taskFlow.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatal(err)
	}
}
