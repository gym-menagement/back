package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type RockerusageController struct {
	controllers.Controller
}

func (c *RockerusageController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *RockerusageController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)

    var args []interface{}
    
    _rocker := c.Geti64("rocker")
    if _rocker != 0 {
        args = append(args, models.Where{Column:"rocker", Value:_rocker, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
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
    _deposit := c.Geti("deposit")
    if _deposit != 0 {
        args = append(args, models.Where{Column:"deposit", Value:_deposit, Compare:"="})    
    }
    _monthlyfee := c.Geti("monthlyfee")
    if _monthlyfee != 0 {
        args = append(args, models.Where{Column:"monthlyfee", Value:_monthlyfee, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
    }
    _assignedby := c.Geti64("assignedby")
    if _assignedby != 0 {
        args = append(args, models.Where{Column:"assignedby", Value:_assignedby, Compare:"="})    
    }
    _startassigneddate := c.Get("startassigneddate")
    _endassigneddate := c.Get("endassigneddate")
    if _startassigneddate != "" && _endassigneddate != "" {        
        var v [2]string
        v[0] = _startassigneddate
        v[1] = _endassigneddate  
        args = append(args, models.Where{Column:"assigneddate", Value:v, Compare:"between"})    
    } else if  _startassigneddate != "" {          
        args = append(args, models.Where{Column:"assigneddate", Value:_startassigneddate, Compare:">="})
    } else if  _endassigneddate != "" {          
        args = append(args, models.Where{Column:"assigneddate", Value:_endassigneddate, Compare:"<="})            
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
                    str += ", ru_" + strings.Trim(v, " ")                
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

func (c *RockerusageController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)

    var args []interface{}
    
    _rocker := c.Geti64("rocker")
    if _rocker != 0 {
        args = append(args, models.Where{Column:"rocker", Value:_rocker, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
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
    _deposit := c.Geti("deposit")
    if _deposit != 0 {
        args = append(args, models.Where{Column:"deposit", Value:_deposit, Compare:"="})    
    }
    _monthlyfee := c.Geti("monthlyfee")
    if _monthlyfee != 0 {
        args = append(args, models.Where{Column:"monthlyfee", Value:_monthlyfee, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
        
    }
    _assignedby := c.Geti64("assignedby")
    if _assignedby != 0 {
        args = append(args, models.Where{Column:"assignedby", Value:_assignedby, Compare:"="})    
    }
    _startassigneddate := c.Get("startassigneddate")
    _endassigneddate := c.Get("endassigneddate")

    if _startassigneddate != "" && _endassigneddate != "" {        
        var v [2]string
        v[0] = _startassigneddate
        v[1] = _endassigneddate  
        args = append(args, models.Where{Column:"assigneddate", Value:v, Compare:"between"})    
    } else if  _startassigneddate != "" {          
        args = append(args, models.Where{Column:"assigneddate", Value:_startassigneddate, Compare:">="})
    } else if  _endassigneddate != "" {          
        args = append(args, models.Where{Column:"assigneddate", Value:_endassigneddate, Compare:"<="})            
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

func (c *RockerusageController) Insert(item *models.Rockerusage) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewRockerusageManager(conn)
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

func (c *RockerusageController) Insertbatch(item *[]models.Rockerusage) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewRockerusageManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *RockerusageController) Update(item *models.Rockerusage) {
    
    
	conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *RockerusageController) Delete(item *models.Rockerusage) {
    
    
    conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *RockerusageController) Deletebatch(item *[]models.Rockerusage) {
    
    
    conn := c.NewConnection()

	manager := models.NewRockerusageManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


