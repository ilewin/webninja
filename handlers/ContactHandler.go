package handlers

import (
	"log"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"webp.ninja/utils"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func ContactHandler(c *fiber.Ctx) error {

	config := utils.GetConfig()

	form := &ContactForm{}

	if err := c.BodyParser(form); err != nil {
		panic(err)
	}

	to := config.Email_To
	subject := config.Email_Subject
	body := "To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + "Name: " + form.Name + "\r\n\r\n" + "Email: " + form.Email + "\r\n\r\n" + "Message: " + form.Message + "\r\n\r\n"
	auth := smtp.PlainAuth("", config.Smtp_Login, config.Smtp_Pass, config.Smtp_Host)
	err := smtp.SendMail(config.Smtp_Server, auth, config.Email_From, []string{to}, []byte(body))
	if err != nil {
		log.Print("ERROR: attempting to send a mail ", err)
		code := fiber.StatusInternalServerError
		return c.Status(code).SendString(err.Error())
	}

	type SentStruct struct {
		Message string `json:"successMessage"`
	}
	sent := SentStruct{
		Message: "Mail was sent",
	}

	return c.JSON(sent)
}
