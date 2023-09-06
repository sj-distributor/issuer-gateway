package utils

import (
	"reflect"
	"testing"
)

var key = "66d2e42661bc292f8237b4736a423a36"
var text = "anson.test.com"

func TestDecrypt(t *testing.T) {

	encrypt, _ := Encrypt(text, key)
	type args struct {
		ciphertext string
		key        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can decrypt",
			args: struct {
				ciphertext string
				key        string
			}{ciphertext: encrypt, key: key},
			want:    text,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.ciphertext, tt.args.key)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}
			if got != string(tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can encrypt",
			args: struct {
				plaintext string
				key       string
			}{plaintext: text, key: key},
			want:    text,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decrypt, err := Decrypt(got, key)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
			}

			if !reflect.DeepEqual(decrypt, tt.want) {
				t.Errorf("Encrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
