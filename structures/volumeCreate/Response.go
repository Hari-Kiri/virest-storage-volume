package volumeCreate

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Name         `json:"data"`
	Error    virest.Error `json:"error"`
}

type Name struct {
	Name string `json:"Name"`
}
