package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type StopController struct {
	controllers.Controller
}

func (c *StopController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *StopController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)

    var args []interface{}
    
    _usehealth := c.Geti64("usehealth")
    if _usehealth != 0 {
        args = append(args, models.Where{Column:"usehealth", Value:_usehealth, Compare:"="})    
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
    _count := c.Geti("count")
    if _count != 0 {
        args = append(args, models.Where{Column:"count", Value:_count, Compare:"="})    
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
                    str += ", s_" + strings.Trim(v, " ")                
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

func (c *StopController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)

    var args []interface{}
    
    _usehealth := c.Geti64("usehealth")
    if _usehealth != 0 {
        args = append(args, models.Where{Column:"usehealth", Value:_usehealth, Compare:"="})    
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
    _count := c.Geti("count")
    if _count != 0 {
        args = append(args, models.Where{Column:"count", Value:_count, Compare:"="})    
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

func (c *StopController) Insert(item *models.Stop) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewStopManager(conn)
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

func (c *StopController) Insertbatch(item *[]models.Stop) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewStopManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *StopController) Update(item *models.Stop) {
    
    
	conn := c.NewConnection()

	manager := models.NewStopManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *StopController) Delete(item *models.Stop) {
    
    
    conn := c.NewConnection()

	manager := models.NewStopManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *StopController) Deletebatch(item *[]models.Stop) {
    
    
    conn := c.NewConnection()

	manager := models.NewStopManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


