package sorting

import (
	"sort"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func SortByCommits(answer []stats.UserData) []stats.UserData {
	sort.SliceStable(answer, func(i, j int) bool {
		if len(answer[i].Commits) == len(answer[j].Commits) && answer[i].Lines != answer[j].Lines {
			return answer[i].Lines > answer[j].Lines
		} else if len(answer[i].Commits) == len(answer[j].Commits) && answer[i].Lines == answer[j].Lines && answer[i].Files != answer[j].Files {
			return answer[i].Files > answer[j].Files
		} else if len(answer[i].Commits) == len(answer[j].Commits) && answer[i].Lines == answer[j].Lines && answer[i].Files == answer[j].Files {
			return strings.Compare(answer[i].Name, answer[j].Name) == -1
		} else {
			return len(answer[i].Commits) > len(answer[j].Commits)
		}
	})
	return answer
}
