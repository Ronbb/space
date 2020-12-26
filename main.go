package main

import (
	"fmt"

	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/model"
	"github.com/ronbb/space/internal/runner"
	"github.com/ronbb/space/internal/server"
	"gopkg.in/toast.v1"
)

func main() {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go runner.SaveSpace(db, func(record *model.SpaceRecord, e error) {
		notification := toast.Notification{
			AppID:   "Space",
			Title:   "Space limit warning",
			Message: "Some message about how important something is...",
			Audio:   toast.IM,
			Loop:    true,
		}
		if err != nil {
			println(err.Error())
			notification.Message = fmt.Sprintf("Failed to update space record %s", err.Error())
			notification.Push()
			return
		}
		println("updated")
		errString := ""
		for _, s := range record.DirectoriesSpace {
			if s.Limit != 0 && s.UsedSpace > uint64(s.Limit) {
				errString += fmt.Sprintf("The size of %s is over limit\n", s.Directory)
			}
		}
		for _, s := range record.VolumesSpace {
			if s.Limit != 0 && (s.TotalSpace-s.FreeSpace) > uint64(s.Limit) {
				errString += fmt.Sprintf("The size of %s is over limit\n", s.Volume)
			}
		}
		if errString != "" {
			println(errString)
			notification.Message = errString
			err = notification.Push()
			if err != nil {
				println(err.Error())
			}
		}
	})

	err = server.Run(db)
	if err != nil {
		panic(err)
	}
}
