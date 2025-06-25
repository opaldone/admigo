package common

import (
	"crypto/rand"
	"fmt"

	"admigo/applog"

	"golang.org/x/crypto/bcrypt"
)

const RLEN = 7

func Hapas(pas string) (hapa string) {
	bh, err := bcrypt.GenerateFromPassword([]byte(pas), bcrypt.DefaultCost)
	if err != nil {
		applog.Danger("crypto.go Hapas", err)
		return
	}

	hapa = string(bh)
	return
}

func CheckPas(hapas, pas string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hapas), []byte(pas))
	return err == nil
}

func CreateUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		applog.Danger("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func RandUID() string {
	b := make([]byte, RLEN)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
