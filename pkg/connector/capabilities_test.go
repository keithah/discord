package connector

import (
	"testing"

	"maunium.net/go/mautrix/event"
)

func TestDiscordMediaCapabilitiesAllowThirtyMiB(t *testing.T) {
	const want = 30 * 1024 * 1024

	if MaxFileSize != want {
		t.Fatalf("MaxFileSize = %d, want %d", MaxFileSize, want)
	}

	for msgType, feature := range discordCaps.File {
		if feature.MaxSize != want {
			t.Fatalf("%s MaxSize = %d, want %d", msgType, feature.MaxSize, want)
		}
	}

	if _, ok := discordCaps.File[event.MsgImage]; !ok {
		t.Fatal("image media capability is missing")
	}
}
