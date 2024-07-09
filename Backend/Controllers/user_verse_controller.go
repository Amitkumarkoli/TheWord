package Controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "backend/Models"
)

func CreateVerse(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        var verse Models.UserVerse
        if err := c.ShouldBindJSON(&verse); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        verse.UserID = userID
        db.Create(&verse)
        c.JSON(http.StatusOK, verse)
    }
}

func GetVerse(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var verse Models.UserVerse
        if err := db.First(&verse, "user_verse_id = ?", c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Verse not found"})
            return
        }
        c.JSON(http.StatusOK, verse)
    }
}

func DeleteVerse(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        verseID := c.Param("id")

        var verse Models.UserVerse
        if err := db.First(&verse, "user_verse_id = ? AND user_id = ?", verseID, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Verse not found or you don't have permission to delete this verse"})
            return
        }

        if err := db.Delete(&verse).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Verse deleted successfully"})
    }
}

func ToggleLike(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        userVerseIDStr := c.Param("id")
        userVerseID, _ := strconv.Atoi(userVerseIDStr)

        var like Models.Like
        if err := db.Where("user_id = ? AND user_verse_id = ?", userID, userVerseID).First(&like).Error; err == nil {
            db.Delete(&like)
            c.JSON(http.StatusOK, gin.H{"message": "Like removed"})
            return
        }

        newLike := Models.Like{
            UserID:      userID,
            UserVerseID: uint(userVerseID),
        }
        db.Create(&newLike)
        c.JSON(http.StatusOK, gin.H{"message": "Verse liked"})
    }
}

func AddComment(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uint)
        userVerseIDStr := c.Param("id")
        userVerseID, _ := strconv.Atoi(userVerseIDStr)

        var comment Models.Comment
        if err := c.ShouldBindJSON(&comment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        comment.UserID = userID
        comment.UserVerseID = uint(userVerseID)
        db.Create(&comment)
        c.JSON(http.StatusOK, comment)
    }
}

func GetComments(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userVerseIDStr := c.Param("id")
        userVerseID, _ := strconv.Atoi(userVerseIDStr)

        var comments []Models.Comment
        if err := db.Where("user_verse_id = ?", userVerseID).Find(&comments).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, comments)
    }
}

func GetLikesCount(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userVerseIDStr := c.Param("id")
        userVerseID, _ := strconv.Atoi(userVerseIDStr)

        var count int64
        if err := db.Model(&Models.Like{}).Where("user_verse_id = ?", userVerseID).Count(&count).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"likes_count": count})
    }
}
