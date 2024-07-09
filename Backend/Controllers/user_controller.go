package Controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "backend/Models"
)

func GetUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        var user Models.User
        if err := db.First(&user, "user_id = ?", userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, user)
    }
}

func GetUserSettings(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        var user Models.User
        if err := db.First(&user, "user_id = ?", userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "primary_color":    user.PrimaryColor,
            "highlight_color":  user.HighlightColor,
            "dark_mode":        user.DarkMode,
            "public_profile":   user.PublicProfile,
            "translation_id":   user.TranslationId,
            "translation_name": user.TranslationName,
        })
    }
}

func UpdateUserSettings(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        var req struct {
            PrimaryColor    int    `json:"primary_color"`
            HighlightColor  int    `json:"highlight_color"`
            DarkMode        bool   `json:"dark_mode"`
            PublicProfile   bool   `json:"public_profile"`
            TranslationId   string `json:"translation_id"`
            TranslationName string `json:"translation_name"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        updates := map[string]interface{}{
            "primary_color":    req.PrimaryColor,
            "highlight_color":  req.HighlightColor,
            "dark_mode":        req.DarkMode,
            "public_profile":   req.PublicProfile,
            "translation_id":   req.TranslationId,
            "translation_name": req.TranslationName,
        }

        if err := db.Model(&Models.User{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
    }
}
