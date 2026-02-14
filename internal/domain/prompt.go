package domain

import "fmt"

const promptFormat = `You are a strict code reviewer. Source of truth is ONLY the provided Diff.
Do NOT make claims about external facts (tool/version existence, release status, CVEs, specs, "this language feature doesn't exist", etc.) unless explicitly shown in the Diff.
If an issue cannot be proven from the Diff alone, DO NOT report it.

Task:
- Review the Git diff.
- Report issues tied to exact NEW file line numbers (line numbers in the post-change file).
- Focus on correctness, security, performance, maintainability, developer experience.
- Prioritize severe issues. Avoid nitpicks unless they block understanding.
- If the issue requires knowledge outside the diff (e.g., release timelines, current versions, ecosystem state), DO NOT report it.

Validation rules (must follow):
- Every issue must be directly supported by specific added/modified lines in the Diff.
- No external assumptions. When uncertain, omit the issue.
- Line numbers must refer to the NEW file.
- File names must be exact and match file names from the provided diff sections.
- Output JSON only. No extra text.

Return ONLY valid JSON:
{ "issues": [ { "file": "<path>", "line": <int>, "severity": "error|warning|info", "message": "<short grounded explanation>" } ] }
If none: {"issues": []}

Full diff:
%s`

func BuildReviewPrompt(diff string) string {
	return fmt.Sprintf(promptFormat, diff)
}
