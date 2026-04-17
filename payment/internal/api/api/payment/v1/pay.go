package v1

import (
	"context"

	"github.com/danilfaer/golang/payment/internal/converter"
	"github.com/danilfaer/golang/payment/internal/model"
	paymentV1 "github.com/danilfaer/golang/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	payment := converter.ConvertFromGRPC(req)
	transactionUUID, err := a.paymentService.PayOrder(ctx, payment)

	if err != nil {
		return nil, model.ErrPayment
	}

	// Возвращаем gRPC ответ
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
