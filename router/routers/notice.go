package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupNoticeRoutes sets up routes for notice domain
func SetupNoticeRoutes(group fiber.Router) {

	group.Get("/notice", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.NoticeController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/notice/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.NoticeController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/notice", func(c *fiber.Ctx) error {
		item_ := &models.Notice{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NoticeController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/notice/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Notice
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NoticeController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/notice/count", func(c *fiber.Ctx) error {

		var controller rest.NoticeController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/notice/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notice{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NoticeController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/notice/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notice{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NoticeController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/notice/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Notice{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.NoticeController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/notice/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Notice{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.NoticeController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}
