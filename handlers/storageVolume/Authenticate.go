package storageVolume

import (
	"net/http"
	"os"
	"time"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/structures/authenticate"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/auth"
	"github.com/golang-jwt/jwt"
)

func Authenticate(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		httpBody authenticate.Response
		result   string
	)

	jwtLifetimeSeconds, errorParseJwtLifetimeSeconds, isError := utils.StringToUint32(os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_JWT_LIFETIME_SECONDS"))
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorParseJwtLifetimeSeconds.Code)
		httpBody.Error = errorParseJwtLifetimeSeconds
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed parsing JWT lifetime seconds environment to uint32 [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorParseJwtLifetimeSeconds.Message,
		)
		return
	}

	result, errorBasicAuth, isErrorBasicAuth := auth.BasicAuth(
		request,
		os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_BA_USER"),
		os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_BA_PASSWORD"),
		os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_NAME"),
		time.Second*time.Duration(jwtLifetimeSeconds),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_JWT_SIGNATURE_KEY")),
	)
	if isErrorBasicAuth {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorBasicAuth.Code)
		httpBody.Error = errorBasicAuth
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed authenticate with basic auth [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorBasicAuth.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Token = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("authenticate user", "[", request.URL.Path, "]")
}
