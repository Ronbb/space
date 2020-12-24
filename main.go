package main

import (
	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/space"
)

func main() {
	p := "C:\\old"
	s, err := space.GetDirectorySpace(p)
	if err != nil {
		println(err.Error())
		return
	}
	println(p, s.UsedSpace, s.Percentage)

	db, err := database.Open()
	if err != nil {
		println(err.Error())
		return
	}

	defer db.Close()
}
