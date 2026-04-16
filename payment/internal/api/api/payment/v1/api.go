package v1

import (
	"github.com/danilfaer/golang/payment/internal/service"
	paymentV1 "github.com/danilfaer/golang/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{paymentService: paymentService}
}
