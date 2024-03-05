package featureflag

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type mapReader[A, B any] struct {
	source SourceReader[A]
	f      func(A) (B, error)
}

func (m *mapReader[A, B]) Read(ctx context.Context) (B, error) {
	var b B

	value, err := m.source.Read(ctx)
	if err != nil {
		return b, fmt.Errorf("inner source: %w", err)
	}

	return m.f(value)
}

// MapFunc is a function that converts a value of type A to a value of type B.
type MapFunc[A, B any] func(A) (B, error)

// MapSourceReader helps to convert the result from one source to another, for example, from a JSON file to a struct.
func MapSourceReader[A, B any](s SourceReader[A], f MapFunc[A, B]) SourceReader[B] {
	return &mapReader[A, B]{source: s, f: f}
}

// MapJSON is a MapFunc to be used with MapSourceReader to convert from a JSON to a struct.
func MapJSON[T any](value []byte) (T, error) {
	var t T
	if err := json.Unmarshal(value, &t); err != nil {
		return t, fmt.Errorf("json unmarshal: %w", err)
	}
	return t, nil
}

// MapXML is a MapFunc to be used with MapSourceReader to convert from an XML to a struct.
func MapXML[T any](value []byte) (T, error) {
	var t T
	if err := xml.Unmarshal(value, &t); err != nil {
		return t, fmt.Errorf("xml unmarshal: %w", err)
	}
	return t, nil
}
