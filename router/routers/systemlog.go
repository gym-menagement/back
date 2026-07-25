package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupSystemlogRoutes sets up routes for systemlog domain
func SetupSystemlogRoutes(group fiber.Router) {

	group.Get("/systemlog", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/systemlog/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/systemlog", func(c *fiber.Ctx) error {
		item_ := &models.Systemlog{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/systemlog/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Systemlog
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/systemlog/count", func(c *fiber.Ctx) error {

		var controller rest.SystemlogController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/systemlog/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Systemlog{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/systemlog/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Systemlog{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/systemlog/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Systemlog{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/systemlog/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Systemlog{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.SystemlogController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}
