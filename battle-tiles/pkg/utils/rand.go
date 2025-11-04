package utils

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

func GenerateUsername() string {
	return "user_" + uuid.New().String()[:6]
}
func GenerateNickName() string {
	return fmt.Sprintf("用户%d", rand.Intn(100000))
}
