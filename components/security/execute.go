package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"io"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, CredentialSaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password, saltHex string) (string, error) {
	if err := ValidatePasswordInput(password); err != nil {
		return "", err
	}
	salt, err := hex.DecodeString(saltHex)
	if err != nil || len(salt) != CredentialSaltLength {
		return "", ErrInvalidPasswordRecord
	}
	tag := pbkdf2HMACSHA256([]byte(password), salt, CredentialIterations, CredentialTagLength)
	return hex.EncodeToString(tag), nil
}

func VerifyPassword(record PasswordRecord, password string) (bool, error) {
	if err := ValidatePasswordRecord(record); err != nil {
		return false, err
	}
	hash, err := HashPassword(password, record.Salt)
	if err != nil {
		return false, err
	}
	expected, err := hex.DecodeString(record.PasswordHash)
	if err != nil {
		return false, ErrInvalidPasswordRecord
	}
	actual, err := hex.DecodeString(hash)
	if err != nil {
		return false, ErrInvalidPasswordRecord
	}
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

func NewPasswordRecord(password string) (PasswordRecord, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return PasswordRecord{}, err
	}
	hash, err := HashPassword(password, salt)
	if err != nil {
		return PasswordRecord{}, err
	}
	return PasswordRecord{
		PasswordHash: hash,
		Salt:         salt,
		Algorithm:    CredentialAlgorithmPBKDF2HMACSHA256V1,
		Iterations:   CredentialIterations,
	}, nil
}

func pbkdf2HMACSHA256(password, salt []byte, iterations, keyLength int) []byte {
	output := make([]byte, 0, keyLength)
	prf := newHMACSHA256PRF(password)
	var blockIndex uint32 = 1
	for len(output) < keyLength {
		block := pbkdf2Block(prf, salt, iterations, blockIndex)
		output = append(output, block...)
		blockIndex++
	}
	return output[:keyLength]
}

func pbkdf2Block(prf hmacSHA256PRF, salt []byte, iterations int, blockIndex uint32) []byte {
	input := make([]byte, len(salt)+4)
	copy(input, salt)
	binary.BigEndian.PutUint32(input[len(salt):], blockIndex)
	u := prf.sum(input)
	out := append([]byte{}, u[:]...)
	for i := 1; i < iterations; i++ {
		u = prf.sum(u[:])
		for j := range out {
			out[j] ^= u[j]
		}
	}
	return out
}

type hmacSHA256PRF struct {
	ipad [64]byte
	opad [64]byte
}

func newHMACSHA256PRF(key []byte) hmacSHA256PRF {
	if len(key) > 64 {
		sum := sha256.Sum256(key)
		key = sum[:]
	}
	var prf hmacSHA256PRF
	for i := range prf.ipad {
		prf.ipad[i] = 0x36
		prf.opad[i] = 0x5c
	}
	for i, b := range key {
		prf.ipad[i] ^= b
		prf.opad[i] ^= b
	}
	return prf
}

func (p hmacSHA256PRF) sum(data []byte) [sha256.Size]byte {
	inner := sha256.New()
	_, _ = inner.Write(p.ipad[:])
	_, _ = inner.Write(data)
	innerSum := inner.Sum(nil)

	outer := sha256.New()
	_, _ = outer.Write(p.opad[:])
	_, _ = outer.Write(innerSum)
	var out [sha256.Size]byte
	copy(out[:], outer.Sum(nil))
	return out
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
