package main

import (
	"fmt"
	"os"

	"vul-parser/internal/checker"
	"vul-parser/internal/config"
	"vul-parser/internal/output"
	"vul-parser/internal/parser"
	"vul-parser/internal/rules"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var configData []byte

	if cfg.Stdin {
		configData, err = parser.ReadFromReader(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		
		if len(configData) == 0 {
			fmt.Fprintf(os.Stderr, "Error: empty input from stdin\n")
			os.Exit(1)
		}
	} else {
		configData, err = os.ReadFile(cfg.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	}

	parsedConfig, err := parser.Parse(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		os.Exit(1)
	}

	rulesList, err := rules.LoadRules(cfg.RulesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading rules: %v\n", err)
		os.Exit(1)
	}

	checkerEngine := checker.NewChecker(rulesList)
	vulnerabilities := checkerEngine.Check(parsedConfig)

	printer := output.NewPrinter(cfg.Silent)
	printer.Print(vulnerabilities)
}