package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     string    `gorm:"uniqueIndex"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AuthResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

var (
	db    *gorm.DB
	store *session.Store
)

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	db.AutoMigrate(&User{})

	seedDemoUser()

	store = session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSecure:   false,
		CookieSameSite: "Lax",
	})

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	app.Get("/login", authRedirect, renderLogin)
	app.Post("/login", handleLogin)

	app.Get("/register", authRedirect, renderRegister)
	app.Post("/register", handleRegister)

	app.Get("/dashboard", requireAuth, renderDashboard)
	app.Post("/logout", handleLogout)

	app.Post("/api/validate/email", validateEmail)
	app.Post("/api/validate/password", validatePassword)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🔮 Hologram Auth Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func seedDemoUser() {
	var count int64
	db.Model(&User{}).Where("phone = ?", "+1 (555) 123-4567").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
		demoUser := User{
			Email:    "demo@hologram.io",
			Phone:    "+1 (555) 123-4567",
			Password: string(hashedPassword),
		}
		db.Create(&demoUser)
		log.Println("📱 Demo user created: demo@hologram.io / demo123")
	}
}

func authRedirect(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		if userID := sess.Get("user_id"); userID != nil {
			return c.Redirect("/dashboard")
		}
	}
	return c.Next()
}

func requireAuth(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Redirect("/login")
	}

	if userID := sess.Get("user_id"); userID == nil {
		return c.Redirect("/login")
	}

	return c.Next()
}

func renderLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Hologram Access",
	})
}

func renderRegister(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Initialize Profile",
	})
}

func renderDashboard(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	userID := sess.Get("user_id")

	var user User
	db.First(&user, userID)

	return c.Render("dashboard", fiber.Map{
		"Title": "Hologram Hub",
		"Email": user.Email,
		"Phone": user.Phone,
	})
}

func handleLogin(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")

	var errors []ValidationError

	if email == "" {
		errors = append(errors, ValidationError{Field: "email", Message: "Email frequency required"})
	}
	if password == "" {
		errors = append(errors, ValidationError{Field: "password", Message: "Access cipher required"})
	}

	if len(errors) > 0 {
		return c.Status(400).JSON(AuthResponse{Success: false, Errors: errors})
	}

	var user User
	result := db.Where("email = ?", strings.ToLower(email)).First(&user)

	if result.Error != nil {
		return c.Status(401).JSON(AuthResponse{
			Success: false,
			Errors:  []ValidationError{{Field: "email", Message: "Identity not found in holomatrix"}},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(401).JSON(AuthResponse{
			Success: false,
			Errors:  []ValidationError{{Field: "password", Message: "Cipher mismatch detected"}},
		})
	}

	sess, _ := store.Get(c)
	sess.Set("user_id", user.ID)
	sess.Save()

	return c.JSON(AuthResponse{Success: true, Message: "Hologram link established"})
}

func handleRegister(c *fiber.Ctx) error {
	email := strings.TrimSpace(strings.ToLower(c.FormValue("email")))
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	var errors []ValidationError

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if email == "" {
		errors = append(errors, ValidationError{Field: "email", Message: "Email frequency required"})
	} else if !emailRegex.MatchString(email) {
		errors = append(errors, ValidationError{Field: "email", Message: "Invalid frequency pattern"})
	} else {
		var count int64
		db.Model(&User{}).Where("email = ?", email).Count(&count)
		if count > 0 {
			errors = append(errors, ValidationError{Field: "email", Message: "Frequency already registered"})
		}
	}

	if len(password) < 6 {
		errors = append(errors, ValidationError{Field: "password", Message: "Cipher must be 6+ characters"})
	} else {
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

		if !hasUpper || !hasLower || !hasNumber {
			errors = append(errors, ValidationError{Field: "password", Message: "Cipher needs uppercase, lowercase, and number"})
		}
	}

	if password != confirmPassword {
		errors = append(errors, ValidationError{Field: "confirm_password", Message: "Cipher confirmation failed"})
	}

	if len(errors) > 0 {
		return c.Status(400).JSON(AuthResponse{Success: false, Errors: errors})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(AuthResponse{
			Success: false,
			Message: "Holomatrix encryption failed",
		})
	}

	user := User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if result := db.Create(&user); result.Error != nil {
		return c.Status(500).JSON(AuthResponse{
			Success: false,
			Message: "Failed to materialize identity",
		})
	}

	sess, _ := store.Get(c)
	sess.Set("user_id", user.ID)
	sess.Save()

	return c.JSON(AuthResponse{Success: true, Message: "Identity materialized successfully"})
}

func handleLogout(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	sess.Destroy()
	return c.Redirect("/login")
}

func validateEmail(c *fiber.Ctx) error {
	email := strings.TrimSpace(strings.ToLower(c.FormValue("email")))

	if email == "" {
		return c.JSON(fiber.Map{"valid": false, "message": "Email frequency required"})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return c.JSON(fiber.Map{"valid": false, "message": "Invalid frequency pattern"})
	}

	var count int64
	db.Model(&User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return c.JSON(fiber.Map{"valid": false, "message": "Frequency already in holomatrix"})
	}

	return c.JSON(fiber.Map{"valid": true, "message": "Frequency available"})
}

func validatePassword(c *fiber.Ctx) error {
	password := c.FormValue("password")

	strength := 0
	messages := []string{}

	if len(password) >= 6 {
		strength++
	} else {
		messages = append(messages, "Minimum 6 characters")
	}

	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		strength++
	} else {
		messages = append(messages, "Add uppercase letter")
	}

	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		strength++
	} else {
		messages = append(messages, "Add lowercase letter")
	}

	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		strength++
	} else {
		messages = append(messages, "Add number")
	}

	if regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		strength++
	}

	if len(password) >= 12 {
		strength++
	}

	return c.JSON(fiber.Map{
		"strength": strength,
		"max":      6,
		"messages": messages,
		"valid":    strength >= 4,
	})
}
