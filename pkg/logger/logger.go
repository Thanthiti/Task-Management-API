package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// 0755 = owner can edit other can read
		if err := os.Mkdir("logs", 0755); err != nil {
			logrus.Fatalf("Unable to make directory log: %v", err)
		}
	}
	// Create or open file
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// 0666 = Everyone can write
	// if not file will create
	if err != nil {
		logrus.Warn("Unable to write log to file, will fallback to stdout.")
		Log.Out = os.Stdout
	} else {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		Log.SetOutput(multiWriter)
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
	})
	// Set log level
	Log.SetLevel(logrus.InfoLevel)
	// Show which file did the log come from
	Log.SetReportCaller(true)
}
