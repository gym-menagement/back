package rest

import (
	"gym/controllers"
	"gym/models"
)

type HelthController struct {
	controllers.Controller
}

func (c *HelthController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewHelthManager(conn)

    var args []interface{}

	category := c.Query("category")
    if category != "" {
        args = append(args, models.Where{Column:"category", Value:category, Compare:"="})
    }

	term := c.Query("term")
    if term != "" {
        args = append(args, models.Where{Column:"term", Value:term, Compare:"="})
    }

    name := c.Query("name")
    if name != "" {
        args = append(args, models.Where{Column:"name", Value:name, Compare:"="})
    }

	count := c.Query("count")
    if count != "" {
        args = append(args, models.Where{Column:"count", Value:count, Compare:"="})
    }

	cost := c.Query("cost")
    if cost != "" {
        args = append(args, models.Where{Column:"cost", Value:cost, Compare:"="})
    }

	discount := c.Query("discount")
    if discount != "" {
        args = append(args, models.Where{Column:"discount", Value:discount, Compare:"="})
    }

	costdiscount := c.Query("costdiscount")
    if costdiscount != "" {
        args = append(args, models.Where{Column:"costdiscount", Value:costdiscount, Compare:"="})
    }

	content := c.Query("content")
    if content != "" {
        args = append(args, models.Where{Column:"content", Value:content, Compare:"="})
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

func (c *HelthController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewHelthManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *HelthController) Insert(item *models.Helth) {
	conn := c.NewConnection()

	manager := models.NewHelthManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *HelthController) Update(item *models.Helth) {
	conn := c.NewConnection()

	manager := models.NewHelthManager(conn)
	manager.Update(item)
}

func (c *HelthController) Delete(item *models.Helth) {
	conn := c.NewConnection()

	manager := models.NewHelthManager(conn)
	manager.Delete(item.Id)
}