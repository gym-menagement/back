package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type TrainermemberController struct {
	controllers.Controller
}

func (c *TrainermemberController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *TrainermemberController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)

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
    _startstartdate := c.Get("startstartdate")
    _endstartdate := c.Get("endstartdate")
    if _startstartdate != "" && _endstartdate != "" {        
        var v [2]string
        v[0] = _startstartdate
        v[1] = _endstartdate  
        args = append(args, models.Where{Column:"startdate", Value:v, Compare:"between"})    
    } else if  _startstartdate != "" {          
        args = append(args, models.Where{Column:"startdate", Value:_startstartdate, Compare:">="})
    } else if  _endstartdate != "" {          
        args = append(args, models.Where{Column:"startdate", Value:_endstartdate, Compare:"<="})            
    }
    _startenddate := c.Get("startenddate")
    _endenddate := c.Get("endenddate")
    if _startenddate != "" && _endenddate != "" {        
        var v [2]string
        v[0] = _startenddate
        v[1] = _endenddate  
        args = append(args, models.Where{Column:"enddate", Value:v, Compare:"between"})    
    } else if  _startenddate != "" {          
        args = append(args, models.Where{Column:"enddate", Value:_startenddate, Compare:">="})
    } else if  _endenddate != "" {          
        args = append(args, models.Where{Column:"enddate", Value:_endenddate, Compare:"<="})            
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
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
                    str += ", tm_" + strings.Trim(v, " ")                
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

func (c *TrainermemberController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)

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
    _startstartdate := c.Get("startstartdate")
    _endstartdate := c.Get("endstartdate")

    if _startstartdate != "" && _endstartdate != "" {        
        var v [2]string
        v[0] = _startstartdate
        v[1] = _endstartdate  
        args = append(args, models.Where{Column:"startdate", Value:v, Compare:"between"})    
    } else if  _startstartdate != "" {          
        args = append(args, models.Where{Column:"startdate", Value:_startstartdate, Compare:">="})
    } else if  _endstartdate != "" {          
        args = append(args, models.Where{Column:"startdate", Value:_endstartdate, Compare:"<="})            
    }
    _startenddate := c.Get("startenddate")
    _endenddate := c.Get("endenddate")

    if _startenddate != "" && _endenddate != "" {        
        var v [2]string
        v[0] = _startenddate
        v[1] = _endenddate  
        args = append(args, models.Where{Column:"enddate", Value:v, Compare:"between"})    
    } else if  _startenddate != "" {          
        args = append(args, models.Where{Column:"enddate", Value:_startenddate, Compare:">="})
    } else if  _endenddate != "" {          
        args = append(args, models.Where{Column:"enddate", Value:_endenddate, Compare:"<="})            
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
        
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

func (c *TrainermemberController) Insert(item *models.Trainermember) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewTrainermemberManager(conn)
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

func (c *TrainermemberController) Insertbatch(item *[]models.Trainermember) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewTrainermemberManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *TrainermemberController) Update(item *models.Trainermember) {
    
    
	conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *TrainermemberController) Delete(item *models.Trainermember) {
    
    
    conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *TrainermemberController) Deletebatch(item *[]models.Trainermember) {
    
    
    conn := c.NewConnection()

	manager := models.NewTrainermemberManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


