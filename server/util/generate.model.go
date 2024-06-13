package util

import (
	"math/rand"

	"github.com/google/uuid"
)

type Generate struct{}

func (ref Generate) OTP(length int) string {
	const otpChars = "1234567890"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}
	return string(otp)
}

func (ref Generate) UUIDv4() string {
	return uuid.New().String()
}
