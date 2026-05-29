package service

import (
	"testing"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

func TestBuildReferralTree(t *testing.T) {
	// Setup data sesuai dengan contoh kasus dari user:
	// - User A (123)
	//   - User B (456)
	//     - User C (789)
	//   - User D (012)

	profiles := []model.UserProfile{
		{UserID: 1, FullName: "User A", PhoneNumber: "123", ReferralNumber: ""},
		{UserID: 2, FullName: "User B", PhoneNumber: "456", ReferralNumber: "123"},
		{UserID: 3, FullName: "User C", PhoneNumber: "789", ReferralNumber: "456"},
		{UserID: 4, FullName: "User D", PhoneNumber: "012", ReferralNumber: "123"},
	}

	roots := BuildReferralTree(profiles)

	// Verifikasi jumlah root node
	if len(roots) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(roots))
	}

	root := roots[0]
	if root.WhatsAppNumber != "123" || root.FullName != "User A" {
		t.Errorf("Expected root to be User A (123), got %s (%s)", root.FullName, root.WhatsAppNumber)
	}

	// Verifikasi anak-anak dari Root (User A)
	if len(root.Children) != 2 {
		t.Fatalf("Expected User A to have 2 children, got %d", len(root.Children))
	}

	// Cari User B dan User D di dalam anak User A
	var childB, childD *model.ReferralNode
	for _, child := range root.Children {
		if child.WhatsAppNumber == "456" {
			childB = child
		} else if child.WhatsAppNumber == "012" {
			childD = child
		}
	}

	if childB == nil {
		t.Error("Expected to find User B (456) as child of User A")
	} else {
		// Verifikasi anak dari User B (yaitu User C)
		if len(childB.Children) != 1 {
			t.Fatalf("Expected User B to have 1 child, got %d", len(childB.Children))
		}
		childC := childB.Children[0]
		if childC.WhatsAppNumber != "789" || childC.FullName != "User C" {
			t.Errorf("Expected child of User B to be User C (789), got %s (%s)", childC.FullName, childC.WhatsAppNumber)
		}
		if len(childC.Children) != 0 {
			t.Errorf("Expected User C to have 0 children, got %d", len(childC.Children))
		}
	}

	if childD == nil {
		t.Error("Expected to find User D (012) as child of User A")
	} else {
		if len(childD.Children) != 0 {
			t.Errorf("Expected User D to have 0 children, got %d", len(childD.Children))
		}
	}
}

func TestBuildReferralTree_DisconnectedParent(t *testing.T) {
	// Mengetes kasus ketika anak memiliki referral number tetapi
	// nomor referral tersebut tidak terdaftar di database.
	// Seharusnya anak tersebut terangkat menjadi root node sendiri.

	profiles := []model.UserProfile{
		{UserID: 1, FullName: "User A", PhoneNumber: "123", ReferralNumber: ""},
		{UserID: 2, FullName: "User B", PhoneNumber: "456", ReferralNumber: "999"},
	}

	roots := BuildReferralTree(profiles)

	// Seharusnya ada 2 root karena User B terputus dari parent-nya
	if len(roots) != 2 {
		t.Fatalf("Expected 2 root nodes, got %d", len(roots))
	}
}
