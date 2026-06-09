package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrExpiredToken       = errors.New("token has expired")
    ErrEmailInUse         = errors.New("email already in use")
)

type AuthService struct {
	userRepo *UserRepository
	refreshTokenRepo *RefreshTokenRepository
	jwtSecret []byte
	accessTokenTTL time.Duration
}

func NewAuthService(userRepo *UserRepository, refreshTokenRepo *RefreshTokenRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService {
		userRepo: userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret: []byte(jwtSecret),
		accessTokenTTL: accessTokenTTL,
	}
}


func (s *AuthService) Register(email, name, username, password string) (*User, error) {
    // Check if user already exists
    _, err := s.userRepo.GetUserByEmail(email)
    if err == nil {
        return nil, ErrEmailInUse
    }
    // Only proceed if the error was "user not found"
    if !errors.Is(err, sql.ErrNoRows) {
        return nil, err
    }
    // Hash the password
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }
    // Create the user
    user, err := s.userRepo.CreateUser(email, name, username, hashedPassword)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

func (s *AuthService) generateAccessToken(user *User) (string, error) {
	expirationTime := time.Now().Add(s.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"username": user.Username,
		"email": user.Email,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthService) LoginWithRefresh(email, password string, refreshTokenTTL time.Duration) (accessToken string, refreshToken string, err error) {
    // Get the user from the database
    user, err := s.userRepo.GetUserByEmail(email)
    if err != nil {
        return "", "", ErrInvalidCredentials
    }
    // Verify the password
    if err := VerifyPassword(user.PasswordHash, password); err != nil {
        return "", "", ErrInvalidCredentials
    }
    // Generate an access token
    accessToken, err = s.generateAccessToken(user)
    if err != nil {
        return "", "", err
    }
    // Create a refresh token
    token, err := s.refreshTokenRepo.CreateRefreshToken(user.ID, refreshTokenTTL)
    if err != nil {
        return "", "", err
    }
    return accessToken, token.Token, nil
}

func (s *AuthService) RefreshAccessToken(refreshTokenString string) (string,error) {
	token, err := s.refreshTokenRepo.GetRefreshToken(refreshTokenString)
	if err != nil {
		return "", ErrInvalidToken
	}

	if token.Revoked {
		return "", ErrInvalidToken
	}

	if time.Now().After(token.ExpiredAt) {
		return "", ErrExpiredToken
	}

	user, err := s.userRepo.GetUserById(token.UserID)
	if err != nil {
		return "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}


func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}

