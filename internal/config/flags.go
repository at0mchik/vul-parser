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
	Recursive bool
}

func ParseFlags() (*Config, error) {
	var silent bool
	var stdin bool
	var rulesPath string
	var recursive bool

	flag.BoolVar(&silent, "s", false, "silent mode - don't exit with error")
	flag.BoolVar(&silent, "silent", false, "silent mode - don't exit with error")
	flag.BoolVar(&stdin, "stdin", false, "read config from stdin")
	flag.StringVar(&rulesPath, "rules", "", "custom rules file path")
	flag.BoolVar(&recursive, "r", false, "recursive directory analysis")
	flag.BoolVar(&recursive, "recursive", false, "recursive directory analysis")

	flag.Parse()

	var filePath string
	if !stdin {
		if flag.NArg() < 1 {
			return nil, fmt.Errorf("config file or directory path required")
		}
		filePath = flag.Arg(0)
	}

	return &Config{
		FilePath:  filePath,
		Silent:    silent,
		Stdin:     stdin,
		RulesPath: rulesPath,
		Recursive: recursive,
	}, nil
}