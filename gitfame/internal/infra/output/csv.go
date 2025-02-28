package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func WriteCSV(userStats []stats.UserData) error {
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	header := []string{"Name", "Lines", "Commits", "Files"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, userData := range userStats {
		row := []string{
			userData.Name,
			fmt.Sprint(userData.Lines),
			fmt.Sprint(len(userData.Commits)),
			fmt.Sprint(userData.Files),
		}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row for user '%s': %w", userData.Name, err)
		}
	}

	csvWriter.Flush()
	return nil
}
