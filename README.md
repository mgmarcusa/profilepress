# profilepress-printed

This directory is the Printing Press generated `profilepress` CLI, augmented with the safety workflow from the CEE plan.

Generated with:

```bash
cli-printing-press generate --plan /home/mgmarcusa/docs/plans/2026-06-26-001-feat-linkedin-private-profile-cli-printing-press-plan.md --name profilepress --output /home/mgmarcusa/profilepress-printed --json
```

Implemented commands:

- `snapshot`
- `privacy-check`
- `propose-for-job`
- `diff`
- `packet export`
- `apply-packet`
- `auth status`
- `doctor`

Live LinkedIn mutation is intentionally blocked until a tested user-controlled browser adapter is added.
