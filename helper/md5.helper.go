package helper

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/rafzarf/system-otp-go/models"
	"net/http"
)

// HashMD5 hashes the input string using MD5
func HashMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CompareHashedPasswords(ctx *gin.Context, hashedInputPassword string, user models.User) bool {
	if hashedInputPassword != user.Password {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return false
	}
	return true
}
