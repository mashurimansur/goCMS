package token

import (
	"errors"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// PasetoMaker is a PASETO token maker using V4
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKeyHex string) (Maker, error) {
	if len(symmetricKeyHex) != 32 {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 characters")
	}

	symmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to create symmetric key: %w", err)
	}

	maker := &PasetoMaker{
		symmetricKey: symmetricKey,
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token := paseto.NewToken()
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	token.SetString("id", payload.ID.String())
	token.SetString("username", payload.Username)

	encrypted := token.V4Encrypt(maker.symmetricKey, nil)
	return encrypted, payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Extract payload from claims
	id, err := parsedToken.GetString("id")
	if err != nil {
		return nil, errors.New("missing id in token")
	}

	username, err := parsedToken.GetString("username")
	if err != nil {
		return nil, errors.New("missing username in token")
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, errors.New("missing issued_at in token")
	}

	expiration, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, errors.New("missing expiration in token")
	}

	payload := &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiration,
	}

	return payload, nil
}
