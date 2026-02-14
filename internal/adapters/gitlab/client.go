package gitlab

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

type Client struct {
	config    domain.ConfigPort
	git       *gitlab.Client
	projectID string
	iid       int
}

func NewClient(url, token, projectID string, iid int, config domain.ConfigPort) (*Client, error) {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}

	return &Client{
		git:       git,
		projectID: projectID,
		iid:       iid,
		config:    config,
	}, nil
}

func (c *Client) GetMergeRequestChanges() ([]domain.Diff, error) {
	changes, _, err := c.git.MergeRequests.ListMergeRequestDiffs(
		c.projectID,
		int64(c.iid),
		&gitlab.ListMergeRequestDiffsOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get MR changes: %w", err)
	}

	var diffs []domain.Diff
	for _, d := range changes {
		diffs = append(diffs, domain.Diff{
			NewPath: d.NewPath,
			OldPath: d.OldPath,
			Content: d.Diff,
		})
	}

	return diffs, nil
}

func (c *Client) GetExistingComments() (map[string][]string, error) {
	notes, _, err := c.git.Notes.ListMergeRequestNotes(
		c.projectID,
		int64(c.iid),
		&gitlab.ListMergeRequestNotesOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	existing := make(map[string][]string)

	for _, note := range notes {
		if note.System {
			continue
		}

		if note.Position == nil {
			continue
		}

		key := fmt.Sprintf("%s:%d", note.Position.NewPath, note.Position.NewLine)
		existing[key] = append(existing[key], note.Body)
	}

	return existing, nil
}

func (c *Client) AddMergeRequestDiscussion(file string, line int, note string) error {
	commitSHA := c.config.GetCommitSHA()
	baseSHA := c.config.GetMergeRequestDiffBaseSHA()
	positionType := "line"

	line64 := int64(line)

	_, _, err := c.git.Discussions.CreateMergeRequestDiscussion(
		c.projectID,
		int64(c.iid),
		&gitlab.CreateMergeRequestDiscussionOptions{
			Body: &note,
			Position: &gitlab.PositionOptions{
				PositionType: &positionType,
				NewLine:      &line64,
				HeadSHA:      &commitSHA,
				BaseSHA:      &baseSHA,
				StartSHA:     &commitSHA,
				OldPath:      &file,
				NewPath:      &file,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add discussion: %w", err)
	}

	return nil
}

func (c *Client) DeleteBotCommentsExceptResolved() error {
	notes, _, err := c.git.Notes.ListMergeRequestNotes(
		c.projectID,
		int64(c.iid),
		&gitlab.ListMergeRequestNotesOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	for _, note := range notes {
		if note.System || note.Resolved {
			continue
		}

		_, err := c.git.Notes.DeleteMergeRequestNote(c.projectID, int64(c.iid), note.ID)
		if err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
	}

	return nil
}

var _ domain.MRProviderPort = (*Client)(nil)
