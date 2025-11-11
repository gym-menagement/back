package routers

import (

	"strconv"

	"gym/global/log"


	"gym/controllers/rest"


	"gym/models"
	"github.com/gofiber/fiber/v2"
)

// SetupPtreservationRoutes sets up routes for ptreservation domain
func SetupPtreservationRoutes(group fiber.Router) {

	group.Post("/ptreservation", func(c *fiber.Ctx) error {
			item_ := &models.Ptreservation{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/ptreservation/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Ptreservation{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/ptreservation", func(c *fiber.Ctx) error {
			item_ := &models.Ptreservation{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/ptreservation", func(c *fiber.Ctx) error {
			item_ := &models.Ptreservation{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/ptreservation/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Ptreservation{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/ptreservation/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/ptreservation", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.PtreservationController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/ptreservation/count", func(c *fiber.Ctx) error {

		var controller rest.PtreservationController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

}