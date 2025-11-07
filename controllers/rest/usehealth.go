package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type UsehealthController struct {
	controllers.Controller
}

func (c *UsehealthController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *UsehealthController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)

    var args []interface{}
    
    _order := c.Geti64("order")
    if _order != 0 {
        args = append(args, models.Where{Column:"order", Value:_order, Compare:"="})    
    }
    _health := c.Geti64("health")
    if _health != 0 {
        args = append(args, models.Where{Column:"health", Value:_health, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _rocker := c.Geti64("rocker")
    if _rocker != 0 {
        args = append(args, models.Where{Column:"rocker", Value:_rocker, Compare:"="})    
    }
    _term := c.Geti64("term")
    if _term != 0 {
        args = append(args, models.Where{Column:"term", Value:_term, Compare:"="})    
    }
    _discount := c.Geti64("discount")
    if _discount != 0 {
        args = append(args, models.Where{Column:"discount", Value:_discount, Compare:"="})    
    }
    _startstartday := c.Get("startstartday")
    _endstartday := c.Get("endstartday")
    if _startstartday != "" && _endstartday != "" {        
        var v [2]string
        v[0] = _startstartday
        v[1] = _endstartday  
        args = append(args, models.Where{Column:"startday", Value:v, Compare:"between"})    
    } else if  _startstartday != "" {          
        args = append(args, models.Where{Column:"startday", Value:_startstartday, Compare:">="})
    } else if  _endstartday != "" {          
        args = append(args, models.Where{Column:"startday", Value:_endstartday, Compare:"<="})            
    }
    _startendday := c.Get("startendday")
    _endendday := c.Get("endendday")
    if _startendday != "" && _endendday != "" {        
        var v [2]string
        v[0] = _startendday
        v[1] = _endendday  
        args = append(args, models.Where{Column:"endday", Value:v, Compare:"between"})    
    } else if  _startendday != "" {          
        args = append(args, models.Where{Column:"endday", Value:_startendday, Compare:">="})
    } else if  _endendday != "" {          
        args = append(args, models.Where{Column:"endday", Value:_endendday, Compare:"<="})            
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
                    str += ", uh_" + strings.Trim(v, " ")                
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

func (c *UsehealthController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)

    var args []interface{}
    
    _order := c.Geti64("order")
    if _order != 0 {
        args = append(args, models.Where{Column:"order", Value:_order, Compare:"="})    
    }
    _health := c.Geti64("health")
    if _health != 0 {
        args = append(args, models.Where{Column:"health", Value:_health, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _rocker := c.Geti64("rocker")
    if _rocker != 0 {
        args = append(args, models.Where{Column:"rocker", Value:_rocker, Compare:"="})    
    }
    _term := c.Geti64("term")
    if _term != 0 {
        args = append(args, models.Where{Column:"term", Value:_term, Compare:"="})    
    }
    _discount := c.Geti64("discount")
    if _discount != 0 {
        args = append(args, models.Where{Column:"discount", Value:_discount, Compare:"="})    
    }
    _startstartday := c.Get("startstartday")
    _endstartday := c.Get("endstartday")

    if _startstartday != "" && _endstartday != "" {        
        var v [2]string
        v[0] = _startstartday
        v[1] = _endstartday  
        args = append(args, models.Where{Column:"startday", Value:v, Compare:"between"})    
    } else if  _startstartday != "" {          
        args = append(args, models.Where{Column:"startday", Value:_startstartday, Compare:">="})
    } else if  _endstartday != "" {          
        args = append(args, models.Where{Column:"startday", Value:_endstartday, Compare:"<="})            
    }
    _startendday := c.Get("startendday")
    _endendday := c.Get("endendday")

    if _startendday != "" && _endendday != "" {        
        var v [2]string
        v[0] = _startendday
        v[1] = _endendday  
        args = append(args, models.Where{Column:"endday", Value:v, Compare:"between"})    
    } else if  _startendday != "" {          
        args = append(args, models.Where{Column:"endday", Value:_startendday, Compare:">="})
    } else if  _endendday != "" {          
        args = append(args, models.Where{Column:"endday", Value:_endendday, Compare:"<="})            
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

func (c *UsehealthController) Insert(item *models.Usehealth) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewUsehealthManager(conn)
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

func (c *UsehealthController) Insertbatch(item *[]models.Usehealth) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewUsehealthManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *UsehealthController) Update(item *models.Usehealth) {
    
    
	conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *UsehealthController) Delete(item *models.Usehealth) {
    
    
    conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *UsehealthController) Deletebatch(item *[]models.Usehealth) {
    
    
    conn := c.NewConnection()

	manager := models.NewUsehealthManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


