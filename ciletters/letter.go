//go:build !solution

package ciletters

import (
	"strings"
	"text/template"
)

var letterTemplate = `Your pipeline #{{ .Pipeline.ID}} {{ if ne .Pipeline.Status "ok" }}has failed{{ else }}passed{{ end }}!
    Project:      {{ .Project.GroupID }}/{{ .Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ slice .Commit.Hash 0 8 }} {{ .Commit.Message }}
    CommitAuthor: {{ .Commit.Author }}{{ range $job := .Pipeline.FailedJobs }}
        Stage: {{ $job.Stage }}, Job {{ $job.Name }}{{ range prepareLog $job.RunnerLog }}
            {{ . }}{{ end }}
{{ end }}`

const skipLines = 9

func prepareLog(s string) []string {
	lines := strings.Split(s, "\n")
	if len(lines) > skipLines {
		return lines[skipLines:]
	}
	return lines
}

func createTemplate() (*template.Template, error) {
	return template.
		New("email").
		Funcs(template.FuncMap{"prepareLog": prepareLog}).
		Parse(letterTemplate)
}

func executeTemplate(templt *template.Template, n *Notification) (string, error) {
	builder := strings.Builder{}
	if err := templt.Execute(&builder, n); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func MakeLetter(n *Notification) (string, error) {
	templt, err := createTemplate()
	if err != nil {
		return "", err
	}
	return executeTemplate(templt, n)
}
