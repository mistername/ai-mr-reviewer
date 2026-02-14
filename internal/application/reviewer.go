package application

import (
	"encoding/json"
	"errors"
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
	if r.config.GetDeleteBotComments() {
		if err := r.mrProvider.DeleteBotCommentsExceptResolved(); err != nil {
			r.logger.Warn("cannot delete bot comments", zap.Error(err))
		}
	}

	existing, err := r.mrProvider.GetExistingComments()
	if err != nil {
		r.logger.Warn("cannot read existing comments", zap.Error(err))

		existing = map[string][]string{}
	}

	diffs, err := r.mrProvider.GetMergeRequestChanges()
	if err != nil {
		return fmt.Errorf("get MR changes: %w", err)
	}

	var diffErrors error

	for _, d := range r.filterNewDiffs(diffs, existing) {
		if err := r.reviewDiff(d); err != nil {
			diffErrors = errors.Join(diffErrors, fmt.Errorf("review failed %s: %w", d.NewPath, err))
			r.logger.Warn("review failed", zap.String("path", d.NewPath), zap.Error(err))
		}
	}

	return diffErrors
}

func (r *Reviewer) filterNewDiffs(diffs []domain.Diff, existing map[string][]string) []domain.Diff {
	filtered := make([]domain.Diff, 0, len(diffs))

	for _, d := range diffs {
		if !hasExistingComments(d.NewPath, existing) {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

func hasExistingComments(path string, existing map[string][]string) bool {
	for key := range existing {
		if strings.HasPrefix(key, path+":") {
			return true
		}
	}

	return false
}

const commentPrefix = "ai-mr-reviewer:"

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
		body := fmt.Sprintf("%s**%s**: %s", commentPrefix, strings.ToUpper(issue.Severity), issue.Message)
		if err := r.mrProvider.AddMergeRequestDiscussion(d.NewPath, issue.Line, body); err != nil {
			r.logger.Warn("failed to add comment", zap.String("path", d.NewPath), zap.Int("line", issue.Line), zap.Error(err))
		}
	}

	return nil
}

func parseReviewResponse(response string) ([]domain.ReviewIssue, error) {
	trimmed := strings.TrimSpace(response)

	jsonStr := extractJSON(trimmed)
	if jsonStr == "" {
		return nil, nil
	}

	var parsed reviewResponse
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, nil
	}

	return parsed.Issues, nil
}

func extractJSON(s string) string {
	if idx := strings.Index(s, "```json"); idx != -1 {
		end := strings.Index(s[idx+7:], "```")
		if end != -1 {
			return s[idx+7 : idx+7+end]
		}
	}

	if idx := strings.Index(s, "```"); idx != -1 {
		content := s[idx+3:]
		if end := strings.Index(content, "```"); end != -1 {
			extracted := strings.TrimSpace(content[:end])
			if isJSON(extracted) {
				return extracted
			}
		}
	}

	start := findJSONStart(s)
	if start == -1 {
		return ""
	}

	end := findMatchingBracket(s, start)
	if end == -1 {
		return ""
	}

	return s[start : end+1]
}

func findJSONStart(s string) int {
	openBrace := strings.Index(s, "{")
	openBracket := strings.Index(s, "[")

	switch {
	case openBrace != -1 && openBracket != -1:
		return min(openBrace, openBracket)
	case openBrace != -1:
		return openBrace
	case openBracket != -1:
		return openBracket
	default:
		return -1
	}
}

func isJSON(s string) bool {
	s = strings.TrimSpace(s)

	return (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"))
}

func findMatchingBracket(s string, start int) int {
	openChar := rune(s[start])
	closeChar := openChar

	if openChar == '{' {
		closeChar = '}'
	}

	count := 0

	for i, r := range s[start:] {
		switch r {
		case openChar:
			count++
		case closeChar:
			count--
			if count == 0 {
				return start + i
			}
		}
	}

	return -1
}

func detectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	return "Unknown"
}
