package utils

import (
	"math/rand"
	"time"
)

type OtpType int

const (
	// Standard 6-digit OTPs
	TWO_FA_LOGIN_OTP OtpType = 6
	CHANGE_PASS_OTP  OtpType = 6
	ADD_EMAIL_OTP    OtpType = 6

	// Longer 8-digit OTPs
	FORGOT_PASS_OTP  OtpType = 8
	REMOVE_EMAIL_OTP OtpType = 8
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateOTP(length int) string {
	if length <= 0 {
		return ""
	}

	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range length {
		randomIndex := seededRand.Intn(len(digits))
		otp[i] = digits[randomIndex]
	}

	return string(otp)
}
