package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestRateLimiterStart(t *testing.T) {
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
		rateConfig    *RateLimiterRateConfig
		expectedError error
		expectations  func(storage *mock.MockRateLimiterStorageInterface)
	}{
		{
			name: "should allow request",
			key:  "127.0.0.1",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: nil,
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, nil).Times(1)
				storage.EXPECT().Increment(suite.ctx, "127.0.0.1", int64(5)).Return(true, nil).Times(1)
			},
		},
		{
			name: "should allow request when has no key",
			key:  "",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: nil,
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "").Times(0)
			},
		},
		{
			name: "should return error on storage get block",
			key:  "127.0.0.1",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: errors.New("error"),
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, errors.New("error")).Times(1)
			},
		},
		{
			name: "should return error on storage increment",
			key:  "127.0.0.1",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: errors.New("error"),
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, nil).Times(1)
				storage.EXPECT().Increment(suite.ctx, "127.0.0.1", int64(5)).Return(false, errors.New("error")).Times(1)
			},
		},
		{
			name: "should add block when increment return false on success",
			key:  "127.0.0.1",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: nil,
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, nil).Times(1)
				storage.EXPECT().Increment(suite.ctx, "127.0.0.1", int64(5)).Return(false, nil).Times(1)

				expiration := time.Duration(int64(time.Millisecond) * int64(500))
				blockedTimeReturn := time.Now().Add(expiration)
				storage.EXPECT().AddBlock(suite.ctx, "127.0.0.1", int64(500)).Return(&blockedTimeReturn, nil).Times(1)
			},
		},
		{
			name: "should return error on add block",
			key:  "127.0.0.1",
			rateConfig: &RateLimiterRateConfig{
				MaxRequestsPerSecond:    int64(5),
				BlockTimeInMilliseconds: int64(500),
			},
			expectedError: errors.New("error"),
			expectations: func(storage *mock.MockRateLimiterStorageInterface) {
				storage.EXPECT().GetBlock(suite.ctx, "127.0.0.1").Return(nil, nil).Times(1)
				storage.EXPECT().Increment(suite.ctx, "127.0.0.1", int64(5)).Return(false, nil).Times(1)
				storage.EXPECT().AddBlock(suite.ctx, "127.0.0.1", int64(500)).Return(nil, errors.New("error")).Times(1)
			},
		},
	}

	for _, testCase := range testCases {
		testCase.expectations(suite.storageMock)
		_, err := checkRateLimit(suite.ctx, testCase.key, suite.storageMock, testCase.rateConfig)
		suite.Equal(testCase.expectedError, err, testCase.name)
	}
}
