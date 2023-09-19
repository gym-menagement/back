package rest

import (
	"gym/controllers"
	"gym/models"
)

type UseHelthController struct {
	controllers.Controller
}

func (c *UseHelthController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewUseHelthManager(conn)

    var args []interface{}

    order := c.Query("order")
    if order != "" {
        args = append(args, models.Where{Column:"order", Value:order, Compare:"="})
    }

	helth := c.Query("helth")
    if helth != "" {
        args = append(args, models.Where{Column:"helth", Value:helth, Compare:"="})
    }

    user := c.Query("user")
    if user != "" {
        args = append(args, models.Where{Column:"user", Value:user, Compare:"="})
    }

	rocker := c.Query("rocker")
    if rocker != "" {
        args = append(args, models.Where{Column:"rocker", Value:rocker, Compare:"="})
    }

	term := c.Query("term")
    if term != "" {
        args = append(args, models.Where{Column:"term", Value:term, Compare:"="})
    }

	discount := c.Query("discount")
    if discount != "" {
        args = append(args, models.Where{Column:"discount", Value:discount, Compare:"="})
    }

    startday := c.Query("startday")
    if startday != "" {
        args = append(args, models.Where{Column:"startday", Value:startday, Compare:"="})
    }

    endday := c.Query("endday")
    if endday != "" {
        args = append(args, models.Where{Column:"endday", Value:endday, Compare:"="})
    }
	
    startdate := c.Query("startdate")
    enddate := c.Query("enddate")
    if startdate != "" && enddate != "" {
        var v [2]string
        v[0] = startdate
        v[1] = enddate
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})
    } else if  startdate != "" {
        args = append(args, models.Where{Column:"date", Value:startdate, Compare:">="})
    } else if  enddate != "" {
        args = append(args, models.Where{Column:"date", Value:enddate, Compare:"<="})
    }
    
    if page != 0 && pagesize != 0 {
        args = append(args, models.Paging(page, pagesize))
    }

    orderby := c.Query("orderby")
    if orderby == "desc" {
        // if page != 0 && pagesize != 0 {
            orderby = "id desc"
        // }
    } else {
		orderby = ""
	}

    if orderby != "" {
        args = append(args, models.Ordering(orderby))
    }

	items := manager.Find(args)
	c.Set("items", items)

    total := manager.Count(args)
	c.Set("total", total)
}

func (c *UseHelthController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewUseHelthManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *UseHelthController) Insert(item *models.UseHelth) {
	conn := c.NewConnection()

	manager := models.NewUseHelthManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *UseHelthController) Update(item *models.UseHelth) {
	conn := c.NewConnection()

	manager := models.NewUseHelthManager(conn)
	manager.Update(item)
}

func (c *UseHelthController) Delete(item *models.UseHelth) {
	conn := c.NewConnection()

	manager := models.NewUseHelthManager(conn)
	manager.Delete(item.Id)
}