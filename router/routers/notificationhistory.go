package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupNotificationhistoryRoutes sets up routes for notificationhistory domain
func SetupNotificationhistoryRoutes(group fiber.Router) {

	group.Get("/notificationhistory", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/notificationhistory/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/notificationhistory", func(c *fiber.Ctx) error {
		item_ := &models.Notificationhistory{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/notificationhistory/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Notificationhistory
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/notificationhistory/count", func(c *fiber.Ctx) error {

		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/notificationhistory/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notificationhistory{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/notificationhistory/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notificationhistory{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/notificationhistory/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notificationhistory{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/notificationhistory/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Notificationhistory{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NotificationhistoryController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}
