package utils

import "golang.org/x/crypto/bcrypt"

// ComparePwd used to check if the password is correct.
// str is the plaintext, pwd is the encrypted ciphertext.
func ComparePwd(str, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(str))
	return err == nil
}

// GenPwd is used to turn str into an encrypted string.
func GenPwd(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hash)
}
