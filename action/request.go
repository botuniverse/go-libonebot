package action

import (
	"github.com/tidwall/gjson"
)

type Request struct {
	Action Action
	Params gjson.Result
	Echo   interface{}
}
