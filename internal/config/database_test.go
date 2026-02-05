package config

import (
	"os"
	"os/exec"
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	val := getEnv("ENV_NOT_EXIST", "default")
	if val != "default" {
		t.Fatalf("expected default, got %s", val)
	}
}

func TestGetEnv_FromEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "value")
	defer os.Unsetenv("TEST_ENV")

	val := getEnv("TEST_ENV", "default")
	if val != "value" {
		t.Fatalf("expected value, got %s", val)
	}
}
func TestInitDB_Fail(t *testing.T) {
	if os.Getenv("TEST_INITDB_FAIL") == "1" {
		os.Setenv("DB_HOST", "invalid_host")
		os.Setenv("DB_PORT", "9999")
		os.Setenv("DB_USER", "invalid")
		os.Setenv("DB_PASSWORD", "invalid")
		os.Setenv("DB_NAME", "invalid")

		InitDB()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitDB_Fail")
	cmd.Env = append(os.Environ(), "TEST_INITDB_FAIL=1")

	err := cmd.Run()

	if err == nil {
		t.Fatalf("expected InitDB to exit with error")
	}
}
