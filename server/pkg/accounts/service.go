package accounts

import "github.com/google/uuid"

type Account struct {
	ID       string
	Username string
	Salt     []byte
	Password string
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

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	acc := Account{
		ID:       id.String(),
		Username: username,
		Salt:     salt,
		Password: hash,
	}

	return &acc, nil
}

func VerifyLogin(pass string, account Account) (bool, error) {
	return verifyPassword(pass, account.Password)
}
