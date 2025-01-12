package volumeDelete

import "libvirt.org/go/libvirtxml"

type Request struct {
	PoolUuid      string                   `json:"poolUuid"`
	Option        uint                     `json:"option"`
	StorageVolume libvirtxml.StorageVolume `json:"storageVolume"`
}
