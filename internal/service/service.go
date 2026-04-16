package service

import (
	"vul-parser/internal/domain/dto"
	"vul-parser/internal/domain/models"
)

type AnalyzisHTTP interface {
	Analyze(req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	AnalyzeWithFile(filePath string, req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	loadRulesFromJSON(data []byte) ([]models.Rule, error)
}

type AnalyzisGRPC interface {
	Analyze(req *AnalyzeRequest) (*AnalyzeResult, error)
	AnalyzeFile(req *AnalyzeRequest) (*AnalyzeResult, error)
	Health() map[string]string
}

type Service struct {
	AnalyzisHTTP
	AnalyzisGRPC
}

func NewService() *Service {
	return &Service{
		AnalyzisHTTP: NewAnalysisHTTPService(),
		AnalyzisGRPC: NewAnalyzerGRPCService(),
	}
}
