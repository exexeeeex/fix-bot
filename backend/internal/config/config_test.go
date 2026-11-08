package config

import "testing"

func TestEnvFileFor(t *testing.T) {
	tests := []struct {
		appEnv string
		want   string
	}{
		{appEnv: "development", want: ".env.dev"},
		{appEnv: "dev", want: ".env.dev"},
		{appEnv: "production", want: ".env.prod"},
		{appEnv: "prod", want: ".env.prod"},
	}

	for _, tt := range tests {
		t.Run(tt.appEnv, func(t *testing.T) {
			got, err := envFileFor(tt.appEnv)
			if err != nil {
				t.Fatalf("envFileFor(%q) error = %v", tt.appEnv, err)
			}
			if got != tt.want {
				t.Errorf("envFileFor(%q) = %q, want %q", tt.appEnv, got, tt.want)
			}
		})
	}
}

func TestEnvFileForRejectsUnknownEnvironment(t *testing.T) {
	if _, err := envFileFor("staging"); err == nil {
		t.Fatal("envFileFor(staging) error = nil, want error")
	}
}
