package cli

import (
	"context"
	"time"

	"profilepress-pp-cli/internal/browser"
	"profilepress-pp-cli/internal/linkedin"
	"profilepress-pp-cli/internal/profile"
)

func waitForManagedProfile(ctx context.Context, client *browser.Client, expectName string, wait time.Duration) ([]profile.Section, string, error) {
	deadline := time.Now().Add(wait)
	var lastErr error
	for time.Now().Before(deadline) {
		page, err := client.CapturePageText(ctx, "", "linkedin.com/in/", 0)
		if err == nil {
			sections, source, parseErr := linkedin.ParseProfilePageText(page, expectName)
			if parseErr == nil {
				return sections, source, nil
			}
			lastErr = parseErr
		} else {
			lastErr = err
		}
		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return nil, "", lastErr
}
