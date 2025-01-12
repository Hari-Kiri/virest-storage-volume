package volumeDelete

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Error    virest.Error `json:"error"`
}
