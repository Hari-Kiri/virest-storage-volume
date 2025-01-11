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
	Key                string                 `json:"key"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Type               libvirt.StorageVolType `json:"Type"`
	Capacity           uint64                 `json:"capacity"`
	CurrentAllocation  uint64                 `json:"currentAllocation"`
	PhysicalAllocation uint64                 `json:"physicalAllocation"`
}
