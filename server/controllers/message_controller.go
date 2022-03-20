package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chayapol-p/cloud-midterm-server/database"
	"github.com/chayapol-p/cloud-midterm-server/models"
	"github.com/gofiber/fiber/v2"
)

var DB *database.Queries
var err error

func InitialDB() {
	// Create database connection.
	DB, err = database.OpenDBConnection()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetMessages(c *fiber.Ctx) error {

	// Get all messages.
	messages, err := DB.GetMessages(c.Params("timestamp"), c.Query("offset"), c.Query("limit"))
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
	updates, err := DB.GetUpdates(c.Params("timestamp"), c.Query("offset"), c.Query("limit"))
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

	formatted_messages := make([]models.OutputMessage, len(messages))
	for j := 0; j < len(messages); j++ {
		formatted_message := models.OutputMessage{messages[j].UUID, messages[j].Author, messages[j].Message, messages[j].Likes}
		formatted_messages[j] = formatted_message
		if !messages[j].UpdatedAuthor.Valid {
			// This mean this id have no update
			continue
		} else {
			if messages[j].UpdatedAuthor.String != "" {
				fmt.Println("Changed Author")
				formatted_messages[j].Author = messages[j].UpdatedAuthor.String
			}
			if messages[j].UpdatedMessage.String != "" {
				formatted_messages[j].Message = messages[j].UpdatedMessage.String
			}
			if messages[j].UpdatedLikes.Int32 != -1 {
				formatted_messages[j].Likes = int(messages[j].UpdatedLikes.Int32)
			}
		}
	}

	fmt.Printf("Query pass")

	// Return status 200 OK.
	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(fiber.Map{
		"error":         false,
		"msg":           nil,
		"count_message": len(formatted_messages),
		"count_update":  len(updates),
		"messages":      formatted_messages,
		"updates":       updates,
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

	// Set initialized default data for message:
	message.Timestamp = time.Now().UTC()

	// // Checking, if message with given ID is exists.
	// if _, err := DB.GetMessage(message.UUID); err == nil {
	// 	// Return status 404 and message not found error.
	// 	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "message with this UUID is already existed",
	// 	})
	// }

	// // Checking, if message with given ID is exists.
	// if _, err := DB.GetUpdate(message.UUID); err == nil {
	// 	// Return status 404 and message not found error.
	// 	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "message with this UUID is already existed and deleted.",
	// 	})
	// }

	if err := DB.CreateMessage(message); err != nil {
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

	// Checking, if message with given ID is exists.
	foundedMessage, err := DB.GetMessage(uuid)
	if err != nil {
		// Return status 404 and message not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "message with this UUID not found",
		})
	}

	update_founded := true
	// Merging if already changed
	foundedUpdate, err := DB.GetUpdate(uuid)
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
		if err := DB.UpdateUpdate(uuid, &new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	} else {
		// Update message by given ID.
		if err := DB.CreateUpdate(&new_update); err != nil {
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

	fmt.Printf("Deleting\n")
	// Delete message by given ID.
	if err := DB.DeleteMessage(uuid); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	fmt.Printf("Deleted\n")

	//Set Delete flag to update
	new_update := models.UpdatedMessage{uuid, time.Now(), "", "", -1, true}
	if _, err := DB.GetUpdate(uuid); err != nil {
		if err := DB.CreateUpdate(&new_update); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		// If not been updated, Create blank update with deleted
	} else {
		if err := DB.UpdateUpdate(uuid, &new_update); err != nil {
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
