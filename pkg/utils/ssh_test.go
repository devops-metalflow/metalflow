package utils

import (
	"testing"
)

func TestIsSafetyCmd(t *testing.T) {
	tests := []struct {
		name, cmd string
		wantErr   bool
	}{
		{
			name:    "testErr",
			cmd:     "rm /",
			wantErr: true,
		},
		{
			name:    "testSuccess",
			cmd:     "cd ..",
			wantErr: false,
		},
		{
			name:    "testNormal",
			cmd:     "rm /home/user/sss",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSafetyCmd(tt.cmd); (got != nil) != tt.wantErr {
				t.Errorf(`[%s] wantErr %v, but Contains(%q) = %v`, tt.name, tt.wantErr, tt.cmd, got)
			}
		})
	}
}
