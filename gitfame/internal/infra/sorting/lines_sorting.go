package sorting

import (
	"sort"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func SortByLines(answer []stats.UserData) []stats.UserData {
	sort.SliceStable(answer, func(i, j int) bool {
		if answer[i].Lines == answer[j].Lines && len(answer[i].Commits) != len(answer[j].Commits) {
			return len(answer[i].Commits) > len(answer[j].Commits)
		} else if answer[i].Lines == answer[j].Lines && len(answer[i].Commits) == len(answer[j].Commits) && answer[i].Files != answer[j].Files {
			return answer[i].Files > answer[j].Files
		} else if answer[i].Lines == answer[j].Lines && len(answer[i].Commits) == len(answer[j].Commits) && answer[i].Files == answer[j].Files {
			return strings.Compare(answer[i].Name, answer[j].Name) == -1
		} else {
			return answer[i].Lines > answer[j].Lines
		}
	})
	return answer
}
