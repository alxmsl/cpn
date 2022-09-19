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
	log.Println(args(n, ss...)...)
	return
}

func Logf(format string, n string, ss ...interface{}) {
	if !NeedLog(n) {
		return
	}
	log.Printf(format, args(n, ss...)...)
}

func args(n string, ss ...interface{}) []interface{} {
	var aa = make([]interface{}, 0, len(ss)+1)
	aa = append(aa, n)
	for _, s := range ss {
		aa = append(aa, s)
	}
	return aa
}

func NeedLog(n string) bool {
	var (
		_, ok1 = nn["*"]
		_, ok2 = nn[n]
	)
	return ok1 || ok2
}
