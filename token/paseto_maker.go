package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	paseto       paseto.Token
	symmetricKey paseto.V4AsymmetricPublicKey
}

func CreateToken(Email string, role string, duration time.Duration) (*Payload, error) {
	payload, err := NewPayload(Email, role, duration)
	if err != nil {
		return payload, err
	}

	token := paseto.NewToken()

	token.SetString("email", payload.Email)
	token.SetString("role", payload.Role)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	// key := paseto.NewV4SymmetricKey()

	// return token, payload, err
}

// func VerifyToken(token string) (*Payload, error) {
// 	payload := &Payload{}
// 	err := paseto.Decrypt(token, symmetricKey, payload, nil)
// 	if err != nil {
// 		return nil, ErrInvalidToken
// 	}
// 	err = payload.Valid()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return payload, nil
// }
