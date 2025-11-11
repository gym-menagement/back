package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type HealthController struct {
	controllers.Controller
}

func (c *HealthController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewHealthManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *HealthController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewHealthManager(conn)

    var args []interface{}
    
    _category := c.Geti64("category")
    if _category != 0 {
        args = append(args, models.Where{Column:"category", Value:_category, Compare:"="})    
    }
    _term := c.Geti64("term")
    if _term != 0 {
        args = append(args, models.Where{Column:"term", Value:_term, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"like"})
        
    }
    _count := c.Geti("count")
    if _count != 0 {
        args = append(args, models.Where{Column:"count", Value:_count, Compare:"="})    
    }
    _cost := c.Geti("cost")
    if _cost != 0 {
        args = append(args, models.Where{Column:"cost", Value:_cost, Compare:"="})    
    }
    _discount := c.Geti64("discount")
    if _discount != 0 {
        args = append(args, models.Where{Column:"discount", Value:_discount, Compare:"="})    
    }
    _costdiscount := c.Geti("costdiscount")
    if _costdiscount != 0 {
        args = append(args, models.Where{Column:"costdiscount", Value:_costdiscount, Compare:"="})    
    }
    _content := c.Get("content")
    if _content != "" {
        args = append(args, models.Where{Column:"content", Value:_content, Compare:"like"})
        
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
                    str += ", h_" + strings.Trim(v, " ")                
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

func (c *HealthController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewHealthManager(conn)

    var args []interface{}
    
    _category := c.Geti64("category")
    if _category != 0 {
        args = append(args, models.Where{Column:"category", Value:_category, Compare:"="})    
    }
    _term := c.Geti64("term")
    if _term != 0 {
        args = append(args, models.Where{Column:"term", Value:_term, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"like"})
        
        
    }
    _count := c.Geti("count")
    if _count != 0 {
        args = append(args, models.Where{Column:"count", Value:_count, Compare:"="})    
    }
    _cost := c.Geti("cost")
    if _cost != 0 {
        args = append(args, models.Where{Column:"cost", Value:_cost, Compare:"="})    
    }
    _discount := c.Geti64("discount")
    if _discount != 0 {
        args = append(args, models.Where{Column:"discount", Value:_discount, Compare:"="})    
    }
    _costdiscount := c.Geti("costdiscount")
    if _costdiscount != 0 {
        args = append(args, models.Where{Column:"costdiscount", Value:_costdiscount, Compare:"="})    
    }
    _content := c.Get("content")
    if _content != "" {
        args = append(args, models.Where{Column:"content", Value:_content, Compare:"like"})
        
        
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

func (c *HealthController) Insert(item *models.Health) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewHealthManager(conn)
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

func (c *HealthController) Insertbatch(item *[]models.Health) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewHealthManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *HealthController) Update(item *models.Health) {
    
    
	conn := c.NewConnection()

	manager := models.NewHealthManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *HealthController) Delete(item *models.Health) {
    
    
    conn := c.NewConnection()

	manager := models.NewHealthManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *HealthController) Deletebatch(item *[]models.Health) {
    
    
    conn := c.NewConnection()

	manager := models.NewHealthManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


