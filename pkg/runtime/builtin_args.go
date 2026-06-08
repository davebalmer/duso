package runtime

import "strconv"

// GetArg retrieves an argument by name or position (0-indexed)
// Checks named arg first, then positional
func GetArg(args map[string]any, index int, name string) any {
	if v, ok := args[name]; ok {
		return v
	}
	if v, ok := args[strconv.Itoa(index)]; ok {
		return v
	}
	return nil
}
