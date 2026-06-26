package store

import (
	"testing"

	"profilepress-pp-cli/internal/packet"
	"profilepress-pp-cli/internal/profile"
)

func TestSnapshotPacketAndApplyLogRoundTrip(t *testing.T) {
	s, err := Open(t.TempDir() + "/profilepress.db")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	snap, err := s.CreateSnapshot("test", []profile.Section{{Name: "headline", RawText: "Old", Normalized: "Old", Source: "test"}, {Name: "about", RawText: "", Source: "test"}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.GetSnapshot(snap.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Sections) != 2 {
		t.Fatalf("sections=%d", len(got.Sections))
	}

	pkt, err := s.CreatePacket(packet.Packet{SnapshotID: snap.ID, Opportunity: "job", Changes: []packet.Change{{Section: "headline", Before: "Old", After: "New", SourceNote: "resume"}}})
	if err != nil {
		t.Fatal(err)
	}
	gotPkt, err := s.GetPacket(pkt.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotPkt.Changes[0].After != "New" {
		t.Fatalf("packet did not round trip: %#v", gotPkt.Changes)
	}

	log, err := s.AddApplyLog(packet.ApplyLog{PacketID: pkt.ID, PrivacyStatus: "passed", SensitiveStatus: "confirmed", Result: "dry-run-passed", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if log.ID == "" {
		t.Fatal("missing apply log id")
	}

	again, err := s.GetPacket(pkt.ID)
	if err != nil {
		t.Fatal(err)
	}
	if again.Changes[0].Before != "Old" {
		t.Fatal("apply log mutated packet history")
	}
}
