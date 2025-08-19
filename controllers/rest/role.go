package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type RoleController struct {
	controllers.Controller
}

func (c *RoleController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *RoleController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _role := c.Geti("role")
    if _role != 0 {
        args = append(args, models.Where{Column:"role", Value:_role, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _startdate := c.Get("startdate")
    _enddate := c.Get("enddate")
    if _startdate != "" && _enddate != "" {        
        var v [2]string
        v[0] = _startdate
        v[1] = _enddate  
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})    
    } else if  _startdate != "" {          
        args = append(args, models.Where{Column:"date", Value:_startdate, Compare:">="})
    } else if  _enddate != "" {          
        args = append(args, models.Where{Column:"date", Value:_enddate, Compare:"<="})            
    }
    

    
    
    if page != 0 && pagesize != 0 {
        args = append(args, models.Paging(page, pagesize))
    }
    
    orderby := c.Get("orderby")
    if orderby == "" {
        if page != 0 && pagesize != 0 {
            orderby = "id desc"
            args = append(args, models.Ordering(orderby))
        }
    } else {
        orderbys := strings.Split(orderby, ",")

        str := ""
        for i, v := range orderbys {
            if i == 0 {
                str += v
            } else {
                if strings.Contains(v, "_") {                   
                    str += ", " + strings.Trim(v, " ")
                } else {
                    str += ", r_" + strings.Trim(v, " ")                
                }
            }
        }
        
        args = append(args, models.Ordering(str))
    }
    
	items := manager.Find(args)
	c.Set("items", items)

    if page == 1 {
       total := manager.Count(args)
	   c.Set("total", total)
    }
}

func (c *RoleController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _role := c.Geti("role")
    if _role != 0 {
        args = append(args, models.Where{Column:"role", Value:_role, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
        
    }
    _startdate := c.Get("startdate")
    _enddate := c.Get("enddate")

    if _startdate != "" && _enddate != "" {        
        var v [2]string
        v[0] = _startdate
        v[1] = _enddate  
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})    
    } else if  _startdate != "" {          
        args = append(args, models.Where{Column:"date", Value:_startdate, Compare:">="})
    } else if  _enddate != "" {          
        args = append(args, models.Where{Column:"date", Value:_enddate, Compare:"<="})            
    }
    
    
    
    
    total := manager.Count(args)
	c.Set("total", total)
}

func (c *RoleController) Insert(item *models.Role) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewRoleManager(conn)
	err := manager.Insert(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *RoleController) Insertbatch(item *[]models.Role) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewRoleManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *RoleController) Update(item *models.Role) {
    
    
	conn := c.NewConnection()

	manager := models.NewRoleManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *RoleController) Delete(item *models.Role) {
    
    
    conn := c.NewConnection()

	manager := models.NewRoleManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *RoleController) Deletebatch(item *[]models.Role) {
    
    
    conn := c.NewConnection()

	manager := models.NewRoleManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


