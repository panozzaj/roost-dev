package server

import "testing"

func TestParseServiceName(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantServiceName string
		wantAppName     string
		wantOk          bool
	}{
		{
			name:            "service-app format",
			input:           "web-myapp",
			wantServiceName: "web",
			wantAppName:     "myapp",
			wantOk:          true,
		},
		{
			name:            "multi-dash app name",
			input:           "api-my-cool-app",
			wantServiceName: "api",
			wantAppName:     "my-cool-app",
			wantOk:          true,
		},
		{
			name:            "simple app name (no dash)",
			input:           "myapp",
			wantServiceName: "",
			wantAppName:     "myapp",
			wantOk:          false,
		},
		{
			name:            "empty string",
			input:           "",
			wantServiceName: "",
			wantAppName:     "",
			wantOk:          false,
		},
		{
			name:            "leading dash",
			input:           "-myapp",
			wantServiceName: "",
			wantAppName:     "myapp",
			wantOk:          true,
		},
		{
			name:            "trailing dash",
			input:           "myapp-",
			wantServiceName: "myapp",
			wantAppName:     "",
			wantOk:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceName, appName, ok := parseServiceName(tt.input)
			if serviceName != tt.wantServiceName {
				t.Errorf("parseServiceName(%q) serviceName = %q, want %q", tt.input, serviceName, tt.wantServiceName)
			}
			if appName != tt.wantAppName {
				t.Errorf("parseServiceName(%q) appName = %q, want %q", tt.input, appName, tt.wantAppName)
			}
			if ok != tt.wantOk {
				t.Errorf("parseServiceName(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
		})
	}
}
