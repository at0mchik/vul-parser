package dto

type AnalyzeRequest struct {
	Config      interface{} `json:"config"`
	Rules       interface{} `json:"rules,omitempty"`
	CheckPermissions bool   `json:"check_permissions,omitempty"`
}

type VulnerabilityResponse struct {
	RuleID       string      `json:"rule_id"`
	Severity     string      `json:"severity"`
	Description  string      `json:"description"`
	Recommendation string   `json:"recommendation"`
	Path         string      `json:"path"`
	Value        interface{} `json:"value"`
	FilePath     string      `json:"file_path,omitempty"`
}

type PermissionResponse struct {
	FilePath    string `json:"file_path"`
	Permission  string `json:"permission"`
	Recommended string `json:"recommended"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type AnalyzeResponse struct {
	Vulnerabilities []VulnerabilityResponse `json:"vulnerabilities"`
	Permissions     []PermissionResponse    `json:"permissions,omitempty"`
	TotalCount      int                     `json:"total_count"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}