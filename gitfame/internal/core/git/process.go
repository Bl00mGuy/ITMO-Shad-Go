package git

import (
	"os/exec"
	"strconv"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func process(fileName string, info stats.RepoFlags) (map[string]stats.UserData, error) {
	cmd := exec.Command("git", "blame", fileName, "--porcelain", info.Revision)
	cmd.Dir = info.Repository
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	blame := string(cmdOutput)

	if len(blame) == 0 {
		cmd := exec.Command("git", "log", info.Revision, "-1", "--pretty=format:%H %an", "--", fileName)
		cmd.Dir = info.Repository
		cmdOutput, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		log := string(cmdOutput)

		lines := strings.Split(log, "\n")

		result := make(stats.UserDataSet)

		for _, line := range lines {
			var singleFile stats.UserData
			currentFields := strings.Fields(line)
			prefix := currentFields[0] + " "
			name := strings.TrimPrefix(line, prefix)
			singleFile.Commits = make(stats.IntSet)
			singleFile.Commits[currentFields[0]] = 1
			result[name] = singleFile
		}

		return result, nil
	}

	lines := strings.Split(blame, "\n")
	var who string
	if info.UseCommitter {
		who = "committer"
	} else {
		who = "author"
	}

	ccommitAmount := make(stats.IntSet)
	commitsFromUser := make(stats.StringSet)

	var currentHash string
	for _, line := range lines {
		currentFields := strings.Fields(line)

		if len(currentFields) == 0 {
			continue
		}

		if len(currentFields) == 4 {
			currentHash = currentFields[0]
			count, _ := strconv.Atoi(currentFields[3])
			ccommitAmount[currentHash] += count
			continue
		}
		if currentFields[0] == who {
			_, ok := commitsFromUser[currentHash]
			if !ok {
				prefix := who + " "
				commitsFromUser[currentHash] = strings.TrimPrefix(line, prefix)
			}
		}
	}

	result := make(stats.UserDataSet)

	for commitHash, user := range commitsFromUser {
		stat, ok := result[user]
		if !ok {
			stat.Commits = make(stats.IntSet)
		}
		stat.Commits[commitHash] = 1
		stat.Lines += ccommitAmount[commitHash]

		result[user] = stat
	}

	return result, nil
}
