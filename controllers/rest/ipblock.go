package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type IpblockController struct {
	controllers.Controller
}

func (c *IpblockController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *IpblockController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)

    var args []interface{}
    
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _policy := c.Geti("policy")
    if _policy != 0 {
        args = append(args, models.Where{Column:"policy", Value:_policy, Compare:"="})    
    }
    _use := c.Geti("use")
    if _use != 0 {
        args = append(args, models.Where{Column:"use", Value:_use, Compare:"="})    
    }
    _order := c.Geti("order")
    if _order != 0 {
        args = append(args, models.Where{Column:"order", Value:_order, Compare:"="})    
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
                    str += ", ib_" + strings.Trim(v, " ")                
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

func (c *IpblockController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)

    var args []interface{}
    
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
        
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _policy := c.Geti("policy")
    if _policy != 0 {
        args = append(args, models.Where{Column:"policy", Value:_policy, Compare:"="})    
    }
    _use := c.Geti("use")
    if _use != 0 {
        args = append(args, models.Where{Column:"use", Value:_use, Compare:"="})    
    }
    _order := c.Geti("order")
    if _order != 0 {
        args = append(args, models.Where{Column:"order", Value:_order, Compare:"="})    
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

func (c *IpblockController) Insert(item *models.Ipblock) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewIpblockManager(conn)
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

func (c *IpblockController) Insertbatch(item *[]models.Ipblock) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewIpblockManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *IpblockController) Update(item *models.Ipblock) {
    
    
	conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *IpblockController) Delete(item *models.Ipblock) {
    
    
    conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *IpblockController) Deletebatch(item *[]models.Ipblock) {
    
    
    conn := c.NewConnection()

	manager := models.NewIpblockManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


