package errors

import "log"

// logs a fatal error and causes the process to exit
func CheckFatal(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

// logs a recoverable error
func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}