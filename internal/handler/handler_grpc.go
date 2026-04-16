package handler

import (
	"vul-parser/gen/proto/analyzer"
	"vul-parser/internal/service"
)

type HandlerGRPC struct {
	analyzer.UnimplementedAnalyzerServiceServer
	service *service.Service
}

func NewHandlerGRPC(service *service.Service) *HandlerGRPC {
	return &HandlerGRPC{
		service: service,
	}
}
