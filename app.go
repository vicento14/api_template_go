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

	// db, err := connectToDatabase()
	// if err != nil {
	// 	log.Print(err)
	// }

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Middleware for database connection
	app.Use(func(c *fiber.Ctx) error {
		db, err := connectToDatabase()
		if err != nil {
			c.Status(500).SendString("Error connecting to the database")
			return err
		}
		c.Locals("db", db)
		return c.Next()
	})

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
		db := c.Locals("db").(*sql.DB)
		user_accounts, err := getUserAccounts(db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(user_accounts)
	})

	app.Get("/UserAccounts/HTML", func(c *fiber.Ctx) error {
		htmlElement := ""

		db := c.Locals("db").(*sql.DB)

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

	app.Get("/UserAccounts/Search/:id_number?/:full_name?/:role?", func(c *fiber.Ctx) error {
		id_number := c.Params("id_number")
		full_name := c.Params("full_name")
		role := c.Params("role")

		db := c.Locals("db").(*sql.DB)
		user_accounts, err := getUserAccountsSearch(id_number, full_name, role, db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(user_accounts)
	})

	app.Get("/UserAccounts/Count/:id_number?/:full_name?/:role?", func(c *fiber.Ctx) error {
		id_number := c.Params("id_number")
		full_name := c.Params("full_name")
		role := c.Params("role")

		db := c.Locals("db").(*sql.DB)
		count, err := countUserAccounts(id_number, full_name, role, db)
		if err != nil {
			log.Print(err)
		}

		return c.SendString(strconv.Itoa(count))
	})

	app.Get("/UserAccounts/Id/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			log.Print(err)
			c.Status(400).SendString("Incorrect ID")
		}

		db := c.Locals("db").(*sql.DB)

		user_account, err := getUserAccountsById(id, db)
		if err != nil {
			log.Print(err)
		}

		return c.JSON(user_account)
	})

	app.Post("/UserAccounts/Insert", func(c *fiber.Ctx) error {
		user_account := new(UserAccounts)

		if err := c.BodyParser(user_account); err != nil {
			log.Print(err)
			c.Status(400).SendString("Incorrect UserAccounts")
			return err
		}

		message := ""

		db := c.Locals("db").(*sql.DB)

		inserted, err := insertUserAccount(UserAccounts{
			Id:       0,
			IdNumber: user_account.IdNumber,
			FullName: user_account.FullName,
			Username: user_account.Username,
			Password: user_account.Password,
			Section:  user_account.Section,
			Role:     user_account.Role,
		}, db)
		if err != nil {
			log.Print(err)
		}

		if inserted > 0 {
			message = "success"
		}

		return c.SendString(message)
	})

	app.Post("/UserAccounts/Update", func(c *fiber.Ctx) error {
		user_account := new(UserAccounts)

		if err := c.BodyParser(user_account); err != nil {
			log.Print(err)
			c.Status(400).SendString("Incorrect UserAccounts")
			return err
		}

		message := ""

		db := c.Locals("db").(*sql.DB)

		updated, err := updateUserAccount(UserAccounts{
			Id:       user_account.Id,
			IdNumber: user_account.IdNumber,
			FullName: user_account.FullName,
			Username: user_account.Username,
			Password: user_account.Password,
			Section:  user_account.Section,
			Role:     user_account.Role,
		}, db)
		if err != nil {
			log.Print(err)
		}

		if updated > 0 {
			message = "success"
		}

		return c.SendString(message)
	})

	app.Post("/UserAccounts/Delete", func(c *fiber.Ctx) error {
		user_account := new(UserAccounts)

		if err := c.BodyParser(user_account); err != nil {
			log.Print(err)
			c.Status(400).SendString("Incorrect UserAccounts")
			return err
		}

		message := ""

		db := c.Locals("db").(*sql.DB)

		deleted, err := deleteUserAccount(user_account.Id, db)
		if err != nil {
			log.Print(err)
		}

		if deleted > 0 {
			message = "success"
		}

		return c.SendString(message)
	})

	// Randoms

	app.Post("/sample", func(c *fiber.Ctx) error {
		sample := new(Sample)

		if err := c.BodyParser(sample); err != nil {
			log.Print(err)
			c.Status(400).SendString("Incorrect Sample")
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

func getUserAccountsSearch(id_number string, full_name string, role string, db *sql.DB) ([]UserAccounts, error) {
	// A user_accounts slice to hold data from returned rows.
	var user_accounts []UserAccounts

	// Start building the query.
	query := "SELECT * FROM user_accounts WHERE 1=1"

	// Slice to hold the arguments for the query.
	var args []interface{}

	// Add conditions to the query based on the parameters.
	if id_number != "" {
		query += " AND id_number LIKE ?"
		id_number = id_number + "%"
		args = append(args, id_number)
	}
	if full_name != "" {
		query += " AND full_name LIKE ?"
		full_name = full_name + "%"
		args = append(args, full_name)
	}
	if role != "" {
		query += " AND role = ?"
		args = append(args, role)
	}

	// Prepare the statement.
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("getUserAccountsSearch: %v", err)
	}
	defer stmt.Close()

	// Execute the query.
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("getUserAccountsSearch : %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var ua UserAccounts
		if err := rows.Scan(&ua.Id, &ua.IdNumber, &ua.FullName, &ua.Username, &ua.Password, &ua.Section, &ua.Role); err != nil {
			return nil, fmt.Errorf("getUserAccountsSearch : %v", err)
		}
		user_accounts = append(user_accounts, ua)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getUserAccountsSearch : %v", err)
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

func countUserAccounts(id_number string, full_name string, role string, db *sql.DB) (int, error) {
	// A variable to hold the count.
	var count int

	// Start building the query.
	query := "SELECT COUNT(*) FROM user_accounts WHERE 1=1"

	// Slice to hold the arguments for the query.
	var args []interface{}

	// Add conditions to the query based on the parameters.
	if id_number != "" {
		query += " AND id_number LIKE ?"
		id_number = id_number + "%"
		args = append(args, id_number)
	}
	if full_name != "" {
		query += " AND full_name LIKE ?"
		full_name = full_name + "%"
		args = append(args, full_name)
	}
	if role != "" {
		query += " AND role = ?"
		args = append(args, role)
	}

	// Prepare the statement.
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("countUserAccounts: %v", err)
	}
	defer stmt.Close()

	// Execute the query.
	row := stmt.QueryRow(args...)
	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countUserAccounts: %v", err)
	}

	return count, nil
}

func insertUserAccount(ua UserAccounts, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO user_accounts (id_number, full_name, username, password, section, role) VALUES (?, ?, ?, ?, ?, ?)", ua.IdNumber, ua.FullName, ua.Username, ua.Password, ua.Section, ua.Role)
	if err != nil {
		return 0, fmt.Errorf("insertUserAccount: %v", err)
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insertUserAccount: %v", err)
	}
	return inserted, nil
}

func updateUserAccount(ua UserAccounts, db *sql.DB) (int64, error) {
	result, err := db.Exec("UPDATE user_accounts SET id_number = ?, full_name = ?, username = ?, password = ?, section = ?, role = ? WHERE id = ?", ua.IdNumber, ua.FullName, ua.Username, ua.Password, ua.Section, ua.Role, ua.Id)
	if err != nil {
		return 0, fmt.Errorf("updateUserAccount: %v", err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateUserAccount: %v", err)
	}
	return updated, nil
}

func deleteUserAccount(id int, db *sql.DB) (int64, error) {
	result, err := db.Exec("DELETE FROM user_accounts WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("deleteUserAccount: %v", err)
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteUserAccount: %v", err)
	}
	return deleted, nil
}
