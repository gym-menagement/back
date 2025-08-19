package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type SettingController struct {
	controllers.Controller
}

func (c *SettingController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewSettingManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *SettingController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewSettingManager(conn)

    var args []interface{}
    
    _category := c.Get("category")
    if _category != "" {
        args = append(args, models.Where{Column:"category", Value:_category, Compare:"like"})
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _key := c.Get("key")
    if _key != "" {
        args = append(args, models.Where{Column:"key", Value:_key, Compare:"like"})
    }
    _value := c.Get("value")
    if _value != "" {
        args = append(args, models.Where{Column:"value", Value:_value, Compare:"like"})
    }
    _remark := c.Get("remark")
    if _remark != "" {
        args = append(args, models.Where{Column:"remark", Value:_remark, Compare:"like"})
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _data := c.Get("data")
    if _data != "" {
        args = append(args, models.Where{Column:"data", Value:_data, Compare:"like"})
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
                    str += ", se_" + strings.Trim(v, " ")                
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

func (c *SettingController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewSettingManager(conn)

    var args []interface{}
    
    _category := c.Get("category")
    if _category != "" {
        args = append(args, models.Where{Column:"category", Value:_category, Compare:"like"})
        
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
        
    }
    _key := c.Get("key")
    if _key != "" {
        args = append(args, models.Where{Column:"key", Value:_key, Compare:"like"})
        
    }
    _value := c.Get("value")
    if _value != "" {
        args = append(args, models.Where{Column:"value", Value:_value, Compare:"like"})
        
    }
    _remark := c.Get("remark")
    if _remark != "" {
        args = append(args, models.Where{Column:"remark", Value:_remark, Compare:"like"})
        
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _data := c.Get("data")
    if _data != "" {
        args = append(args, models.Where{Column:"data", Value:_data, Compare:"like"})
        
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

func (c *SettingController) Insert(item *models.Setting) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewSettingManager(conn)
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

func (c *SettingController) Insertbatch(item *[]models.Setting) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewSettingManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *SettingController) Update(item *models.Setting) {
    
    
	conn := c.NewConnection()

	manager := models.NewSettingManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *SettingController) Delete(item *models.Setting) {
    
    
    conn := c.NewConnection()

	manager := models.NewSettingManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *SettingController) Deletebatch(item *[]models.Setting) {
    
    
    conn := c.NewConnection()

	manager := models.NewSettingManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


