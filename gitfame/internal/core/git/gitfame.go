package git

import (
	"os/exec"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/core/filters"
	"gitlab.com/slon/shad-go/gitfame/internal/core/languages"
	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func Gitfame(request stats.RepoFlags) ([]stats.UserData, error) {
	var fileTree []string

	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", request.Revision)
	cmd.Dir = request.Repository
	fileNames, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	if len(fileNames) == 0 {
		fileTree = nil
	} else {
		fileTree = strings.Split(strings.TrimSpace(string(fileNames)), "\n")
	}

	userStats := make(stats.UserDataSet)

	for _, fileName := range fileTree {
		if filters.Excludes(request, fileName) || filters.Restricted(request, fileName) ||
			filters.ExtensionCheck(request, fileName) || languages.LanguageCheck(request, fileName) {
			continue
		}

		processedFile, err := process(fileName, request)
		if err != nil {
			return nil, err
		}

		for name, singleFile := range processedFile {
			userStat, ok := userStats[name]
			if !ok {
				userStat.Commits = make(stats.IntSet)
			}
			userStat.Files++
			userStat.Lines += singleFile.Lines

			for commitHash := range singleFile.Commits {
				userStat.Commits[commitHash] = 1
			}

			userStats[name] = userStat
		}
	}

	var userInfo []stats.UserData

	for name, stat := range userStats {
		var info stats.UserData
		info.Name = name
		info.Files = stat.Files
		info.Lines = stat.Lines
		info.Commits = stat.Commits

		userInfo = append(userInfo, info)
	}

	return userInfo, nil
}
