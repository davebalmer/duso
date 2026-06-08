package runtime

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

// builtinParseCSV parses a CSV string to array of arrays
func builtinParseCSV(evaluator *Evaluator, args map[string]any) (any, error) {
	csvString := GetArg(args, 0, "str")
	if csvString == nil {
		return nil, fmt.Errorf("parse_csv() requires a string as first argument")
	}

	stringValue, ok := csvString.(string)
	if !ok {
		return nil, fmt.Errorf("parse_csv() requires a string as first argument")
	}

	delimiter := ","
	if delimiterArg := GetArg(args, 1, "delimiter"); delimiterArg != nil {
		if delimiterStr, ok := delimiterArg.(string); ok {
			if len(delimiterStr) > 0 {
				delimiter = delimiterStr
			}
		}
	}

	reader := csv.NewReader(bytes.NewReader([]byte(stringValue)))
	if len(delimiter) > 0 {
		reader.Comma = rune(delimiter[0])
	}

	var result []any
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse_csv() error: %v", err)
		}

		recordArray := make([]any, len(record))
		for i, field := range record {
			recordArray[i] = field
		}
		result = append(result, recordArray)
	}

	return result, nil
}

// builtinFormatCSV formats array of arrays to CSV string
func builtinFormatCSV(evaluator *Evaluator, args map[string]any) (any, error) {
	arrayArg := GetArg(args, 0, "array")
	if arrayArg == nil {
		return nil, fmt.Errorf("format_csv() requires an array as first argument")
	}

	arrayPtr, ok := arrayArg.(*[]Value)
	if !ok {
		return nil, fmt.Errorf("format_csv() requires an array as first argument")
	}

	delimiter := ","
	if delimiterArg := GetArg(args, 1, "delimiter"); delimiterArg != nil {
		if delimiterStr, ok := delimiterArg.(string); ok {
			if len(delimiterStr) > 0 {
				delimiter = delimiterStr
			}
		}
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if len(delimiter) > 0 {
		writer.Comma = rune(delimiter[0])
	}

	for _, row := range *arrayPtr {
		var record []string

		// Handle both array of values and array of arrays
		if rowArray := row.AsArray(); rowArray != nil {
			for _, field := range rowArray {
				record = append(record, field.String())
			}
		} else {
			record = append(record, row.String())
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("format_csv() error: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("format_csv() error: %v", err)
	}

	return buffer.String(), nil
}
