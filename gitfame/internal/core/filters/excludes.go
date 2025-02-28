package filters

import (
	"path"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func hasExclusionPatterns(request stats.RepoFlags) bool {
	return len(request.Exclude) > 0
}

func isFileExcluded(patterns []string, fileName string) bool {
	for _, pattern := range patterns {
		if matchesPattern(pattern, fileName) {
			return true
		}
	}
	return false
}

func matchesPattern(pattern, fileName string) bool {
	match, err := path.Match(pattern, fileName)
	if err != nil {
		return true
	}
	return match
}

func Excludes(request stats.RepoFlags, fileName string) bool {
	if !hasExclusionPatterns(request) {
		return false
	}
	return isFileExcluded(request.Exclude, fileName)
}
