package out

import "log"

func printf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Success(format string, a ...interface{}) {
	format = "\033[32m" + format + "\033[0m \U0001F604\n"
	printf(format, a...)
}

func Start(format string, a ...interface{}) {
	format = "\033[1;32m" + format + "\033[0m \U0001F7E2\n"
	printf(format, a...)
}

func Running(format string, a ...interface{}) {
	format = "\033[1;33m" + format + "\033[0m \U0001F7E1\n"
	printf(format, a...)
}

func Error(format string, a ...interface{}) {
	format = "\033[31m" + format + "\033[0m \U0001F534\n"
	printf(format, a...)
}
