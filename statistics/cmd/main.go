package main

import (
	"log"
	"statistics/web/app"
)

func main() {
	app := app.Setup()
	log.Fatal(app.Listen(":3000"))
}
