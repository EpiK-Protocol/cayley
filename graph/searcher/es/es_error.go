package es

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type EsError struct {
	Reason string
	Status float64
}

func (e *EsError) Error() string {
	return fmt.Sprintf(`{"status: %f, "reason": "%s"}`, e.Status, e.Reason)
}
