package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/idtoken"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const GOOGLE_CLIENT_ID = "351570360340-c2fh8e6t265el6d73a0n25m95s9uq49j.apps.googleusercontent.com"
const GOOGLE_ISSUER_1 = "accounts.google.com"
const GOOGLE_ISSUER_2 = "https://accounts.google.com"

type jwtCustomClaims struct {
	Email        string `json:"email"`
	GoogleIdUser string `json:"google_id_user"`
	UserId       int    `json:"user_id"`
	jwt.StandardClaims
}

func getToken(email, googleUserId string, userId int) (string, error) {
	claims := &jwtCustomClaims{
		email,
		googleUserId,
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}
func verifyIdToken(token string) (idtoken.Payload, error) {
	user, err := idtoken.Validate(context.Background(), token, GOOGLE_CLIENT_ID)
	if err != nil {
		return idtoken.Payload{}, err
	}
	if user.Audience != GOOGLE_CLIENT_ID {
		return idtoken.Payload{}, errors.New("google client id not match")
	}
	if user.Issuer != GOOGLE_ISSUER_1 && user.Issuer != GOOGLE_ISSUER_2 {
		return idtoken.Payload{}, errors.New("google issuer not match")
	}
	return *user, nil
}

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "root", "localhost", "googleauth")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	e := echo.New()
	e.Use(middleware.CORS())
	//e.Use(middleware.Recover())
	e.POST("/google", func(c echo.Context) error {
		var req struct {
			Token string `json:"token"`
		}
		err := c.Bind(&req)
		if err != nil {
			 return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		user, err := verifyIdToken(req.Token)
		if err != nil {
			 return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Println("google response", user)


		emailResp := user.Claims["email"].(string)

		var userdb User
		if err := db.First(&userdb, "email = ?", emailResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userdb = User{
					GoogleID: user.Subject,
					Email:    emailResp,
				}
				if err := db.Create(&userdb).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Println(userdb)

		token, err := getToken(emailResp, user.Subject, userdb.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, echo.Map{
			"user": echo.Map{
				"user_id":        userdb.ID,
				"google_id_user": user.Subject,
				"email":          emailResp,
			},
			"token": token,
		})
	})
	r := e.Group("/me")

	//Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		userId := claims.UserId
		fmt.Println(userId)

		var userdb User
		if err := db.First(&userdb, userId).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Println(userdb)
		return c.JSON(http.StatusOK, userdb)
	})
	e.Logger.Fatal(e.Start(":5000"))
}
