package routers

import (
	"strconv"

	"gym/global/log"

	"gym/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"gym/models"
)

// SetupGymadminRoutes sets up routes for gymadmin domain
func SetupGymadminRoutes(group fiber.Router) {

	group.Get("/gymadmin", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.GymadminController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/gymadmin/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.GymadminController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/gymadmin", func(c *fiber.Ctx) error {
		item_ := &models.Gymadmin{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.GymadminController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/gymadmin/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Gymadmin
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.GymadminController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/gymadmin/count", func(c *fiber.Ctx) error {

		var controller rest.GymadminController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/gymadmin/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Gymadmin{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.GymadminController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Patch("/gymadmin/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Gymadmin{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.GymadminController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/gymadmin/:id<int>", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		item_ := &models.Gymadmin{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		item_.Id = id_
		var controller rest.GymadminController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/gymadmin/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Gymadmin{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		var controller rest.GymadminController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	// 프론트 호환 검색: bare 배열 반환 (res.data 가 배열)
	group.Get("/gymadmin/search/gym", func(c *fiber.Ctx) error {
		gym_, _ := strconv.ParseInt(c.Query("gym"), 10, 64)
		var controller rest.GymadminController
		controller.Init(c)
		controller.FindByGym(gym_)
		controller.Close()
		return c.JSON(controller.Result["items"])
	})

	group.Get("/gymadmin/search/user", func(c *fiber.Ctx) error {
		user_, _ := strconv.ParseInt(c.Query("user"), 10, 64)
		var controller rest.GymadminController
		controller.Init(c)
		controller.FindByUser(user_)
		controller.Close()
		return c.JSON(controller.Result["items"])
	})

}
