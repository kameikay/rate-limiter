package ratelimiter

import (
	"github.com/goccy/go-json"
	"github.com/kameikay/rate-limiter/ratelimiter/storage"
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

type RateLimiterRateConfig struct {
	MaxRequestsPerSecond    int64 `json:"max_requests_per_second"`
	BlockTimeInMilliseconds int64 `json:"block_time_in_milliseconds"`
}

type RateLimiterConfig struct {
	IP           *RateLimiterRateConfig              `json:"ip"`
	CustomTokens map[string]*RateLimiterRateConfig   `json:"custom_tokens"`
	Storage      storage.RateLimiterStorageInterface `json:"-"`
}

func (rlc *RateLimiterConfig) GetRateLimiterRateConfigForToken(token string) *RateLimiterRateConfig {
	customTokenConfig, ok := rlc.CustomTokens[token]
	if ok {
		return customTokenConfig
	} else {
		return rlc.IP
	}
}

func setConfig() *RateLimiterConfig {
	defaultMaxRequests := viper.GetInt64("RATE_LIMITER_DEFAULT_MAX_REQUESTS")
	defaultBlockTime := viper.GetInt64("RATE_LIMITER_DEFAULT_BLOCK_TIME")
	customTokensEnvs := viper.GetString("RATE_LIMITER_CUSTOM_TOKENS")

	if defaultMaxRequests == 0 {
		defaultMaxRequests = 5
	}

	if defaultBlockTime == 0 {
		defaultBlockTime = 500
	}

	if customTokensEnvs == "" {
		customTokensEnvs = `[]`
	}

	// IP rate limiter
	IP := &RateLimiterRateConfig{
		MaxRequestsPerSecond:    defaultMaxRequests,
		BlockTimeInMilliseconds: defaultBlockTime,
	}

	// Token rate limiter
	type token struct {
		Name                    string `json:"name"`
		MaxRequestsPerSecond    int64  `json:"max_requests_per_second"`
		BlockTimeInMilliseconds int64  `json:"block_time_in_milliseconds"`
	}
	var tokens []token
	err := json.Unmarshal([]byte(customTokensEnvs), &tokens)
	if err != nil {
		panic(err)
	}

	// Set custom tokens
	var customTokens = make(map[string]*RateLimiterRateConfig)
	for _, tkn := range tokens {
		if tkn.MaxRequestsPerSecond == 0 {
			tkn.MaxRequestsPerSecond = defaultMaxRequests
		}

		if tkn.BlockTimeInMilliseconds == 0 {
			tkn.BlockTimeInMilliseconds = defaultBlockTime
		}

		customTokens[tkn.Name] = &RateLimiterRateConfig{
			MaxRequestsPerSecond:    tkn.MaxRequestsPerSecond,
			BlockTimeInMilliseconds: tkn.BlockTimeInMilliseconds,
		}
	}

	return &RateLimiterConfig{
		IP:           IP,
		CustomTokens: customTokens,
		Storage:      storage.NewRedisStorage(),
	}
}
