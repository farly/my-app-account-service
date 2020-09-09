package accounts

import (
	models "accounts/models"

	"github.com/go-redis/redis"
)

type Context struct {
	UserModel *models.UserModel
	Rdb       *redis.Client
}

// leave this for now
type Success struct {
	Ok bool `json:"ok"`
}

// leave it for now
type Fail struct {
	Ok     bool              `json:"ok"`
	Errors map[string]string `json:"errors,omitempty"`
}
