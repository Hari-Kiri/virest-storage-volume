package storageVolume

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeCreate"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Create a storage volume within a pool based on an XML description. Not all pools support creation of volumes.
//
// Option:
//
//   - 1 = (0x1; 1 << 0) : Prealloc metadata in flags can be used to get higher performance with qcow2 image files
//     which don't support full preallocation, by creating a sparse image file with metadata.
//
//   - 2 = (0x2; 1 << 1) : perform a btrfs lightweight copy.
//
//   - 4 = (0x4; 1 << 2) : Validate the XML document against schema.
func VolumeCreate(connection virest.Connection, poolUuid string, xmlConfig libvirtxml.StorageVolume, option uint) (volumeCreate.Key, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return volumeCreate.Key{}, virestError, isError
	}
	defer storagePoolObject.Free()

	marshalXmlConfig, errorMarshalXmlConfig := xmlConfig.Marshal()
	virestError.Error, isError = errorMarshalXmlConfig.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed marshaling supplied storage volume xml config: %s", virestError.Message)
		return volumeCreate.Key{}, virestError, isError
	}

	storageVolumeObject, errorCreateStorageVolumeObject := storagePoolObject.StorageVolCreateXML(marshalXmlConfig, libvirt.StorageVolCreateFlags(option))
	virestError.Error, isError = errorCreateStorageVolumeObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed create storage volume: %s", virestError.Message)
		return volumeCreate.Key{}, virestError, isError
	}
	defer storageVolumeObject.Free()

	storageVolumeName, errorGetStorageVolumeName := storageVolumeObject.GetKey()
	virestError.Error, isError = errorGetStorageVolumeName.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("Failed to get the name of the created volume: %s", virestError.Message)
		return volumeCreate.Key{}, virestError, isError
	}

	return volumeCreate.Key{
		Key: storageVolumeName,
	}, virestError, false
}
