package api

import (
	"fmt"
	"time"

	"dynexo.de/pkg/log"

	"github.com/labstack/echo"
	//sha "golang.org/x/crypto/sha3"
	sha "crypto/sha256"
)

const (
	debug_api_auth = true
)

func ApiAuthMiddleware(AccessKey, SecretKey []byte, logger log.ILogger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// return next(c) for successful authentications
			if ok := apiAuthHandler(c, AccessKey, SecretKey, logger); ok {
				return next(c)
			}

			c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Access restricted")
			return echo.ErrUnauthorized
		}
	}
}

func apiAuthHandler(c echo.Context, accessKey, secretKey []byte, logger log.ILogger) bool {

	req := c.Request()
	if req == nil {
		if debug_api_auth {
			logger.Error("API::apiAuthHandler: Invalid request")
		}
		return false
	}

	reqAuthStr := req.Header.Get("Authorization")
	if len(reqAuthStr) < 32 {
		if debug_api_auth {
			logger.Error("API::apiAuthHandler: Invalid or missing authorization header")
		}
		return false
	}

	reqDateStr := req.Header.Get("x-wix-date")
	if len(reqDateStr) < 14 {
		if debug_api_auth {
			logger.Error("API::apiAuthHandler: Invalid or missing request date")
		}
		return false
	}

	dateStr := reqDateStr

	// generate authorization based on AWS S3 authz definition
	//dateStr := time.Now().UTC().Format("20060102T150405Z")

	time_parsed, err := time.Parse("20060102T150405Z", dateStr)
	if err != nil {
		if debug_api_auth {
			logger.Error("API::apiAuthHandler: Invalid date string:", err)
		}
		return false
	}

	// validate remote time is close to local time -> prevent replays
	if time_parsed.Before(time.Now().Add(5*time.Minute)) || time_parsed.After(time.Now().Add(5*time.Minute)) {
		return false
	}

	reqStr := fmt.Sprintf("%s %s\r\nHost: %s\r\nx-wix-date: %s",
		req.Method, req.URL.Path, req.URL.Host, dateStr,
	)

	// cretae auth string
	dateSign := sha.Sum256(append(secretKey, []byte(dateStr)...))
	authSign := sha.Sum256(append(dateSign[:], []byte(reqStr)...))
	authStr := fmt.Sprintf("id=%s,sign=%64x", accessKey, authSign)

	if authStr == reqAuthStr {
		return true
	}

	//req.Header.Set("x-wix-date", dateStr)
	//req.Header.Set("x-wix-endp", endpointStr)
	//req.Header.Set("Authorization", authStr)
	//req.Header.Set("User-Agent", "dxoFiles/1.0 SecureDav/1.2.24")

	//log.Printf("ERR: API::apiAuthHandler: Invalid authorization")

	return false
}
