package handler

import (
	"time"

	"github.com/alexliesenfeld/health"
)

func HealthCheck() health.Checker {
	checker := health.NewChecker(
		health.WithTimeout(10 * time.Second),
	)
	return checker
}
