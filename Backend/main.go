package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "backend/Auth"
    "backend/Controllers"
    "backend/Models"
)

var db *gorm.DB

func main() {
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
        dbHost, dbUser, dbPassword, dbName)

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    db.AutoMigrate(&Models.User{}, &Models.UserVerse{}, &Models.Like{}, &Models.Comment{})
    log.Println("Database tables created or already exist.")

    createAdminUser()

    r := gin.Default()

    r.POST("/register", Auth.RegisterUser(db))
    r.POST("/login", Auth.LoginUser(db))

    auth := r.Group("/")
    auth.Use(Auth.AuthMiddleware(db))
    {
        auth.GET("/user/:id", Controllers.GetUser(db))
        auth.POST("/verse", Controllers.CreateVerse(db))
        auth.GET("/verse/:id", Controllers.GetVerse(db))
        auth.DELETE("/verses/:id", Controllers.DeleteVerse(db))
        auth.POST("/verse/:id/toggle-like", Controllers.ToggleLike(db))
        auth.POST("/verse/:id/comment", Controllers.AddComment(db))
        auth.GET("/user/settings", Controllers.GetUserSettings(db))
        auth.POST("/user/settings", Controllers.UpdateUserSettings(db))
        auth.GET("/verses/public", Controllers.GetPublicVerses(db))
        auth.POST("/verses/save", Controllers.SaveVerse(db))
        auth.GET("/verses/saved", Controllers.GetSavedVerses(db))
        auth.PUT("/verses/:id", Controllers.UpdateVerse(db))
        auth.GET("/verse/:id/comments", Controllers.GetComments(db))
        auth.GET("/verse/:id/likes", Controllers.GetLikesCount(db))
    }

    r.Run()
}

func createAdminUser() {
    var user Models.User
    if err := db.First(&user, "email = ?", "admin@example.com").Error; err == nil {
        return
    }

    passwordHash, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
    admin := Models.User{
        Email:           "admin@example.com",
        Username:        "AdminUser",
        PasswordHash:    string(passwordHash),
        PublicProfile:   true,
        PrimaryColor:    0xFF000000, // ARGB for black
        HighlightColor:  0xFFFF0000, // ARGB for red
        DarkMode:        true,
        TranslationId:   "bba9f40183526463-01",
        TranslationName: "Berean Standard Bible",
    }
    db.Create(&admin)
    log.Println("Admin user created or already exists.")
}
