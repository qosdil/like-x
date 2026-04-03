package model

import "testing"

func TestModelConstants(t *testing.T) {
	if FullNameMinLength != 8 {
		t.Fatalf("expected FullNameMinLength 8, got %d", FullNameMinLength)
	}
	if FullNameMaxlength != 32 {
		t.Fatalf("expected FullNameMaxlength 32, got %d", FullNameMaxlength)
	}
	if PasswordMinLength != 8 {
		t.Fatalf("expected PasswordMinLength 8, got %d", PasswordMinLength)
	}
	if PasswordMaxLength != 16 {
		t.Fatalf("expected PasswordMaxLength 16, got %d", PasswordMaxLength)
	}
}
