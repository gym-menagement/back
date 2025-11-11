package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type PtreservationController struct {
	controllers.Controller
}

func (c *PtreservationController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *PtreservationController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)

    var args []interface{}
    
    _trainer := c.Geti64("trainer")
    if _trainer != 0 {
        args = append(args, models.Where{Column:"trainer", Value:_trainer, Compare:"="})    
    }
    _member := c.Geti64("member")
    if _member != 0 {
        args = append(args, models.Where{Column:"member", Value:_member, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _startreservationdate := c.Get("startreservationdate")
    _endreservationdate := c.Get("endreservationdate")
    if _startreservationdate != "" && _endreservationdate != "" {        
        var v [2]string
        v[0] = _startreservationdate
        v[1] = _endreservationdate  
        args = append(args, models.Where{Column:"reservationdate", Value:v, Compare:"between"})    
    } else if  _startreservationdate != "" {          
        args = append(args, models.Where{Column:"reservationdate", Value:_startreservationdate, Compare:">="})
    } else if  _endreservationdate != "" {          
        args = append(args, models.Where{Column:"reservationdate", Value:_endreservationdate, Compare:"<="})            
    }
    _starttime := c.Get("starttime")
    if _starttime != "" {
        args = append(args, models.Where{Column:"starttime", Value:_starttime, Compare:"like"})
    }
    _endtime := c.Get("endtime")
    if _endtime != "" {
        args = append(args, models.Where{Column:"endtime", Value:_endtime, Compare:"like"})
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
    _cancelreason := c.Get("cancelreason")
    if _cancelreason != "" {
        args = append(args, models.Where{Column:"cancelreason", Value:_cancelreason, Compare:"like"})
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
                    str += ", pr_" + strings.Trim(v, " ")                
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

func (c *PtreservationController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)

    var args []interface{}
    
    _trainer := c.Geti64("trainer")
    if _trainer != 0 {
        args = append(args, models.Where{Column:"trainer", Value:_trainer, Compare:"="})    
    }
    _member := c.Geti64("member")
    if _member != 0 {
        args = append(args, models.Where{Column:"member", Value:_member, Compare:"="})    
    }
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _startreservationdate := c.Get("startreservationdate")
    _endreservationdate := c.Get("endreservationdate")

    if _startreservationdate != "" && _endreservationdate != "" {        
        var v [2]string
        v[0] = _startreservationdate
        v[1] = _endreservationdate  
        args = append(args, models.Where{Column:"reservationdate", Value:v, Compare:"between"})    
    } else if  _startreservationdate != "" {          
        args = append(args, models.Where{Column:"reservationdate", Value:_startreservationdate, Compare:">="})
    } else if  _endreservationdate != "" {          
        args = append(args, models.Where{Column:"reservationdate", Value:_endreservationdate, Compare:"<="})            
    }
    _starttime := c.Get("starttime")
    if _starttime != "" {
        args = append(args, models.Where{Column:"starttime", Value:_starttime, Compare:"like"})
        
    }
    _endtime := c.Get("endtime")
    if _endtime != "" {
        args = append(args, models.Where{Column:"endtime", Value:_endtime, Compare:"like"})
        
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
    _cancelreason := c.Get("cancelreason")
    if _cancelreason != "" {
        args = append(args, models.Where{Column:"cancelreason", Value:_cancelreason, Compare:"like"})
        
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

func (c *PtreservationController) Insert(item *models.Ptreservation) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewPtreservationManager(conn)
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

func (c *PtreservationController) Insertbatch(item *[]models.Ptreservation) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewPtreservationManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *PtreservationController) Update(item *models.Ptreservation) {
    
    
	conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *PtreservationController) Delete(item *models.Ptreservation) {
    
    
    conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *PtreservationController) Deletebatch(item *[]models.Ptreservation) {
    
    
    conn := c.NewConnection()

	manager := models.NewPtreservationManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


