package item

import (
	"context"
	"errors"

	pb "go-sample-http/gen/item"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler는 HTTP handler와 동일한 UseCase를 재사용합니다.
type GRPCHandler struct {
	pb.UnimplementedItemServiceServer
	uc UseCase
}

func NewGRPCHandler(uc UseCase) *GRPCHandler {
	return &GRPCHandler{uc: uc}
}

func (h *GRPCHandler) ListItems(ctx context.Context, _ *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	items, err := h.uc.List()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.ListItemsResponse{}
	for _, item := range items {
		resp.Items = append(resp.Items, toProto(item))
	}
	return resp, nil
}

func (h *GRPCHandler) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.ItemResponse, error) {
	item, err := h.uc.Get(req.Id)
	if errors.Is(err, ErrNotFound) {
		return nil, status.Error(codes.NotFound, "item not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProto(*item), nil
}

func (h *GRPCHandler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.ItemResponse, error) {
	item, err := h.uc.Create(CreateRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProto(*item), nil
}

func (h *GRPCHandler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.ItemResponse, error) {
	item, err := h.uc.Update(req.Id, UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if errors.Is(err, ErrNotFound) {
		return nil, status.Error(codes.NotFound, "item not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProto(*item), nil
}

func (h *GRPCHandler) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	if err := h.uc.Delete(req.Id); errors.Is(err, ErrNotFound) {
		return nil, status.Error(codes.NotFound, "item not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteItemResponse{}, nil
}

func toProto(item Item) *pb.ItemResponse {
	return &pb.ItemResponse{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.UTC().String(),
	}
}
