package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/kameikay/rate-limiter/ratelimiter/storage"
)

func checkRateLimit(
	ctx context.Context,
	key string,
	storage storage.RateLimiterStorageInterface,
	rateConfig *RateLimiterRateConfig,
) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := storage.GetBlock(ctx, key)
	if err != nil {
		return nil, err
	}

	if block == nil {
		success, err := storage.Increment(ctx, key, rateConfig.MaxRequestsPerSecond)
		if err != nil {
			return nil, err
		}

		if !success {
			fmt.Println(fmt.Sprintf("block: %d", rateConfig.BlockTimeInMilliseconds))
			block, err = storage.AddBlock(ctx, key, rateConfig.BlockTimeInMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		fmt.Printf("block time %.2f seconds\n", block.Sub(time.Now()).Seconds())
		return block, nil
	}

	return nil, nil
}
