package api

import (
	"net/http"
	"runtime"
	"time"

	"dynexo.de/ufyle/pkg/model/v1"

	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
	// "github.com/sevenNt/echo-pprof"
)

func apiServiceStart(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServiceStop(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServiceReload(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServicePause(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServiceRun(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServicePurge(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data map[string]interface{}

		runtime.GC()

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}

func apiServiceStats(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		m := &runtime.MemStats{}

		data := map[string]interface{}{
			"time":       time.Now(),
			"goroutines": runtime.NumGoroutine(),
			"mem_acq":    m.Sys,
			"mem_use":    m.Alloc,
			"num_malloc": m.Mallocs,
			"num_free":   m.Frees,
			"gc_enabled": m.EnableGC,
			"num_gc":     m.NumGC,
			"last_gc":    m.LastGC,
			"next_gc":    m.NextGC,
		}

		return c.JSON(http.StatusOK, JsonSuccess("succeeded", data))
	})
}
