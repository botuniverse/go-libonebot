package comm

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// TODO: this should be moved to some other package
func handleAction(request gjson.Result) gjson.Result {
	log.Debugf("Action request: %v", request)
	// TODO: now it simply return the request
	return request
}
