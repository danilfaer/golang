package model

import (
	"errors"

	sharedErrors "github.com/danilfaer/golang/shared/pkg/errors"
)

var ErrPayment = sharedErrors.NewInvalidArgumentError(errors.New("payment error"))
