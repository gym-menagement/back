package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type NotificationhistoryController struct {
	controllers.Controller
}

func (c *NotificationhistoryController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *NotificationhistoryController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)

    var args []interface{}
    
    _sender := c.Geti64("sender")
    if _sender != 0 {
        args = append(args, models.Where{Column:"sender", Value:_sender, Compare:"="})    
    }
    _receiver := c.Geti64("receiver")
    if _receiver != 0 {
        args = append(args, models.Where{Column:"receiver", Value:_receiver, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _title := c.Get("title")
    if _title != "" {
        args = append(args, models.Where{Column:"title", Value:_title, Compare:"like"})
        
    }
    _body := c.Get("body")
    if _body != "" {
        args = append(args, models.Where{Column:"body", Value:_body, Compare:"like"})
    }
    _data := c.Get("data")
    if _data != "" {
        args = append(args, models.Where{Column:"data", Value:_data, Compare:"like"})
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _errormessage := c.Get("errormessage")
    if _errormessage != "" {
        args = append(args, models.Where{Column:"errormessage", Value:_errormessage, Compare:"like"})
    }
    _startsentdate := c.Get("startsentdate")
    _endsentdate := c.Get("endsentdate")
    if _startsentdate != "" && _endsentdate != "" {        
        var v [2]string
        v[0] = _startsentdate
        v[1] = _endsentdate  
        args = append(args, models.Where{Column:"sentdate", Value:v, Compare:"between"})    
    } else if  _startsentdate != "" {          
        args = append(args, models.Where{Column:"sentdate", Value:_startsentdate, Compare:">="})
    } else if  _endsentdate != "" {          
        args = append(args, models.Where{Column:"sentdate", Value:_endsentdate, Compare:"<="})            
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
                    str += ", nh_" + strings.Trim(v, " ")                
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

func (c *NotificationhistoryController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)

    var args []interface{}
    
    _sender := c.Geti64("sender")
    if _sender != 0 {
        args = append(args, models.Where{Column:"sender", Value:_sender, Compare:"="})    
    }
    _receiver := c.Geti64("receiver")
    if _receiver != 0 {
        args = append(args, models.Where{Column:"receiver", Value:_receiver, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _title := c.Get("title")
    if _title != "" {
        args = append(args, models.Where{Column:"title", Value:_title, Compare:"like"})
        
        
    }
    _body := c.Get("body")
    if _body != "" {
        args = append(args, models.Where{Column:"body", Value:_body, Compare:"like"})
        
    }
    _data := c.Get("data")
    if _data != "" {
        args = append(args, models.Where{Column:"data", Value:_data, Compare:"like"})
        
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _errormessage := c.Get("errormessage")
    if _errormessage != "" {
        args = append(args, models.Where{Column:"errormessage", Value:_errormessage, Compare:"like"})
        
    }
    _startsentdate := c.Get("startsentdate")
    _endsentdate := c.Get("endsentdate")

    if _startsentdate != "" && _endsentdate != "" {        
        var v [2]string
        v[0] = _startsentdate
        v[1] = _endsentdate  
        args = append(args, models.Where{Column:"sentdate", Value:v, Compare:"between"})    
    } else if  _startsentdate != "" {          
        args = append(args, models.Where{Column:"sentdate", Value:_startsentdate, Compare:">="})
    } else if  _endsentdate != "" {          
        args = append(args, models.Where{Column:"sentdate", Value:_endsentdate, Compare:"<="})            
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

func (c *NotificationhistoryController) Insert(item *models.Notificationhistory) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewNotificationhistoryManager(conn)
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

func (c *NotificationhistoryController) Insertbatch(item *[]models.Notificationhistory) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewNotificationhistoryManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *NotificationhistoryController) Update(item *models.Notificationhistory) {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *NotificationhistoryController) Delete(item *models.Notificationhistory) {
    
    
    conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *NotificationhistoryController) Deletebatch(item *[]models.Notificationhistory) {
    
    
    conn := c.NewConnection()

	manager := models.NewNotificationhistoryManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


