package storageVolume

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func VolumeDelete(connection virest.Connection, poolUuid, volumeName string, option uint) (virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	storageVolumeObject, errorCreateStorageVolumeObject := storagePoolObject.LookupStorageVolByName(volumeName)
	virestError.Error, isError = errorCreateStorageVolumeObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage volume object: %s", virestError.Message)
		return virestError, isError
	}

	virestError.Error, isError = storageVolumeObject.Delete(libvirt.StorageVolDeleteFlags(option)).(libvirt.Error)

	return virestError, isError
}
