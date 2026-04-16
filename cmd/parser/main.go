package main

import (
	"fmt"
	"os"
	"path/filepath"

	"vul-parser/internal/checker"
	"vul-parser/internal/config"
	"vul-parser/internal/domain/models"
	"vul-parser/internal/output"
	"vul-parser/internal/parser"
	"vul-parser/internal/permission"
	"vul-parser/internal/rules"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	rulesList, err := rules.LoadRules(cfg.RulesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading rules: %v\n", err)
		os.Exit(1)
	}

	checkerEngine := checker.NewChecker(rulesList)
	permChecker := permission.NewPermissionChecker()
	printer := output.NewPrinter(cfg.Silent)

	var allVulnerabilities []models.Vulnerability
	var allPermissions []models.FilePermission

	if cfg.Stdin {
		configData, err := parser.ReadFromReader(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		
		parsedConfig, err := parser.Parse(configData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
			os.Exit(1)
		}
		
		vulnerabilities := checkerEngine.Check(parsedConfig, "stdin")
		allVulnerabilities = append(allVulnerabilities, vulnerabilities...)
	} else {
		info, err := os.Stat(cfg.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", err)
			os.Exit(1)
		}
		
		if info.IsDir() {
			err := filepath.Walk(cfg.FilePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				
				if info.IsDir() {
					if !cfg.Recursive && path != cfg.FilePath {
						return filepath.SkipDir
					}
					return nil
				}
				
				ext := filepath.Ext(path)
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					return nil
				}
				
				if perm := permChecker.CheckFile(path); perm != nil {
					allPermissions = append(allPermissions, *perm)
				}
				
				configData, err := os.ReadFile(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: cannot read %s: %v\n", path, err)
					return nil
				}
				
				parsedConfig, err := parser.Parse(configData)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: cannot parse %s: %v\n", path, err)
					return nil
				}
				
				vulnerabilities := checkerEngine.Check(parsedConfig, path)
				allVulnerabilities = append(allVulnerabilities, vulnerabilities...)
				
				return nil
			})
			
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			if perm := permChecker.CheckFile(cfg.FilePath); perm != nil {
				allPermissions = append(allPermissions, *perm)
			}
			
			configData, err := os.ReadFile(cfg.FilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
			
			parsedConfig, err := parser.Parse(configData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
				os.Exit(1)
			}
			
			vulnerabilities := checkerEngine.Check(parsedConfig, cfg.FilePath)
			allVulnerabilities = append(allVulnerabilities, vulnerabilities...)
		}
	}
	
	printer.PrintPermissions(allPermissions)
	printer.Print(allVulnerabilities)
}