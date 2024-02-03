package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

//var db *sql.DB

type Sample struct {
	Method string `json:"method" xml:"method" form:"method"`
}

type UserAccounts struct {
	Id       int    `json:"Id"`
	IdNumber string `json:"IdNumber"`
	FullName string `json:"FullName"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Section  string `json:"Section"`
	Role     string `json:"Role"`
}

func main() {
	app := fiber.New()

	db, err := connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		htmlElement := "<p>"
		htmlElement += "API TEMPLATE GO (localhost:914)"
		htmlElement += "</p>"
		htmlElement += "<p>"
		htmlElement += "Developed By : Vince Dale D. Alcantara"
		htmlElement += "</p>"
		htmlElement += "<p>"
		htmlElement += "Version 1.0.0"
		htmlElement += "</p>"

		return c.Type("html", "utf-8").SendString(htmlElement)
	})

	app.Get("/UserAccounts", func(c *fiber.Ctx) error {
		user_accounts, err := getUserAccounts(db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(user_accounts)
	})

	app.Get("/UserAccounts/HTML", func(c *fiber.Ctx) error {
		htmlElement := ""

		user_accounts, err := getUserAccounts(db)
		if err != nil {
			log.Print(err)
		}

		count := 0

		for _, ua := range user_accounts {
			count++
			htmlElement += "<tr style=\"cursor:pointer;\" class=\"modal-trigger\" data-toggle=\"modal\" data-target=\"#update_account\" onclick=\"get_accounts_details(&quot;" + strconv.Itoa(ua.Id) + "~!~" + ua.IdNumber + "~!~" + ua.Username + "~!~" + ua.FullName + "~!~" + ua.Password + "~!~" + ua.Section + "~!~" + ua.Role + "&quot;)\">"
			htmlElement += "<td>" + strconv.Itoa(count) + "</td>"
			htmlElement += "<td>" + ua.IdNumber + "</td>"
			htmlElement += "<td>" + ua.Username + "</td>"
			htmlElement += "<td>" + ua.FullName + "</td>"
			htmlElement += "<td>" + ua.Section + "</td>"
			htmlElement += "<td>" + strings.ToUpper(ua.Role) + "</td>"
			htmlElement += "</tr>"
		}

		return c.Type("html", "utf-8").SendString(htmlElement)
	})

	app.Get("/UserAccounts/Count", func(c *fiber.Ctx) error {
		count, err := countUserAccounts(db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(count)
	})

	app.Get("/UserAccounts/Id/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		user_account, err := getUserAccountsById(id, db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(user_account)
	})

	// Randoms

	app.Post("/sample", func(c *fiber.Ctx) error {
		sample := new(Sample)

		if err := c.BodyParser(sample); err != nil {
			return err
		}

		htmlElement := "<tr>"
		htmlElement += "<td colspan='6' style='text-align:center; color:red;'>TESTING API TEMPLATE GO ON localhost:914 method = " + sample.Method + "</td>"
		htmlElement += "</tr>"

		return c.Type("html", "utf-8").SendString(htmlElement)
	})

	app.Get("/sample2/:method?", func(c *fiber.Ctx) error {
		method := c.Params("method")

		htmlElement := "<tr>"
		htmlElement += "<td colspan='6' style='text-align:center; color:red;'>TESTING API TEMPLATE GO ON localhost:914 method = " + method + "</td>"
		htmlElement += "</tr>"

		return c.Type("html", "utf-8").SendString(htmlElement)
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		str := "TESTING API TEMPLATE GO ON localhost:914"

		return c.SendString(str)
	})

	log.Fatal(app.Listen(":914"))
}

// Database Connection
func connectToDatabase() (*sql.DB, error) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "web_template",
		AllowNativePasswords: true,
	}

	// Open the database connection.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is established.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getUserAccounts(db *sql.DB) ([]UserAccounts, error) {
	// A user_accounts slice to hold data from returned rows.
	var user_accounts []UserAccounts

	rows, err := db.Query("SELECT * FROM user_accounts")
	if err != nil {
		return nil, fmt.Errorf("getUserAccounts : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var ua UserAccounts
		if err := rows.Scan(&ua.Id, &ua.IdNumber, &ua.FullName, &ua.Username, &ua.Password, &ua.Section, &ua.Role); err != nil {
			return nil, fmt.Errorf("getUserAccounts : %v", err)
		}
		user_accounts = append(user_accounts, ua)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getUserAccounts : %v", err)
	}
	return user_accounts, nil
}

func getUserAccountsById(id int64, db *sql.DB) (UserAccounts, error) {
	// A ua to hold data from the returned row.
	var ua UserAccounts

	row := db.QueryRow("SELECT * FROM user_accounts WHERE id = ?", id)
	if err := row.Scan(&ua.Id, &ua.IdNumber, &ua.FullName, &ua.Username, &ua.Password, &ua.Section, &ua.Role); err != nil {
		if err == sql.ErrNoRows {
			return ua, fmt.Errorf("getUserAccountsById %d: no such UserAccount", id)
		}
		return ua, fmt.Errorf("getUserAccountsById %d: %v", id, err)
	}
	return ua, nil
}

func countUserAccounts(db *sql.DB) (int, error) {
	// A variable to hold the count.
	var count int

	// Query the database.
	row := db.QueryRow("SELECT COUNT(*) FROM user_accounts")
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countUserAccounts: %v", err)
	}

	return count, nil
}
