package tools


import (
	"github.com/google/uuid"
)

// 生成唯一的 UUID
func GetUuid() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}