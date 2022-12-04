package hash

import (
	"crypto/md5"
	"fmt"
)

type MD5Hasher struct {
	salt string
}

func NewMD5Hasher(salt string) *MD5Hasher {
	return &MD5Hasher{salt: salt}
}

func (h MD5Hasher) Hash(password string) (string, error) {
	md5h := md5.New()

	if _, err := md5h.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5h.Sum([]byte(h.salt))), nil
}
