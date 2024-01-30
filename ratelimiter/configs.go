package ratelimiter

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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

	// IP rate limiter
	IP := &RateLimiterRateConfig{
		MaxRequestsPerSecond:    defaultMaxRequests,
		BlockTimeInMilliseconds: defaultBlockTime,
	}

	// Token rate limiter
	// Get custom tokens from env that satisfy the regex
	envKeyRegex := regexp.MustCompile("^RATE_LIMITER_TOKEN_(.*)_(MAX_REQUESTS|BLOCK_TIME)$")
	foundTokens := map[string]bool{}

	envs := os.Environ()
	for _, env := range envs {
		envPair := strings.SplitN(env, "=", 2)
		envKey := envPair[0]
		if envKeyRegex.Match([]byte(envKey)) {
			foundTokens[envKeyRegex.FindStringSubmatch(envKey)[1]] = true
		}
	}

	tokens := []string{}
	for t := range foundTokens {
		tokens = append(tokens, t)
	}

	// Set custom tokens
	var customTokens map[string]*RateLimiterRateConfig
	for _, token := range tokens {
		maxRequestsPerSecondEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_MAX_REQUESTS", token)
		maxRequestsPerSecond := viper.GetInt64(maxRequestsPerSecondEnvKey)
		if maxRequestsPerSecond == 0 {
			maxRequestsPerSecond = defaultMaxRequests
		}

		blockTimeEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_BLOCK_TIME", token)
		blockTime := viper.GetInt64(blockTimeEnvKey)
		if blockTime == 0 {
			blockTime = defaultBlockTime
		}

		customTokens[token] = &RateLimiterRateConfig{
			MaxRequestsPerSecond:    maxRequestsPerSecond,
			BlockTimeInMilliseconds: blockTime,
		}
	}

	return &RateLimiterConfig{
		IP:           IP,
		CustomTokens: customTokens,
		Storage:      storage.NewRedisStorage(),
	}
}
