package helper

import (
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"testing"
)


func TestIsStrongPassword(t *testing.T) {
	cfg.LoadEnvVariables()
	logger.NewZlog()
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "pw: abc", input: "abc", want: false},
		{name: "pw: abcFIEJ", input: "abcFIEJ", want: false},
		{name: "pw: isej@Ifie", input: "isej@Ifie", want: true},
		{name: "pw: ise3j@Ifie", input: "ise3j@Ifie", want: true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := isStrongPassword(testCase.input)

			if got != testCase.want {
				t.Errorf("got %+v want %+v", got, testCase.want)
			}
		})
	}
}
