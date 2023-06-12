package user

import (
	"errors"
	"fmt"
	"go_api/auth"
	"go_api/database"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user database.Account
	query := "SELECT * FROM accounts WHERE id = $1"
	database.DB.QueryRowx(query, claims.Issuer).StructScan(&user)

	return c.JSON(user)
}

func GetBalance(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var balance int
	database.DB.QueryRow("SELECT balance FROM accounts WHERE id = $1", claims.Issuer).Scan(&balance)

	return c.JSON(fiber.Map{
		"balance": balance,
	})
}

func Click(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	// Increment balance
	_, err = database.DB.Exec("UPDATE accounts SET balance = balance + 1 WHERE id = $1", claims.Issuer)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update value: %v\n", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func ChangeAttr(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["attrName"] == "" {
		return errors.New("ERROR: attrName does not exist")
	} else if data["value"] == "" {
		return errors.New("ERROR: value does not exist")
	}

	_, err = database.DB.Exec("UPDATE accounts SET $1 = $2 WHERE id = $3", data["attrName"], data["value"], claims.Issuer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update value: %v\n", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
