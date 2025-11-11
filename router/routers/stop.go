package routers

import (

	"strconv"

	"gym/global/log"


	"gym/controllers/rest"


	"gym/models"
	"github.com/gofiber/fiber/v2"
)

// SetupStopRoutes sets up routes for stop domain
func SetupStopRoutes(group fiber.Router) {

	group.Put("/stop", func(c *fiber.Ctx) error {
			item_ := &models.Stop{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.StopController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/stop", func(c *fiber.Ctx) error {
			item_ := &models.Stop{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.StopController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/stop/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Stop{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.StopController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/stop/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.StopController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/stop", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.StopController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/stop/count", func(c *fiber.Ctx) error {

		var controller rest.StopController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/stop", func(c *fiber.Ctx) error {
			item_ := &models.Stop{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.StopController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/stop/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Stop{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.StopController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}