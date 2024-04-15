package jwt

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	// "github.com/golang-jwt/jwt/v4"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/util"
	"gorm.io/gorm"
)

// https://dev.to/techschoolguru/how-to-create-and-verify-jwt-paseto-token-in-golang-1l5j

var userModel *models.User
var orgModel *models.Organization
var userMutex sync.Mutex

func InitModels() {
	userModel = &models.User{}
	orgModel = &models.Organization{}
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type ClaimStrings []string

type NumericDate struct {
	time.Time
}

type RegisteredClaims struct {
	// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Issuer string `json:"iss,omitempty"`

	// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Subject string `json:"sub,omitempty"`

	// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	Audience ClaimStrings `json:"aud,omitempty"`

	// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	ExpiredAt *NumericDate `json:"exp,omitempty"`

	// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	NotBefore *NumericDate `json:"nbf,omitempty"`

	// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	IssuedAt *NumericDate `json:"iat,omitempty"`

	// the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
	ID string `json:"jti,omitempty"`
}

type Payload struct {
	RegisteredClaims
	UserID         int64 `json:"user_id,omitempty"`
	OrganizationId int64 `json:"organization_id,omitempty"`
}

func (c *Payload) Valid() error {
	if time.Now().After(c.ExpiredAt.Time) {
		return fmt.Errorf("payload is invalid")
	}
	return nil
}

type TMSRegisteredClaims struct {
	// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Issuer string `json:"iss,omitempty"`

	// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Subject string `json:"sub,omitempty"`

	// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	Audience string `json:"aud,omitempty"`

	// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	ExpiredAt int64 `json:"exp,omitempty"`

	// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	NotBefore *NumericDate `json:"nbf,omitempty"`

	// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	IssuedAt *NumericDate `json:"iat,omitempty"`

	// the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
	ID string `json:"jti,omitempty"`
}

type PayloadTMS struct {
	TMSRegisteredClaims
	Username string `json:"TglCode,omitempty"`
	Email    string `json:"Email,omitempty"`
	FullName string `json:"FullName,omitempty"`
}

func (c *PayloadTMS) Valid() error {
	if time.Now().UnixMilli() < c.ExpiredAt {
		return fmt.Errorf("payload is invalid")
	}
	return nil
}

type JWTMaker struct {
	SecretKey    string
	TMSSecretKey string
	Lifetime     int
	Issuer       string
}

func (maker *JWTMaker) CreateToken(userID int64) (string, error) {

	payload := &Payload{
		UserID: userID,
		RegisteredClaims: RegisteredClaims{
			Issuer: maker.Issuer,
		},
	}

	payload.ExpiredAt = &NumericDate{
		Time: time.Now().Add(time.Duration(maker.Lifetime) * time.Hour),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.SecretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return maker.VerifyTokenWithTMS(token)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (maker *JWTMaker) VerifyTokenWithTMS(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.TMSSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &PayloadTMS{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payloadTMS, ok := jwtToken.Claims.(*PayloadTMS)
	if !ok {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		RegisteredClaims: RegisteredClaims{
			Issuer:  payloadTMS.Issuer,
			Subject: payloadTMS.Subject,
			// Audience: payloadTMS.Audience,
			ExpiredAt: &NumericDate{
				Time: time.Unix(0, payloadTMS.ExpiredAt*int64(time.Millisecond)),
			},
			NotBefore: payloadTMS.NotBefore,
			IssuedAt:  payloadTMS.IssuedAt,
			ID:        payloadTMS.ID,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*6))
	defer cancel()

	userMutex.Lock()
	defer userMutex.Unlock()
	user, err := userModel.ReadByCondition(ctx, "username", payloadTMS.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// add new user
		data := &entities.User{
			DisplayName:  payloadTMS.FullName,
			Username:     payloadTMS.Username,
			PasswordHash: "",
			Email:        payloadTMS.Email,
		}

		// TGL Solutions
		org, err := orgModel.First(ctx)
		if err == nil {
			data.OrganizationID = org.ID
		} else {
			data.Organization = &entities.Organization{
				DisplayName: "TGL Solutions",
			}
		}

		data.PasswordHash = util.GetMD5Hash(data.PasswordHash)
		data.UUID = uuid.NewString()

		var userErr error
		payload.UserID, userErr = userModel.Create(ctx, data)
		if userErr != nil {
			return nil, userErr
		}
	} else {
		payload.UserID = user.ID
	}

	return payload, nil
}
