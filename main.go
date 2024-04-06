/*
 * Copyright (C) 2024- Germano Rizzo
 *
 * This file is part of MailRG.
 *
 * MailRG is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * MailRG is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with FoodHubber.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/wneessen/go-mail"
)

type body struct {
	Token       string   `json:"token"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Bcc         []string `json:"bcc"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	HTML        bool     `json:"html"`
	Attachments []string `json:"attachments"`
}

func main() {
	smtp, err := mail.NewClient(
		os.Getenv("SMTP_SERVER"),
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("SMTP_USER")),
		mail.WithPassword(os.Getenv("SMTP_PASS")),
	)
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("smtp", smtp)
		c.Locals("token", os.Getenv("MAILRG_TOKEN"))
		c.Locals("dataDir", os.Getenv("DATA_DIR"))
		return c.Next()
	})

	app.Post("/send", send)

	log.Fatal(app.Listen(":2163"))
}

func send(c *fiber.Ctx) error {
	var smtp *mail.Client = c.Locals("smtp").(*mail.Client)
	var token = c.Locals("token").(string)
	var dataDir = c.Locals("dataDir").(string)

	var body body
	err := c.BodyParser(&body)
	if err != nil {
		c.SendString(fmt.Sprintf("error parsing body: %s", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if body.Token != token {
		c.SendString("wrong token")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	m := mail.NewMsg()
	if err := m.From(body.From); err != nil {
		c.SendString(fmt.Sprintf("failed to set From address: %s", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := m.To(body.To...); err != nil {
		c.SendString(fmt.Sprintf("failed to set To address: %s", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := m.Cc(body.Cc...); err != nil {
		c.SendString(fmt.Sprintf("failed to set Cc address: %s", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := m.Bcc(body.Bcc...); err != nil {
		c.SendString(fmt.Sprintf("failed to set Bcc address: %s", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	m.Subject(body.Subject)

	if body.HTML {
		m.SetBodyString(mail.TypeTextHTML, body.Body)
	} else {
		m.SetBodyString(mail.TypeTextPlain, body.Body)
	}

	for _, a := range body.Attachments {
		m.AttachFile(filepath.Join(dataDir, a), mail.WithFileName(a))
	}

	if err := smtp.DialAndSend(m); err != nil {
		c.SendString(fmt.Sprintf("failed to send email: %s", err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendString("Mail sent")
}
