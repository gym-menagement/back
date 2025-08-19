package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type LoginlogController struct {
	controllers.Controller
}

func (c *LoginlogController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *LoginlogController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)

    var args []interface{}
    
    _ip := c.Get("ip")
    if _ip != "" {
        args = append(args, models.Where{Column:"ip", Value:_ip, Compare:"like"})
    }
    _ipvalue := c.Geti64("ipvalue")
    if _ipvalue != 0 {
        args = append(args, models.Where{Column:"ipvalue", Value:_ipvalue, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
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
                    str += ", ll_" + strings.Trim(v, " ")                
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

func (c *LoginlogController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)

    var args []interface{}
    
    _ip := c.Get("ip")
    if _ip != "" {
        args = append(args, models.Where{Column:"ip", Value:_ip, Compare:"like"})
        
    }
    _ipvalue := c.Geti64("ipvalue")
    if _ipvalue != 0 {
        args = append(args, models.Where{Column:"ipvalue", Value:_ipvalue, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
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

func (c *LoginlogController) Insert(item *models.Loginlog) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewLoginlogManager(conn)
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

func (c *LoginlogController) Insertbatch(item *[]models.Loginlog) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewLoginlogManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *LoginlogController) Update(item *models.Loginlog) {
    
    
	conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *LoginlogController) Delete(item *models.Loginlog) {
    
    
    conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *LoginlogController) Deletebatch(item *[]models.Loginlog) {
    
    
    conn := c.NewConnection()

	manager := models.NewLoginlogManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


