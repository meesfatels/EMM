package runtime

import "testing"

func TestNormalizeConversationName_Valid(t *testing.T) {
	tests := []string{"my-chat", "chat_01", "chat.v2", "A-1_2.3"}
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			got, err := NormalizeConversationName(tc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc {
				t.Fatalf("got %q, want %q", got, tc)
			}
		})
	}
}

func TestNormalizeConversationName_Invalid(t *testing.T) {
	tests := []string{"", " ", ".", "..", "../x", "a/b", `a\\b`, "bad:name", "emoji🔥"}
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			if _, err := NormalizeConversationName(tc); err == nil {
				t.Fatalf("expected error for %q", tc)
			}
		})
	}
}
