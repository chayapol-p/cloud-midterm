package controllers

import (
	"fmt"
	"time"

	"github.com/chayapol-p/cloud-midterm-server/database"
	"github.com/chayapol-p/cloud-midterm-server/models"
	"github.com/gofiber/fiber/v2"
)

func GetMessages(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all messages.
	messages, err := db.GetMessages(c.Params("timestamp"))
	if err != nil {
		// Return, if messages not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "messages were not found",
			"count":    0,
			"messages": nil,
			"err":      err.Error(),
		})
	}

	// Get all updates.
	updates, err := db.GetUpdates()
	if err != nil {
		// Return, if messages not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "updates were not found",
			"count":    0,
			"messages": nil,
			"err":      err.Error(),
		})
	}

	for i := len(updates) - 1; i >= 0; i-- {
		if updates[i].IsDeleted {
			continue
		}
		for j := 0; j < len(messages); j++ {
			if updates[i].UUID == messages[j].UUID {
				if updates[i].Author != "" {
					fmt.Println("Changed Author")
					messages[j].Author = updates[i].Author
				}
				if updates[i].Message != "" {
					messages[j].Message = updates[i].Message
				}
				if updates[i].Likes != -1 {
					messages[j].Likes = updates[i].Likes
				}
				updates = append(updates[:i], updates[i+1:]...)
				break
			}
		}
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"count":    len(messages),
		"messages": messages,
		"updates":  updates,
	})
}

func CreateMessage(c *fiber.Ctx) error {
	// Create new Book struct
	message := &models.Message{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(message); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set initialized default data for message:
	message.Timestamp = time.Now().UTC()

	// Checking, if message with given ID is exists.
	if _, err := db.GetMessage(message.UUID); err == nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this UUID is already existed",
		})
	}

	// Checking, if message with given ID is exists.
	if _, err := db.GetUpdate(message.UUID); err == nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this UUID is already existed and deleted.",
		})
	}

	if err := db.CreateMessage(message); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"message": message,
	})
}

func UpdateMessage(c *fiber.Ctx) error {
	// Create new Book struct
	update := &models.Message{}
	uuid := c.Params("uuid")
	// Check, if received JSON data is valid.
	if err := c.BodyParser(update); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if message with given ID is exists.
	foundedMessage, err := db.GetMessage(uuid)
	if err != nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this UUID not found",
		})
	}

	update_founded := true
	// Merging if already changed
	foundedUpdate, err := db.GetUpdate(uuid)
	if err != nil {
		// If not been updated, Create blank update
		update_founded = false
		foundedUpdate = models.UpdatedMessage{uuid, time.Now(), "", "", -1, false}
	}

	new_update := models.UpdatedMessage{uuid, time.Now(), "", "", -1, false}
	if foundedUpdate.Author == "" {
		if foundedMessage.Author != update.Author {
			new_update.Author = update.Author
		}
	} else if foundedUpdate.Author != "" {
		if update.Author == foundedMessage.Author {
			new_update.Author = ""
		} else {
			new_update.Author = update.Author
		}
	}
	if foundedUpdate.Message == "" {
		if foundedMessage.Message != update.Message {
			new_update.Message = update.Message
		}
	} else if foundedUpdate.Message != "" {
		if update.Message == foundedMessage.Message {
			new_update.Message = ""
		} else {
			new_update.Message = update.Message
		}
	}
	if foundedUpdate.Likes == -1 {
		if foundedMessage.Likes != update.Likes {
			new_update.Likes = update.Likes
		}
	} else if foundedUpdate.Likes != -1 {
		if update.Likes == foundedMessage.Likes {
			new_update.Likes = -1
		} else {
			new_update.Likes = update.Likes
		}
	}

	if update_founded {
		if err := db.UpdateUpdate(uuid, &new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	} else {
		// Update message by given ID.
		if err := db.CreateUpdate(&new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	// Return status 204.
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteMessage(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if message with given ID is exists.
	if _, err := db.GetMessage(uuid); err != nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this ID not found",
		})
	}

	// Delete message by given ID.
	if err := db.DeleteMessage(uuid); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Set Delete flag to update
	new_update := models.UpdatedMessage{uuid, time.Now(), "", "", -1, true}
	if _, err := db.GetUpdate(uuid); err != nil {
		if err := db.CreateUpdate(&new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		// If not been updated, Create blank update with deleted
	} else {
		if err := db.UpdateUpdate(uuid, &new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
