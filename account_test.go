package main

import (
	"path"
	"testing"
)

type T = *testing.T

func init() {
	rsaKeyLength = 1024 // for faster account generation
}

func TestNewAccount(t *testing.T) {
	t.Run("Creating an account with valid parameters", func(t T) {
		dir := t.TempDir()
		account, err := NewAccount(dir, "Ahmed Mohamed", "ahmed@example.com", "strong password")

		ASSERT(t, err == nil, "Error was returned when creating an account: %s", err)
		ASSERT(t, account != nil, "Account value is nil, expected a value")

		t.Run("Include correct information", func(t T) {
			identity, _ := account.Identity()
			pgpkey, _ := account.Export()

			ASSERT_EQUAL(t, "ahmed@example.com", account.Email())
			ASSERT_EQUAL(t, "Ahmed Mohamed", account.Name())
			ASSERT_EQUAL(t, "Ahmed Mohamed <ahmed@example.com>", identity)
			REFUTE_EQUAL(t, 0, len(pgpkey))
		})

		t.Run("Creates the correct file structure", func(t T) {
			ASSERT_DIR_EXISTS(t, path.Join(dir, ".mau"))
			ASSERT_FILE_EXISTS(t, path.Join(dir, ".mau", "account.pgp"))
		})
	})

	t.Run("Creating an account without a password", func(t T) {
		dir := t.TempDir()
		account, err := NewAccount(dir, "Ahmed Mohamed", "ahmed@example.com", "")

		ASSERT_ERROR(t, ErrPassphraseRequired, err)
		ASSERT_EQUAL(t, nil, account)
	})

	t.Run("Creating an account in an existing account directory", func(t T) {
		dir := t.TempDir()
		NewAccount(dir, "Ahmed Mohamed", "ahmed@example.com", "password")
		account, err := NewAccount(dir, "Ahmed Mahmoud", "ahmed.mahmoud@example.com", "password")

		ASSERT(t, err == ErrAccountAlreadyExists, "Expected an error: %s Got: %s", ErrAccountAlreadyExists, err)
		ASSERT(t, account == nil, "Expected the account to be nil value got : %v", account)
	})

	t.Run("Two accounts with same identity", func(t T) {
		account1, _ := NewAccount(t.TempDir(), "Ahmed Mohamed", "ahmed@example.com", "password")
		account2, _ := NewAccount(t.TempDir(), "Ahmed Mohamed", "ahmed@example.com", "password")

		REFUTE_EQUAL(t, account1.Fingerprint(), account2.Fingerprint())
	})
}

func TestOpenAccount(t *testing.T) {
	dir := t.TempDir()
	account, _ := NewAccount(dir, "Ahmed Mohamed", "ahmed@example.com", "strong password")

	t.Run("Using same password", func(t T) {
		opened, err := OpenAccount(dir, "strong password")
		ASSERT_ERROR(t, nil, err)
		ASSERT_EQUAL(t, "ahmed@example.com", opened.Email())
		ASSERT_EQUAL(t, "Ahmed Mohamed", opened.Name())
		ASSERT_EQUAL(t, account.Fingerprint(), opened.Fingerprint())
	})
}
