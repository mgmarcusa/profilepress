package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureRun(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmdErr := func() error {
		// Execute uses os.Stdout through cobra/default helpers.
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = append([]string{"profilepress"}, args...)
		return Execute()
	}()
	w.Close()
	os.Stdout = oldOut
	_, _ = buf.ReadFrom(r)
	return buf.String(), cmdErr
}

func TestPrintedCLIWorkflow(t *testing.T) {
	dir := t.TempDir()
	db := filepath.Join(dir, "pp.db")
	fixture := filepath.Join(dir, "profile.json")
	job := filepath.Join(dir, "job.txt")
	if err := os.WriteFile(fixture, []byte(`{"section_map":{"headline":"Old headline","about":"Old about"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(job, []byte("Research Engineer role"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := captureRun(t, "snapshot", "--db", db, "--fixture", fixture)
	if err != nil {
		t.Fatal(err)
	}
	var snap map[string]any
	if err := json.Unmarshal([]byte(out), &snap); err != nil {
		t.Fatal(err)
	}

	out, err = captureRun(t, "propose-for-job", "--db", db, "--snapshot", snap["snapshot_id"].(string), "--job-file", job, "--change", "headline=AI Evaluation and Privacy Research Leader", "--source-note", "resume")
	if err != nil {
		t.Fatal(err)
	}
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}

	out, err = captureRun(t, "diff", "--db", db, "--packet", pkt["packet_id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- Old headline") || !strings.Contains(out, "+ AI Evaluation and Privacy Research Leader") {
		t.Fatalf("bad diff: %s", out)
	}

	_, err = captureRun(t, "apply-packet", "--db", db, "--packet", pkt["packet_id"].(string), "--privacy-status", "unknown", "--dry-run")
	if err == nil {
		t.Fatal("unknown privacy should block")
	}

	out, err = captureRun(t, "apply-packet", "--db", db, "--packet", pkt["packet_id"].(string), "--privacy-status", "disabled", "--dry-run", "--confirm-sensitive", "APPLY-SENSITIVE")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "dry-run-passed") {
		t.Fatalf("bad apply: %s", out)
	}
}
