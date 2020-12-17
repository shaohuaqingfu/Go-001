package service

import (
	"Week04/api/order"
	"Week04/internal/Week04/data"
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	Dao *data.OrderData
}

func (svc *OrderService) GetOrderByCreateUserId(ctx context.Context, req *order.OrderRequest) (*order.OrderListResponse, error) {
	createUserId := req.GetCreateUserId()
	orders, err := svc.Dao.GetByCreateUserId(createUserId)
	if err != nil {
		return nil, err
	}
	var orderResponses []*order.OrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, &order.OrderResponse{
			Id:           o.Id,
			OrderNo:      o.OrderNo,
			CreateTime:   timestamppb.New(o.CreateTime),
			CreateUserId: o.CreateUserId,
		})
	}
	return &order.OrderListResponse{
		OrderResponse: orderResponses,
	}, nil
}

func (svc *OrderService) GetOrderByNo(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	//id := req.GetOrderNo()
	//for _, response := range responses {
	//	if response.GetOrderNo() == id {
	//		return &response, nil
	//	}
	//}
	return nil, errors.New("没有找到订单")
}

func (svc *OrderService) GetOrderById(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	id := req.GetId()
	o, err := svc.Dao.GetById(id)
	if err != nil {
		fmt.Println(err)
		if svc.Dao.IsNotExists(err) {
			return nil, errors.New("没有找到对应订单")
		}
		return nil, err
	}
	orderResponse := &order.OrderResponse{
		Id:           o.Id,
		OrderNo:      o.OrderNo,
		CreateTime:   timestamppb.New(o.CreateTime),
		CreateUserId: o.CreateUserId,
	}
	return orderResponse, nil
}
