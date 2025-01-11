package volumeListAll

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     []Data       `json:"data"`
	Error    virest.Error `json:"error"`
}

type Data struct {
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Type       libvirt.StorageVolType `json:"Type"`
	Capacity   uint64                 `json:"capacity"`
	Allocation uint64                 `json:"allocation"`
	Physical   uint64                 `json:"physical"`
}
