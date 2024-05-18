package main

import (
	"log"

	"github.com/alanpramil7/go-chat/db"
	"github.com/alanpramil7/go-chat/internal/router"
	"github.com/alanpramil7/go-chat/internal/user"
)

func main () {
  db, err := db.NewDatabase()
  if err != nil {
    log.Fatalf("Error connecting database: %v", err)
  }

  userRep := user.NewRepository(db.GetDB())
  userService := user.NewService(userRep)
  userHadler := user.NewHandler(userService)

  router.InitRouter(userHadler)
  router.Start("0.0.0.0:8080")
}
