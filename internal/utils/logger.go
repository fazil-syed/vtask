package utils

// add a simple logger that logs messages to the console with different log levels (info, warning, error)
import (
	"log"
)
func Info(msg string) {
	log.Printf("[INFO] %s", msg)
}
func Warning(msg string) {
	log.Printf("[WARNING] %s", msg)
}
func Error(msg string) {
	log.Printf("[ERROR] %s", msg)
}
