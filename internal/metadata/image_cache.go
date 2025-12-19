package metadata

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const tmdbImageBaseURL = "https://image.tmdb.org/t/p/original"

func cacheTMDBImage(
	ctx context.Context,
	baseDir string,
	kind string,
	id int64,
	name string,
	tmdbPath string,
) (string, error) {

	base := filepath.Join(baseDir, kind, fmt.Sprint(id))

	ext := filepath.Ext(tmdbPath)
	if ext == "" {
		ext = ".jpg"
	}

	fullPath := filepath.Join(base, name+ext)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Idempotent
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, nil
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		tmdbImageBaseURL+tmdbPath,
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image download failed: %s", resp.Status)
	}

	tmp := fullPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return "", err
	}
	f.Close()

	if err := os.Rename(tmp, fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
