package application

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/mistername/ai-mr-reviewer/internal/domain"
	"go.uber.org/zap"
)

type Reviewer struct {
	logger     *zap.Logger
	mrProvider domain.MRProviderPort
	aiProvider domain.AIProviderPort
	runtime    domain.RuntimeConfig
}

type reviewResponse struct {
	Issues []reviewIssuePayload `json:"issues"`
}

type reviewIssuePayload struct {
	FilePath string `json:"file"`
	Path     string `json:"path"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
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

func NewReviewer(runtime domain.RuntimeConfig, mrProvider domain.MRProviderPort, aiProvider domain.AIProviderPort, logger *zap.Logger) *Reviewer {
	return &Reviewer{runtime: runtime, mrProvider: mrProvider, aiProvider: aiProvider, logger: logger}
}

func (r *Reviewer) Run(ctx context.Context) error {
	if r.runtime.DeleteBotComments {
		if err := r.mrProvider.DeleteBotCommentsExceptResolved(ctx); err != nil {
			r.logger.Warn("cannot delete bot comments", zap.Error(err))
		}
	}

	existing, err := r.mrProvider.GetExistingComments(ctx)
	if err != nil {
		r.logger.Warn("cannot read existing comments", zap.Error(err))

		existing = map[string][]string{}
	}

	diffs, err := r.mrProvider.GetMergeRequestChanges(ctx)
	if err != nil {
		return fmt.Errorf("get MR changes: %w", err)
	}

	prefix := r.runtime.CommentPrefix + ":"

	filteredDiffs := r.filterNewDiffs(diffs, existing, prefix)
	if len(filteredDiffs) == 0 {
		return nil
	}

	if err := r.reviewDiffs(ctx, filteredDiffs); err != nil {
		return fmt.Errorf("review diffs: %w", err)
	}

	return nil
}

func (r *Reviewer) filterNewDiffs(diffs []domain.Diff, existing map[string][]string, prefix string) []domain.Diff {
	filtered := make([]domain.Diff, 0, len(diffs))

	for _, d := range diffs {
		if !hasExistingComments(d.NewPath, existing, prefix) {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

func hasExistingComments(path string, existing map[string][]string, prefix string) bool {
	for key, bodies := range existing {
		if strings.HasPrefix(key, path+":") {
			for _, body := range bodies {
				if strings.HasPrefix(body, prefix) {
					return true
				}
			}
		}
	}

	return false
}

func (r *Reviewer) reviewDiffs(ctx context.Context, diffs []domain.Diff) error {
	combinedDiff := buildCombinedDiff(diffs)

	reviewText, err := r.aiProvider.ReviewCode(ctx, combinedDiff)
	if err != nil {
		return fmt.Errorf("review code: %w", err)
	}

	issues, err := parseReviewResponse(reviewText)
	if err != nil {
		return fmt.Errorf("parse review response: %w", err)
	}

	knownFiles := make(map[string]struct{}, len(diffs))
	for _, d := range diffs {
		knownFiles[d.NewPath] = struct{}{}
	}

	prefix := r.runtime.CommentPrefix

	for _, issue := range issues {
		filePath := issue.FilePath
		if filePath == "" && len(knownFiles) == 1 {
			for onlyPath := range knownFiles {
				filePath = onlyPath
			}
		}

		if _, ok := knownFiles[filePath]; !ok {
			r.logger.Warn("skip issue for unknown file", zap.String("file", issue.FilePath), zap.Int("line", issue.Line))
			continue
		}

		body := fmt.Sprintf("%s:**%s**: %s", prefix, strings.ToUpper(issue.Severity), issue.Message)
		if err := r.mrProvider.AddMergeRequestDiscussion(ctx, filePath, issue.Line, body); err != nil {
			r.logger.Warn("failed to add comment", zap.String("path", filePath), zap.Int("line", issue.Line), zap.Error(err))
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
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	issues := make([]domain.ReviewIssue, 0, len(parsed.Issues))
	for _, issue := range parsed.Issues {
		filePath := issue.FilePath
		if filePath == "" {
			filePath = issue.Path
		}

		issues = append(issues, domain.ReviewIssue{
			FilePath: filePath,
			Severity: issue.Severity,
			Message:  issue.Message,
			Line:     issue.Line,
		})
	}

	return issues, nil
}

func buildCombinedDiff(diffs []domain.Diff) string {
	uniquePaths := make(map[string]struct{}, len(diffs))
	for _, d := range diffs {
		uniquePaths[d.NewPath] = struct{}{}
	}

	paths := slices.Collect(maps.Keys(uniquePaths))
	sort.Strings(paths)

	var builder strings.Builder

	for _, path := range paths {
		for _, d := range diffs {
			if d.NewPath != path {
				continue
			}

			builder.WriteString("File: ")
			builder.WriteString(d.NewPath)
			builder.WriteString("\nLanguage: ")
			builder.WriteString(detectLanguage(d.NewPath))
			builder.WriteString("\nDiff:\n")
			builder.WriteString(d.Content)
			builder.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(builder.String())
}

func extractJSON(s string) string {
	if idx := strings.Index(s, "```json"); idx != -1 {
		end := strings.Index(s[idx+7:], "```")
		if end != -1 {
			return s[idx+7 : idx+7+end]
		}
	}

	if _, after, ok := strings.Cut(s, "```"); ok {
		content := after
		if before, _, ok := strings.Cut(content, "```"); ok {
			extracted := strings.TrimSpace(before)
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
