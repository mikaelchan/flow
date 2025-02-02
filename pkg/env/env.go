package env

import "os"

func IsRelease() bool {
	return os.Getenv("APP_MODE") == "release"
}
