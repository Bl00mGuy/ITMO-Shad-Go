//go:build !solution

package main

import (
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/core/git"
	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
	"gitlab.com/slon/shad-go/gitfame/internal/infra/output"
	"gitlab.com/slon/shad-go/gitfame/internal/infra/sorting"

	"github.com/spf13/pflag"
)

func main() {
	request := parseFlags()

	answer, err := git.Gitfame(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get statistics: ", err)
		os.Exit(1)
	}

	sorting.SortStatistics(answer, request.OrderBy)

	if err := output.OutputResults(answer, request.Format); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing result: ", err)
		os.Exit(1)
	}
}

func parseFlags() stats.RepoFlags {
	var request stats.RepoFlags
	pflag.StringVar(&request.Repository, "repository", ".", "path to the Git repository; defaults to the current directory")
	pflag.StringVar(&request.Revision, "revision", "HEAD", "commit reference; defaults to HEAD")
	pflag.StringVar(&request.OrderBy, "order-by", "lines", "key for sorting results; one of lines (default), commits, or files")
	pflag.BoolVar(&request.UseCommitter, "use-committer", false, "uses committer instead of author in calculations")
	pflag.StringVar(&request.Format, "format", "tabular", "output format: tabular, csv, json, or json-lines")
	pflag.StringSliceVar(&request.Extensions, "extensions", []string{}, "limits files by extension, e.g., '.go,.md'")
	pflag.StringSliceVar(&request.Languages, "languages", []string{}, "limits files by language, e.g., 'go,markdown'")
	pflag.StringSliceVar(&request.Exclude, "exclude", []string{}, "excludes files matching Glob patterns, e.g., 'foo/*,bar/*'")
	pflag.StringSliceVar(&request.RestrictTo, "restrict-to", []string{}, "includes only files matching Glob patterns")
	pflag.Parse()
	return request
}
