package include

import "github.com/google/uuid"

func createHash() string {
	return uuid.New().String()
}
