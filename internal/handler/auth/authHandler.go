package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GoogleUser struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

//	var googleOAuthConfig = &oauth2.Config{
//		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
//		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
//		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
//		Scopes:       []string{"profile", "email"},
//		Endpoint:     google.Endpoint,
//	}
var googleOAuthConfig *oauth2.Config

func RegisterHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Printf("register request recived")
	var req RegisterRequest
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cheching if user already exists..
	_, err := queries.GetUserByEmail(dbCtx, req.Email)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		log.Printf("register request failed user already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing the password: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := queries.CreateUserWithEmail(dbCtx, repository.CreateUserWithEmailParams{
		Email:        req.Email,
		PasswordHash: sql.NullString{String: string(hashedPassword), Valid: true},
		Name:         sql.NullString{String: string(req.Name), Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server err"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created sucessfully",
		"token":   token,
		"user": gin.H{
			"Id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})

}

func LoginHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Printf("login request recived")

	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		log.Println("error binding json: ", err)
		return
	}

	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !user.PasswordHash.Valid {
		// Handle NULL case
		log.Println("password hash is NULL in the data base")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(string(user.ID.String()), user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"Id":    user.ID,
			"email": user.Email,
			"name":  string(user.Name.String),
		},
	})

}

func GoogleAuthhandler(ctx *gin.Context, queries *repository.Queries) {

	log.Printf("Googleauth request recived")

	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "google authorizationcode missing"})
		log.Println("google authorizationcode missing")
		return
	}

	log.Printf("Google OAuthConfig: %+v\n", googleOAuthConfig)
	Gtoken, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange google token"})
		log.Printf("failed to exchange google token :%v ", err)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), Gtoken)
	// https://www.googleapis.com/oauth2/v1/userinfo?alt=json
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
		log.Printf("failed to fetch user info :%v ", err)
		return
	}

	defer resp.Body.Close()

	// decode user info
	var googleuser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleuser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		log.Printf("failed to decode user info : %v ", err)
		return
	}

	// check if user exist
	user, err := queries.GetUserByGoogleID(ctx, sql.NullString{String: string(googleuser.ID), Valid: true})
	// the below condition is for creating new google user
	if err != nil {
		if err == sql.ErrNoRows {
			newuser, err := queries.CreateUserWithGoogle(ctx, repository.CreateUserWithGoogleParams{
				GoogleID: sql.NullString{String: string(googleuser.ID), Valid: true},
				Email:    googleuser.Email,
				Name:     sql.NullString{String: string(googleuser.Name), Valid: true},
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				log.Printf("failed to create user :%vv", err)
				return
			}
			// user = newuser
			token, err := utils.GenerateJWT(newuser.ID.String(), newuser.Email)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			// Send response
			ctx.JSON(http.StatusOK, gin.H{
				"message": "new user login successful",
				"token":   token,
				"user": gin.H{
					"Id":    newuser.ID,
					"email": newuser.Email,
					"name":  newuser.Name,
				},
			})
		}
		return
	}
	// the below condition is for creating new google user
	token, err := utils.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"Id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})

}
