package handler

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"shelter-it-be/internal/database"
	"shelter-it-be/internal/model/dto"
	"shelter-it-be/internal/model/request"
	"shelter-it-be/internal/utils"
	"time"
)

type AuthHandler struct {
	db              *database.Database
	jwtSecret       []byte
	tokenExpiration time.Duration
}

func NewAuthHandler(db *database.Database, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtSecret:       jwtSecret,
		tokenExpiration: 24 * time.Hour,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user request.RegisterReq
	req := utils.RegisterReq{RegisterReq: &user}

	// Validate input JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	// Additional validation
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var isUserExists bool
	err := h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&isUserExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if isUserExists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database with transaction
	tx, err := h.db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction start failed"})
		return
	}

	var id int
	err = tx.QueryRow(`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`, user.Email, hashedPassword).Scan(&id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"uid":     id,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var login request.LoginReq
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Get user from database
	var user dto.User
	err := h.db.DB.QueryRow(`SELECT EXISTS(SELECT id, email, password_hash FROM users WHERE email = $1)`, user.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Gemerate JWT with claims
	now := time.Now()
	claims := jwt.Claims(jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(h.tokenExpiration).Unix(),
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token with expiration
	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"exp_in":     h.tokenExpiration.Seconds(),
		"token_type": "Bearer",
	})
}

// RefreshToken generates token for valid user
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get UID from Context (set by auth middleware)
	uid, isUserExist := c.Get("uid")
	if !isUserExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate new token
	now := time.Now()
	claims := jwt.MapClaims{
		"uid": uid,
		"iat": now.Unix(),
		"exp": now.Add(h.tokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"exp_in":     h.tokenExpiration.Seconds(),
		"token_type": "Bearer",
	})
}

// Logout endpoint
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":      "Logout successful",
		"instruction:": "Please remove token from Authorization header",
	})
}
