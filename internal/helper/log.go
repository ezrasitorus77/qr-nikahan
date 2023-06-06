package helper

import (
	"fmt"
	"log"
	"os"
	"time"

	"qr-nikahan/domain"
)

var logger *domain.Log

func init() {
	file, e := os.OpenFile(fmt.Sprintf("%s.txt", time.Now().Format("01-02-2006")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if e != nil {
		log.Fatal(e)
	}

	logger = &domain.Log{
		File:          file,
		InfoLogger:    log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		PanicLogger:   log.New(file, "PANIC: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func INFO(text interface{}) {
	logger.InfoLogger.Println(text)
}

func WARNING(text interface{}) {
	logger.WarningLogger.Println(text)
}

func ERROR(text interface{}) {
	logger.ErrorLogger.Println(text)
}

func PANIC(text interface{}) {
	logger.PanicLogger.Panicln(text)
}
