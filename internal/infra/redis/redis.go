package redis

import (
	"github.com/spf13/viper"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

type RedisUserSessionDTO struct {
}

func NewRedis() *Client {
	redisAddr := viper.GetString("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       0,
		Password: "",
	})

	return &Client{
		client: client,
	}
}

// func (c *Client) SetUser(ctx context.Context, user *entity.User, expiresIn time.Duration) error {
// 	jsonBytes, err := json.Marshal(user)
// 	if err != nil {
// 		return fmt.Errorf("error on encode token")
// 	}
// 	if err := c.client.Set(ctx, "users", string(jsonBytes), expiresIn).Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Client) GetUser(ctx context.Context, userId string) (entity.User, error) {
// 	jsonString := c.client.Get(ctx, "users").Val()

// 	var userDTO entity.User
// 	err := json.Unmarshal([]byte(jsonString), &userDTO)
// 	if err != nil {
// 		return entity.User{}, err
// 	}
// 	return userDTO, nil
// }
