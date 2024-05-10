package plugin

import "testing"

func TestEnsurePrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		input  string
		want   string
	}{
		{
			name:   "empty input",
			prefix: "pre_",
			input:  "",
			want:   "",
		},
		{
			name:   "input already has prefix",
			prefix: "pre_",
			input:  "pre_value",
			want:   "pre_value",
		},
		{
			name:   "input needs prefix",
			prefix: "pre_",
			input:  "value",
			want:   "pre_value",
		},
		{
			name:   "input with leading/trailing spaces",
			prefix: "pre_",
			input:  "  value  ",
			want:   "pre_  value  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnsurePrefix(tt.prefix, tt.input)
			if got != tt.want {
				t.Errorf("EnsurePrefix(%q, %q) = %q, want %q", tt.prefix, tt.input, got, tt.want)
			}
		})
	}
}
