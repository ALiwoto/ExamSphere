package hashing

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CompareSHA256(hash, value string) bool {
	return hash == HashSHA256(value)
}

func CompareSHA512(hash, value string) bool {
	return hash == HashSHA512(value)
}

func HashSHA256(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashSHA512(value string) string {
	hash := sha512.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetToken(size, id int64) string {
	return FormatToken(RandomString(size), id)
}

func GenerateAuthHash() string {
	b := make([]byte, AuthHashSize)
	for i := range b {
		b[i] = authHashChars[rand.Intn(len(authHashChars))]
	}
	return string(b)
}

func GenerateAgentAuthKey() string {
	b := make([]byte, AgentAuthKeySize)
	for i := range b {
		b[i] = authHashChars[rand.Intn(len(authHashChars))]
	}
	return string(b)
}

func FormatToken(hash string, id int64) string {
	return strconv.FormatInt(id, 10) + ":" + hash
}

func GetUserToken(id int64) string {
	return GetToken(8, id)
}

func GetIdFromToken(value string) int64 {
	if !strings.Contains(value, ":") {
		return 0
	}

	id, _ := strconv.ParseInt(strings.Split(value, ":")[0], 10, 64)
	return id
}

func GetIdAndHashFromToken(value string) (int64, string) {
	if !strings.Contains(value, ":") {
		return 0, ""
	}

	strs := strings.Split(value, ":")
	id, _ := strconv.ParseInt(strs[0], 10, 64)
	return id, strs[1]
}

// RandomString generates a random string of n length
func RandomString(n int64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

func RandomCommonString(n int64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = commonCharacterRunes[rand.Intn(len(commonCharacterRunes))]
	}
	return string(b)
}

func EncodeValueToBase64(value any) string {
	b, _ := json.Marshal(value)
	if len(b) == 0 {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}

func DecodeValueFromBase64(value string, result any) error {
	b, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, result)
}
