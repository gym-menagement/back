package rest

import (
	"gym/controllers"
	"gym/models"
)

type TermController struct {
	controllers.Controller
}

func (c *TermController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewTermManager(conn)

    var args []interface{}

    gym := c.Query("gym")
    if gym != "" {
        args = append(args, models.Where{Column:"gym", Value:gym, Compare:"="})
    }

	day_type := c.Query("day_type")
    if day_type != "" {
        args = append(args, models.Where{Column:"day_type", Value:day_type, Compare:"="})
    }

    name := c.Query("name")
    if name != "" {
        args = append(args, models.Where{Column:"name", Value:name, Compare:"="})
    }

	term := c.Query("term")
    if term != "" {
        args = append(args, models.Where{Column:"term", Value:term, Compare:"="})
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

func (c *TermController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewTermManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *TermController) Insert(item *models.Term) {
	conn := c.NewConnection()

	manager := models.NewTermManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *TermController) Update(item *models.Term) {
	conn := c.NewConnection()

	manager := models.NewTermManager(conn)
	manager.Update(item)
}

func (c *TermController) Delete(item *models.Term) {
	conn := c.NewConnection()

	manager := models.NewTermManager(conn)
	manager.Delete(item.Id)
}