package volumeCreate

import "libvirt.org/go/libvirtxml"

type Request struct {
	PoolUuid      string                   `json:"poolUuid"`
	StorageVolume libvirtxml.StorageVolume `json:"storageVolume"`
}
