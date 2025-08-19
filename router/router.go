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
		return ctx.JSON(JwtAuth(ctx, loginid, passwd))
	})
	apiGroup := app.Group("/api")
	apiGroup.Use(JwtAuthRequired)
	{
		
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
			item_ := &models.UserUpdate{}
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
			item_ := &models.UserUpdate{}
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
			var controller rest.DaytypeController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/daytype", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.DaytypeController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/daytype", func(ctx *fiber.Ctx) error {
			item_ := &models.Daytype{}
			ctx.BodyParser(item_)
			var controller rest.DaytypeController
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
			item_ := &models.Daytype{}
			ctx.BodyParser(item_)
			var controller rest.DaytypeController
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
			item_ := &models.Daytype{}
			ctx.BodyParser(item_)
			var controller rest.DaytypeController
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
			var controller rest.HelthcategoryController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/helthcategory", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.HelthcategoryController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/helthcategory", func(ctx *fiber.Ctx) error {
			item_ := &models.Helthcategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthcategoryController
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
			item_ := &models.Helthcategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthcategoryController
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
			item_ := &models.Helthcategory{}
			ctx.BodyParser(item_)
			var controller rest.HelthcategoryController
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

		apiGroup.Get("/order/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.OrderController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/order", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.OrderController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/order", func(ctx *fiber.Ctx) error {
			item_ := &models.Order{}
			ctx.BodyParser(item_)
			var controller rest.OrderController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/order", func(ctx *fiber.Ctx) error {
			item_ := &models.Order{}
			ctx.BodyParser(item_)
			var controller rest.OrderController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/order", func(ctx *fiber.Ctx) error {
			item_ := &models.Order{}
			ctx.BodyParser(item_)
			var controller rest.OrderController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/paymenttype/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.PaymenttypeController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/paymenttype", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.PaymenttypeController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/paymenttype", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymenttype{}
			ctx.BodyParser(item_)
			var controller rest.PaymenttypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/paymenttype", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymenttype{}
			ctx.BodyParser(item_)
			var controller rest.PaymenttypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/paymenttype", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymenttype{}
			ctx.BodyParser(item_)
			var controller rest.PaymenttypeController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/paymentform/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.PaymentformController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/paymentform", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.PaymentformController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/paymentform", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymentform{}
			ctx.BodyParser(item_)
			var controller rest.PaymentformController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/paymentform", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymentform{}
			ctx.BodyParser(item_)
			var controller rest.PaymentformController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/paymentform", func(ctx *fiber.Ctx) error {
			item_ := &models.Paymentform{}
			ctx.BodyParser(item_)
			var controller rest.PaymentformController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/payment/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.PaymentController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/payment", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.PaymentController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/payment", func(ctx *fiber.Ctx) error {
			item_ := &models.Payment{}
			ctx.BodyParser(item_)
			var controller rest.PaymentController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/payment", func(ctx *fiber.Ctx) error {
			item_ := &models.Payment{}
			ctx.BodyParser(item_)
			var controller rest.PaymentController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/payment", func(ctx *fiber.Ctx) error {
			item_ := &models.Payment{}
			ctx.BodyParser(item_)
			var controller rest.PaymentController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/rockergroup/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.RockergroupController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/rockergroup", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.RockergroupController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/rockergroup", func(ctx *fiber.Ctx) error {
			item_ := &models.Rockergroup{}
			ctx.BodyParser(item_)
			var controller rest.RockergroupController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/rockergroup", func(ctx *fiber.Ctx) error {
			item_ := &models.Rockergroup{}
			ctx.BodyParser(item_)
			var controller rest.RockergroupController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/rockergroup", func(ctx *fiber.Ctx) error {
			item_ := &models.Rockergroup{}
			ctx.BodyParser(item_)
			var controller rest.RockergroupController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/rocker/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.RockerController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/rocker", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.RockerController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/rocker", func(ctx *fiber.Ctx) error {
			item_ := &models.Rocker{}
			ctx.BodyParser(item_)
			var controller rest.RockerController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/rocker", func(ctx *fiber.Ctx) error {
			item_ := &models.Rocker{}
			ctx.BodyParser(item_)
			var controller rest.RockerController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/rocker", func(ctx *fiber.Ctx) error {
			item_ := &models.Rocker{}
			ctx.BodyParser(item_)
			var controller rest.RockerController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/role/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.RoleController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/role", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.RoleController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/role", func(ctx *fiber.Ctx) error {
			item_ := &models.Role{}
			ctx.BodyParser(item_)
			var controller rest.RoleController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/role", func(ctx *fiber.Ctx) error {
			item_ := &models.Role{}
			ctx.BodyParser(item_)
			var controller rest.RoleController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/role", func(ctx *fiber.Ctx) error {
			item_ := &models.Role{}
			ctx.BodyParser(item_)
			var controller rest.RoleController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/stop/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.StopController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/stop", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.StopController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/stop", func(ctx *fiber.Ctx) error {
			item_ := &models.Stop{}
			ctx.BodyParser(item_)
			var controller rest.StopController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/stop", func(ctx *fiber.Ctx) error {
			item_ := &models.Stop{}
			ctx.BodyParser(item_)
			var controller rest.StopController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/stop", func(ctx *fiber.Ctx) error {
			item_ := &models.Stop{}
			ctx.BodyParser(item_)
			var controller rest.StopController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/term/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.TermController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/term", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.TermController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/term", func(ctx *fiber.Ctx) error {
			item_ := &models.Term{}
			ctx.BodyParser(item_)
			var controller rest.TermController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/term", func(ctx *fiber.Ctx) error {
			item_ := &models.Term{}
			ctx.BodyParser(item_)
			var controller rest.TermController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/term", func(ctx *fiber.Ctx) error {
			item_ := &models.Term{}
			ctx.BodyParser(item_)
			var controller rest.TermController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/usehelth/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.UsehelthController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/usehelth", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.UsehelthController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/usehelth", func(ctx *fiber.Ctx) error {
			item_ := &models.Usehelth{}
			ctx.BodyParser(item_)
			var controller rest.UsehelthController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/usehelth", func(ctx *fiber.Ctx) error {
			item_ := &models.Usehelth{}
			ctx.BodyParser(item_)
			var controller rest.UsehelthController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/usehelth", func(ctx *fiber.Ctx) error {
			item_ := &models.Usehelth{}
			ctx.BodyParser(item_)
			var controller rest.UsehelthController
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