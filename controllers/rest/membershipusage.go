package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type MembershipusageController struct {
	controllers.Controller
}

func (c *MembershipusageController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *MembershipusageController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)

    var args []interface{}
    
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _totaldays := c.Geti("totaldays")
    if _totaldays != 0 {
        args = append(args, models.Where{Column:"totaldays", Value:_totaldays, Compare:"="})    
    }
    _useddays := c.Geti("useddays")
    if _useddays != 0 {
        args = append(args, models.Where{Column:"useddays", Value:_useddays, Compare:"="})    
    }
    _remainingdays := c.Geti("remainingdays")
    if _remainingdays != 0 {
        args = append(args, models.Where{Column:"remainingdays", Value:_remainingdays, Compare:"="})    
    }
    _totalcount := c.Geti("totalcount")
    if _totalcount != 0 {
        args = append(args, models.Where{Column:"totalcount", Value:_totalcount, Compare:"="})    
    }
    _usedcount := c.Geti("usedcount")
    if _usedcount != 0 {
        args = append(args, models.Where{Column:"usedcount", Value:_usedcount, Compare:"="})    
    }
    _remainingcount := c.Geti("remainingcount")
    if _remainingcount != 0 {
        args = append(args, models.Where{Column:"remainingcount", Value:_remainingcount, Compare:"="})    
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
    _pausedays := c.Geti("pausedays")
    if _pausedays != 0 {
        args = append(args, models.Where{Column:"pausedays", Value:_pausedays, Compare:"="})    
    }
    _startlastuseddate := c.Get("startlastuseddate")
    _endlastuseddate := c.Get("endlastuseddate")
    if _startlastuseddate != "" && _endlastuseddate != "" {        
        var v [2]string
        v[0] = _startlastuseddate
        v[1] = _endlastuseddate  
        args = append(args, models.Where{Column:"lastuseddate", Value:v, Compare:"between"})    
    } else if  _startlastuseddate != "" {          
        args = append(args, models.Where{Column:"lastuseddate", Value:_startlastuseddate, Compare:">="})
    } else if  _endlastuseddate != "" {          
        args = append(args, models.Where{Column:"lastuseddate", Value:_endlastuseddate, Compare:"<="})            
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
                    str += ", mu_" + strings.Trim(v, " ")                
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

func (c *MembershipusageController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)

    var args []interface{}
    
    _membership := c.Geti64("membership")
    if _membership != 0 {
        args = append(args, models.Where{Column:"membership", Value:_membership, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _totaldays := c.Geti("totaldays")
    if _totaldays != 0 {
        args = append(args, models.Where{Column:"totaldays", Value:_totaldays, Compare:"="})    
    }
    _useddays := c.Geti("useddays")
    if _useddays != 0 {
        args = append(args, models.Where{Column:"useddays", Value:_useddays, Compare:"="})    
    }
    _remainingdays := c.Geti("remainingdays")
    if _remainingdays != 0 {
        args = append(args, models.Where{Column:"remainingdays", Value:_remainingdays, Compare:"="})    
    }
    _totalcount := c.Geti("totalcount")
    if _totalcount != 0 {
        args = append(args, models.Where{Column:"totalcount", Value:_totalcount, Compare:"="})    
    }
    _usedcount := c.Geti("usedcount")
    if _usedcount != 0 {
        args = append(args, models.Where{Column:"usedcount", Value:_usedcount, Compare:"="})    
    }
    _remainingcount := c.Geti("remainingcount")
    if _remainingcount != 0 {
        args = append(args, models.Where{Column:"remainingcount", Value:_remainingcount, Compare:"="})    
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
    _pausedays := c.Geti("pausedays")
    if _pausedays != 0 {
        args = append(args, models.Where{Column:"pausedays", Value:_pausedays, Compare:"="})    
    }
    _startlastuseddate := c.Get("startlastuseddate")
    _endlastuseddate := c.Get("endlastuseddate")

    if _startlastuseddate != "" && _endlastuseddate != "" {        
        var v [2]string
        v[0] = _startlastuseddate
        v[1] = _endlastuseddate  
        args = append(args, models.Where{Column:"lastuseddate", Value:v, Compare:"between"})    
    } else if  _startlastuseddate != "" {          
        args = append(args, models.Where{Column:"lastuseddate", Value:_startlastuseddate, Compare:">="})
    } else if  _endlastuseddate != "" {          
        args = append(args, models.Where{Column:"lastuseddate", Value:_endlastuseddate, Compare:"<="})            
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

func (c *MembershipusageController) Insert(item *models.Membershipusage) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewMembershipusageManager(conn)
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

func (c *MembershipusageController) Insertbatch(item *[]models.Membershipusage) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewMembershipusageManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *MembershipusageController) Update(item *models.Membershipusage) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *MembershipusageController) Delete(item *models.Membershipusage) {
    
    
    conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *MembershipusageController) Deletebatch(item *[]models.Membershipusage) {
    
    
    conn := c.NewConnection()

	manager := models.NewMembershipusageManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


