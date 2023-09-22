package utils

import "testing"

func TestGenJwt(t *testing.T) {
	type args struct {
		id        string
		name      string
		jwtSecret string
		secret    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "can gen jwt token",
			args: struct {
				id        string
				name      string
				jwtSecret string
				secret    string
			}{id: "ansonansonanson", name: "ansonansonanson", jwtSecret: "F2NFhljyALXM86x", secret: "66d2e42661bc292f8237b4736a423a36"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenJwt(tt.args.id, tt.args.name, tt.args.jwtSecret, tt.args.secret)
			if err != nil {
				t.Error(err)
			}
			if err := ParseJwt(got, tt.args.jwtSecret, tt.args.secret, tt.args.id, tt.args.name); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParseJwt(t *testing.T) {
	type args struct {
		jwtSecret  string
		secret     string
		targetId   string
		targetName string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "can parse token",
			args: struct {
				jwtSecret  string
				secret     string
				targetId   string
				targetName string
			}{jwtSecret: "F2NFhljyALXM86x", secret: "66d2e42661bc292f8237b4736a423a36", targetId: "ansonansonanson", targetName: "ansonansonanson"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenJwt(tt.args.targetId, tt.args.targetName, tt.args.jwtSecret, tt.args.secret)
			if err != nil {
				t.Error(err)
			}

			if err := ParseJwt(got, tt.args.jwtSecret, tt.args.secret, tt.args.targetId, tt.args.targetName); err != nil {
				t.Errorf("ParseJwt() error = %v", err)
			}
		})
	}
}
