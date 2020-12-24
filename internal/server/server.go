package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/model"
)

// Run .
func Run(db database.DB) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/directories", func(c echo.Context) error {
		s, err := db.GetDirectories()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/volumes", func(c echo.Context) error {
		s, err := db.GetVolumes()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/space/directory/:path", func(c echo.Context) error {
		start, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		end, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		b, err := base64.RawStdEncoding.DecodeString(c.Param("path"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		path := string(b)

		s, err := db.GetDirectorySpace(path, start, end)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/space/volume/:path", func(c echo.Context) error {
		start, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		end, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		b, err := base64.RawStdEncoding.DecodeString(c.Param("path"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		path := string(b)

		s, err := db.GetVolumeSpace(path, start, end)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/last_record_time", func(c echo.Context) error {
		t, err := db.GetLastRecordTime()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("%d", t))
	})

	e.PUT("/directory", func(c echo.Context) error {
		v := model.DirectoryHash{}
		err := c.Bind(&v)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = db.PutDirectory(v.Directory)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "success")
	})

	e.PUT("/volume", func(c echo.Context) error {
		v := model.VolumeHash{}
		err := c.Bind(&v)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = db.PutVolume(v.Volume)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "success")
	})

	return e.Start(":8005")
}
