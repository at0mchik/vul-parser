package permission

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vul-parser/internal/domain/models"
)

type PermissionChecker struct{}

func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{}
}

func (p *PermissionChecker) CheckFile(filePath string) *models.FilePermission {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	mode := info.Mode()
	perm := mode.Perm()
	
	// Проверка слишком широких прав
	if perm&0077 != 0 {
		permStr := fmt.Sprintf("%o", perm)
		var recommended string
		var severity models.Severity
		
		if perm&0007 != 0 {
			recommended = "chmod 640 или chmod 600"
			severity = models.High
		} else if perm&0070 != 0 {
			recommended = "chmod 640 или chmod 600"
			severity = models.Medium
		} else {
			recommended = "chmod 644"
			severity = models.Medium
		}
		
		return &models.FilePermission{
			FilePath:    filePath,
			Permission:  permStr,
			Recommended: recommended,
			Severity:    severity,
			Description: "Слишком широкие права доступа на файл конфигурации",
		}
	}
	
	return nil
}

func (p *PermissionChecker) CheckDirectory(dirPath string, recursive bool) []models.FilePermission {
	var results []models.FilePermission
	
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			if !recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}
		
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			return nil
		}
		
		if perm := p.CheckFile(path); perm != nil {
			results = append(results, *perm)
		}
		
		return nil
	}
	
	filepath.Walk(dirPath, walkFn)
	return results
}