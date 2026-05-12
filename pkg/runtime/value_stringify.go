package runtime

import (
	"fmt"
	"sort"
	"strings"

	"github.com/duso-org/duso/pkg/script"
)

// ValueForDisplay converts a value to a display string for print/output.
// Strings are printed as-is without quotes (for readability).
// Other types use Duso syntax so arrays/objects display correctly.
func ValueForDisplay(val any) string {
	// Strings print as-is without quotes
	if str, ok := val.(string); ok {
		return str
	}
	// Everything else uses Duso syntax
	return ValueToDusoString(val)
}

// ValueToDusoString converts any value to a Duso-parsable string representation.
// This is used for tostring(), embedded values, and any place Duso syntax is needed.
// The result is valid Duso syntax that can be parsed back with parse().
func ValueToDusoString(val any) string {
	return valueToDusoStringInner(val, 0)
}

// valueToDusoStringInner is the recursive implementation with depth tracking.
func valueToDusoStringInner(val any, depth int) string {
	switch v := val.(type) {
	case nil:
		return "nil"

	case bool:
		if v {
			return "true"
		}
		return "false"

	case float64:
		// Format numbers without unnecessary decimals
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)

	case string:
		// Escape special characters and wrap in quotes
		return fmt.Sprintf("\"%s\"", escapeString(v))

	case *[]Value:
		// Arrays come in as *[]Value from ValueToInterface
		if v == nil || len(*v) == 0 {
			return "[]"
		}
		var parts []string
		for _, item := range *v {
			parts = append(parts, valueToDusoStringInner(ValueToInterface(item), depth+1))
		}
		return "[" + strings.Join(parts, ", ") + "]"

	case map[string]Value:
		// Objects come in as map[string]Value from ValueToInterface
		if len(v) == 0 {
			return "{}"
		}

		// Sort keys for consistent output
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var parts []string
		for _, k := range keys {
			valStr := valueToDusoStringInner(ValueToInterface(v[k]), depth+1)
			parts = append(parts, fmt.Sprintf("%s=%s", k, valStr))
		}
		return "{" + strings.Join(parts, ", ") + "}"

	case *ValueRef:
		// ValueRef wraps functions, code, errors, binaries
		if v == nil {
			return "<nil>"
		}
		switch v.Val.Type {
		case script.VAL_CODE:
			return "<code>"
		case script.VAL_ERROR:
			// Error is wrapped in ValueRef
			errVal := v.Val.AsErrorVal()
			if errVal != nil {
				msgStr := valueToDusoStringInner(ValueToInterface(errVal.Message), depth+1)
				return fmt.Sprintf("error(%s)", msgStr)
			}
			return "error(<unknown>)"
		case VAL_FUNCTION:
			return "<function>"
		case script.VAL_BINARY:
			return "<binary>"
		default:
			return fmt.Sprintf("<%s>", v.Val.Type.String())
		}

	case *script.DusoError:
		// Runtime errors (from throw() or evaluator)
		msgStr := valueToDusoStringInner(v.Message, depth+1)
		return fmt.Sprintf("error(%s)", msgStr)

	case []any:
		// Fallback for []any (shouldn't normally happen)
		if len(v) == 0 {
			return "[]"
		}
		var parts []string
		for _, item := range v {
			parts = append(parts, valueToDusoStringInner(item, depth+1))
		}
		return "[" + strings.Join(parts, ", ") + "]"

	case map[string]any:
		// Fallback for map[string]any (shouldn't normally happen)
		if len(v) == 0 {
			return "{}"
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			valStr := valueToDusoStringInner(v[k], depth+1)
			parts = append(parts, fmt.Sprintf("%s=%s", k, valStr))
		}
		return "{" + strings.Join(parts, ", ") + "}"

	default:
		// For unknown types, use a generic representation
		return fmt.Sprintf("<%T>", v)
	}
}

// escapeString escapes special characters in a string for Duso syntax
// Only escapes quotes and backslashes that would break Duso parsing
// Preserves actual newlines/tabs/etc for display purposes
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
