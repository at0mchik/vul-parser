package handler

import (
	"context"
	"fmt"

	"vul-parser/gen/proto/analyzer"
	"vul-parser/internal/service"

	"github.com/sirupsen/logrus"
)

func (h *HandlerGRPC) Analyze(ctx context.Context, req *analyzer.AnalyzeRequest) (*analyzer.AnalyzeResponse, error) {
	logrus.Info("Received Analyze request")
	logrus.Infof("Config is nil: %v", req.Config == nil)

	// if req.Config != nil {
	// 	configJSON, _ := json.Marshal(req.Config.AsMap())
	// 	logrus.Printf("Config: %s", string(configJSON))
	// }

	// if req.Rules != nil {
	// 	rulesJSON, _ := json.Marshal(req.Rules.AsMap())
	// 	logrus.Printf("Rules: %s", string(rulesJSON))
	// }

	// if req.Config == nil {
	// 	logrus.Printf("ERROR: Config is required but was nil")
	// 	return nil, fmt.Errorf("config is required")
	// }

	var rulesMap map[string]interface{}
	if req.Rules != nil {
		rulesMap = req.Rules.AsMap()
	}

	result, err := h.service.AnalyzisGRPC.Analyze(&service.AnalyzeRequest{
		Config: req.Config.AsMap(),
		Rules:  rulesMap,
	})
	if err != nil {
		logrus.Printf("Service error: %v", err)
		return nil, err
	}

	logrus.Infof("Found %d vulnerabilities", result.TotalCount)

	// respJSON, _ := json.Marshal(resp)
	// logrus.Info("Response: %s", string(respJSON))

	return h.toProtoResponse(result), nil
}

func (h *HandlerGRPC) AnalyzeFile(ctx context.Context, req *analyzer.AnalyzeFileRequest) (*analyzer.AnalyzeResponse, error) {
	logrus.Info("Received Analyze request")

	if req.FilePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	var rulesMap map[string]interface{}

	if req.Rules != nil {
		rulesMap = req.Rules.AsMap()
	}

	result, err := h.service.AnalyzisGRPC.AnalyzeFile(&service.AnalyzeRequest{
		FilePath:         req.FilePath,
		Rules:            rulesMap,
		CheckPermissions: req.CheckPermissions,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("Found %d vulnerabilities", result.TotalCount)

	return h.toProtoResponse(result), nil
}

func (h *HandlerGRPC) Health(ctx context.Context, req *analyzer.HealthRequest) (*analyzer.HealthResponse, error) {
	health := h.service.AnalyzisGRPC.Health()
	return &analyzer.HealthResponse{
		Status:  health["status"],
		Version: health["version"],
	}, nil
}

func (h *HandlerGRPC) toProtoResponse(result *service.AnalyzeResult) *analyzer.AnalyzeResponse {
	resp := &analyzer.AnalyzeResponse{
		Vulnerabilities: make([]*analyzer.Vulnerability, 0, len(result.Vulnerabilities)),
		Permissions:     make([]*analyzer.Permission, 0, len(result.Permissions)),
		TotalCount:      int32(result.TotalCount),
	}

	for _, v := range result.Vulnerabilities {
		resp.Vulnerabilities = append(resp.Vulnerabilities, &analyzer.Vulnerability{
			RuleId:         v.RuleID,
			Severity:       string(v.Severity),
			Description:    v.Description,
			Recommendation: v.Recommendation,
			Path:           v.Path,
			Value:          fmt.Sprintf("%v", v.Value),
			FilePath:       v.FilePath,
		})
	}

	for _, p := range result.Permissions {
		resp.Permissions = append(resp.Permissions, &analyzer.Permission{
			FilePath:    p.FilePath,
			Permission:  p.Permission,
			Recommended: p.Recommended,
			Severity:    string(p.Severity),
			Description: p.Description,
		})
	}

	return resp
}
