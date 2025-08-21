package routers

import (

	"strconv"

	"gym/global/log"


	"gym/controllers/rest"


	"gym/models"
	"github.com/gofiber/fiber/v2"
)

// SetupDiscountRoutes sets up routes for discount domain
func SetupDiscountRoutes(group fiber.Router) {

	group.Post("/discount", func(c *fiber.Ctx) error {
			item_ := &models.Discount{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.DiscountController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/discount/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Discount{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.DiscountController
		controller.Init(c)
		controller.Insertbatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/discount", func(c *fiber.Ctx) error {
			item_ := &models.Discount{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.DiscountController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/discount", func(c *fiber.Ctx) error {
			item_ := &models.Discount{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.DiscountController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/discount/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.Discount{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.DiscountController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/discount/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.DiscountController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/discount", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.DiscountController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/discount/count", func(c *fiber.Ctx) error {

		var controller rest.DiscountController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

}