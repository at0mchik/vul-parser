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
		if v.FilePath != "" {
			fmt.Printf("File: %s\n", v.FilePath)
		}
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

func (p *Printer) PrintPermissions(permissions []models.FilePermission) {
	if len(permissions) == 0 {
		return
	}
	
	for _, perm := range permissions {
		fmt.Printf("%s: %s\n", perm.Severity, perm.Description)
		fmt.Printf("  File: %s\n", perm.FilePath)
		fmt.Printf("  Current permissions: %s\n", perm.Permission)
		fmt.Printf("  Recommended: %s\n", perm.Recommended)
		fmt.Println()
	}
}