package rest

import (
	"gym/controllers"
	"gym/models"
)

type StopController struct {
	controllers.Controller
}

func (c *StopController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)

    var args []interface{}

    use_helth := c.Query("use_helth")
    if use_helth != "" {
        args = append(args, models.Where{Column:"use_helth", Value:use_helth, Compare:"="})
    }    

    startday := c.Query("startday")
    if startday != "" {
        args = append(args, models.Where{Column:"startday", Value:startday, Compare:"="})
    }

    endday := c.Query("endday")
    if endday != "" {
        args = append(args, models.Where{Column:"endday", Value:endday, Compare:"="})
    }

    count := c.Query("count")
    if count != "" {
        args = append(args, models.Where{Column:"count", Value:count, Compare:"="})
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

func (c *StopController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *StopController) Insert(item *models.Stop) {
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *StopController) Update(item *models.Stop) {
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
	manager.Update(item)
}

func (c *StopController) Delete(item *models.Stop) {
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
	manager.Delete(item.Id)
}