package term_test

import (
	"testing"

	"github.com/batmac/ccat/term"
)

func TestIsStdoutTerminal(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.IsStdoutTerminal(); got != tt.want {
				t.Errorf("IsStdoutTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTerminalSize(t *testing.T) {
	tests := []struct {
		name       string
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{"", 80, 24, false}, // default
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight, err := term.GetTerminalSize()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTerminalSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("GetTerminalSize() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("GetTerminalSize() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
		})
	}
}

func TestClearScreen(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		term.ClearScreen()
	})
}

func TestSupportedColors(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		{"", 8}, // default
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.SupportedColors(); got < tt.want {
				t.Errorf("SupportedColors() = %v, want >= %v", got, tt.want)
			}
		})
	}
}
