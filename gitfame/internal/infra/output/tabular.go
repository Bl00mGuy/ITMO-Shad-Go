package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func WriteTabular(userStatistics []stats.UserData) error {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(tabWriter, "Name\tLines\tCommits\tFiles")

	for _, userData := range userStatistics {
		fmt.Fprintf(
			tabWriter,
			"%s\t%d\t%d\t%d\n",
			userData.Name,
			userData.Lines,
			len(userData.Commits),
			userData.Files,
		)
	}

	tabWriter.Flush()
	return nil
}
