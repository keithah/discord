package bridgev2

import (
	"testing"

	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

func TestShouldSkipLoginForRelayModeSkipsRelayLogin(t *testing.T) {
	relay := &UserLogin{UserLogin: &database.UserLogin{ID: networkid.UserLoginID("relay-login")}}
	candidate := &UserLogin{UserLogin: &database.UserLogin{ID: networkid.UserLoginID("relay-login")}}

	if !shouldSkipLoginForRelayMode(candidate, relay) {
		t.Fatal("relay login candidate was not skipped")
	}
}

func TestShouldSkipLoginForRelayModeKeepsDifferentLogin(t *testing.T) {
	relay := &UserLogin{UserLogin: &database.UserLogin{ID: networkid.UserLoginID("relay-login")}}
	candidate := &UserLogin{UserLogin: &database.UserLogin{ID: networkid.UserLoginID("personal-login")}}

	if shouldSkipLoginForRelayMode(candidate, relay) {
		t.Fatal("personal login candidate was skipped")
	}
}
