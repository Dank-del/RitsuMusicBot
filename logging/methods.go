package logging

import (
	"fmt"
	"time"
)

func display(data string) {
	fmt.Printf("%s %s===== [%s]\n", botLabel, data, time.Now())
}

// Error Log Error
// Usage : logging.Error("Yamete Kudasai")
func Error(data string) {
	display(fmt.Sprintf("%s ==== [%s] ==", errorLabel, data))
}

// Info Log Info
// Usage : logging.Info("Kimochiii")
func Info(data string) {
	display(fmt.Sprintf("%s ==== [%s] ==", infoLabel, data))
}

// Warn Log warning
// Usage : logging.Warn("UOHHHHHH")
func Warn(data string) {
	display(fmt.Sprintf("%s ==== [%s] ==", warnLabel, data))
}
