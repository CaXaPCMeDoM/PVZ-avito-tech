package security

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, inputPassword string) error
}
