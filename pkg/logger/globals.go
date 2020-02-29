package logger

import (
	"fmt"
	"log"
)

// Dummy global logger implementation
// Possible improvements - implement as a interface / struct, prefixes, different log levels, push logs to any monitoring tool, ...

func Errorf(err error, format string, v ...interface{}) {
	m := fmt.Errorf(fmt.Sprintf(format, v)+": %w", err)
	log.Println(m)
}
