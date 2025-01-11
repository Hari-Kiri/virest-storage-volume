package storageVolume

import (
	"fmt"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-volume/structures/volumeListAll"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Collect the list of all storage volumes inside a pool. The physical on disk usage
// can be different than the calculated allocation value as is the case with qcow2
// files.
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
	waitGroup.Add(len(storageVolumes) * 4)
	for i := 0; i < len(storageVolumes); i++ {
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

			storageVolumeInfo, errorGetStorageVolumeInfo := storageVolumes[index].GetInfo()
			if errorGetStorageVolumeInfo != nil {
				temboLog.ErrorLogging("failed get storage volume info:", errorGetStorageVolumeInfo)
				return
			}

			result[index].Type = storageVolumeInfo.Type
			result[index].Capacity = storageVolumeInfo.Capacity
			result[index].Allocation = storageVolumeInfo.Allocation
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStorageVolumeRef := storageVolumes[index].Ref()
			if errorGetStorageVolumeRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage volume:", errorGetStorageVolumeRef)
				return
			}
			defer storageVolumes[index].Free()

			storageVolumeInfo, errorGetStorageVolumeInfo := storageVolumes[index].GetInfoFlags(1)
			if errorGetStorageVolumeInfo != nil {
				temboLog.ErrorLogging("failed get storage volume info:", errorGetStorageVolumeInfo)
				return
			}

			result[index].Physical = storageVolumeInfo.Allocation
		}(i)
	}
	waitGroup.Wait()

	return result, virestError, false
}
