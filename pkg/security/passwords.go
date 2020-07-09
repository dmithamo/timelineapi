// package security handles encryption/decryption of pwds and auth tokens
package security

import "golang.org/x/crypto/bcrypt"

// pwdHashingStrength sets the level of rigour employed in hashing the pwd
const pwdHashingStrength = 12

// GeneratePasswordHash encrypts a pwd for storage in the db
func GeneratePasswordHash(password *string) (string, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(*password), pwdHashingStrength)
	if err != nil {
		return "", nil
	}

	return string(pwdHash), nil
}

// VerifyPassword compares a pwd with the stored hash
func VerifyPassword(pwdHash, password *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*pwdHash), []byte(*password))
	return err == nil
}
