package email

import (
	"errors"
	"testing"
)

type emailValidationTests struct {
	name  string
	input EmailInfo
	want  error
}

func TestCheckEmailInfo(t *testing.T) {
	tests := []emailValidationTests{
		{
			name: "all fields pass the validation",
			input: EmailInfo{
				To:         []string{"user1@example.com", "user2@example.ai"},
				MsgMeta:    map[string]interface{}{"subject": "This is a test mail"},
				MsgContent: "This is a test mail from Go-Fiber API",
			},
			want: nil,
		},
		{
			name: "check 0 recipients",
			input: EmailInfo{
				To: []string{},
			},
			want: errors.New("Recipient is required"),
		},
		{
			name: "check empty subject",
			input: EmailInfo{
				To:      []string{"user1@example.com", "user2@example.ai"},
				MsgMeta:    map[string]interface{}{"subject": ""},
			},
			want: errors.New("Subject is required"),
		},
		{
			name: "check empty message",
			input: EmailInfo{
				To:      []string{"user1@example.com", "user2@example.ai"},
				MsgMeta:    map[string]interface{}{"subject": "This is a test mail"},
				MsgContent: "",
			},
			want: errors.New("Cannot send empty message"),
		},
		{
			name: "check email is valid",
			input: EmailInfo{
				To:      []string{"user1@example.com", "user2example.ai"},
				MsgMeta:    map[string]interface{}{"subject": "This is a test mail"},
				MsgContent: "This is a test mail from Go-Fiber API",
			},
			want: errors.New("user2example.ai is not a valid email..."),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.input.checkInfo()

			if got != nil && testCase.want == nil {
				t.Errorf("got an error but didn't want one")
				return
			}

			if got != nil {
				assertError(t, got, testCase.want)
			}

		})
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want.Error() {
		t.Errorf("got %q, want %q", got, want)
	}
}
