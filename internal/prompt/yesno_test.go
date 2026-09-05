package prompt

import (
	"os"
	"strings"
	"testing"
)

func TestYesNo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr string
	}{
		{name: "y", input: "y\n", want: true},
		{name: "yes", input: "YES\n", want: true},
		{name: "n", input: "n\n", want: false},
		{name: "no", input: "No\n", want: false},
		{name: "retry then yes", input: "maybe\ny\n", want: true},
		{name: "eof", input: "", wantErr: "error reading input"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			orig := os.Stdin
			os.Stdin = r
			t.Cleanup(func() {
				os.Stdin = orig
				_ = r.Close()
			})

			go func() {
				_, _ = w.Write([]byte(tc.input))
				_ = w.Close()
			}()

			got, err := YesNo("confirm?")
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("YesNo: %v", err)
			}
			if got != tc.want {
				t.Fatalf("YesNo = %v, want %v", got, tc.want)
			}
		})
	}
}
