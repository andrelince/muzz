package logger

import (
	"github.com/sirupsen/logrus"
)

type Settings struct {
}

func New(settings Settings) *logrus.Logger {
	return logrus.New()
}
