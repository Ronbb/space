package main

import (
	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/size"
)

func main() {
	s, err := size.GetSpaceOfDir("C:\\Software\\ARCTIME_PRO_2.3_WIN64")
	if err != nil {
		println(err.Error())
		return
	}
	println("a", s)

	db, err := database.Open()
	if err != nil {
		println(err.Error())
		return
	}

	db.Test()

}
