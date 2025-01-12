package storageVolume

import (
	"fmt"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeListAll"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Collect the list of all storage volumes inside a pool. The physical on disk usage
// can be different than the calculated allocation value as is the case with qcow2
// files.
//
// Storage Volume Type:
//
//   - 0 = (0x0) : Regular file based volumes
//
//   - 1 = (0x1) : Block based volumes
//
//   - 2 = (0x2) : Directory-passthrough based volume
//
//   - 3 = (0x3) : Network volumes like RBD (RADOS Block Device)
//
//   - 4 = (0x4) : Network accessible directory that can contain other network volumes
//
//   - 5 = (0x5) : Ploop based volumes
//
//   - 6 = (0x6)
func VolumeListAll(connection virest.Connection, poolUuid string) ([]volumeListAll.Data, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return []volumeListAll.Data{}, virestError, isError
	}
	defer storagePoolObject.Free()

	// extra flags; not used yet, so callers should always pass 0
	storageVolumes, errorGetListOfAllStorageVolumes := storagePoolObject.ListAllStorageVolumes(0)
	virestError.Error, isError = errorGetListOfAllStorageVolumes.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get list of all storage volumes: %s", virestError.Message)
		return []volumeListAll.Data{}, virestError, isError
	}

	var waitGroup sync.WaitGroup
	result := make([]volumeListAll.Data, len(storageVolumes))
	waitGroup.Add(len(storageVolumes) * 3)
	for i := 0; i < len(storageVolumes); i++ {
		defer storageVolumes[i].Free()

		go func(index int) {
			defer waitGroup.Done()

			errorGetStorageVolumeRef := storageVolumes[index].Ref()
			if errorGetStorageVolumeRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage volume:", errorGetStorageVolumeRef)
				return
			}
			defer storageVolumes[index].Free()

			storageVolumeName, errorGetStorageVolumeName := storageVolumes[index].GetName()
			if errorGetStorageVolumeName != nil {
				temboLog.ErrorLogging("failed get storage volume name:", errorGetStorageVolumeName)
				return
			}

			result[index].Name = storageVolumeName
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStorageVolumeRef := storageVolumes[index].Ref()
			if errorGetStorageVolumeRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage volume:", errorGetStorageVolumeRef)
				return
			}
			defer storageVolumes[index].Free()

			storageVolumePath, errorGetStorageVolumePath := storageVolumes[index].GetPath()
			if errorGetStorageVolumePath != nil {
				temboLog.ErrorLogging("failed get storage volume path:", errorGetStorageVolumePath)
				return
			}

			result[index].Path = storageVolumePath
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStorageVolumeRef := storageVolumes[index].Ref()
			if errorGetStorageVolumeRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage volume:", errorGetStorageVolumeRef)
				return
			}
			defer storageVolumes[index].Free()

			// extra flags; not used yet, so callers should always pass 0
			storageVolumeXmlDesc, errorGetStorageVolumeXmlDesc := storageVolumes[index].GetXMLDesc(0)
			if errorGetStorageVolumeXmlDesc != nil {
				temboLog.ErrorLogging("failed get storage volume detail:", errorGetStorageVolumeXmlDesc)
				return
			}

			var storageVolume libvirtxml.StorageVolume
			errorUnmarshal := storageVolume.Unmarshal(storageVolumeXmlDesc)
			if errorUnmarshal != nil {
				temboLog.ErrorLogging("failed unmarshal storage volume xml desc:", errorGetStorageVolumeXmlDesc)
				return
			}

			result[index].Type = storageVolume.Type
			result[index].Capacity = *storageVolume.Capacity
			result[index].Allocation = *storageVolume.Allocation
			result[index].Physical = *storageVolume.Physical
		}(i)
	}
	waitGroup.Wait()

	return result, virestError, false
}
