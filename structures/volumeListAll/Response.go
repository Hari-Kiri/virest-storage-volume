package volumeListAll

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     []Data       `json:"data"`
	Error    virest.Error `json:"error"`
}

type Data struct {
	Name       string                       `json:"name"`
	Path       string                       `json:"path"`
	Type       string                       `json:"Type"`
	Capacity   libvirtxml.StorageVolumeSize `json:"capacity"`
	Allocation libvirtxml.StorageVolumeSize `json:"allocation"`
	Physical   libvirtxml.StorageVolumeSize `json:"physical"`
}
