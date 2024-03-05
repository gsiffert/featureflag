package featureflag

import (
	"context"
)

// SourceReader defines the interface used to read the Value from the Source of the FeatureFlag.
type SourceReader[T any] interface {
	Read(ctx context.Context) (T, error)
}
