package routers

import (

	"strconv"

	"gym/global/log"


	"gym/controllers/rest"


	"gym/models"
	"github.com/gofiber/fiber/v2"
)

// SetupMemberbodyRoutes sets up routes for memberbody domain
func SetupMemberbodyRoutes(group fiber.Router) {

	group.Delete("/memberbody/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Memberbody{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/memberbody/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/memberbody", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/memberbody/count", func(c *fiber.Ctx) error {

		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/memberbody", func(c *fiber.Ctx) error {
			item_ := &models.Memberbody{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/memberbody/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Memberbody{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/memberbody", func(c *fiber.Ctx) error {
			item_ := &models.Memberbody{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/memberbody", func(c *fiber.Ctx) error {
			item_ := &models.Memberbody{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.MemberbodyController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}