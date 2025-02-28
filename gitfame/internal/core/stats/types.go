package stats

type IntSet map[string]int

type StringSet map[string]string

type UserData struct {
	Files   int
	Lines   int
	Commits IntSet
	Name    string
}

type UserDataSet map[string]UserData

type RepoFlags struct {
	UseCommitter bool
	Repository   string
	Revision     string
	OrderBy      string
	Format       string
	Extensions   []string
	Languages    []string
	Exclude      []string
	RestrictTo   []string
}

type Language struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}
