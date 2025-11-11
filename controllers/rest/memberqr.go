package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type MemberqrController struct {
	controllers.Controller
}

func (c *MemberqrController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *MemberqrController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _code := c.Get("code")
    if _code != "" {
        args = append(args, models.Where{Column:"code", Value:_code, Compare:"like"})
    }
    _imageurl := c.Get("imageurl")
    if _imageurl != "" {
        args = append(args, models.Where{Column:"imageurl", Value:_imageurl, Compare:"like"})
    }
    _isactive := c.Geti("isactive")
    if _isactive != 0 {
        args = append(args, models.Where{Column:"isactive", Value:_isactive, Compare:"="})    
    }
    _startexpiredate := c.Get("startexpiredate")
    _endexpiredate := c.Get("endexpiredate")
    if _startexpiredate != "" && _endexpiredate != "" {        
        var v [2]string
        v[0] = _startexpiredate
        v[1] = _endexpiredate  
        args = append(args, models.Where{Column:"expiredate", Value:v, Compare:"between"})    
    } else if  _startexpiredate != "" {          
        args = append(args, models.Where{Column:"expiredate", Value:_startexpiredate, Compare:">="})
    } else if  _endexpiredate != "" {          
        args = append(args, models.Where{Column:"expiredate", Value:_endexpiredate, Compare:"<="})            
    }
    _startgenerateddate := c.Get("startgenerateddate")
    _endgenerateddate := c.Get("endgenerateddate")
    if _startgenerateddate != "" && _endgenerateddate != "" {        
        var v [2]string
        v[0] = _startgenerateddate
        v[1] = _endgenerateddate  
        args = append(args, models.Where{Column:"generateddate", Value:v, Compare:"between"})    
    } else if  _startgenerateddate != "" {          
        args = append(args, models.Where{Column:"generateddate", Value:_startgenerateddate, Compare:">="})
    } else if  _endgenerateddate != "" {          
        args = append(args, models.Where{Column:"generateddate", Value:_endgenerateddate, Compare:"<="})            
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
    _usecount := c.Geti("usecount")
    if _usecount != 0 {
        args = append(args, models.Where{Column:"usecount", Value:_usecount, Compare:"="})    
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
                    str += ", mq_" + strings.Trim(v, " ")                
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

func (c *MemberqrController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _code := c.Get("code")
    if _code != "" {
        args = append(args, models.Where{Column:"code", Value:_code, Compare:"like"})
        
    }
    _imageurl := c.Get("imageurl")
    if _imageurl != "" {
        args = append(args, models.Where{Column:"imageurl", Value:_imageurl, Compare:"like"})
        
    }
    _isactive := c.Geti("isactive")
    if _isactive != 0 {
        args = append(args, models.Where{Column:"isactive", Value:_isactive, Compare:"="})    
    }
    _startexpiredate := c.Get("startexpiredate")
    _endexpiredate := c.Get("endexpiredate")

    if _startexpiredate != "" && _endexpiredate != "" {        
        var v [2]string
        v[0] = _startexpiredate
        v[1] = _endexpiredate  
        args = append(args, models.Where{Column:"expiredate", Value:v, Compare:"between"})    
    } else if  _startexpiredate != "" {          
        args = append(args, models.Where{Column:"expiredate", Value:_startexpiredate, Compare:">="})
    } else if  _endexpiredate != "" {          
        args = append(args, models.Where{Column:"expiredate", Value:_endexpiredate, Compare:"<="})            
    }
    _startgenerateddate := c.Get("startgenerateddate")
    _endgenerateddate := c.Get("endgenerateddate")

    if _startgenerateddate != "" && _endgenerateddate != "" {        
        var v [2]string
        v[0] = _startgenerateddate
        v[1] = _endgenerateddate  
        args = append(args, models.Where{Column:"generateddate", Value:v, Compare:"between"})    
    } else if  _startgenerateddate != "" {          
        args = append(args, models.Where{Column:"generateddate", Value:_startgenerateddate, Compare:">="})
    } else if  _endgenerateddate != "" {          
        args = append(args, models.Where{Column:"generateddate", Value:_endgenerateddate, Compare:"<="})            
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
    _usecount := c.Geti("usecount")
    if _usecount != 0 {
        args = append(args, models.Where{Column:"usecount", Value:_usecount, Compare:"="})    
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

func (c *MemberqrController) Insert(item *models.Memberqr) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewMemberqrManager(conn)
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

func (c *MemberqrController) Insertbatch(item *[]models.Memberqr) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewMemberqrManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *MemberqrController) Update(item *models.Memberqr) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *MemberqrController) Delete(item *models.Memberqr) {
    
    
    conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *MemberqrController) Deletebatch(item *[]models.Memberqr) {
    
    
    conn := c.NewConnection()

	manager := models.NewMemberqrManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


