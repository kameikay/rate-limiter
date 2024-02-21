package ratelimiter

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
}

func TestConfigsStart(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	viper.Reset()
}

func (suite *ConfigTestSuite) ConfigTestSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	tmpfile, err := os.Create(".env")
	if err != nil {
		suite.T().Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	text := []byte("RATE_LIMITER_DEFAULT_MAX_REQUESTS=10\nRATE_LIMITER_DEFAULT_BLOCK_TIME=500")
	if _, err := tmpfile.Write(text); err != nil {
		suite.T().Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		suite.T().Fatal(err)
	}

	err = LoadConfig(".")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(10), viper.GetInt64("RATE_LIMITER_DEFAULT_MAX_REQUESTS"))
	assert.Equal(suite.T(), int64(500), viper.GetInt64("RATE_LIMITER_DEFAULT_BLOCK_TIME"))
}

func (suite *ConfigTestSuite) TestGetRateLimiterRateConfigForToken() {
	testCases := []struct {
		name     string
		expected *RateLimiterRateConfig
		token    string
	}{
		{
			name: "should return custom token config",
			expected: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    10,
				BlockTimeInMilliseconds: 1000,
			},
			token: "token",
		},
		{
			name: "should return default config when no custom token is provided",
			expected: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    5,
				BlockTimeInMilliseconds: 500,
			},
			token: "",
		},
		{
			name: "should return default config when wrong token is provided",
			expected: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    5,
				BlockTimeInMilliseconds: 500,
			},
			token: "wrong-token",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			rateLimiterConfig := &RateLimiterConfig{
				IP: &RateLimiterRateConfig{
					MaxRequestsPerSecond:    5,
					BlockTimeInMilliseconds: 500,
				},
				CustomTokens: map[string]*RateLimiterRateConfig{
					"token": {
						MaxRequestsPerSecond:    10,
						BlockTimeInMilliseconds: 1000,
					},
				},
			}

			result := rateLimiterConfig.GetRateLimiterRateConfigForToken(tc.token)
			suite.Equal(tc.expected, result)
		})
	}
}

func (suite *ConfigTestSuite) TestSetConfig() {
	testCases := []struct {
		name                             string
		envRateLimiterDefaultMaxRequests int64
		envRateLimiterDefaultBlockTime   int64
		envRateLimiterCustomTokens       string
		expected                         *RateLimiterConfig
	}{
		{
			name: "should use default values when no env is provided",
			expected: &RateLimiterConfig{
				IP: &RateLimiterRateConfig{
					MaxRequestsPerSecond:    5,
					BlockTimeInMilliseconds: 500,
				},
				CustomTokens: map[string]*RateLimiterRateConfig{},
			},
		},
		{
			name: "should use env values when provided",
			expected: &RateLimiterConfig{
				IP: &RateLimiterRateConfig{
					MaxRequestsPerSecond:    15,
					BlockTimeInMilliseconds: 1500,
				},
				CustomTokens: map[string]*RateLimiterRateConfig{
					"ABC123": {
						MaxRequestsPerSecond:    1,
						BlockTimeInMilliseconds: 10000,
					},
					"DEF321": {
						MaxRequestsPerSecond:    15,
						BlockTimeInMilliseconds: 1500,
					},
				},
			},
			envRateLimiterDefaultMaxRequests: 15,
			envRateLimiterDefaultBlockTime:   1500,
			envRateLimiterCustomTokens:       `[{"name":"ABC123","max_requests_per_second":1,"block_time_in_milliseconds":10000},{"name":"DEF321","max_requests_per_second":0,"block_time_in_milliseconds":0}]`,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			viper.Set("RATE_LIMITER_DEFAULT_MAX_REQUESTS", tc.envRateLimiterDefaultMaxRequests)
			viper.Set("RATE_LIMITER_DEFAULT_BLOCK_TIME", tc.envRateLimiterDefaultBlockTime)
			viper.Set("RATE_LIMITER_CUSTOM_TOKENS", tc.envRateLimiterCustomTokens)
			result := setConfig()
			// We don't want to compare the storage field
			result.Storage = nil
			suite.Equal(tc.expected, result)
		})
	}
}
