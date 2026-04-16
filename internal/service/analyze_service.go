package service

import (
	"encoding/json"
	"fmt"
	"vul-parser/internal/checker"
	"vul-parser/internal/domain/dto"
	"vul-parser/internal/domain/models"
	"vul-parser/internal/parser"
	"vul-parser/internal/permission"
	"vul-parser/internal/rules"

	"github.com/sirupsen/logrus"
)

type AnalysisHTTPService struct {
	defaultRules []models.Rule
}

func NewAnalysisHTTPService() *AnalysisHTTPService {
	defaultRules, err := rules.LoadRules("")
	if err != nil {
		logrus.Errorf("failed to load default rules: %v", err)
		return nil
	}

	return &AnalysisHTTPService{
		defaultRules: defaultRules,
	}
}

func (s *AnalysisHTTPService) Analyze(req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error) {
	var rulesList []models.Rule
	var err error

	if req.Rules != nil {
		rulesData, err := json.Marshal(req.Rules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal rules: %w", err)
		}
		rulesList, err = s.loadRulesFromJSON(rulesData)
		if err != nil {
			return nil, fmt.Errorf("failed to load custom rules: %w", err)
		}
	} else {
		rulesList = s.defaultRules
	}

	checkerEngine := checker.NewChecker(rulesList)

	var vulnerabilities []models.Vulnerability

	configData, err := json.Marshal(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	parsedConfig, err := parser.Parse(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	vulnerabilities = checkerEngine.Check(parsedConfig, "request")

	response := &dto.AnalyzeResponse{
		Vulnerabilities: make([]dto.VulnerabilityResponse, 0, len(vulnerabilities)),
		TotalCount:      len(vulnerabilities),
	}

	for _, v := range vulnerabilities {
		response.Vulnerabilities = append(response.Vulnerabilities, dto.VulnerabilityResponse{
			RuleID:         v.RuleID,
			Severity:       string(v.Severity),
			Description:    v.Description,
			Recommendation: v.Recommendation,
			Path:           v.Path,
			Value:          v.Value,
			FilePath:       v.FilePath,
		})
	}

	if req.CheckPermissions {
		// permChecker := permission.NewPermissionChecker()
		response.Permissions = make([]dto.PermissionResponse, 0)
	}

	return response, nil
}

func (s *AnalysisHTTPService) AnalyzeWithFile(filePath string, req *dto.AnalyzeRequest) (*dto.AnalyzeResponse, error) {
	var rulesList []models.Rule
	var err error

	if req.Rules != nil {
		rulesData, err := json.Marshal(req.Rules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal rules: %w", err)
		}
		rulesList, err = s.loadRulesFromJSON(rulesData)
		if err != nil {
			return nil, fmt.Errorf("failed to load custom rules: %w", err)
		}
	} else {
		rulesList = s.defaultRules
	}

	checkerEngine := checker.NewChecker(rulesList)
	permChecker := permission.NewPermissionChecker()

	var vulnerabilities []models.Vulnerability
	var permissions []models.FilePermission

	if req.CheckPermissions {
		if perm := permChecker.CheckFile(filePath); perm != nil {
			permissions = append(permissions, *perm)
		}
	}

	configData, err := parser.ReadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parsedConfig, err := parser.Parse(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	vulnerabilities = checkerEngine.Check(parsedConfig, filePath)

	response := &dto.AnalyzeResponse{
		Vulnerabilities: make([]dto.VulnerabilityResponse, 0, len(vulnerabilities)),
		Permissions:     make([]dto.PermissionResponse, 0, len(permissions)),
		TotalCount:      len(vulnerabilities) + len(permissions),
	}

	for _, v := range vulnerabilities {
		response.Vulnerabilities = append(response.Vulnerabilities, dto.VulnerabilityResponse{
			RuleID:         v.RuleID,
			Severity:       string(v.Severity),
			Description:    v.Description,
			Recommendation: v.Recommendation,
			Path:           v.Path,
			Value:          v.Value,
			FilePath:       v.FilePath,
		})
	}

	for _, p := range permissions {
		response.Permissions = append(response.Permissions, dto.PermissionResponse{
			FilePath:    p.FilePath,
			Permission:  p.Permission,
			Recommended: p.Recommended,
			Severity:    string(p.Severity),
			Description: p.Description,
		})
	}

	return response, nil
}

func (s *AnalysisHTTPService) loadRulesFromJSON(data []byte) ([]models.Rule, error) {
	var ruleSet struct {
		Rules []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Severity       string `json:"severity"`
			Description    string `json:"description"`
			Recommendation string `json:"recommendation"`
			Conditions     []struct {
				Path              string      `json:"path"`
				Operator          string      `json:"operator"`
				Value             interface{} `json:"value"`
				AndValueNotEmpty  bool        `json:"and_value_not_empty"`
				ExcludeValueRegex string      `json:"exclude_value_regex,omitempty"`
			} `json:"conditions"`
		} `json:"rules"`
	}

	if err := json.Unmarshal(data, &ruleSet); err != nil {
		return nil, err
	}

	rules := make([]models.Rule, 0, len(ruleSet.Rules))
	for _, r := range ruleSet.Rules {
		conditions := make([]models.Condition, 0, len(r.Conditions))
		for _, c := range r.Conditions {
			conditions = append(conditions, models.Condition{
				Path:              c.Path,
				Operator:          models.ConditionOperator(c.Operator),
				Value:             c.Value,
				AndValueNotEmpty:  c.AndValueNotEmpty,
				ExcludeValueRegex: c.ExcludeValueRegex,
			})
		}
		rules = append(rules, models.Rule{
			ID:             r.ID,
			Name:           r.Name,
			Severity:       models.Severity(r.Severity),
			Description:    r.Description,
			Recommendation: r.Recommendation,
			Conditions:     conditions,
		})
	}

	return rules, nil
}
