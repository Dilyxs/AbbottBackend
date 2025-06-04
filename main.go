package main

import (
	"adsayan.com/test/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4321, http://abbotthustlers.com, https://abbotthustlers.com, https://abbott.adsayan.com",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Post("/ClientAddition", utils.AddClient)
	app.Get("/AllClient", utils.GetAllClients)
	app.Post("/Register", utils.RegisterUser)
	app.Post("/Login", utils.LoginUser)
	app.Get("/isCookie", utils.ValidateCookie)
	app.Post("/userInformation", utils.ReturnUserInformation)
	app.Post("/AddCalendarInfo", utils.PostCalendarInfo)
	app.Get("/CalendarInfo", utils.FetchCalendarInformation)
	app.Delete("/DeleteCalendarEvent", utils.DeleteCalendarEvent)
	app.Post("/SingleClient", utils.ReturnASingleClient)
	app.Get("/AllNote", utils.FetchAllNotes)
	app.Post("/PostNote", utils.UploadANote)
	app.Delete("/DeleteNote", utils.DeleteANote)
	app.Put("/UpdateNote", utils.UpdateANote)
	app.Delete("/DeleteFinance", utils.DeleteAFinanceDetail)
	app.Post("/ChangeFinance", utils.ChangeHasBeenTaken)

	app.Listen("0.0.0.0:4330")
}
