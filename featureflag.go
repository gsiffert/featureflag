// Package featureflag implements a modular feature flag client,
// each feature-flag can define its own source and refresh interval.
//
// There are multiple packages intended to cover common feature-flag servers:
// - github.com/gsiffert/ffaws/ffsecretmanager: AWS SecretManager
// - github.com/gsiffert/ffaws/ffappconfig: AWS AppConfig
//
// If unfortunately your Source is not implemented,
// you can implement your own feature-flag source by implementing the SourceReader interface.
package featureflag

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FeatureFlag define a property which should be refreshed in the background at a given interval.
// Error while refreshing the value will be logged using the logger.
type FeatureFlag[T any] interface {
	// Value returns the current value of the feature flag.
	// This function is safe to be called concurrently.
	Value() T
}

// New creates a new FeatureFlag with the given refresh interval and source.
// If the initial download fails, the function will return a nil FeatureFlag and an error.
// Otherwise, the FeatureFlag will be refreshed in the background at the given refreshInterval.
// Canceling the context will stop the FeatureFlag from refreshing its value.
func New[T any](ctx context.Context, refreshInterval time.Duration, source SourceReader[T]) (FeatureFlag[T], error) {
	f := &featureFlag[T]{
		refreshInterval: refreshInterval,
		source:          source,
	}
	err := f.refresh(ctx)
	if err != nil {
		return nil, err
	}

	go f.run(ctx)
	return f, nil
}

type featureFlag[T any] struct {
	refreshInterval time.Duration
	source          SourceReader[T]
	value           T
	m               sync.RWMutex
}

func (f *featureFlag[T]) Value() T {
	f.m.RLock()
	defer f.m.RUnlock()
	return f.value
}

func (f *featureFlag[T]) setValue(value T) {
	f.m.Lock()
	defer f.m.Unlock()
	f.value = value
}

func (f *featureFlag[T]) refresh(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorContext(ctx, "Panic while refreshing the feature flag: %v.", r)
		}
	}()

	value, err := f.source.Read(ctx)
	if err != nil {
		return fmt.Errorf("retrieve feature-flag from source: %w", err)
	}

	f.setValue(value)
	return nil
}

func (f *featureFlag[T]) run(ctx context.Context) {
	ticker := time.NewTicker(f.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := f.refresh(ctx); err != nil {
				logger.ErrorContext(ctx, "Refreshing the feature flag failed because: %v.", err)
			}
		}
	}
}
