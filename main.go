package main

import (
	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/runner"
	"github.com/ronbb/space/internal/server"
)

func main() {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go runner.SaveSpace(db, func(e error) {
		if err != nil {
			println(err.Error())
		}
		println("updated")
	})
	err = server.Run(db)
	if err != nil {
		panic(err)
	}
}
