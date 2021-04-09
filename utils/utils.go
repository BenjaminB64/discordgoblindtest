package utils

import "os"

func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
