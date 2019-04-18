package battleship

import "testing"

func TestPrompt(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Prompt(); got != tt.want {
				t.Errorf("Prompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
