package service

import (
	"encoding/json"
	"fmt"

	"vul-parser/internal/checker"
	"vul-parser/internal/domain/models"
	"vul-parser/internal/parser"
	"vul-parser/internal/permission"
	"vul-parser/internal/rules"
)

type AnalyzeRequest struct {
	Config           map[string]interface{}
	Rules            map[string]interface{}
	CheckPermissions bool
	FilePath         string
}

type AnalyzeResult struct {
	Vulnerabilities []models.Vulnerability
	Permissions     []models.FilePermission
	TotalCount      int
}

type AnalyzerGRPCService struct {
	defaultRules []models.Rule
}

func NewAnalyzerGRPCService() *AnalyzerGRPCService {
	defaultRules, _ := rules.LoadRules("")

	return &AnalyzerGRPCService{
		defaultRules: defaultRules,
	}
}

func (s *AnalyzerGRPCService) Analyze(req *AnalyzeRequest) (*AnalyzeResult, error) {
	rulesList, err := s.loadRulesFromMap(req.Rules)
	if err != nil {
		return nil, err
	}

	checkerEngine := checker.NewChecker(rulesList)

	// Конвертируем map в JSON
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	parsedConfig, err := parser.Parse(configJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	vulnerabilities := checkerEngine.Check(parsedConfig, "grpc_request")

	return &AnalyzeResult{
		Vulnerabilities: vulnerabilities,
		TotalCount:      len(vulnerabilities),
	}, nil
}

func (s *AnalyzerGRPCService) AnalyzeFile(req *AnalyzeRequest) (*AnalyzeResult, error) {
	rulesList, err := s.loadRulesFromMap(req.Rules)
	if err != nil {
		return nil, err
	}

	checkerEngine := checker.NewChecker(rulesList)
	permChecker := permission.NewPermissionChecker()

	var vulnerabilities []models.Vulnerability
	var permissions []models.FilePermission

	if req.CheckPermissions {
		if perm := permChecker.CheckFile(req.FilePath); perm != nil {
			permissions = append(permissions, *perm)
		}
	}

	configData, err := parser.ReadFromFile(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parsedConfig, err := parser.Parse(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	vulnerabilities = checkerEngine.Check(parsedConfig, req.FilePath)

	return &AnalyzeResult{
		Vulnerabilities: vulnerabilities,
		Permissions:     permissions,
		TotalCount:      len(vulnerabilities) + len(permissions),
	}, nil
}

func (s *AnalyzerGRPCService) Health() map[string]string {
	return map[string]string{
		"status":  "ok",
		"version": "1.0.0",
	}
}

func (s *AnalyzerGRPCService) loadRulesFromMap(rulesMap map[string]interface{}) ([]models.Rule, error) {
	if len(rulesMap) == 0 {
		return s.defaultRules, nil
	}

	// Конвертируем map в YAML (для совместимости с существующим loader)
	rulesYAML, err := json.Marshal(rulesMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules: %w", err)
	}

	// Правила ожидают YAML, конвертируем JSON в YAML
	var data interface{}
	if err := json.Unmarshal(rulesYAML, &data); err != nil {
		return nil, err
	}

	return s.loadRulesFromInterface(data)
}

func (s *AnalyzerGRPCService) loadRulesFromInterface(data interface{}) ([]models.Rule, error) {
	// Получаем map правил
	rulesMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid rules format")
	}

	rulesArray, ok := rulesMap["rules"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("rules.rules must be an array")
	}

	var rules []models.Rule

	for _, r := range rulesArray {
		ruleMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		rule := models.Rule{
			ID:             s.getString(ruleMap, "id"),
			Name:           s.getString(ruleMap, "name"),
			Severity:       models.Severity(s.getString(ruleMap, "severity")),
			Description:    s.getString(ruleMap, "description"),
			Recommendation: s.getString(ruleMap, "recommendation"),
		}

		conditions, ok := ruleMap["conditions"].([]interface{})
		if ok {
			for _, c := range conditions {
				condMap, ok := c.(map[string]interface{})
				if !ok {
					continue
				}

				cond := models.Condition{
					Path:              s.getString(condMap, "path"),
					Operator:          models.ConditionOperator(s.getString(condMap, "operator")),
					Value:             condMap["value"],
					AndValueNotEmpty:  s.getBool(condMap, "and_value_not_empty"),
					ExcludeValueRegex: s.getString(condMap, "exclude_value_regex"),
				}
				rule.Conditions = append(rule.Conditions, cond)
			}
		}

		rules = append(rules, rule)
	}

	if len(rules) == 0 {
		return s.defaultRules, nil
	}

	return rules, nil
}

func (s *AnalyzerGRPCService) getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func (s *AnalyzerGRPCService) getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}