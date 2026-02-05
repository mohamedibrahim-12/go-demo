package tests

import (
	"testing"
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/repositories"
	"go-demo/worker"
)

// TestCleanupWorkerDeletesOldRecords inserts old and new records, runs cleanup
// once with a short retention, and asserts that old records are deleted and
// newer ones remain. No real ticker or long sleeps are used.
func TestCleanupWorkerDeletesOldRecords(t *testing.T) {
	// Ensure created_at column exists (migration v2 runs in TestMain)
	if database.GormDB == nil {
		t.Fatal("database not connected")
	}

	// Create a "new" user via repository (created_at = now)
	newUser := models.User{Name: "New Cleanup User", Role: "Tester"}
	if err := repositories.CreateUser(newUser); err != nil {
		t.Fatalf("create new user: %v", err)
	}

	// Insert an "old" user via raw SQL with created_at in the past (10 days ago)
	// so it will be deleted when retention is 7 days
	err := database.GormDB.Exec(
		`INSERT INTO users (name, role, uuid, created_at) VALUES (?, ?, gen_random_uuid(), NOW() - INTERVAL '10 days')`,
		"Old Cleanup User",
		"Role",
	).Error
	if err != nil {
		t.Fatalf("insert old user: %v", err)
	}

	// Run cleanup once with 7-day retention (no ticker, no sleep)
	worker.RunCleanupOnce(7 * 24 * time.Hour)

	// Assert: old record is gone, new record remains
	users, err := repositories.GetUsers()
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	var foundNew, foundOld bool
	for _, u := range users {
		if u.Name == "New Cleanup User" {
			foundNew = true
		}
		if u.Name == "Old Cleanup User" {
			foundOld = true
		}
	}
	if !foundNew {
		t.Error("expected new user to remain after cleanup")
	}
	if foundOld {
		t.Error("expected old user to be deleted by cleanup")
	}
}

// TestCleanupWorkerKeepsRecentRecords inserts two users with "recent" created_at
// and runs cleanup with a long retention; both should remain.
func TestCleanupWorkerKeepsRecentRecords(t *testing.T) {
	if database.GormDB == nil {
		t.Fatal("database not connected")
	}

	// Create two users via repository (both have created_at = now)
	u1 := models.User{Name: "Recent User 1", Role: "A"}
	u2 := models.User{Name: "Recent User 2", Role: "B"}
	if err := repositories.CreateUser(u1); err != nil {
		t.Fatalf("create user 1: %v", err)
	}
	if err := repositories.CreateUser(u2); err != nil {
		t.Fatalf("create user 2: %v", err)
	}

	// Run cleanup with 1-year retention (nothing should be deleted)
	worker.RunCleanupOnce(365 * 24 * time.Hour)

	users, err := repositories.GetUsers()
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	var found1, found2 bool
	for _, u := range users {
		if u.Name == "Recent User 1" {
			found1 = true
		}
		if u.Name == "Recent User 2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("expected both recent users to remain; found Recent User 1=%v Recent User 2=%v", found1, found2)
	}
}

// TestCleanupProductsDeletesOld inserts an old product via raw SQL, runs
// cleanup once, and asserts the old product is deleted.
func TestCleanupProductsDeletesOld(t *testing.T) {
	if database.GormDB == nil {
		t.Fatal("database not connected")
	}

	// New product via repository
	p := models.Product{Name: "New Product", Price: 99.99}
	if err := repositories.CreateProduct(p); err != nil {
		t.Fatalf("create new product: %v", err)
	}

	// Old product via raw SQL (created_at 10 days ago)
	err := database.GormDB.Exec(
		`INSERT INTO products (name, price, uuid, created_at) VALUES (?, ?, gen_random_uuid(), NOW() - INTERVAL '10 days')`,
		"Old Cleanup Product",
		1.99,
	).Error
	if err != nil {
		t.Fatalf("insert old product: %v", err)
	}

	worker.RunCleanupOnce(7 * 24 * time.Hour)

	products, err := repositories.GetProducts()
	if err != nil {
		t.Fatalf("get products: %v", err)
	}
	var foundNew, foundOld bool
	for _, p := range products {
		if p.Name == "New Product" {
			foundNew = true
		}
		if p.Name == "Old Cleanup Product" {
			foundOld = true
		}
	}
	if !foundNew {
		t.Error("expected new product to remain after cleanup")
	}
	if foundOld {
		t.Error("expected old product to be deleted by cleanup")
	}
}
