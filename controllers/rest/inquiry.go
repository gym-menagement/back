package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type InquiryController struct {
	controllers.Controller
}

func (c *InquiryController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *InquiryController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
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
    _content := c.Get("content")
    if _content != "" {
        args = append(args, models.Where{Column:"content", Value:_content, Compare:"like"})
        
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _answer := c.Get("answer")
    if _answer != "" {
        args = append(args, models.Where{Column:"answer", Value:_answer, Compare:"like"})
    }
    _answeredby := c.Geti64("answeredby")
    if _answeredby != 0 {
        args = append(args, models.Where{Column:"answeredby", Value:_answeredby, Compare:"="})    
    }
    _startanswereddate := c.Get("startanswereddate")
    _endanswereddate := c.Get("endanswereddate")
    if _startanswereddate != "" && _endanswereddate != "" {        
        var v [2]string
        v[0] = _startanswereddate
        v[1] = _endanswereddate  
        args = append(args, models.Where{Column:"answereddate", Value:v, Compare:"between"})    
    } else if  _startanswereddate != "" {          
        args = append(args, models.Where{Column:"answereddate", Value:_startanswereddate, Compare:">="})
    } else if  _endanswereddate != "" {          
        args = append(args, models.Where{Column:"answereddate", Value:_endanswereddate, Compare:"<="})            
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
                    str += ", iq_" + strings.Trim(v, " ")                
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

func (c *InquiryController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)

    var args []interface{}
    
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
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
    _content := c.Get("content")
    if _content != "" {
        args = append(args, models.Where{Column:"content", Value:_content, Compare:"like"})
        
        
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _answer := c.Get("answer")
    if _answer != "" {
        args = append(args, models.Where{Column:"answer", Value:_answer, Compare:"like"})
        
    }
    _answeredby := c.Geti64("answeredby")
    if _answeredby != 0 {
        args = append(args, models.Where{Column:"answeredby", Value:_answeredby, Compare:"="})    
    }
    _startanswereddate := c.Get("startanswereddate")
    _endanswereddate := c.Get("endanswereddate")

    if _startanswereddate != "" && _endanswereddate != "" {        
        var v [2]string
        v[0] = _startanswereddate
        v[1] = _endanswereddate  
        args = append(args, models.Where{Column:"answereddate", Value:v, Compare:"between"})    
    } else if  _startanswereddate != "" {          
        args = append(args, models.Where{Column:"answereddate", Value:_startanswereddate, Compare:">="})
    } else if  _endanswereddate != "" {          
        args = append(args, models.Where{Column:"answereddate", Value:_endanswereddate, Compare:"<="})            
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

func (c *InquiryController) Insert(item *models.Inquiry) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewInquiryManager(conn)
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

func (c *InquiryController) Insertbatch(item *[]models.Inquiry) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewInquiryManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *InquiryController) Update(item *models.Inquiry) {
    
    
	conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *InquiryController) Delete(item *models.Inquiry) {
    
    
    conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *InquiryController) Deletebatch(item *[]models.Inquiry) {
    
    
    conn := c.NewConnection()

	manager := models.NewInquiryManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


