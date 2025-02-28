package filters

import (
	"path"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func hasRestrictions(request stats.RepoFlags) bool {
	return len(request.RestrictTo) > 0
}

func matchesAnyPattern(patterns []string, fileName string) (bool, bool) {
	for _, pattern := range patterns {
		match, err := path.Match(pattern, fileName)
		if err != nil {
			return false, true
		}
		if match {
			return true, false
		}
	}
	return false, false
}

func Restricted(request stats.RepoFlags, fileName string) bool {
	if !hasRestrictions(request) {
		return false
	}

	match, hasError := matchesAnyPattern(request.RestrictTo, fileName)
	if hasError {
		return true
	}

	return !match
}
