package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupRockerusageRoutes sets up routes for rockerusage domain
func SetupRockerusageRoutes(group fiber.Router) {

	group.Get("/rockerusage/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/rockerusage", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/rockerusage/count", func(c *fiber.Ctx) error {

		var controller rest.RockerusageController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/rockerusage", func(c *fiber.Ctx) error {
		item_ := &models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/rockerusage/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/rockerusage/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/rockerusage/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/rockerusage/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/rockerusage/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Rockerusage{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockerusageController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}
