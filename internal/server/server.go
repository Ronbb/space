package server

import (
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

	e.GET("/space/directory", func(c echo.Context) error {
		start, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		end, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		path := c.QueryParam("path")

		s, err := db.GetDirectorySpace(path, start, end)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/space/volume", func(c echo.Context) error {
		start, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		end, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		path := c.QueryParam("path")

		s, err := db.GetVolumeSpace(path, start, end)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &s)
	})

	e.GET("/last_record", func(c echo.Context) error {
		record, err := db.GetLastRecord()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &record)
	})

	e.PUT("/directory", func(c echo.Context) error {
		v := model.DirectoryHash{}
		err := c.Bind(&v)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = db.PutDirectory(v.Directory, v.Limit)
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

		err = db.PutVolume(v.Volume, v.Limit, v.LimitPercentage)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "success")
	})

	e.DELETE("/directory", func(c echo.Context) error {
		v := model.DirectoryHash{}
		err := c.Bind(&v)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = db.RemoveDirectory(v.Directory)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "success")
	})

	e.DELETE("/volume", func(c echo.Context) error {
		v := model.VolumeHash{}
		err := c.Bind(&v)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = db.RemoveVolume(v.Volume)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "success")
	})

	return e.Start(":8005")
}
