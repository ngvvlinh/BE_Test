package api

import (
	//"encoding/json"
	"net/http"

	"github.com/kr/pretty"

	"dynexo.de/ufyle/pkg/err"
	"dynexo.de/ufyle/pkg/model/v1"

	"github.com/labstack/echo"
)

var (
	dbg = pretty.Println
)

type (
	UfylerV1Func map[int]func([]byte) (umodel.IRes, error)
)

func apiUfylerV1Init(s umodel.IController) UfylerV1Func {

	return UfylerV1Func{
		16: s.HandleFile,
		17: s.HandleFolder,
		20: s.HandleSearch,
		22: s.HandleRecenlyUsed,
	}
}

func apiUfylerV1(s umodel.IController, uf UfylerV1Func) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		// var m struct {
		// 	Filename string
		// }

		var m umodel.Message

		if err := c.Bind(&m); err != nil {
			return c.JSON(http.StatusOK, JsonError(http.StatusBadRequest, "Failed to load request"))
		}

		//dbg(m)

		uhandler, ok := uf[m.T]
		if !ok {
			return c.JSON(http.StatusOK, JsonError(http.StatusBadRequest, "Unsupported request type"))
		}

		resp, err := uhandler(m.D)

		dbg(resp)

		data := map[string]interface{}{
			"t": 1,
			"d": resp,
			"e": err,
			"s": 0,
		}

		return c.JSON(http.StatusOK, data)
	})
}

func apiUfylerV1Errors(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		data := map[string]interface{}{
			"t": 1,
			"d": err.ErrorsList(),
			"e": 0,
			"s": 0,
		}

		return c.JSON(http.StatusOK, data)
	})
}

func apiUfylerV1Message(s umodel.IController) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		data := map[string]interface{}{
			"t": 32,
			"d": s.Messages(),
			"e": 0,
			"s": 0,
		}

		return c.JSON(http.StatusOK, data)
	})
}
