package controllers

import (
	"net/http"
	"time"

	"identity-service/config"
	"identity-service/models"
	"identity-service/redisconn"
	"identity-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	c.BindJSON(&input)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	user := models.User{Username: input.Username, Password: string(hashed)}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var input struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		DeviceID  string `json:"device_id"`
		UserAgent string `json:"user_agent"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessJti := uuid.New().String()
	refreshJti := uuid.New().String()

	redisconn.SetChangePassAt(user.ID, user.UpdatedAt.Unix())

	accessToken, err := utils.GenerateAccessToken(user.ID, accessJti, user.UpdatedAt.Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, refreshJti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	refreshRecord := models.RefreshToken{
		ID:        refreshJti,
		UserID:    user.ID,
		Token:     refreshToken,
		DeviceID:  input.DeviceID,
		UserAgent: input.UserAgent,
		ExpiresAt: time.Now().Add(time.Duration(utils.REFRESH_TOKEN_TIME) * time.Hour),
	}
	if err := config.DB.Create(&refreshRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Parse & verify token
	claims, err := utils.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	tokenID := claims["jti"].(string)
	userID := uint(claims["user_id"].(float64))
	exp := int64(claims["exp"].(float64))

	// Check blacklist
	isBlacklisted, _ := redisconn.IsRefreshTokenBlacklisted(tokenID)
	if isBlacklisted {
		// Có thể nghi ngờ bị leak -> xoá toàn bộ token
		config.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
		redisconn.BlacklistAccessToken(tokenID, exp)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token reuse detected"})
		return
	}

	// Kiểm tra DB: refresh token có tồn tại không?
	var storedToken models.RefreshToken
	err = config.DB.Where("id = ? AND user_id = ?", tokenID, userID).First(&storedToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Re-use token đã bị xóa
			redisconn.BlacklistRefreshToken(tokenID, exp)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token invalid or reused"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Delete old refresh token (rotation)
	config.DB.Delete(&storedToken)

	// Tạo access + refresh mới
	newAccessJTI := uuid.New().String()
	newRefreshJTI := uuid.New().String()
	changePassAt, _ := redisconn.GetChangePassAt(userID)

	accessToken, err := utils.GenerateAccessToken(userID, newAccessJTI, changePassAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(userID, newRefreshJTI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Lưu refresh token mới vào DB
	newRecord := models.RefreshToken{
		ID:        newRefreshJTI,
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	config.DB.Create(&newRecord)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
