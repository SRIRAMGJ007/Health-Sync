package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	Role     string `json:"role" binding:"required"`
}

type DoctorRegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	Name            string `json:"name" binding:"required"`
	Specialization  string `json:"specialization" binding:"required"`
	Experience      int32  `json:"experience" binding:"required"`
	Qualification   string `json:"qualification" binding:"required"`
	HospitalName    string `json:"hospital_name" binding:"required"`
	ConsultationFee string `json:"consultation_fee" binding:"required"`
	Role            string `json:"role" binding:"required"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type DoctorLoginRequest struct {
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
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "user" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPasswordStr := string(hashedPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		user, err := queries.CreateUserWithEmail(ctx, repository.CreateUserWithEmailParams{
			Email:        req.Email,
			PasswordHash: &hashedPasswordStr,
			Name:         &req.Name,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server err"})
			return
		}

		token, err := utils.GenerateJWT(user.ID.String(), user.Email, "user")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
			"token":   token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"name":  user.Name,
			},
		})
	} else if req.Role == "doctor" {
		var req DoctorRegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPasswordStr := string(hashedPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var consultationFee pgtype.Numeric
		err = consultationFee.Scan(req.ConsultationFee) // Scan the string into pgtype.Numeric
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consultation fee format"})
			return
		}
		doctor, err := queries.CreateDoctor(ctx, repository.CreateDoctorParams{
			Name:            req.Name,
			Specialization:  req.Specialization,
			Experience:      req.Experience,
			Qualification:   req.Qualification,
			HospitalName:    req.HospitalName,
			ConsultationFee: consultationFee,
			Email:           &req.Email,
			PasswordHash:    &hashedPasswordStr,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server err"})
			return
		}

		token, err := utils.GenerateJWT(doctor.ID.String(), *doctor.Email, "doctor")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Doctor login successful",
			"token":   token,
			"doctor": gin.H{
				"id":               doctor.ID,
				"email":            doctor.Email,
				"name":             doctor.Name,
				"specialization":   doctor.Specialization,
				"experience":       doctor.Experience,
				"qualification":    doctor.Qualification,
				"hospital_name":    doctor.HospitalName,
				"consultation_fee": doctor.ConsultationFee,
			},
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
	}
}

func UserRegisterHandler(ctx *gin.Context, queries *repository.Queries) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := queries.CreateUserWithEmail(ctx, repository.CreateUserWithEmailParams{
		Email:        req.Email,
		PasswordHash: &hashedPasswordStr,
		Name:         &req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server err"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func DoctorRegisterHandler(ctx *gin.Context, queries *repository.Queries) {
	var req DoctorRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var consultationFee pgtype.Numeric
	err = consultationFee.Scan(req.ConsultationFee)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consultation fee format"})
		return
	}

	doctor, err := queries.CreateDoctor(ctx, repository.CreateDoctorParams{
		Name:            req.Name,
		Specialization:  req.Specialization,
		Experience:      req.Experience,
		Qualification:   req.Qualification,
		HospitalName:    req.HospitalName,
		ConsultationFee: consultationFee,
		Email:           &req.Email,
		PasswordHash:    &hashedPasswordStr,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server err"})
		return
	}

	token, err := utils.GenerateJWT(doctor.ID.String(), *doctor.Email, "doctor")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Doctor created successfully",
		"token":   token,
		"doctor": gin.H{
			"id":               doctor.ID,
			"email":            doctor.Email,
			"name":             doctor.Name,
			"specialization":   doctor.Specialization,
			"experience":       doctor.Experience,
			"qualification":    doctor.Qualification,
			"hospital_name":    doctor.HospitalName,
			"consultation_fee": doctor.ConsultationFee,
		},
	})
}

func UserLoginHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Printf("user login request received")
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Println("error binding json: ", err)
		return
	}

	user, err := queries.GetUserByEmail(dbCtx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func DoctorLoginHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Printf("doctor login request received")
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req DoctorLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Println("error binding json: ", err)
		return
	}

	doctor, err := queries.GetDoctorByEmail(dbCtx, &req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*doctor.PasswordHash), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(doctor.ID.String(), *doctor.Email, "doctor")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Doctor login successful",
		"token":   token,
		"doctor": gin.H{
			"id":               doctor.ID,
			"email":            doctor.Email,
			"name":             doctor.Name,
			"specialization":   doctor.Specialization,
			"experience":       doctor.Experience,
			"qualification":    doctor.Qualification,
			"hospital_name":    doctor.HospitalName,
			"consultation_fee": doctor.ConsultationFee,
		},
	})
}

func GoogleAuthhandler(ctx *gin.Context, queries *repository.Queries) {

	log.Printf("Googleauth request recived")
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "google authorizationcode missing"})
		log.Println("google authorizationcode missing")
		return
	}

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
	user, err := queries.GetUserByGoogleID(ctx, &googleuser.ID)
	// the below condition is for creating new google user
	log.Println(" checking if user already exists ")
	if err != nil {
		if err == pgx.ErrNoRows {
			newuser, err := queries.CreateUserWithGoogle(dbCtx, repository.CreateUserWithGoogleParams{
				GoogleID: &googleuser.ID,
				Email:    googleuser.Email,
				Name:     &googleuser.Name,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				log.Printf("failed to create user :%v", err)
				return
			}
			// user = newuser
			token, err := utils.GenerateJWT(newuser.ID.String(), newuser.Email, "user")
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				log.Printf("failed to generate token for the  user :%v", err)
				return
			}

			// Send response
			ctx.JSON(http.StatusOK, gin.H{
				"message": "new user login successful",
				"token":   token,
				"user": gin.H{
					"id":    newuser.ID,
					"email": newuser.Email,
					"name":  newuser.Name,
				},
			})
		}
		return
	}
	// the below condition is for existing google user
	token, err := utils.GenerateJWT(user.ID.String(), user.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		log.Printf("failed to decode user info : %v ", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})

}
