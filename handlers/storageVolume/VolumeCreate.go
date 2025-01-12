package storageVolume

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/modules/storageVolume"
	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeCreate"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func VolumeCreate(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData volumeCreate.Request
		httpBody        volumeCreate.Response
	)

	connection, errorRequestPrecondition, isError := storageVolume.RequestPrecondition(
		request,
		http.MethodPost,
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

	result, errorGetVolumeList, isErrorGetVolumeList := storageVolume.VolumeCreate(
		connection,
		requestBodyData.PoolUuid,
		requestBodyData.StorageVolume,
		requestBodyData.Option,
	)
	if isErrorGetVolumeList {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetVolumeList.Code)
		httpBody.Error = errorGetVolumeList
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed create storage volume [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetVolumeList.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("storage volume created on pool", requestBodyData.PoolUuid, "inside hypervisor", request.Header["Hypervisor-Uri"][0],
		"[", request.URL.Path, "]")
}
