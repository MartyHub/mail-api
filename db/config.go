package db

import "github.com/invopop/validation"

type Config struct {
	URL string `env:"URL"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.URL, validation.Required),
	)
}
