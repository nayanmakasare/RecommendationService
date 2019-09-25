// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: RecommendationService.proto

package RecommendationService

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for RecommendationService service

type RecommendationService interface {
	TileClicked(ctx context.Context, in *TileClickedRequest, opts ...client.CallOption) (*TileClickedResponse, error)
	GetCollabrativeFilteringData(ctx context.Context, in *GetRecommendationRequest, opts ...client.CallOption) (RecommendationService_GetCollabrativeFilteringDataService, error)
	GetContentbasedData(ctx context.Context, in *GetRecommendationRequest, opts ...client.CallOption) (RecommendationService_GetContentbasedDataService, error)
	InitialRecommendationEngine(ctx context.Context, in *InitRecommendationRequest, opts ...client.CallOption) (*InitRecommendationResponse, error)
}

type recommendationService struct {
	c    client.Client
	name string
}

func NewRecommendationService(name string, c client.Client) RecommendationService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "RecommendationService"
	}
	return &recommendationService{
		c:    c,
		name: name,
	}
}

func (c *recommendationService) TileClicked(ctx context.Context, in *TileClickedRequest, opts ...client.CallOption) (*TileClickedResponse, error) {
	req := c.c.NewRequest(c.name, "RecommendationService.TileClicked", in)
	out := new(TileClickedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recommendationService) GetCollabrativeFilteringData(ctx context.Context, in *GetRecommendationRequest, opts ...client.CallOption) (RecommendationService_GetCollabrativeFilteringDataService, error) {
	req := c.c.NewRequest(c.name, "RecommendationService.GetCollabrativeFilteringData", &GetRecommendationRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &recommendationServiceGetCollabrativeFilteringData{stream}, nil
}

type RecommendationService_GetCollabrativeFilteringDataService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*MovieTile, error)
}

type recommendationServiceGetCollabrativeFilteringData struct {
	stream client.Stream
}

func (x *recommendationServiceGetCollabrativeFilteringData) Close() error {
	return x.stream.Close()
}

func (x *recommendationServiceGetCollabrativeFilteringData) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *recommendationServiceGetCollabrativeFilteringData) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *recommendationServiceGetCollabrativeFilteringData) Recv() (*MovieTile, error) {
	m := new(MovieTile)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *recommendationService) GetContentbasedData(ctx context.Context, in *GetRecommendationRequest, opts ...client.CallOption) (RecommendationService_GetContentbasedDataService, error) {
	req := c.c.NewRequest(c.name, "RecommendationService.GetContentbasedData", &GetRecommendationRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &recommendationServiceGetContentbasedData{stream}, nil
}

type RecommendationService_GetContentbasedDataService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*MovieTile, error)
}

type recommendationServiceGetContentbasedData struct {
	stream client.Stream
}

func (x *recommendationServiceGetContentbasedData) Close() error {
	return x.stream.Close()
}

func (x *recommendationServiceGetContentbasedData) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *recommendationServiceGetContentbasedData) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *recommendationServiceGetContentbasedData) Recv() (*MovieTile, error) {
	m := new(MovieTile)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *recommendationService) InitialRecommendationEngine(ctx context.Context, in *InitRecommendationRequest, opts ...client.CallOption) (*InitRecommendationResponse, error) {
	req := c.c.NewRequest(c.name, "RecommendationService.InitialRecommendationEngine", in)
	out := new(InitRecommendationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RecommendationService service

type RecommendationServiceHandler interface {
	TileClicked(context.Context, *TileClickedRequest, *TileClickedResponse) error
	GetCollabrativeFilteringData(context.Context, *GetRecommendationRequest, RecommendationService_GetCollabrativeFilteringDataStream) error
	GetContentbasedData(context.Context, *GetRecommendationRequest, RecommendationService_GetContentbasedDataStream) error
	InitialRecommendationEngine(context.Context, *InitRecommendationRequest, *InitRecommendationResponse) error
}

func RegisterRecommendationServiceHandler(s server.Server, hdlr RecommendationServiceHandler, opts ...server.HandlerOption) error {
	type recommendationService interface {
		TileClicked(ctx context.Context, in *TileClickedRequest, out *TileClickedResponse) error
		GetCollabrativeFilteringData(ctx context.Context, stream server.Stream) error
		GetContentbasedData(ctx context.Context, stream server.Stream) error
		InitialRecommendationEngine(ctx context.Context, in *InitRecommendationRequest, out *InitRecommendationResponse) error
	}
	type RecommendationService struct {
		recommendationService
	}
	h := &recommendationServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&RecommendationService{h}, opts...))
}

type recommendationServiceHandler struct {
	RecommendationServiceHandler
}

func (h *recommendationServiceHandler) TileClicked(ctx context.Context, in *TileClickedRequest, out *TileClickedResponse) error {
	return h.RecommendationServiceHandler.TileClicked(ctx, in, out)
}

func (h *recommendationServiceHandler) GetCollabrativeFilteringData(ctx context.Context, stream server.Stream) error {
	m := new(GetRecommendationRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.RecommendationServiceHandler.GetCollabrativeFilteringData(ctx, m, &recommendationServiceGetCollabrativeFilteringDataStream{stream})
}

type RecommendationService_GetCollabrativeFilteringDataStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*MovieTile) error
}

type recommendationServiceGetCollabrativeFilteringDataStream struct {
	stream server.Stream
}

func (x *recommendationServiceGetCollabrativeFilteringDataStream) Close() error {
	return x.stream.Close()
}

func (x *recommendationServiceGetCollabrativeFilteringDataStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *recommendationServiceGetCollabrativeFilteringDataStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *recommendationServiceGetCollabrativeFilteringDataStream) Send(m *MovieTile) error {
	return x.stream.Send(m)
}

func (h *recommendationServiceHandler) GetContentbasedData(ctx context.Context, stream server.Stream) error {
	m := new(GetRecommendationRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.RecommendationServiceHandler.GetContentbasedData(ctx, m, &recommendationServiceGetContentbasedDataStream{stream})
}

type RecommendationService_GetContentbasedDataStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*MovieTile) error
}

type recommendationServiceGetContentbasedDataStream struct {
	stream server.Stream
}

func (x *recommendationServiceGetContentbasedDataStream) Close() error {
	return x.stream.Close()
}

func (x *recommendationServiceGetContentbasedDataStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *recommendationServiceGetContentbasedDataStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *recommendationServiceGetContentbasedDataStream) Send(m *MovieTile) error {
	return x.stream.Send(m)
}

func (h *recommendationServiceHandler) InitialRecommendationEngine(ctx context.Context, in *InitRecommendationRequest, out *InitRecommendationResponse) error {
	return h.RecommendationServiceHandler.InitialRecommendationEngine(ctx, in, out)
}
