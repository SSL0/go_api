package auth

import (
	"fmt"
	"go_api/database"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"
const emailRegex = "/^[a-zA-Z0-9.!#$%&'*+\\=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$/"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	matched, _ := regexp.MatchString(emailRegex, data["email"])
	fmt.Println(matched)
	if len(data["name"]) < 5 || len(data["password"]) < 1 || matched {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Неверно введены данные",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := database.Account{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	query := `INSERT INTO accounts(name, email, password) VALUES (:name, :email, :password)`
	_, err := database.DB.NamedExec(query, user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Insert Row: %v\n", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Такие логин или пароль уже используются",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if len(data["name"]) < 5 || len(data["password"]) < 1 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Неверно введены данные",
		})
	}

	var user database.Account
	query := "SELECT * FROM accounts WHERE name = $1"
	err := database.DB.QueryRowx(query, data["name"]).StructScan(&user)
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Логин не найден",
		})
	}

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Пользователь не найден",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Неверный пароль",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Внутренняя ошибка сервера",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func GetUser(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}
		fmt.Println(err)
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user database.Account
	query := "SELECT * FROM accounts WHERE id = $1"
	database.DB.QueryRowx(query, claims.Issuer).StructScan(&user)

	return c.JSON(user)
}

func IsAuthorized(tokenString string) bool {
	_, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired {
			return false
		}
		fmt.Println(err)
	}

	return true
}
