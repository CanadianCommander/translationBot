package translationFile

import (
	"fmt"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

const dotPlaceholder = "\\DOT_PLACEHOLDER\\"

// formatKeypath formats the given key and path
func formatKeypath(path string, key string) string {
	key = strings.ReplaceAll(key, ".", dotPlaceholder)
	return strings.TrimLeft(fmt.Sprintf("%s.%s", path, key), ".")
}

// unEscapeKeypath removes escaping characters from the keypath
func unEscapeKeypath(keypath string) string {
	return strings.ReplaceAll(keypath, dotPlaceholder, ".")
}

// splitKeypath in to its path parts (also handles escaping)
func splitKeypath(keypath string) []string {
	pathParts := strings.Split(keypath, ".")

	for idx, _ := range pathParts {
		pathParts[idx] = unEscapeKeypath(pathParts[idx])
	}

	return pathParts
}
