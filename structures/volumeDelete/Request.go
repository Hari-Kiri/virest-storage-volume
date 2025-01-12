package volumeDelete

type Request struct {
	PoolUuid   string `json:"poolUuid"`
	Option     uint   `json:"option"`
	VolumeName string `json:"volumeName"`
}
