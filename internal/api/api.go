package api

import (
	"fmt"
	"net/http"
	"strings"

	"dynexo.de/pkg/log"

	"dynexo.de/ufyle/pkg/model/v1"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//	"github.com/sevenNt/echo-pprof"
)

const (
	debug_api_handler = false
	const_base_path   = "ufyle/lib."

	// index = start pos, "connect/relay." = 14, "api" = 3, ".func1" = 6
	const_base_path_len = 14
)

var (
	name_base_path     string
	name_base_path_len int
)

type ServiceApi struct {
	e      *echo.Echo
	Logger log.ILogger
}

func NewServiceApi(s umodel.IController, logger log.ILogger) *ServiceApi {

	e := echo.New()

	if debug_api_handler {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	}

	api_service := e.Group("/service/")
	{
		api_service.POST("reload", apiServiceReload(s))
		api_service.POST("pause", apiServicePause(s))
		api_service.POST("run", apiServiceRun(s))
		api_service.POST("purge", apiServicePurge(s))
		api_service.POST("stats", apiServiceStats(s))
	}

	api_ufyler := e.Group("/ufyler/")
	{
		uf := apiUfylerV1Init(s)

		api_ufyler.POST("v1", apiUfylerV1(s, uf))
		api_ufyler.POST("v1/err", apiUfylerV1Errors(s))
		api_ufyler.POST("v1/msg", apiUfylerV1Message(s))
	}

	/*
		e.GET("/api/v1/stats/conn", apiV1RelayStatsConn(wix))
		e.GET("/api/v1/stats/endpoint", apiV1RelayStatsEndpoint(wix))
		e.GET("/api/v1/stats/mem", apiV1RelayStatsMem(wix))
		e.GET("/api/v1/stats/relay", apiV1RelayStatsRelay(wix))

		api_relay := e.Group("/api/v1/relay/")
		{
			if !wix.config.Api.NoAuth {
				api_relay.Use(api.ApiAuthMiddleware(wix.config.Api.AccessKey, wix.config.Api.SecretKey, wix.logger))
			}

			api_relay.POST("reload", apiRelayControllerReload(wix))
			api_relay.POST("pause", apiRelayControllerDisable(wix))
			api_relay.POST("run", apiRelayControllerEnable(wix))
			api_relay.POST("disable", apiRelayControllerDisable(wix))
			api_relay.POST("enable", apiRelayControllerEnable(wix))
			api_relay.POST("purge", apiRelayControllerPurge(wix))

			api_relay.GET("endpoints", apiRelayEndpointList(wix))
			api_relay.GET("endpoint/:id", apiRelayEndpointGet(wix))
			api_relay.POST("endpoint/:id", apiRelayEndpointCtrl(wix))
			//api_relay.POST("endpoint", apiRelayEndpointAdd(wix))
			//api_relay.DELETE("endpoint", apiRelayEndpointDel(wix))

			// domain(s) is no longer needed since /api/v1/auth/domain provides some function
			//api_relay.GET("domains", apiRelayDomainList(wix))
			//api_relay.GET("domain/:id", apiRelayDomainGet(wix))
			//api_relay.POST("domain", apiRelayDomainAdd(wix))
			//api_relay.DELETE("domain", apiRelayDomainDel(wix))

			api_relay.POST("config/show", apiRelayConfigGet(wix))
			api_relay.POST("config/load", apiRelayConfigLoad(wix))
			api_relay.POST("config/save", apiRelayConfigSave(wix))
			api_relay.POST("config/reload", apiRelayConfigReload(wix))
		}

		api_db := e.Group("/api/v1/auth/")
		{
			if !wix.config.Api.NoAuth {
				api_db.Use(api.ApiAuthMiddleware(wix.config.Api.AccessKey, wix.config.Api.SecretKey, wix.logger))
			}

			api_db.GET("endpoints", apiDatabaseEndpointList(wix))
			api_db.GET("endpoint/:id", apiDatabaseEndpointGet(wix))
			api_db.POST("endpoint", apiDatabaseEndpointAdd(wix))
			api_db.PUT("endpoint/:id", apiDatabaseEndpointUpd(wix))
			api_db.DELETE("endpoint/:id", apiDatabaseEndpointDel(wix))

			api_db.GET("domains", apiDatabaseDomainList(wix))
			api_db.GET("domain/:id", apiDatabaseDomainGet(wix))
			api_db.POST("domain", apiDatabaseDomainAdd(wix))
			api_db.PUT("domain/:id", apiDatabaseDomainUpd(wix))
			api_db.DELETE("domain/:id", apiDatabaseDomainDel(wix))

			api_db.GET("endpoints/domain/:id", apiDatabaseDomainEndpointsList(wix))
			api_db.GET("domain/:id/endpoints", apiDatabaseDomainEndpointsList(wix))

			api_db.GET("export", apiDatabaseExport(wix))
			api_db.POST("import", apiDatabaseImport(wix))
		}

		// this must be the last entry
		e.GET("/api/v1/routes", apiV1Routes(e.Routes()))

	*/

	api := &ServiceApi{e, logger}

	logger.Debug("Registered API routes")
	for _, entry := range e.Routes() {
		if index := strings.Index(entry.Name, const_base_path); index >= 0 {
			name := entry.Name[index+const_base_path_len+3 : len(entry.Name)-6]
			str := fmt.Sprintf("%-08s %s -> %s", entry.Method, entry.Path, name)
			logger.Info("API::Route:", str)
		}
	}

	return api
}

func apiV1Routes(routes []*echo.Route) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {

		var data []interface{}

		for _, entry := range routes {
			if index := strings.Index(entry.Name, const_base_path); index >= 0 {
				name := entry.Name[index+const_base_path_len+3 : len(entry.Name)-6]
				_data := map[string]interface{}{
					"method": entry.Method,
					"path":   entry.Path,
					"name":   name,
				}

				data = append(data, _data)
			}
		}

		return c.JSON(http.StatusOK, JsonSuccess("success", data))
	})
}

func (c *ServiceApi) Run(server *http.Server) error {

	c.e.HideBanner = true
	if !debug_api_handler {
		c.e.Logger.SetOutput(log.NewNullWriter())
	}
	c.e.Debug = debug_api_handler

	return c.e.StartServer(server)
}

func (c *ServiceApi) Routes() []string {
	var routes []string

	for _, entry := range c.e.Routes() {
		routes = append(routes, entry.Method+" "+entry.Path+" "+entry.Name)
	}
	return routes
}

//
// handler
//

func ListenAndServeApi(addr string, s umodel.IController, logger log.ILogger) error {

	address := ":903"
	if len(addr) > 2 {
		address = addr
	}

	server, err := NewApiServer(address, logger)
	if err != nil {
		return err
	}

	apiService := NewServiceApi(s, logger)
	err = apiService.Run(server)

	return err
}
