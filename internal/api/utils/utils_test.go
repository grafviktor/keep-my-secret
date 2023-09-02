package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPasswordConformsPolicy(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HashedPassword is OK",
			args: args{password: "1"},
			want: true,
		},
		{
			name: "HashedPassword is empty",
			args: args{password: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordConformsPolicy(tt.args.password); got != tt.want {
				t.Errorf("IsPasswordConformsPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUsernameConformsPolicy(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Login is OK",
			args: args{username: "tony.tester@example.com"},
			want: true,
		},
		{
			name: "Login is empty",
			args: args{username: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUsernameConformsPolicy(tt.args.username); got != tt.want {
				t.Errorf("IsUsernameConformsPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		cipherdata []byte
		key        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Can encrypt and decrypt and confirm the same result",
			args: args{
				cipherdata: []byte("6tXPNaEV&!xC?3>#"),
				key:        "12345",
			},
			want:    "6tXPNaEV&!xC?3>#",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.args.cipherdata, tt.args.key)
			require.NoError(t, err)

			got, err := Decrypt(encrypted, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
