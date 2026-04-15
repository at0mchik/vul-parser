package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse(data []byte) (interface{}, error)
}

func Parse(data []byte) (interface{}, error) {
	data = []byte(strings.TrimSpace(string(data)))
	
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	
	if data[0] == '{' || data[0] == '[' {
		return parseJSON(data)
	}
	
	return parseYAML(data)
}

func parseJSON(data []byte) (interface{}, error) {
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return result, nil
}

func parseYAML(data []byte) (interface{}, error) {
	var result interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	return result, nil
}

func ReadFromReader(r io.Reader) ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	
	// Проверяем, является ли stdin терминалом (интерактивный режим)
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "Введите конфигурацию (нажмите Ctrl+D на пустой строке для окончания ввода):")
		
		var lines []string
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}
		
		if len(lines) == 0 {
			return nil, fmt.Errorf("no input provided")
		}
		
		return []byte(strings.Join(lines, "\n")), nil
	}
	
	// Неинтерактивный режим (pipe или redirection)
	return io.ReadAll(r)
}

func ReadFromFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}