package internal

import (
	"log"
	"path/filepath"
	"runtime"
)

func JustLogErr(args ...interface{}) {
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			_, f, l, _ := runtime.Caller(1)
			log.Printf("%s %d: %+v", filepath.Base(f), l, err)
		}
	}
}
