package ratelimiter

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kameikay/rate-limiter/ratelimiter/storage"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	ctx  context.Context
}

func TestMiddlewareStart(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (suite *MiddlewareTestSuite) MiddlewareTestSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.ctx = context.Background()
}

func (suite *MiddlewareTestSuite) TestNewRateLimiter() {
	middleware := NewRateLimiter()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	suite.NotNil(middleware)
	suite.NotNil(middleware(handler))
}

func (suite *MiddlewareTestSuite) TestRateLimiter() {
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:    5,
			BlockTimeInMilliseconds: 500,
		},
	}

	testCases := []struct {
		name                     string
		rateLimiterCheckFunction rateLimiterCheckFunction
		statusCode               int
		message                  []byte
		handler                  http.HandlerFunc
		token                    string
	}{
		{
			name: "should allow request when no token is provided",
			rateLimiterCheckFunction: func(ctx context.Context, key string, storage storage.RateLimiterStorageInterface, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
				return nil, nil
			},
			statusCode: http.StatusOK,
			message:    []byte("OK"),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
		{
			name: "should not allow request when rate limit is reached",
			rateLimiterCheckFunction: func(ctx context.Context, key string, storage storage.RateLimiterStorageInterface, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
				block := time.Now().Add(time.Second * 5)
				return &block, nil
			},
			statusCode: http.StatusTooManyRequests,
			message:    []byte("you have reached the maximum number of requests or actions allowed within a certain time frame"),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
		{
			name: "should allow request when token is provided",
			rateLimiterCheckFunction: func(ctx context.Context, key string, storage storage.RateLimiterStorageInterface, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
				return nil, nil
			},
			statusCode: http.StatusOK,
			message:    []byte("OK"),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			token: "test-token",
		},
		{
			name: "should return error on checkRateLimitFunc",
			rateLimiterCheckFunction: func(ctx context.Context, key string, storage storage.RateLimiterStorageInterface, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
				return nil, errors.New("error")
			},
			statusCode: http.StatusInternalServerError,
			message:    []byte("error"),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			request := httptest.NewRequest("GET", "http://test", nil)
			recorder := httptest.NewRecorder()
			if tc.token != "" {
				request.Header.Set("API_KEY", tc.token)
			}

			rateLimiter(config, tc.handler, tc.rateLimiterCheckFunction).ServeHTTP(recorder, request)

			response := recorder.Result()
			responseStatusCode := response.StatusCode
			responseBody, err := io.ReadAll(response.Body)

			suite.Equal(tc.statusCode, responseStatusCode)
			suite.Equal(tc.message, responseBody)
			suite.Nil(err)
		})
	}
}
