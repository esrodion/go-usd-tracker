package controller

import (
	"context"
	"go-usdtrub/internal/models"
	pb "go-usdtrub/pkg/grpc"
	"log"
	"net"

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

	return c, err
}

func (c *GrpcController) GetRates(ctx context.Context, p *pb.GetRatesParams) (*pb.Rate, error) {
	rate, err := c.service.GetRates(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	server := grpc.NewServer()
	pb.RegisterRatesServer(server, c)

	// TODO: Zap logging
	log.Println("Geo gRPC server started at " + hostName)
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return server, nil
}

func (c *GrpcController) Stop() {
	c.server.GracefulStop()
}
