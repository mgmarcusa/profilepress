package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"profilepress-pp-cli/internal/linkedin"
	"profilepress-pp-cli/internal/packet"
)

func newApplyPacketCmd() *cobra.Command {
	var dbPath, packetID, privacyRaw, confirmSensitive, confirmApply string
	var override, dryRun, simulateLive bool
	cmd := &cobra.Command{
		Use:   "apply-packet",
		Short: "Apply an approved packet only after privacy preflight, sensitive-change confirmation, and final user approval.",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := openStore(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()
			p, err := packetByIDOrLatest(db, packetID)
			if err != nil {
				return err
			}
			privacyStatus, err := linkedin.RequirePrivacyPassed(privacyRaw, override)
			if err != nil {
				return err
			}
			sensitive := packet.SensitiveChanges(p.Changes)
			sensitiveStatus := "none"
			if len(sensitive) > 0 {
				sensitiveStatus = "requires-confirmation"
				if confirmSensitive != "APPLY-SENSITIVE" {
					return fmt.Errorf("packet has %d sensitive change(s); pass --confirm-sensitive APPLY-SENSITIVE", len(sensitive))
				}
				sensitiveStatus = "confirmed"
			}
			if !dryRun && confirmApply != "APPLY" {
				return errors.New("non-dry-run apply requires --confirm-apply APPLY")
			}
			adapter := linkedin.ApplyAdapter(linkedin.NotImplementedAdapter{})
			if dryRun || simulateLive {
				adapter = linkedin.DryRunAdapter{}
			}
			for _, ch := range p.Changes {
				if err := adapter.Apply(ch.Section, ch.After); err != nil {
					return err
				}
			}
			result := "applied"
			if dryRun {
				result = "dry-run-passed"
			}
			log, err := db.AddApplyLog(packet.ApplyLog{PacketID: p.ID, PrivacyStatus: string(privacyStatus), SensitiveStatus: sensitiveStatus, Result: result, DryRun: dryRun, ConfirmationSource: "cli"})
			if err != nil {
				return err
			}
			return printJSON(map[string]any{"apply_log_id": log.ID, "packet_id": p.ID, "result": result, "privacy_status": privacyStatus, "sensitive_status": sensitiveStatus})
		},
	}
	cmd.Flags().StringVar(&dbPath, "db", "", "database path")
	cmd.Flags().StringVar(&packetID, "packet", "", "packet ID")
	cmd.Flags().StringVar(&privacyRaw, "privacy-status", "unknown", "disabled|enabled|unknown")
	cmd.Flags().BoolVar(&override, "override-privacy-risk", false, "override failed/unknown privacy status")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "run checks without writing")
	cmd.Flags().BoolVar(&simulateLive, "simulate-live", false, "test-only live adapter")
	cmd.Flags().StringVar(&confirmSensitive, "confirm-sensitive", "", "must equal APPLY-SENSITIVE for sensitive changes")
	cmd.Flags().StringVar(&confirmApply, "confirm-apply", "", "must equal APPLY for non-dry-run apply")
	return cmd
}
