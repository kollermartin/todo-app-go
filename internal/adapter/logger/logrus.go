package logger

import (
	"github.com/sirupsen/logrus"
)

func Set() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
}