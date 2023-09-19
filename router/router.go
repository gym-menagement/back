package router

import (
	"gym/controllers/rest"
	"gym/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(app *fiber.App) {
	app.Get("/api/jwt", func(ctx *fiber.Ctx) error {
		loginid := ctx.Query("loginid")
		passwd := ctx.Query("passwd")
		return ctx.JSON(JwtAuth(loginid, passwd))
	})
	// app.Get("/api/jwt/token", func(ctx *fiber.Ctx) error {
	// 	token := ctx.Get("Authorization")
	// 	return ctx.JSON(JwtToken(token))
	// })
	apiGroup := app.Group("/api")
	apiGroup.Use(JwtAuthRequired())
	{
		// apiGroup.Get("/me", func(ctx *fiber.Ctx) error {
		// 	token := ctx.Get("Authorization")
		// 	return ctx.JSON(JwtMe(token))
		// })

		apiGroup.Get("/user/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.UserController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/user", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.UserController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/user", func(ctx *fiber.Ctx) error {
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/user", func(ctx *fiber.Ctx) error {
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/user", func(ctx *fiber.Ctx) error {
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/daytype/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.DayTypeController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/daytype", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.DayTypeController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/daytype", func(ctx *fiber.Ctx) error {
			item_ := &models.DayType{}
			ctx.BodyParser(item_)
			var controller rest.DayTypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/daytype", func(ctx *fiber.Ctx) error {
			item_ := &models.DayType{}
			ctx.BodyParser(item_)
			var controller rest.DayTypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/daytype", func(ctx *fiber.Ctx) error {
			item_ := &models.DayType{}
			ctx.BodyParser(item_)
			var controller rest.DayTypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/discount/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.DiscountController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/discount", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.DiscountController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/discount", func(ctx *fiber.Ctx) error {
			item_ := &models.Discount{}
			ctx.BodyParser(item_)
			var controller rest.DiscountController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/discount", func(ctx *fiber.Ctx) error {
			item_ := &models.Discount{}
			ctx.BodyParser(item_)
			var controller rest.DiscountController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/discount", func(ctx *fiber.Ctx) error {
			item_ := &models.Discount{}
			ctx.BodyParser(item_)
			var controller rest.DiscountController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/gym/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.GymController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/gym", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.GymController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/gym", func(ctx *fiber.Ctx) error {
			item_ := &models.Gym{}
			ctx.BodyParser(item_)
			var controller rest.GymController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/gym", func(ctx *fiber.Ctx) error {
			item_ := &models.Gym{}
			ctx.BodyParser(item_)
			var controller rest.GymController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/gym", func(ctx *fiber.Ctx) error {
			item_ := &models.Gym{}
			ctx.BodyParser(item_)
			var controller rest.GymController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/helthcategory/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.HelthCategoryController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/helthcategory", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.HelthCategoryController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/helthcategory", func(ctx *fiber.Ctx) error {
			item_ := &models.HelthCategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthCategoryController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/helthcategory", func(ctx *fiber.Ctx) error {
			item_ := &models.HelthCategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthCategoryController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/helthcategory", func(ctx *fiber.Ctx) error {
			item_ := &models.HelthCategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthCategoryController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/helth/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.HelthController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/helth", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.HelthController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/helth", func(ctx *fiber.Ctx) error {
			item_ := &models.Helth{}
			ctx.BodyParser(item_)
			var controller rest.HelthController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/helth", func(ctx *fiber.Ctx) error {
			item_ := &models.Helth{}
			ctx.BodyParser(item_)
			var controller rest.HelthController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/helth", func(ctx *fiber.Ctx) error {
			item_ := &models.Helth{}
			ctx.BodyParser(item_)
			var controller rest.HelthController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})


		apiGroup.Get("/membership/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.MembershipController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/membership", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.MembershipController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/membership", func(ctx *fiber.Ctx) error {
			item_ := &models.Membership{}
			ctx.BodyParser(item_)
			var controller rest.MembershipController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/membership", func(ctx *fiber.Ctx) error {
			item_ := &models.Membership{}
			ctx.BodyParser(item_)
			var controller rest.MembershipController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/membership", func(ctx *fiber.Ctx) error {
			item_ := &models.Membership{}
			ctx.BodyParser(item_)
			var controller rest.MembershipController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})
	}
}