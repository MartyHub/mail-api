package worker

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/invopop/validation"
)

type Config struct {
	BatchSize  int32         `envDefault:"5"`
	Count      int           `envDefault:"2"`
	Interval   time.Duration `envDefault:"30s"`
	MaxTries   int16         `envDefault:"3"`
	RetryDelay time.Duration `envDefault:"30s"`

	Stopper chan bool
	Waiter  *sync.WaitGroup
}

func (c Config) Stop() {
	// Send one stop signal for each sender
	for i := 0; i < c.Count; i++ {
		c.Stopper <- true
	}
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BatchSize, validation.Required, validation.Min(1)),
		validation.Field(&c.Count, validation.Required, validation.Min(1)),
		validation.Field(&c.Interval, validation.Required),
		validation.Field(&c.MaxTries, validation.Required, validation.Min(1)),
		validation.Field(&c.RetryDelay, validation.Required),
	)
}

func (c Config) String() string {
	sb := strings.Builder{}

	sb.WriteString("Sender Config:\n")
	sb.WriteString(fmt.Sprintf("  - Batch Size: %d\n", c.BatchSize))
	sb.WriteString(fmt.Sprintf("  - Count: %d\n", c.Count))
	sb.WriteString(fmt.Sprintf("  - Interval: %v\n", c.Interval))
	sb.WriteString(fmt.Sprintf("  - Max Tries: %d\n", c.MaxTries))

	return sb.String()
}
