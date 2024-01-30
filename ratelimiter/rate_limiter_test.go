package ratelimiter

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/kameikay/rate-limiter/ratelimiter/mocks"
	"github.com/stretchr/testify/suite"
)

type RateLimiterTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	ctx         context.Context
	storageMock *mock.MockRateLimiterStorageInterface
}

func TestCreateItemInputUseCaseStart(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

func (suite *RateLimiterTestSuite) RateLimiterTestSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *RateLimiterTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.storageMock = mock.NewMockRateLimiterStorageInterface(suite.ctrl)
}

func (suite *RateLimiterTestSuite) TestSuccess() {
	testCases := []struct {
		name          string
		key           string
		config        *RateLimiterConfig
		expectedError error
		expectations  func(storage *mock.MockRateLimiterStorageInterface)
	}{
		{
			name: "should allow request",
			key:  "127.0.0.1",
			config: &RateLimiterConfig{
				IP: &RateLimiterRateConfig{
					MaxRequestsPerSecond:    5,
					BlockTimeInMilliseconds: 500,
				},
			},
			expectedError: nil,
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, nil).Times(1)
				storage.EXPECT().Increment(suite.ctx, "127.0.0.1", 5).Return(true, nil).Times(1)
			},
		},
	}

	for _, testCase := range testCases {
		testCase.expectations(suite.storageMock)
		_, err := checkRateLimit(suite.ctx, testCase.key, testCase.config, testCase.config.IP)
		suite.Equal(testCase.expectedError, err, testCase.name)
	}
}
