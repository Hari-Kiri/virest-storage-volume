package storageVolume

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/modules/storageVolume"
	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeListAll"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func VolumeListAll(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData volumeListAll.Request
		httpBody        volumeListAll.Response
	)

	connection, errorRequestPrecondition, isError := storageVolume.RequestPrecondition(
		request,
		http.MethodGet,
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

	result, errorGetVolumeList, isErrorGetVolumeList := storageVolume.VolumeListAll(connection, requestBodyData.PoolUuid)
	if isErrorGetVolumeList {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetVolumeList.Code)
		httpBody.Error = errorGetVolumeList
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed get list of storage volume [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetVolumeList.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("listing storage volume on pool", requestBodyData.PoolUuid, "inside hypervisor", request.Header["Hypervisor-Uri"][0],
		"[", request.URL.Path, "]")
}
