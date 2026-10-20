package rules

type Issue struct {
	Path    string
	Line    int
	Column  int
	Message string
}
