package simulator

import (
	"alphanonce.com/exchangesimulator/internal/log"
)

var logger *log.Logger

func init() {
	logger = log.NewDefault().With(log.String("package", "simulator"))
}
