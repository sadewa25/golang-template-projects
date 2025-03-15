package main

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/send-email", func(c *fiber.Ctx) error {
		type EmailRequest struct {
			From    string `json:"from"`
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}

		// Parse request body
		var req EmailRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bad Request",
				"msg":   "Invalid request payload",
			})
		}

		// Mailpit server details
		smtpHost := "localhost"
		smtpPort := 1025 // Default Mailpit SMTP port

		// Sender and recipient data
		from := req.From
		to := []string{req.To}

		contentType := "text/html"

		templateHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Email Notification</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333;
					max-width: 600px;
					margin: 0 auto;
				}
				.header {
					background-color: #4285f4;
					color: white;
					padding: 20px;
					text-align: center;
				}
				.content {
					padding: 20px;
					background-color: #f9f9f9;
				}
				.footer {
					text-align: center;
					padding: 10px;
					font-size: 12px;
					color: #666;
					border-top: 1px solid #eee;
				}
				.button {
					display: inline-block;
					background-color: #4285f4;
					color: white;
					text-decoration: none;
					padding: 10px 20px;
					border-radius: 4px;
					margin-top: 15px;
				}
			</style>
		</head>
		<body>
			<div class="header">
				<h1>{{.Subject}}</h1>
			</div>
			<div class="content">
				<p>Hello,</p>
				<p>{{.Body}}</p>
				<p>Thank you for using our service!</p>
				<a href="https://example.com" class="button">Learn More</a>
			</div>
			<div class="footer">
				<p>© 2025 Your Company. All rights reserved.</p>
				<p>You're receiving this email because you signed up for notifications.</p>
			</div>
		</body>
		</html>
		`

		// Replace template placeholders with actual content
		htmlContent := strings.Replace(templateHTML, "{{.Subject}}", req.Subject, -1)
		htmlContent = strings.Replace(htmlContent, "{{.Body}}", req.Body, -1)

		// Email message
		subject := req.Subject
		message := []byte(fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: %s; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n", from, to[0], subject, contentType, htmlContent))

		// Send the email without authentication
		// Mailpit doesn't require authentication for local development
		err := smtp.SendMail(
			fmt.Sprintf("%s:%d", smtpHost, smtpPort),
			nil, // No authentication needed for Mailpit
			from,
			to,
			message,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"msg":   "Failed to send email",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "",
			"msg":   "Sent email successfully",
		})
	})

	app.Listen(":3000")
}
