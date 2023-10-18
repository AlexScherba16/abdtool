package main

import (
	"abdtool/internal/composer"
	"log"
)

func main() {
	app, err := composer.NewComposedApplication()
	if err != nil {
		err.AppendStackTrace("main")
		log.Fatalln(err.Error())
	}

	_ = app
}
