package output

import (
	"fmt"
	"os"
	"vul-parser/internal/domain/models"
)

type Printer struct {
	silent bool
}

func NewPrinter(silent bool) *Printer {
	return &Printer{silent: silent}
}

func (p *Printer) Print(vulnerabilities []models.Vulnerability) {
	if len(vulnerabilities) == 0 {
		fmt.Println("No vulnerabilities found")
		return
	}

	for _, v := range vulnerabilities {
		fmt.Printf("%s: %s\n", v.Severity, v.Description)
		fmt.Printf("  Location: %s\n", v.Path)
		fmt.Printf("  Value: %v\n", v.Value)
		fmt.Printf("  Recommendation: %s\n", v.Recommendation)
		fmt.Println()
	}

	if !p.silent {
		os.Exit(1)
	}
}