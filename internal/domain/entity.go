package domain

type MergeRequest struct {
	ProjectID string
	Title     string
	WebURL    string
	IID       int
}

type Diff struct {
	NewPath string
	OldPath string
	Content string
}

type ReviewIssue struct {
	Severity string
	Message  string
	Line     int
}

type Review struct {
	FilePath string
	Language string
	Issues   []ReviewIssue
}

type Comment struct {
	FilePath string
	Body     string
	Line     int
}
