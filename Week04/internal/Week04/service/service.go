package service

import (
	"Week04/api/order"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
}

var (
	responses = []order.OrderResponse{
		{
			Id:           1,
			OrderNo:      "order1",
			CreateTime:   timestamppb.Now(),
			CreateUserId: 1,
		},
		{
			Id:           2,
			OrderNo:      "order2",
			CreateTime:   timestamppb.Now(),
			CreateUserId: 1,
		},
		{
			Id:           3,
			OrderNo:      "order3",
			CreateTime:   timestamppb.Now(),
			CreateUserId: 2,
		},
	}
)

func (svc *OrderService) GetOrderByCreateUserId(ctx context.Context, req *order.OrderRequest) (*order.OrderListResponse, error) {
	createUserId := req.GetCreateUserId()
	var results []*order.OrderResponse
	for _, response := range responses {
		if response.GetCreateUserId() == createUserId {
			results = append(results, &response)
		}
	}
	return &order.OrderListResponse{
		OrderResponse: results,
	}, nil
}

func (svc *OrderService) GetOrderByNo(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	id := req.GetOrderNo()
	for _, response := range responses {
		if response.GetOrderNo() == id {
			return &response, nil
		}
	}
	return nil, errors.New("没有找到订单")
}

func (svc *OrderService) GetOrderById(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	id := req.GetId()
	for _, response := range responses {
		if response.GetId() == id {
			return &response, nil
		}
	}
	return nil, errors.New("没有找到订单")
}
