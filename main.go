package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/model"
	"github.com/ronbb/space/internal/runner"
	"github.com/ronbb/space/internal/server"
	"gopkg.in/toast.v1"
)

var gen = flag.Bool("gen", false, "generate")

func init() {
	flag.Parse()
}

func main() {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if *gen {
		rs, err := db.GetRecords()
		if err != nil {
			println(err.Error())
			return
		}
		b, err := json.Marshal(&rs)
		if err != nil {
			println(err.Error())
			return
		}
		ioutil.WriteFile("out.json", b, 0644)
		return
	}

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
			if s.Limit != 0 {
				limit := s.Limit
				if s.LimitPercentage {
					limit = int64(float64(limit) / 100 * float64(s.TotalSpace))
				}
				if s.FreeSpace < uint64(limit) {
					errString += fmt.Sprintf("The free size of %s is less than limit\n", s.Volume)
				}
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
