package script

import (
	"fmt"
	"sort"
	"strings"
)

// ValueForDisplay converts a value to a display string for print/output.
// Strings are printed as-is without quotes (for readability).
// Other types use Duso syntax so arrays/objects display correctly.
func ValueForDisplay(val Value) string {
	// Strings print as-is without quotes
	if val.IsString() {
		return val.AsString()
	}
	// Everything else uses Duso syntax
	return ValueToDusoString(val)
}

// ValueToDusoString converts any Value to a Duso-parsable string representation.
// This is used for tostring(), templates, and any place Duso syntax is needed.
// The result is valid Duso syntax that can be parsed back with parse().
func ValueToDusoString(val Value) string {
	return valueToDusoStringInner(val, 0)
}

// valueToDusoStringInner is the recursive implementation with depth tracking.
func valueToDusoStringInner(val Value, depth int) string {
	switch val.Type {
	case VAL_NIL:
		return "nil"

	case VAL_BOOL:
		if val.AsBool() {
			return "true"
		}
		return "false"

	case VAL_NUMBER:
		n := val.AsNumber()
		if n == float64(int64(n)) {
			return fmt.Sprintf("%.0f", n)
		}
		return fmt.Sprintf("%g", n)

	case VAL_STRING:
		// Escape special characters and wrap in quotes
		return fmt.Sprintf("\"%s\"", escapeString(val.AsString()))

	case VAL_ARRAY:
		arr := val.AsArray()
		if len(arr) == 0 {
			return "[]"
		}
		var parts []string
		for _, item := range arr {
			parts = append(parts, valueToDusoStringInner(item, depth+1))
		}
		return "[" + strings.Join(parts, ", ") + "]"

	case VAL_OBJECT:
		obj := val.Data.(map[string]Value)
		if len(obj) == 0 {
			return "{}"
		}

		// Sort keys for consistent output
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var parts []string
		for _, k := range keys {
			valStr := valueToDusoStringInner(obj[k], depth+1)
			parts = append(parts, fmt.Sprintf("%s=%s", k, valStr))
		}
		return "{" + strings.Join(parts, ", ") + "}"

	case VAL_FUNCTION:
		return "<function>"

	case VAL_CODE:
		return "<code>"

	case VAL_ERROR:
		errVal := val.AsErrorVal()
		if errVal != nil {
			msgStr := valueToDusoStringInner(errVal.Message, depth+1)
			return fmt.Sprintf("error(%s)", msgStr)
		}
		return "error(<unknown>)"

	case VAL_BINARY:
		bin := val.AsBinary()
		if bin != nil && bin.Data != nil {
			size := len(*bin.Data)
			filename := ""
			if fn, ok := bin.Metadata["filename"]; ok && fn.IsString() {
				filename = fn.AsString()
			}
			if filename != "" {
				return fmt.Sprintf("<binary: %s (%d bytes)>", filename, size)
			}
			return fmt.Sprintf("<binary: %d bytes>", size)
		}
		return "<binary>"

	default:
		return fmt.Sprintf("<%s>", val.Type.String())
	}
}

// escapeString escapes special characters in a string for Duso syntax
// Escapes newlines, tabs, and other control characters as literals (\n, \t, etc)
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
