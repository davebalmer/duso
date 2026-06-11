package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/duso-org/duso/pkg/script"
)

// generateSyntax generates and outputs the TextMate syntax definition as JSON
func generateSyntax() error {
	// Get keywords from the script package
	keywords := script.GetKeywords()
	keywordNames := make([]string, 0, len(keywords))
	for kw := range keywords {
		keywordNames = append(keywordNames, kw)
	}
	sort.Strings(keywordNames)

	// Get builtins - they're registered at startup
	builtinNames := script.GetAllBuiltinNames()
	sort.Strings(builtinNames)

	// Build keywords regex
	keywordRegex := "\\b(" + strings.Join(keywordNames, "|") + ")\\b"

	// Build builtins regex
	builtinRegex := "\\b(" + strings.Join(builtinNames, "|") + ")\\b"

	// Create the TextMate syntax object
	syntax := map[string]interface{}{
		"$schema":   "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		"name":      "Duso",
		"scopeName": "source.duso",
		"fileTypes": []string{"du"},
		"patterns": []map[string]interface{}{
			{"include": "#comments"},
			{"include": "#multiline-strings"},
			{"include": "#strings"},
			{"include": "#keywords"},
			{"include": "#builtins"},
			{"include": "#structures"},
			{"include": "#functions"},
			{"include": "#constants"},
			{"include": "#numbers"},
			{"include": "#operators"},
			{"include": "#punctuation"},
		},
		"repository": map[string]interface{}{
			"comments": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "comment.block.duso",
						"begin": "/\\*",
						"beginCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.comment.begin.duso",
							},
						},
						"end": "\\*/",
						"endCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.comment.end.duso",
							},
						},
						"patterns": []map[string]interface{}{
							{
								"name":  "comment.block.duso",
								"begin": "/\\*",
								"beginCaptures": map[string]interface{}{
									"0": map[string]string{
										"name": "punctuation.definition.comment.begin.duso",
									},
								},
								"end": "\\*/",
								"endCaptures": map[string]interface{}{
									"0": map[string]string{
										"name": "punctuation.definition.comment.end.duso",
									},
								},
							},
						},
					},
					{
						"name":  "comment.line.duso",
						"match": "//.*?$",
					},
				},
			},
			"strings": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "source.regexp.duso",
						"begin": "~",
						"beginCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.regexp.begin.duso",
							},
						},
						"end": "~",
						"endCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.regexp.end.duso",
							},
						},
						"patterns": []map[string]string{
							{
								"name":  "constant.character.escape.regexp.duso",
								"match": "\\\\.",
							},
						},
					},
					{
						"name":  "string.quoted.double.duso",
						"begin": "\"",
						"end":   "\"",
						"patterns": []map[string]interface{}{
							{
								"name":  "constant.character.escape.duso",
								"match": "\\\\([\\\\\"'ntr{}])",
							},
							{
								"include": "#template-expression",
							},
						},
					},
					{
						"name":  "string.quoted.single.duso",
						"begin": "'",
						"end":   "'",
						"patterns": []map[string]interface{}{
							{
								"name":  "constant.character.escape.duso",
								"match": "\\\\([\\\\\"'ntr{}])",
							},
							{
								"include": "#template-expression",
							},
						},
					},
				},
			},
			"template-expression": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "meta.template.duso",
						"begin": "{{",
						"beginCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.template.begin.duso",
							},
						},
						"end": "}}",
						"endCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.definition.template.end.duso",
							},
						},
						"patterns": []map[string]interface{}{
							{"include": "#builtins"},
							{"include": "#operators"},
							{"include": "#numbers"},
							{"include": "#constants"},
							{"include": "#identifiers"},
							{"include": "#punctuation"},
						},
					},
				},
			},
			"multiline-strings": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"begin": "\"\"\"",
						"beginCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.bracket.multiline.begin.duso",
							},
						},
						"end": "\"\"\"",
						"endCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.bracket.multiline.end.duso",
							},
						},
						"contentName":        "string.quoted.double.duso",
						"applyEndPatternLast": true,
						"patterns": []map[string]interface{}{
							{
								"include": "#template-expression",
							},
						},
					},
					{
						"begin": "'''",
						"beginCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.bracket.multiline.begin.duso",
							},
						},
						"end": "'''",
						"endCaptures": map[string]interface{}{
							"0": map[string]string{
								"name": "punctuation.bracket.multiline.end.duso",
							},
						},
						"contentName":        "string.quoted.single.duso",
						"applyEndPatternLast": true,
						"patterns": []map[string]interface{}{
							{
								"include": "#template-expression",
							},
						},
					},
				},
			},
			"keywords": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "keyword.control.duso",
						"match": keywordRegex,
					},
				},
			},
			"builtins": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "support.function.duso",
						"match": builtinRegex,
					},
				},
			},
			"structures": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "entity.name.function.duso",
						"match": "\\b([A-Z][A-Za-z0-9_]*)\\s*(?=\\()",
					},
				},
			},
			"functions": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "meta.function.duso",
						"match": "\\b(function)\\s+([a-z_][a-z0-9_]*)\\b",
						"captures": map[string]interface{}{
							"1": map[string]string{
								"name": "keyword.control.duso",
							},
							"2": map[string]string{
								"name": "entity.name.function.duso",
							},
						},
					},
				},
			},
			"constants": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "constant.language.duso",
						"match": "\\b(true|false|nil)\\b",
					},
				},
			},
			"numbers": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "constant.numeric.duso",
						"match": "\\b([0-9]+(\\.[0-9]+)?)\\b",
					},
				},
			},
			"operators": map[string]interface{}{
				"patterns": []map[string]interface{}{
					{
						"name":  "keyword.operator.comparison.duso",
						"match": "(==|!=|<=|>=|<|>)",
					},
					{
						"name":  "keyword.operator.assignment.duso",
						"match": "(\\+=|-=|\\*=|/=|%=|=)",
					},
					{
						"name":  "keyword.operator.arithmetic.duso",
						"match": "(\\+|-|\\*|/|%)",
					},
					{
						"name":  "keyword.operator.postfix.duso",
						"match": "(\\+\\+|--)",
					},
					{
						"name":  "keyword.operator.ternary.duso",
						"match": "[?:]",
					},
				},
			},
			"punctuation": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "punctuation.duso",
						"match": "[{}\\[\\]().,;]",
					},
				},
			},
			"identifiers": map[string]interface{}{
				"patterns": []map[string]string{
					{
						"name":  "variable.other.duso",
						"match": "\\b[a-z_][a-z0-9_]*\\b",
					},
				},
			},
		},
	}

	// Output as formatted JSON
	output, err := json.MarshalIndent(syntax, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate syntax JSON: %v", err)
	}

	fmt.Println(string(output))
	return nil
}
