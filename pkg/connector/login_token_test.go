package connector

import "testing"

func TestNormalizeDiscordLoginToken(t *testing.T) {
	userToken := "mfa.abcdefghijklmnopqrstuvwxyz"
	botToken := "not-a-real-token-for-normalizer-test"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "user token unchanged",
			input: userToken,
			want:  userToken,
		},
		{
			name:  "raw bot token unchanged",
			input: botToken,
			want:  botToken,
		},
		{
			name:  "prefixed bot token unchanged",
			input: "Bot " + botToken,
			want:  "Bot " + botToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeDiscordLoginToken(tt.input); got != tt.want {
				t.Fatalf("normalizeDiscordLoginToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
