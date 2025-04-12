package auth

import (
	"fmt"
	"testing"
)

func TestGenerateOtp(t *testing.T) {
	auth := Auth{}
	auth.otpSecondsValid = 10
	pin, err := auth.generatePinForOtp()
	if err != nil {
		t.Fatalf(":Failed to generate a random pin for otp: %v", err.Error())
	}
	fmt.Println("OTP: ", pin)
}
