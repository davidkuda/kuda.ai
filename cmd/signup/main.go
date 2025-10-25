package main

import (
	"flag"
	"log"

	"github.com/davidkuda/kudaai/internal/envcfg"
	"github.com/davidkuda/kudaai/internal/models"
)

func main() {
	var err error

	email := flag.String("email", "", "email address of the new user")
	password := flag.String("password", "", "email address of the new user")
	update := flag.Bool("update", false, "update pw of user")
	flag.Parse()

	if *email == "" {
		log.Fatal("use this cmd with the -email <email> flag.")
	}

	if *password == "" {
		log.Fatal("use this cmd with the -password <password> flag.")
	}

	db, err := envcfg.DB()
	if err != nil {
		log.Fatalf("could not open DB: %v\n", err)
	}
	defer db.Close()
	m := models.UserModel{DB: db}

	if *update {
		err = m.UpdatePassword(*email, *password)
		if err != nil {
			log.Fatalf("m.Update: failed to update password: %v\n", err)
		}
		log.Printf("success: updated password of %s\n", *email)
	} else {
		err = m.Insert(*email, *password)
		if err != nil {
			log.Fatalf("m.Insert(): failed to create new user: %v\n", err)
		}
		log.Println("Inserted new user to DB!")
	}
}
