package controllers

import (
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
	messages, err := db.GetMessages()
	if err != nil {
		// Return, if messages not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "messages were not found",
			"count":    0,
			"messages": nil,
		})
	}

	// Get all updates.
	updates, err := db.GetUpdates()
	if err != nil {
		// Return, if messages not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "messages were not found",
			"count":    0,
			"messages": nil,
		})
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

// func GetMessage(c *fiber.Ctx) error {
// 	// Catch message ID from URL.
// 	uuid, err := uuid.Parse(c.Params("uuid"))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	// Create database connection.
// 	db, err := database.OpenDBConnection()
// 	if err != nil {
// 		// Return status 500 and database connection error.
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	// Get message by UUID.
// 	message, err := db.GetMessage(uuid)
// 	if err != nil {
// 		// Return, if message not found.
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error":   true,
// 			"msg":     "message with the given ID is not found",
// 			"message": nil,
// 		})
// 	}

// 	// Return status 200 OK.
// 	return c.JSON(fiber.Map{
// 		"error":   false,
// 		"msg":     nil,
// 		"message": message,
// 	})
// }

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
	foundedMessage, err := db.GetMessage(update.UUID)
	if err != nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this ID not found",
		})
	}

	// Merging if already changed
	foundedUpdate, err := db.GetUpdate(update.UUID)
	if err != nil {
		// If not been updated, Create blank update
		foundedUpdate = models.UpdatedMessage{update.UUID, time.Now(), "", "", -1, false}
	}

	new_update := models.UpdatedMessage{update.UUID, time.Now(), "", "", -1, false}
	if foundedUpdate.Author == "" && foundedMessage.Author != update.Author {
		new_update.Author = update.Author
	} else if update.Author != foundedUpdate.Author {
		if update.Author == foundedMessage.Author {
			new_update.Author = ""
		} else {
			new_update.Author = update.Author
		}
	}
	if foundedUpdate.Message == "" && foundedMessage.Message != update.Message {
		new_update.Message = update.Message
	} else if update.Message != foundedUpdate.Message {
		if update.Message == foundedMessage.Message {
			new_update.Message = ""
		} else {
			new_update.Message = update.Message
		}
	}
	if foundedUpdate.Likes == -1 && foundedMessage.Likes != update.Likes {
		new_update.Likes = update.Likes
	} else if update.Likes != foundedUpdate.Likes {
		if update.Likes == foundedMessage.Likes {
			new_update.Likes = -1
		} else {
			new_update.Likes = update.Likes
		}
	}

	// Update message by given ID.
	if err := db.UpdateUpdate(foundedMessage.UUID, &new_update); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204.
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteMessage(c *fiber.Ctx) error {
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

	// Checking, if message with given ID is exists.
	foundedMessage, err := db.GetMessage(message.UUID)
	if err != nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this ID not found",
		})
	}

	// Delete message by given ID.
	if err := db.DeleteMessage(foundedMessage.UUID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Noted in update table as delete
	new_update := models.UpdatedMessage{foundedMessage.UUID, time.Now(), "", "", -1, true}
	if err := db.UpdateUpdate(foundedMessage.UUID, &new_update); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
