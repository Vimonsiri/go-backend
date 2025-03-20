package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.com/ployMatsuri/go-backend/config"
	"gitlab.com/ployMatsuri/go-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Register function สำหรับสมัครสมาชิก
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// แฮชรหัสผ่าน
	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Password = hashedPassword

	// บันทึกลงฐานข้อมูล
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login function สำหรับการเข้าสู่ระบบ
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ตรวจสอบผู้ใช้ในฐานข้อมูล
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// ตรวจสอบรหัสผ่าน
	if err := CheckPasswordHash(user.Password, existingUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials - pws"})

		return
	}

	// สร้าง JWT Token
	token, err := GenerateJWT(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// ฟังก์ชันสำหรับการสร้าง JWT Token
func GenerateJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET") // ใช้ JWT_SECRET จาก .env
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not set in environment")
	}

	// สร้าง token และกำหนด claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // หมดอายุใน 24 ชั่วโมง

	// เซ็นต์และสร้าง token
	tk, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tk, nil
}

// ฟังก์ชันสำหรับการตรวจสอบรหัสผ่าน
func CheckPasswordHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
