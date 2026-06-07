package cli

import (
	"github.com/duso-org/duso/pkg/script"
)

// LintScript analyzes a Duso script and returns diagnostics
func LintScript(filename string, source string) ([]*script.LintDiagnostic, error) {
	// Tokenize
	lexer := script.NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	// Parse
	parser := script.NewParserWithFile(tokens, filename)
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Analyze AST
	analyzer := script.NewLintAnalyzer(program, filename)
	return analyzer.Analyze(), nil
}
