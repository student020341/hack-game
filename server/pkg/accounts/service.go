package accounts

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID       string
	Username string `gorm:"unique;not null"`
	Salt     []byte
	Password string `gorm:"not null"`
	// 0 = admin, 1 = mod, 2 = user?
	Level int
}

func CreateAccount(username string, pass string) (*Account, error) {
	p := &params{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 1,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := generateSalt(p.saltLength)
	if err != nil {
		return nil, err
	}

	hash, err := generateHashFromPassword(pass, salt, p)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	acc := Account{
		ID:       id.String(),
		Username: username,
		Salt:     salt,
		Password: hash,
		Level:    2,
	}

	return &acc, nil
}

func Login(pass string, account Account) (*AuthSession, error) {
	ok, err := verifyPassword(pass, account.Password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	id := uuid.New()

	token := uuid.New()

	auth := AuthSession{
		ID:        id.String(),
		Token:     token.String(),
		CreatedAt: time.Now(),
		AccountID: account.ID,
	}

	return &auth, nil
}
