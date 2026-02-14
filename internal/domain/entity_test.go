package domain

import (
	"testing"
)

func TestEntitiesCanBeConstructed(t *testing.T) {
	mr := MergeRequest{IID: 10, ProjectID: "42", Title: "MR", WebURL: "https://gitlab/mr/10"}
	diff := Diff{NewPath: "a.go", OldPath: "a.go", Content: "@@ -1 +1 @@"}
	issue := ReviewIssue{Line: 12, Severity: "warning", Message: "msg"}
	review := Review{FilePath: "a.go", Language: "Go", Issues: []ReviewIssue{issue}}
	comment := Comment{FilePath: "a.go", Line: 12, Body: "text"}

	if mr.IID != 10 || diff.NewPath == "" || review.Language != "Go" || comment.Line != 12 {
		t.Fatal("entity values are not set as expected")
	}
}
