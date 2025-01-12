package storageVolume

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/modules/storageVolume"
	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeDelete"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func VolumeDelete(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData volumeDelete.Request
		httpBody        volumeDelete.Response
	)

	connection, errorRequestPrecondition, isError := storageVolume.RequestPrecondition(
		request,
		http.MethodDelete,
		&requestBodyData,
		os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_NAME"),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_VOLUME_APPLICATION_JWT_SIGNATURE_KEY")),
	)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorRequestPrecondition.Code)
		httpBody.Error = errorRequestPrecondition
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorRequestPrecondition.Message,
		)
		return
	}
	defer connection.Close()

	errorDeleteVolume, isErrorDeleteVolume := storageVolume.VolumeDelete(
		connection,
		requestBodyData.PoolUuid,
		requestBodyData.VolumeName,
		requestBodyData.Option,
	)
	if isErrorDeleteVolume {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorDeleteVolume.Code)
		httpBody.Error = errorDeleteVolume
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed delete storage volume [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorDeleteVolume.Message,
		)
		return
	}

	utils.NoContentResponseBuilder(responseWriter)
	temboLog.InfoLogging("storage volume deleted on pool", requestBodyData.PoolUuid, "inside hypervisor", request.Header["Hypervisor-Uri"][0],
		"[", request.URL.Path, "]")
}
