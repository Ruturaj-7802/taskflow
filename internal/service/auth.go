package service

import (
    "context"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "golang.org/x/crypto/bcrypt"

    "github.com/Ruturaj-7802/taskflow/internal/dto"
    "github.com/Ruturaj-7802/taskflow/internal/model"
    "github.com/Ruturaj-7802/taskflow/internal/repository"
)

var (
    ErrEmailTaken      = errors.New("email already registered")
    ErrInvalidCreds    = errors.New("invalid email or password")
)

type AuthService struct {
    userRepo  *repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
    // Check duplicate email
    existing, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, pgx.ErrNoRows) {
        return nil, err
    }
    if existing != nil {
        return nil, ErrEmailTaken
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        return nil, err
    }

    user := &model.User{
        ID:        uuid.New(),
        Name:      req.Name,
        Email:     req.Email,
        Password:  string(hash),
        CreatedAt: time.Now(),
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        Token: token,
        User:  dto.UserDTO{ID: user.ID, Name: user.Name, Email: user.Email},
    }, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrInvalidCreds
        }
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, ErrInvalidCreds
    }

    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        Token: token,
        User:  dto.UserDTO{ID: user.ID, Name: user.Name, Email: user.Email},
    }, nil
}

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    jwt.RegisteredClaims
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtSecret), nil
    })
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    return claims, nil
}