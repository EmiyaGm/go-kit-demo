package service

import (
	"context"
	"errors"
	"strings"
	"fmt"
)

// Service describes a service that adds things together.
type Service interface {
	Uppercase(ctx context.Context, a string) (string, error)
	Create(ctx context.Context, ID string, FlowID uint32, Source string, Type string) (error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New() Service {
	var svc Service
	{
		svc = NewBasicService()
	}
	return svc
}

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}


func (s basicService) Uppercase(_ context.Context, a string) (string, error) {
	return strings.ToUpper(a), nil
}

func (s basicService) Create(_ context.Context, ID string, FlowID uint32, Source string, Type string) error{
	fmt.Print("get alarm data")
	return nil
}
