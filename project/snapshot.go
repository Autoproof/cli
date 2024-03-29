package project

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type ignoreList struct {
	ignores map[string]struct{}
}

func newIgnoreList(ignores ...string) *ignoreList {
	il := ignoreList{
		ignores: make(map[string]struct{}),
	}
	for _, ignore := range ignores {
		il.ignores[ignore] = struct{}{}
	}
	return &il
}

func (il *ignoreList) has(s string) bool {
	_, exists := il.ignores[s]
	return exists
}

func (p *Project) Snapshot(ctx context.Context) (*Snapshot, error) {
	var snapshot Snapshot

	il := newIgnoreList(p.Config().GetStringSlice("ignore")...)
	err := filepath.Walk(p.path, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while walking file/directory: %w", err)
		}

		// interrupt snapshotting if context has been cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if il.has(filepath.Base(filePath)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer func() { _ = file.Close() }()

		h := newKeccakHash()
		if _, err := io.Copy(h, file); err != nil {
			return fmt.Errorf("calculate file hash for %s: %w", filePath, err)
		}

		relFilePath, err := filepath.Rel(p.path, filePath)
		if err != nil {
			return fmt.Errorf("file path relative to project root: %w", err)
		}
		snapshot.Items = append(snapshot.Items, SnapshotItem{
			Filename: relFilePath,
			Hash:     h.String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
