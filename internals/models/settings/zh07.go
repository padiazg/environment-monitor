package settings

import (
	"fmt"
	"os"
)

type ZH07Settings struct {
	Mode string
	Port string
}

func (z *ZH07Settings) Validate() error {
	if z == nil {
		return fmt.Errorf("ZH7 settings not set")
	}

	if z.Mode != "qa" && z.Mode != "initiative" {
		return fmt.Errorf("invalid mode '%s': must be 'qa' or 'initiative'", z.Mode)
	}

	if z.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	// Check if port exists (basic validation)
	if _, err := os.Stat(z.Port); os.IsNotExist(err) {
		return fmt.Errorf("port '%s' does not exist", z.Port)
	}

	return nil
}
