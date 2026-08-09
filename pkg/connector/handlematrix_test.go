package connector

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestShouldBridgeMatrixReadReceipts(t *testing.T) {
	tests := []struct {
		name   string
		client DiscordClient
		want   bool
	}{
		{
			name: "missing session",
			want: false,
		},
		{
			name: "bot session",
			client: DiscordClient{
				Session: &discordgo.Session{IsUser: false},
			},
			want: false,
		},
		{
			name: "user session",
			client: DiscordClient{
				Session: &discordgo.Session{IsUser: true},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.client.shouldBridgeMatrixReadReceipts(); got != tc.want {
				t.Fatalf("unexpected read receipt support: got %t, want %t", got, tc.want)
			}
		})
	}
}
