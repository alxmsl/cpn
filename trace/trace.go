package trace

import (
	"log"
	"os"
	"strings"
)

var nn = map[string]struct{}{}

func init() {
	var ss = strings.Split(os.Getenv("CPN_DEBUG_NAMES"), ",")
	for _, s := range ss {
		var k = strings.TrimSpace(s)
		if k == "" {
			continue
		}
		nn[k] = struct{}{}
	}
	if len(nn) > 0 {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}
}

func Log(n string, ss ...interface{}) {
	if !NeedLog(n) {
		return
	}

	var aa = make([]interface{}, 0, len(ss)+1)
	aa = append(aa, n)
	for _, s := range ss {
		aa = append(aa, s)
	}
	log.Println(aa...)
	return
}

func NeedLog(n string) bool {
	var (
		_, ok1 = nn["*"]
		_, ok2 = nn[n]
	)
	return ok1 || ok2
}
