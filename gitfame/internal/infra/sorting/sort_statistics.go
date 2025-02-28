package sorting

import (
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func SortStatistics(answer []stats.UserData, orderBy string) {
	switch orderBy {
	case "lines":
		SortByLines(answer)
	case "commits":
		SortByCommits(answer)
	case "files":
		SortByFiles(answer)
	default:
		fmt.Fprintln(os.Stderr, "Invalid order-by value: ", orderBy)
		os.Exit(1)
	}
}
