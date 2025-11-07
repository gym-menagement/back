package routers

import (

	"encoding/json"

	"strconv"

	"gym/global/log"


	"gym/controllers/rest"


	"gym/models/user"

	"gym/models"
	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up routes for user domain
func SetupUserRoutes(group fiber.Router) {

	group.Get("/user/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.UserController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/user", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.UserController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/user/count", func(c *fiber.Ctx) error {

		var controller rest.UserController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/user", func(c *fiber.Ctx) error {
			item_ := &models.UserUpdate{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.UserController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/user/batch", func(c *fiber.Ctx) error {
			var results map[string]interface{}
			jsonData := c.Body()
			jsonErr := json.Unmarshal(jsonData, &results)
			if jsonErr != nil {
			    log.Error().Msg(jsonErr.Error())
			}
			var items_ *[]models.UserUpdate
			items__ref := &items_
			err := c.BodyParser(items__ref)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.UserController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/user", func(c *fiber.Ctx) error {
			item_ := &models.UserUpdate{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.UserController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/user", func(c *fiber.Ctx) error {
			item_ := &models.User{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.UserController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/user/count/loginid/:loginid", func(c *fiber.Ctx) error {
			loginid_ := c.Params("loginid")
		var controller rest.UserController
		controller.Init(c)
		controller.CountByLoginid(loginid_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/user/batch", func(c *fiber.Ctx) error {
			item_ := &[]models.User{}
			err := c.BodyParser(item_)
			if err != nil {
			    log.Error().Msg(err.Error())
			}
		var controller rest.UserController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/user/get/loginid/:loginid", func(c *fiber.Ctx) error {
			loginid_ := c.Params("loginid")
		var controller rest.UserController
		controller.Init(c)
		controller.GetByLoginid(loginid_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/user/get/connectid/:connectid", func(c *fiber.Ctx) error {
			connectid_ := c.Params("connectid")
		var controller rest.UserController
		controller.Init(c)
		controller.GetByConnectid(connectid_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/user/logindatebyid", func(c *fiber.Ctx) error {
			var results map[string]interface{}
			jsonData := c.Body()
			jsonErr := json.Unmarshal(jsonData, &results)
			if jsonErr != nil {
			    log.Error().Msg(jsonErr.Error())
			}
			var logindate_ string
			if v, flag := results["logindate"]; flag {
				logindate_ = v.(string)
			}
			var id_ int64
			if v, flag := results["id"]; flag {
				id_ = int64(v.(float64))
			}
		var controller rest.UserController
		controller.Init(c)
		controller.UpdateLogindateById(logindate_, id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/user/find/level/:level", func(c *fiber.Ctx) error {
			var level_ user.Level
			level__, _ := strconv.Atoi(c.Params("level"))
			level_ = user.Level(level__)
		var controller rest.UserController
		controller.Init(c)
		controller.FindByLevel(level_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}