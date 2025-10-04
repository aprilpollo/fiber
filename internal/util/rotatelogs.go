package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"time"
)

func RotateLogs() (io.Writer, error) {
	_ = os.MkdirAll("logs", 0755)

	rl, err := rotatelogs.New(
		"logs/app.%Y-%m-%d.log",                
		rotatelogs.WithLinkName("logs/app.log"), 
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour), 
	)
	if err != nil {
		log.Fatalf("rotatelogs init error: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, rl)

	return mw, nil

}
