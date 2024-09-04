package debug

import "fmt"

var is_debug bool

func DebugInit(debug bool) {
	is_debug = debug
}

func Logging(str string) {
	if is_debug {
		fmt.Println(str)
	}
}
