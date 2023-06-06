package domain

import (
	"log"
	"os"
)

type (
	Log struct {
		File          *os.File
		InfoLogger    *log.Logger
		WarningLogger *log.Logger
		ErrorLogger   *log.Logger
		PanicLogger   *log.Logger
	}
)
