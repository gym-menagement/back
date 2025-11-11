package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type AttendanceController struct {
	controllers.Controller
}

func (c *AttendanceController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *AttendanceController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _method := c.Geti("method")
    if _method != 0 {
        args = append(args, models.Where{Column:"method", Value:_method, Compare:"="})    
    }
    _startcheckintime := c.Get("startcheckintime")
    _endcheckintime := c.Get("endcheckintime")
    if _startcheckintime != "" && _endcheckintime != "" {        
        var v [2]string
        v[0] = _startcheckintime
        v[1] = _endcheckintime  
        args = append(args, models.Where{Column:"checkintime", Value:v, Compare:"between"})    
    } else if  _startcheckintime != "" {          
        args = append(args, models.Where{Column:"checkintime", Value:_startcheckintime, Compare:">="})
    } else if  _endcheckintime != "" {          
        args = append(args, models.Where{Column:"checkintime", Value:_endcheckintime, Compare:"<="})            
    }
    _startcheckouttime := c.Get("startcheckouttime")
    _endcheckouttime := c.Get("endcheckouttime")
    if _startcheckouttime != "" && _endcheckouttime != "" {        
        var v [2]string
        v[0] = _startcheckouttime
        v[1] = _endcheckouttime  
        args = append(args, models.Where{Column:"checkouttime", Value:v, Compare:"between"})    
    } else if  _startcheckouttime != "" {          
        args = append(args, models.Where{Column:"checkouttime", Value:_startcheckouttime, Compare:">="})
    } else if  _endcheckouttime != "" {          
        args = append(args, models.Where{Column:"checkouttime", Value:_endcheckouttime, Compare:"<="})            
    }
    _duration := c.Geti("duration")
    if _duration != 0 {
        args = append(args, models.Where{Column:"duration", Value:_duration, Compare:"="})    
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
    }
    _ip := c.Get("ip")
    if _ip != "" {
        args = append(args, models.Where{Column:"ip", Value:_ip, Compare:"like"})
    }
    _device := c.Get("device")
    if _device != "" {
        args = append(args, models.Where{Column:"device", Value:_device, Compare:"like"})
    }
    _createdby := c.Geti64("createdby")
    if _createdby != 0 {
        args = append(args, models.Where{Column:"createdby", Value:_createdby, Compare:"="})    
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
                    str += ", at_" + strings.Trim(v, " ")                
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

func (c *AttendanceController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _method := c.Geti("method")
    if _method != 0 {
        args = append(args, models.Where{Column:"method", Value:_method, Compare:"="})    
    }
    _startcheckintime := c.Get("startcheckintime")
    _endcheckintime := c.Get("endcheckintime")

    if _startcheckintime != "" && _endcheckintime != "" {        
        var v [2]string
        v[0] = _startcheckintime
        v[1] = _endcheckintime  
        args = append(args, models.Where{Column:"checkintime", Value:v, Compare:"between"})    
    } else if  _startcheckintime != "" {          
        args = append(args, models.Where{Column:"checkintime", Value:_startcheckintime, Compare:">="})
    } else if  _endcheckintime != "" {          
        args = append(args, models.Where{Column:"checkintime", Value:_endcheckintime, Compare:"<="})            
    }
    _startcheckouttime := c.Get("startcheckouttime")
    _endcheckouttime := c.Get("endcheckouttime")

    if _startcheckouttime != "" && _endcheckouttime != "" {        
        var v [2]string
        v[0] = _startcheckouttime
        v[1] = _endcheckouttime  
        args = append(args, models.Where{Column:"checkouttime", Value:v, Compare:"between"})    
    } else if  _startcheckouttime != "" {          
        args = append(args, models.Where{Column:"checkouttime", Value:_startcheckouttime, Compare:">="})
    } else if  _endcheckouttime != "" {          
        args = append(args, models.Where{Column:"checkouttime", Value:_endcheckouttime, Compare:"<="})            
    }
    _duration := c.Geti("duration")
    if _duration != 0 {
        args = append(args, models.Where{Column:"duration", Value:_duration, Compare:"="})    
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
        
    }
    _ip := c.Get("ip")
    if _ip != "" {
        args = append(args, models.Where{Column:"ip", Value:_ip, Compare:"like"})
        
    }
    _device := c.Get("device")
    if _device != "" {
        args = append(args, models.Where{Column:"device", Value:_device, Compare:"like"})
        
    }
    _createdby := c.Geti64("createdby")
    if _createdby != 0 {
        args = append(args, models.Where{Column:"createdby", Value:_createdby, Compare:"="})    
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

func (c *AttendanceController) Insert(item *models.Attendance) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewAttendanceManager(conn)
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

func (c *AttendanceController) Insertbatch(item *[]models.Attendance) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewAttendanceManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *AttendanceController) Update(item *models.Attendance) {
    
    
	conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *AttendanceController) Delete(item *models.Attendance) {
    
    
    conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *AttendanceController) Deletebatch(item *[]models.Attendance) {
    
    
    conn := c.NewConnection()

	manager := models.NewAttendanceManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


