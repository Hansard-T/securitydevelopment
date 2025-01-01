package vars

import "sync"

var (
	Result    *sync.Map
)

func init() {
	Result = &sync.Map{}
}