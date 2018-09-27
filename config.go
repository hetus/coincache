package coincache

import (
	"time"
)

type Config struct {
	Database string
	Debug    bool
	Market   string
	Interval time.Duration
}
