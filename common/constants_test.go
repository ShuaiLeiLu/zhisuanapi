package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThemeAwarePathMapsLegacyConsoleRoutes(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "topup",
			path: "/console/topup",
			want: "/wallet",
		},
		{
			name: "log",
			path: "/console/log?type=all",
			want: "/usage-logs?type=all",
		},
		{
			name: "personal",
			path: "/console/personal",
			want: "/profile",
		},
		{
			name: "unknown",
			path: "/console/unknown",
			want: "/console/unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ThemeAwarePath(tt.path))
		})
	}
}
