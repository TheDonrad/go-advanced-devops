package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSettingsList_setConfigFlags(t *testing.T) {

	settings := SettingsList{
		Addr:           "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
		Key:            "",
		RateLimit:      5,
	}

	os.Args = []string{"test", "-a", "localhost:1234", "-p", "30s", "-r", "60s", "-k", "hash", "-l", "10"}

	settings.setConfigFlags()

	want := SettingsList{
		Addr:           "localhost:1234",
		PollInterval:   30 * time.Second,
		ReportInterval: 60 * time.Second,
		Key:            "hash",
		RateLimit:      10,
	}

	t.Run("config from env", func(t *testing.T) {
		if !reflect.DeepEqual(settings, want) {
			t.Errorf("Config = %v, want %v", settings, &want)
		}
	})
}

func TestSettingsList_setConfigEnv(t *testing.T) {

	settings := SettingsList{
		Addr:           "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
		Key:            "",
		RateLimit:      5,
	}

	os.Setenv("ADDRESS", "localhost:123")
	defer os.Unsetenv("ADDRESS")

	os.Setenv("REPORT_INTERVAL", "60s")
	defer os.Unsetenv("REPORT_INTERVAL")

	os.Setenv("POLL_INTERVAL", "20s")
	defer os.Unsetenv("POLL_INTERVAL")

	os.Setenv("KEY", "hash")
	defer os.Unsetenv("Key")

	os.Setenv("RATE_LIMIT", "10")
	defer os.Unsetenv("RateLimit")

	settings.setConfigEnv()

	want := SettingsList{
		Addr:           "localhost:123",
		PollInterval:   20 * time.Second,
		ReportInterval: 60 * time.Second,
		Key:            "hash",
		RateLimit:      10,
	}

	t.Run("config from env", func(t *testing.T) {
		if got := Config(false); !reflect.DeepEqual(got, &want) {
			t.Errorf("Config() = %v, want %v", got, &want)
		}
	})
}

func TestConfig(t *testing.T) {
	want := SettingsList{
		Addr:           "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
		Key:            "",
		RateLimit:      5,
	}
	os.Setenv("ADDRESS", "")
	defer os.Unsetenv("ADDRESS")

	os.Setenv("REPORT_INTERVAL", "")
	defer os.Unsetenv("REPORT_INTERVAL")

	os.Setenv("POLL_INTERVAL", "")
	defer os.Unsetenv("POLL_INTERVAL")

	os.Setenv("KEY", "")
	defer os.Unsetenv("Key")

	os.Setenv("RATE_LIMIT", "")
	defer os.Unsetenv("RateLimit")

	t.Run("config test", func(t *testing.T) {
		if got := Config(false); !reflect.DeepEqual(got, &want) {
			t.Errorf("Config() = %v, want %v", got, &want)
		}
	})

}
