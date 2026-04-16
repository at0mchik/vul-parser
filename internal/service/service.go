package service

import (
	"vul-parser/internal/domain/dto"
	"vul-parser/internal/domain/models"
)

type AnalysisHTTP interface {
	Analyze(req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	AnalyzeWithFile(filePath string, req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	loadRulesFromJSON(data []byte) ([]models.Rule, error)
}

type AnalysisGRPC interface {
	Analyze(req *AnalyzeRequest) (*AnalyzeResult, error)
	AnalyzeFile(req *AnalyzeRequest) (*AnalyzeResult, error)
	Health() map[string]string
}

type Service struct {
	AnalysisHTTP
	AnalysisGRPC
}

func NewService() *Service {
	return &Service{
		AnalysisHTTP: NewAnalysisHTTPService(),
		AnalysisGRPC: NewAnalyzerGRPCService(),
	}
}
