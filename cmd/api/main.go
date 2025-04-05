package main

import (
	"log"
	"net/http"
	"secretWall/internal/handler"
	"secretWall/internal/infra"
	"secretWall/internal/repository"
	"secretWall/internal/service"
)

func main() {
	mux := http.NewServeMux()

	db, err := infra.InitDB()

	if err != nil {
		log.Fatal(err)
	}

	var userRepo = repository.NewUserRepo(db)
	var authService = service.NewAuthService(userRepo)
	handler.NewAuthHandler(authService, mux)

	log.Fatal(http.ListenAndServe("192.168.1.76:8080", mux))
}
