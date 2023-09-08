package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/grafviktor/keep-my-secret/internal/api/utils"
)

func hashString(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	return string(bytes), err
}

type User struct {
	ID             int64  `json:"id"`
	Login          string `json:"login,omitempty"`
	HashedPassword string `json:"-"`
	// RestorePassword was not implemented and not used anywhere
	RestorePassword string `json:"-"`
	DataKey         string `json:"-"`
}

// NewUser creates a new New User model with a random data key. The key should never be given to a user.
// They will be automatically restored from the database when user logs in.
func NewUser(login, password string) (*User, error) {
	hashedPassword, err := hashString(password)
	if err != nil {
		return nil, err
	}

	// When we create a new user, we generate a new random password
	// this password is used to encrypt user's data internally. For security
	// reasons, the user never knows his own 'data' password.
	// 'data' password is stored in the database and encrypted by the user's
	// original password. When user logs in, we decrypt 'data' password and save
	// it into RAM. This process is transparent for the users.
	key := utils.GenerateRandomPassword()
	encryptedKey, err := utils.Encrypt([]byte(key), password)
	if err != nil {
		return nil, err
	}

	u := User{
		Login:           login,
		HashedPassword:  hashedPassword,
		RestorePassword: hashedPassword,
		DataKey:         string(encryptedKey),
	}

	return &u, nil
}

// PasswordMatches check if password which was provided by the user during login process is correct
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) GetDataKey(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("cannot decrypt data key - no password set")
	}

	key, err := utils.Decrypt([]byte(u.DataKey), password)
	if err != nil {
		return "", err
	}

	return string(key), nil
}
