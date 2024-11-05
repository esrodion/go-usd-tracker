package controller

import (
	"context"
	"fmt"
	"go-usdtrub/internal/metrics"
	"go-usdtrub/internal/models"
	"go-usdtrub/internal/traces"
	pb "go-usdtrub/pkg/grpc"
	"go-usdtrub/pkg/logger"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Service interface {
	GetRates(ctx context.Context) (models.CurrenceyRate, error)
}

type GrpcController struct {
	pb.UnimplementedRatesServer
	service Service
	server  *grpc.Server
}

func NewGrpcController(service Service, hostName string) (*GrpcController, error) {
	var err error
	c := &GrpcController{}

	c.service = service
	c.server, err = c.start(hostName)

	if err != nil {
		return nil, fmt.Errorf("controller.NewGrpcController: %w", err)
	}

	return c, nil
}

func (c *GrpcController) GetRates(ctx context.Context, p *pb.GetRatesParams) (*pb.Rate, error) {
	ctx = metrics.GRPCMethodHead(ctx, "GetRates")
	defer metrics.GRPCMethodTail(ctx, "GetRates", time.Now())

	ctx, span := traces.Start(ctx, "GRPCGetRates")
	defer span.End()

	rate, err := c.service.GetRates(ctx)
	if err != nil {
		return nil, fmt.Errorf("controller.GrpcController.GetRates: %w", err)
	}

	result := &pb.Rate{}
	result.Ask = rate.Ask
	result.Bid = rate.Bid
	result.Timestamp = rate.Timestamp.UnixMicro()

	return result, nil
}

func (c *GrpcController) start(hostName string) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", hostName)
	if err != nil {
		return nil, fmt.Errorf("controller.GrpcController.start: %w", err)
	}

	log := logger.Logger().Sugar().Named("gRPC")

	server := grpc.NewServer()
	pb.RegisterRatesServer(server, c)

	log.Info("gRPC server started at ", hostName)
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.DPanic(err)
		}
	}()

	return server, nil
}

func (c *GrpcController) Stop() {
	c.server.GracefulStop()
}
