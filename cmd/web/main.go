package main

import (
	"log"
	"net/http"
)


func main() {
	log.Print("Starting web server, listening on port 8873")
	err := http.ListenAndServe(":8873", routes())
	log.Fatal(err)
}
