package application

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"go.uber.org/zap"
)

type Reviewer struct {
	config     domain.ConfigPort
	mrProvider domain.MRProviderPort
	aiProvider domain.AIProviderPort
	logger     *zap.Logger
}

type reviewResponse struct {
	Issues []domain.ReviewIssue `json:"issues"`
}

var languageMap = map[string]string{
	".go":    "Go",
	".js":    "JavaScript/TypeScript",
	".jsx":   "JavaScript/TypeScript",
	".ts":    "JavaScript/TypeScript",
	".tsx":   "JavaScript/TypeScript",
	".py":    "Python",
	".java":  "Java",
	".rs":    "Rust",
	".c":     "C",
	".h":     "C",
	".cpp":   "C++",
	".hpp":   "C++",
	".cc":    "C++",
	".rb":    "Ruby",
	".php":   "PHP",
	".swift": "Swift",
	".kt":    "Kotlin",
	".kts":   "Kotlin",
	".scala": "Scala",
	".sh":    "Shell",
	".bash":  "Shell",
	".sql":   "SQL",
	".yaml":  "YAML",
	".yml":   "YAML",
	".json":  "JSON",
	".xml":   "XML",
	".md":    "Markdown",
}

func NewReviewer(config domain.ConfigPort, mrProvider domain.MRProviderPort, aiProvider domain.AIProviderPort, logger *zap.Logger) *Reviewer {
	return &Reviewer{config: config, mrProvider: mrProvider, aiProvider: aiProvider, logger: logger}
}

func (r *Reviewer) Run() error {
	existing, err := r.mrProvider.GetExistingComments()
	if err != nil {
		r.logger.Warn("cannot read existing comments", zap.Error(err))

		existing = map[string][]string{}
	}

	diffs, err := r.mrProvider.GetMergeRequestChanges()
	if err != nil {
		return fmt.Errorf("get MR changes: %w", err)
	}

	for _, d := range r.filterNewDiffs(diffs, existing) {
		if err := r.reviewDiff(d); err != nil {
			r.logger.Warn("review failed", zap.String("path", d.NewPath), zap.Error(err))
		}
	}

	return nil
}

func (r *Reviewer) filterNewDiffs(diffs []domain.Diff, existing map[string][]string) []domain.Diff {
	filtered := make([]domain.Diff, 0, len(diffs))

	for _, d := range diffs {
		key := fmt.Sprintf("%s:1", d.NewPath)
		if _, ok := existing[key]; !ok {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

func (r *Reviewer) reviewDiff(d domain.Diff) error {
	reviewText, err := r.aiProvider.ReviewCode(filepath.Base(d.NewPath), d.Content, detectLanguage(d.NewPath))
	if err != nil {
		return fmt.Errorf("review code: %w", err)
	}

	issues, err := parseReviewResponse(reviewText)
	if err != nil {
		return fmt.Errorf("parse AI review response: %w", err)
	}

	for _, issue := range issues {
		body := fmt.Sprintf("**%s**: %s", strings.ToUpper(issue.Severity), issue.Message)
		if err := r.mrProvider.AddMergeRequestDiscussion(d.NewPath, issue.Line, body); err != nil {
			r.logger.Warn("failed to add comment", zap.String("path", d.NewPath), zap.Int("line", issue.Line), zap.Error(err))
		}
	}

	return nil
}

func parseReviewResponse(response string) ([]domain.ReviewIssue, error) {
	trimmed := strings.TrimSpace(response)
	start := strings.Index(trimmed, "{")

	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("json object not found")
	}

	var parsed reviewResponse
	if err := json.Unmarshal([]byte(trimmed[start:end+1]), &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return parsed.Issues, nil
}

func detectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	return "Unknown"
}
