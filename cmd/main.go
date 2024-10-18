package main

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/handler"
	"log"
)

func main() {
	srv := new(taskFlow.Server)
	handlers := new(handler.Handler)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
