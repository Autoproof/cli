package autoproofapi

import (
	"context"
	"net/http"
)

type SnapshotMode string

const (
	ProductionSnapshotMode SnapshotMode = "live"
	TestingSnapshotMode    SnapshotMode = "test"
)

type UploadSnapshotResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type SnapshotItem struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
}

type Snapshot struct {
	Project     string         `json:"project"`
	Description string         `json:"description,omitempty"`
	Mode        SnapshotMode   `json:"mode"`
	Data        []SnapshotItem `json:"data"`
}

func (c *Client) UploadSnapshot(ctx context.Context, snapshot *Snapshot) (*UploadSnapshotResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "v1/upload-snapshot", snapshot)
	if err != nil {
		return nil, err
	}

	var snapshotResponse UploadSnapshotResponse
	if err := c.do(req, &snapshotResponse); err != nil {
		return nil, err
	}

	return &snapshotResponse, nil
}
