package rest

import (
	"gym/controllers"
	"gym/models"
)

type RoleController struct {
	controllers.Controller
}

func (c *RoleController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)

    var args []interface{}

    gym := c.Query("gym")
    if gym != "" {
        args = append(args, models.Where{Column:"gym", Value:gym, Compare:"="})
    }

	role := c.Query("role")
    if role != "" {
        args = append(args, models.Where{Column:"role", Value:role, Compare:"="})
    }    

    name := c.Query("name")
    if name != "" {
        args = append(args, models.Where{Column:"name", Value:name, Compare:"="})
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

func (c *RoleController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *RoleController) Insert(item *models.Role) {
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *RoleController) Update(item *models.Role) {
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
	manager.Update(item)
}

func (c *RoleController) Delete(item *models.Role) {
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
	manager.Delete(item.Id)
}