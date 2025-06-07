package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AddClient(c *fiber.Ctx) error {
	var client Client
	if err := c.BodyParser(&client); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := InsertClient(client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert client"})
	}

	return c.Status(200).JSON(fiber.Map{"result": "success"})

}

func RegisterUser(c *fiber.Ctx) error {
	var user RegisterIngo

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := InsertUser(user.Email, user.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"message": "ok"})
}

func GetAllClients(c *fiber.Ctx) error {
	data, err := GetEveryClientDB()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "could not fetch client"})
	} else {
		return c.Status(200).JSON(data)
	}
}

func LoginUser(c *fiber.Ctx) error {
	var user RegisterIngo

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "wrong format data"})
	}
	if err, id := VerifyUser(user.Email, user.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "wrong password"})
	} else { // correct data
		c.Cookie(&fiber.Cookie{
			Name:     "session_token",
			Value:    fmt.Sprintf("%v", id),
			Path:     "/",
			HTTPOnly: true,
			Domain:   "abbotthustlers.com",
			Secure:   true, //change to true if in production
			SameSite: fiber.CookieSameSiteNoneMode,
			MaxAge:   3600 * 24 * 7, //once a week

		})
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}

func ValidateCookie(c *fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	if cookie == "" {
		return c.Status(400).JSON(fiber.Map{
			"Auth": "false",
			"id":   "none",
		})
	}

	userID := cookie

	return c.Status(200).JSON(fiber.Map{
		"Auth": true,
		"id":   userID,
	})
}

func ReturnUserInformation(c *fiber.Ctx) error {
	var id SendUserID

	if err := c.BodyParser(&id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "wrong format"})
	}
	if info, err := FetchUserData(id.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	} else {
		return c.Status(200).JSON(info)
	}

}

func PostCalendarInfo(c *fiber.Ctx) error {
	var data CalenderInfoReceived

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "wrong format"})
	}
	if err := AddACalendarInfo(data); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "ok"})

}

func FetchCalendarInformation(c *fiber.Ctx) error {
	if info, err := FetchCalendarInfo(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	} else {
		return c.Status(200).JSON(info)
	}

}

func DeleteCalendarEvent(c *fiber.Ctx) error {
	var id SendUserID
	if err := c.BodyParser(&id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := DeleteACalendarElementDB(id.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})

	}
	return c.Status(200).JSON(fiber.Map{"message": "ok"})
}

func ReturnASingleClient(c *fiber.Ctx) error {
	var id SendUserID
	if err := c.BodyParser(&id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if data, err := ReturnDataSingleClient(id.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	} else {
		return c.Status(200).JSON(data)
	}
}

func FetchAllNotes(c *fiber.Ctx) error {
	if data, err := FetchAllNotesDB(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(data)
	}
}

func UploadANote(c *fiber.Ctx) error {
	var SingleNote NotesDetailSender
	if err := c.BodyParser(&SingleNote); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "wrong format"})
	}
	var has_sent_client bool = false
	if SingleNote.ClientRelation != 0 {
		has_sent_client = true
	}
	if err := UploadANoteDB(SingleNote, has_sent_client); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}

func DeleteANote(c *fiber.Ctx) error {
	var id SendUserID
	if err := c.BodyParser(&id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	if err := DeleteASpecificNoteDB(id.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}

func UpdateANote(c *fiber.Ctx) error {
	var SingleNote NotesDetailUpdater
	if err := c.BodyParser(&SingleNote); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := UpdateANoteDB(SingleNote); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}

}

func FetchAllFinanceDetail(c *fiber.Ctx) error {
	if AllNotes, err := FetchAllFinanceDetailDB(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": `${err}`})
	} else {
		return c.Status(200).JSON(AllNotes)
	}
}

func InsertAFinanceDetail(c *fiber.Ctx) error {
	var Details FinanceDetailSend
	if err := c.BodyParser(&Details); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := InsertAFinanceDetailDB(Details); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}

func DeleteAFinanceDetail(c *fiber.Ctx) error {
	var Id SendUserID
	if err := c.BodyParser(&Id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := DeleteAFinanceDetailDB(Id.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}

}

func ChangeHasBeenTaken(c *fiber.Ctx) error {
	var Ids FinanceDetailChangeHasBeenTaken
	if err := c.BodyParser(&Ids); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := UpdateUtilizedFinanceDetailDB(Ids.Ids); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}

func InsertAnImage(c *fiber.Ctx) error {
	var data ImageDataSent
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}
	if err := UploadAnImageDB(data.Url, data.Title, data.Description); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	} else {
		return c.Status(200).JSON(fiber.Map{"message": "ok"})

	}
}
