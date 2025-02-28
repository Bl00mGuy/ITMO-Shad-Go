package output

import (
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func OutputResults(answer []stats.UserData, format string) error {
	switch format {
	case "tabular":
		return WriteTabular(answer)
	case "csv":
		return WriteCSV(answer)
	case "json":
		return WriteJSON(answer)
	case "json-lines":
		return WriteJSONLines(answer)
	default:
		fmt.Fprintln(os.Stderr, "Invalid format value: ", format)
		os.Exit(1)
		return nil
	}
}
