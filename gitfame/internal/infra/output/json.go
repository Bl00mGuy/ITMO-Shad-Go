package output

import (
	"fmt"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func WriteJSON(userStatistics []stats.UserData) error {
	fmt.Print("[")

	for index, userData := range userStatistics {
		fmt.Printf(
			"{\"name\":\"%s\",\"lines\":%d,\"commits\":%d,\"files\":%d}",
			userData.Name,
			userData.Lines,
			len(userData.Commits),
			userData.Files,
		)

		if index != len(userStatistics)-1 {
			fmt.Print(",")
		}
	}

	fmt.Println("]")
	return nil
}
