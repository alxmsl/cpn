package transition

import (
	"log"

	"github.com/alxmsl/cpn"
)

func First(mm []*cpn.M) *cpn.M {
	log.Println("mm", mm[0])
	return mm[0]
}
