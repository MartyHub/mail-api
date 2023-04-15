package smtp

import (
	"fmt"
	"strings"

	"github.com/MartyHub/mail-api/utils"
	"github.com/invopop/validation"
)

type Config struct {
	AuthType     string
	AuthUsername string
	AuthPassword string

	From string
	Host string `envDefault:"localhost"`
	Port int    `envDefault:"25"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AuthType, validation.In(AuthCRAMMD5, AuthPlain)),
		validation.Field(&c.From, validation.Required),
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required, validation.Min(1)),
	)
}

func (c Config) String() string {
	sb := strings.Builder{}

	sb.WriteString("SMTP Config:\n")
	sb.WriteString(fmt.Sprintf("  - Auth Type: %s\n", c.AuthType))
	sb.WriteString(fmt.Sprintf("  - Auth Username: %s\n", c.AuthUsername))
	sb.WriteString(fmt.Sprintf("  - Auth Password: %s\n", utils.Mask(c.AuthPassword)))
	sb.WriteString(fmt.Sprintf("  - From: %s\n", c.From))
	sb.WriteString(fmt.Sprintf("  - Host: %s\n", c.Host))
	sb.WriteString(fmt.Sprintf("  - Port: %d\n", c.Port))

	return sb.String()
}
