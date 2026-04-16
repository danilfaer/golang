package model

import (
	"errors"

	sharedErrors "github.com/danilfaer/golang/shared/pkg/errors"
)

var (
	ErrPartNotFound = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrInvalidUUID  = sharedErrors.NewInvalidArgumentError(errors.New("invalid uuid"))
)
