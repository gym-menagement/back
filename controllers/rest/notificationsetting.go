package rest


import (
	"gym/controllers"
	
	"gym/models"

    "strings"
)

type NotificationsettingController struct {
	controllers.Controller
}

func (c *NotificationsettingController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *NotificationsettingController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _enabled := c.Geti("enabled")
    if _enabled != 0 {
        args = append(args, models.Where{Column:"enabled", Value:_enabled, Compare:"="})    
    }
    _membershipexpiry := c.Geti("membershipexpiry")
    if _membershipexpiry != 0 {
        args = append(args, models.Where{Column:"membershipexpiry", Value:_membershipexpiry, Compare:"="})    
    }
    _membershipnear := c.Geti("membershipnear")
    if _membershipnear != 0 {
        args = append(args, models.Where{Column:"membershipnear", Value:_membershipnear, Compare:"="})    
    }
    _attendanceenc := c.Geti("attendanceenc")
    if _attendanceenc != 0 {
        args = append(args, models.Where{Column:"attendanceenc", Value:_attendanceenc, Compare:"="})    
    }
    _gymannounce := c.Geti("gymannounce")
    if _gymannounce != 0 {
        args = append(args, models.Where{Column:"gymannounce", Value:_gymannounce, Compare:"="})    
    }
    _systemnotice := c.Geti("systemnotice")
    if _systemnotice != 0 {
        args = append(args, models.Where{Column:"systemnotice", Value:_systemnotice, Compare:"="})    
    }
    _paymentconfirm := c.Geti("paymentconfirm")
    if _paymentconfirm != 0 {
        args = append(args, models.Where{Column:"paymentconfirm", Value:_paymentconfirm, Compare:"="})    
    }
    _pauseexpiry := c.Geti("pauseexpiry")
    if _pauseexpiry != 0 {
        args = append(args, models.Where{Column:"pauseexpiry", Value:_pauseexpiry, Compare:"="})    
    }
    _weeklygoal := c.Geti("weeklygoal")
    if _weeklygoal != 0 {
        args = append(args, models.Where{Column:"weeklygoal", Value:_weeklygoal, Compare:"="})    
    }
    _personalrecord := c.Geti("personalrecord")
    if _personalrecord != 0 {
        args = append(args, models.Where{Column:"personalrecord", Value:_personalrecord, Compare:"="})    
    }
    _quietenabled := c.Geti("quietenabled")
    if _quietenabled != 0 {
        args = append(args, models.Where{Column:"quietenabled", Value:_quietenabled, Compare:"="})    
    }
    _quietstart := c.Get("quietstart")
    if _quietstart != "" {
        args = append(args, models.Where{Column:"quietstart", Value:_quietstart, Compare:"like"})
    }
    _quietend := c.Get("quietend")
    if _quietend != "" {
        args = append(args, models.Where{Column:"quietend", Value:_quietend, Compare:"like"})
    }
    _startcreateddate := c.Get("startcreateddate")
    _endcreateddate := c.Get("endcreateddate")
    if _startcreateddate != "" && _endcreateddate != "" {        
        var v [2]string
        v[0] = _startcreateddate
        v[1] = _endcreateddate  
        args = append(args, models.Where{Column:"createddate", Value:v, Compare:"between"})    
    } else if  _startcreateddate != "" {          
        args = append(args, models.Where{Column:"createddate", Value:_startcreateddate, Compare:">="})
    } else if  _endcreateddate != "" {          
        args = append(args, models.Where{Column:"createddate", Value:_endcreateddate, Compare:"<="})            
    }
    _startupdateddate := c.Get("startupdateddate")
    _endupdateddate := c.Get("endupdateddate")
    if _startupdateddate != "" && _endupdateddate != "" {        
        var v [2]string
        v[0] = _startupdateddate
        v[1] = _endupdateddate  
        args = append(args, models.Where{Column:"updateddate", Value:v, Compare:"between"})    
    } else if  _startupdateddate != "" {          
        args = append(args, models.Where{Column:"updateddate", Value:_startupdateddate, Compare:">="})
    } else if  _endupdateddate != "" {          
        args = append(args, models.Where{Column:"updateddate", Value:_endupdateddate, Compare:"<="})            
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
                    str += ", ns_" + strings.Trim(v, " ")                
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

func (c *NotificationsettingController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _enabled := c.Geti("enabled")
    if _enabled != 0 {
        args = append(args, models.Where{Column:"enabled", Value:_enabled, Compare:"="})    
    }
    _membershipexpiry := c.Geti("membershipexpiry")
    if _membershipexpiry != 0 {
        args = append(args, models.Where{Column:"membershipexpiry", Value:_membershipexpiry, Compare:"="})    
    }
    _membershipnear := c.Geti("membershipnear")
    if _membershipnear != 0 {
        args = append(args, models.Where{Column:"membershipnear", Value:_membershipnear, Compare:"="})    
    }
    _attendanceenc := c.Geti("attendanceenc")
    if _attendanceenc != 0 {
        args = append(args, models.Where{Column:"attendanceenc", Value:_attendanceenc, Compare:"="})    
    }
    _gymannounce := c.Geti("gymannounce")
    if _gymannounce != 0 {
        args = append(args, models.Where{Column:"gymannounce", Value:_gymannounce, Compare:"="})    
    }
    _systemnotice := c.Geti("systemnotice")
    if _systemnotice != 0 {
        args = append(args, models.Where{Column:"systemnotice", Value:_systemnotice, Compare:"="})    
    }
    _paymentconfirm := c.Geti("paymentconfirm")
    if _paymentconfirm != 0 {
        args = append(args, models.Where{Column:"paymentconfirm", Value:_paymentconfirm, Compare:"="})    
    }
    _pauseexpiry := c.Geti("pauseexpiry")
    if _pauseexpiry != 0 {
        args = append(args, models.Where{Column:"pauseexpiry", Value:_pauseexpiry, Compare:"="})    
    }
    _weeklygoal := c.Geti("weeklygoal")
    if _weeklygoal != 0 {
        args = append(args, models.Where{Column:"weeklygoal", Value:_weeklygoal, Compare:"="})    
    }
    _personalrecord := c.Geti("personalrecord")
    if _personalrecord != 0 {
        args = append(args, models.Where{Column:"personalrecord", Value:_personalrecord, Compare:"="})    
    }
    _quietenabled := c.Geti("quietenabled")
    if _quietenabled != 0 {
        args = append(args, models.Where{Column:"quietenabled", Value:_quietenabled, Compare:"="})    
    }
    _quietstart := c.Get("quietstart")
    if _quietstart != "" {
        args = append(args, models.Where{Column:"quietstart", Value:_quietstart, Compare:"like"})
        
    }
    _quietend := c.Get("quietend")
    if _quietend != "" {
        args = append(args, models.Where{Column:"quietend", Value:_quietend, Compare:"like"})
        
    }
    _startcreateddate := c.Get("startcreateddate")
    _endcreateddate := c.Get("endcreateddate")

    if _startcreateddate != "" && _endcreateddate != "" {        
        var v [2]string
        v[0] = _startcreateddate
        v[1] = _endcreateddate  
        args = append(args, models.Where{Column:"createddate", Value:v, Compare:"between"})    
    } else if  _startcreateddate != "" {          
        args = append(args, models.Where{Column:"createddate", Value:_startcreateddate, Compare:">="})
    } else if  _endcreateddate != "" {          
        args = append(args, models.Where{Column:"createddate", Value:_endcreateddate, Compare:"<="})            
    }
    _startupdateddate := c.Get("startupdateddate")
    _endupdateddate := c.Get("endupdateddate")

    if _startupdateddate != "" && _endupdateddate != "" {        
        var v [2]string
        v[0] = _startupdateddate
        v[1] = _endupdateddate  
        args = append(args, models.Where{Column:"updateddate", Value:v, Compare:"between"})    
    } else if  _startupdateddate != "" {          
        args = append(args, models.Where{Column:"updateddate", Value:_startupdateddate, Compare:">="})
    } else if  _endupdateddate != "" {          
        args = append(args, models.Where{Column:"updateddate", Value:_endupdateddate, Compare:"<="})            
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

func (c *NotificationsettingController) Insert(item *models.Notificationsetting) {
    
    
    

	conn := c.NewConnection()
    
	manager := models.NewNotificationsettingManager(conn)
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

func (c *NotificationsettingController) Insertbatch(item *[]models.Notificationsetting) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewNotificationsettingManager(conn)

    for i := 0; i < rows; i++ {
        
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *NotificationsettingController) Update(item *models.Notificationsetting) {
    
    
    

	conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *NotificationsettingController) Delete(item *models.Notificationsetting) {
    
    
    conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *NotificationsettingController) Deletebatch(item *[]models.Notificationsetting) {
    
    
    conn := c.NewConnection()

	manager := models.NewNotificationsettingManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


