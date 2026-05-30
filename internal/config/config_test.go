package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.App.Name != "petverse-be" {
		t.Fatalf("App.Name = %q, want petverse-be", cfg.App.Name)
	}

	if cfg.HTTP.Port != 8080 {
		t.Fatalf("HTTP.Port = %d, want 8080", cfg.HTTP.Port)
	}

	if cfg.App.Env != "local" {
		t.Fatalf("App.Env = %q, want local", cfg.App.Env)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "70000")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want invalid port error")
	}
}

func TestLoadRequiresCriticalConfigOutsideLocal(t *testing.T) {
	t.Setenv("APP_ENV", "staging")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing staging configuration error")
	}
}

func TestLoadRejectsUnknownEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "qa")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want invalid environment error")
	}
}
