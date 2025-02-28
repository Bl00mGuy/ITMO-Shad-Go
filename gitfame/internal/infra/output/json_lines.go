package output

import (
	"fmt"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func WriteJSONLines(userStatistics []stats.UserData) error {
	for _, userData := range userStatistics {
		fmt.Printf(
			"{\"name\":\"%s\",\"lines\":%d,\"commits\":%d,\"files\":%d}\n",
			userData.Name,
			userData.Lines,
			len(userData.Commits),
			userData.Files,
		)
	}
	return nil
}
