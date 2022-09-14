package service

import "errors"

// HeavyComputation runs nohting and is absolutely not heavy.
func HeavyComputation() (string, error) {
	return "", errors.New("oh no")
}
