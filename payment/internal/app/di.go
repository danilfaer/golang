package app

import (
	"context"

	paymentV1API "github.com/danilfaer/golang/payment/internal/api/api/payment/v1"
	"github.com/danilfaer/golang/payment/internal/service"
	paymentService "github.com/danilfaer/golang/payment/internal/service/payment"
	paymentV1 "github.com/danilfaer/golang/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentGRPCAPI paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

// PaymentGRPCAPI — реализация gRPC PaymentService для регистрации на сервере.
func (d *diContainer) PaymentGRPCAPI(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentGRPCAPI == nil {
		d.paymentGRPCAPI = paymentV1API.NewAPI(d.PaymentService(ctx))
	}
	return d.paymentGRPCAPI
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
