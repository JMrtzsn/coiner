package export
import (
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
}

// TODO should be loaded from a central location?
// LoadEnv loads all possible env variables
func (c *Config) LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
