package utils

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

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

	// Test wrong AES key kength
	aesKeyLength = 11
	_, err := Decrypt([]byte("6tXPNaEV&!xC?3>#"), "12345")
	require.Error(t, err)
}

func TestGenerateRandomPassword(t *testing.T) {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random password
	password := GenerateRandomPassword()

	// Check if the generated password has the correct length
	expectedLength := aesKeyLength // Assuming aesKeyLength is defined in the same package
	if len(password) != expectedLength {
		t.Errorf("Generated password has incorrect length. Expected %d, but got %d", expectedLength, len(password))
	}

	// Check if the generated password only contains valid validRandomChars
	for _, char := range password {
		if !strings.ContainsRune(validRandomChars, char) {
			t.Errorf("Generated password contains invalid character: %s", string(char))
		}
	}
}

func BenchmarkGenerateRandomPassword(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		_ = GenerateRandomPassword()
	}
}

func TestNormalizeAESKey(t *testing.T) {
	// aesKeyLength is defined in utils.go. Redefining, to simplify debugging
	aesKeyLength = 16
	// Test cases for key normalization
	testCases := []struct {
		input    string
		expected string
	}{
		{"1234", "1234000000000000"},                     // Shorter key, should be padded with zeros
		{"abcdefghijklmnop", "abcdefghijklmnop"},         // Correct length key, should remain unchanged
		{"abcdefghijklmnopqrstuvwx", "abcdefghijklmnop"}, // Longer key, should be truncated
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			normalizedKey := normalizeAESKey(tc.input)

			if normalizedKey != tc.expected {
				t.Errorf("Expected normalized key: %s, but got: %s", tc.expected, normalizedKey)
			}
		})
	}
}
