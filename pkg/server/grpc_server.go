package server

import (
	"net"

	"vul-parser/gen/proto/analyzer"
	"vul-parser/internal/handler"
	"vul-parser/internal/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerGRPC struct {
	grpcServer *grpc.Server
	port       string
}

func NewServer(port string) (*ServerGRPC, error) {
	service := service.NewService()


	handler := handler.NewHandlerGRPC(service)

	grpcServer := grpc.NewServer()
	analyzer.RegisterAnalyzerServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return &ServerGRPC{
		grpcServer: grpcServer,
		port:       port,
	}, nil
}

func (s *ServerGRPC) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	logrus.Infof("gRPC server starting on port %s", s.port)
	return s.grpcServer.Serve(lis)
}

func (s *ServerGRPC) Stop() {
	logrus.Infof("Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
}