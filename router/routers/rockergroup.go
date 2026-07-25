package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupRockergroupRoutes sets up routes for rockergroup domain
func SetupRockergroupRoutes(group fiber.Router) {

	group.Post("/rockergroup/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/rockergroup/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/rockergroup/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/rockergroup/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/rockergroup/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/rockergroup/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/rockergroup", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/rockergroup/count", func(c *fiber.Ctx) error {

		var controller rest.RockergroupController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/rockergroup", func(c *fiber.Ctx) error {
		item_ := &models.Rockergroup{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.RockergroupController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}
