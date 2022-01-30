package support

import (
	"os"
)

func WebServerPort() string {
	return os.Getenv("PORT")
}
