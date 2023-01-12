package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	max = 999_999
	min = 100_000
)

type OtpConfig struct {
	Enabled           bool   `help:"listen string" env:"ENABLED" default:"false"`
	Issuer            string `help:"listen string" env:"ISSUER"`
	RecoveryCodeCount int    `help:"listen string" env:"RECOVERY_CODE_COUNT"`
}

type Otp struct {
	opt       *OtpConfig
	otpAlgo   otp.Algorithm
	otpDigits otp.Digits
}

func NewOtpConfig(opt *OtpConfig) *Otp {
	return &Otp{
		opt:       opt,
		otpAlgo:   otp.AlgorithmSHA512,
		otpDigits: otp.DigitsSix,
	}
}

func (o *Otp) Generate(_ context.Context, uid string) (*OtpGenerated, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      o.opt.Issuer,
		AccountName: uid,
		Algorithm:   o.otpAlgo,
		Digits:      o.otpDigits,
	})
	if err != nil {
		return nil, err
	}

	return &OtpGenerated{
		URL:    key.URL(),
		Secret: key.Secret(),
	}, nil
}

func (o *Otp) ValidateCode(_ context.Context, secret, code string) (bool, error) {
	nowInUTC := time.Now().UTC()
	isCorrect, err := totp.ValidateCustom(code, secret, nowInUTC, totp.ValidateOpts{
		Algorithm: o.otpAlgo,
		Digits:    o.otpDigits,
	})
	if err != nil {
		return false, err
	}

	return isCorrect, nil
}

func getRandom() (int64, error) {
	bigInt := big.NewInt(max - min)
	v, err := rand.Int(rand.Reader, bigInt)
	if err != nil {
		return 0, err
	}

	return v.Int64() + max, nil
}

func (o *Otp) GenerateRecoveryCodes(_ context.Context) ([]string, error) {
	codes := make([]string, 0)
	uniqCodes := make(map[string]struct{})
	cnt := o.opt.RecoveryCodeCount
	for cnt > 0 {
		randomNum, err := getRandom()
		if err != nil {
			return nil, err
		}
		randInt := fmt.Sprintf("%d", randomNum)
		if _, ok := uniqCodes[randInt]; ok {
			continue
		}
		uniqCodes[randInt] = struct{}{}
		cnt -= 1
	}

	for k := range uniqCodes {
		codes = append(codes, k)
	}
	return codes, nil
}
