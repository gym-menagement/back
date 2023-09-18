package rest

import (
	"gym/controllers"
	"gym/models"
)

type UserController struct {
	controllers.Controller
}

func (c *UserController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    var args []interface{}

    gym := c.Query("gym")
    if gym != "" {
        args = append(args, models.Where{Column:"gym", Value:gym, Compare:"="})
    }    

    loginid := c.Query("loginid")
    if loginid != "" {
        args = append(args, models.Where{Column:"loginid", Value:loginid, Compare:"like"})
    }

    passwd := c.Query("passwd")
    if passwd != "" {
        args = append(args, models.Where{Column:"passwd", Value:passwd, Compare:"like"})
    }

    name := c.Query("name")
    if name != "" {
        args = append(args, models.Where{Column:"name", Value:name, Compare:"="})
    }

    role := c.Query("role")
    if role != "" {
        args = append(args, models.Where{Column:"role", Value:role, Compare:"="})
    }

    image := c.Query("image")
    if image != "" {
        args = append(args, models.Where{Column:"image", Value:image, Compare:"="})
    }

    sex := c.Query("sex")
    if sex != "" {
        args = append(args, models.Where{Column:"sex", Value:sex, Compare:"="})
    }

    birth := c.Query("birth")
    if birth != "" {
        args = append(args, models.Where{Column:"birth", Value:birth, Compare:"="})
    }

    phonenum := c.Query("phonenum")
    if phonenum != "" {
        args = append(args, models.Where{Column:"phonenum", Value:phonenum, Compare:"="})
    }

    address := c.Query("address")
    if address != "" {
        args = append(args, models.Where{Column:"address", Value:address, Compare:"="})
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

func (c *UserController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *UserController) Insert(item *models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *UserController) Update(item *models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Update(item)
}

func (c *UserController) Delete(item *models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Delete(item.Id)
}