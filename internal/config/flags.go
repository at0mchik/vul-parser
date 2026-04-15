package config

import (
	"flag"
	"fmt"
)

type Config struct {
	FilePath  string
	Silent    bool
	Stdin     bool
	RulesPath string
}

func ParseFlags() (*Config, error) {
	var silent bool
	var stdin bool
	var rulesPath string

	flag.BoolVar(&silent, "s", false, "silent mode - don't exit with error")
	flag.BoolVar(&silent, "silent", false, "silent mode - don't exit with error")
	flag.BoolVar(&stdin, "stdin", false, "read config from stdin")
	flag.StringVar(&rulesPath, "rules", "", "custom rules file path")

	flag.Parse()

	var filePath string
	if !stdin {
		if flag.NArg() < 1 {
			return nil, fmt.Errorf("config file path required")
		}
		filePath = flag.Arg(0)
	}

	return &Config{
		FilePath:  filePath,
		Silent:    silent,
		Stdin:     stdin,
		RulesPath: rulesPath,
	}, nil
}