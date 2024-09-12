package simulator

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/types"
)

//go:generate mockery --name Simulator --output ../mocks --outpkg mocks

type Simulator interface {
	Process(request types.Request, startTime time.Time) (types.Response, time.Time)
}
