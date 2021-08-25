package reflect

import (
	"runtime"
	"strconv"
)

func At(skip int) string {
	pc, _, lineNumber, _ := runtime.Caller(skip + 1)
	return runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(lineNumber)
}
