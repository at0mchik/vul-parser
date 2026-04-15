package service

import (
	"vul-parser/internal/domain/dto"
	"vul-parser/internal/domain/models"
)

type Analyzis interface {
	Analyze(req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	AnalyzeWithFile(filePath string, req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error)
	loadRulesFromJSON(data []byte) ([]models.Rule, error)
}

type Service struct {
	Analyzis
}

func NewService() *Service {
	return &Service{
		Analyzis: NewAnalysisService(),
	}
}
