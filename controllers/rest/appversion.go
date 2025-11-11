package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type AppversionController struct {
	controllers.Controller
}

func (c *AppversionController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *AppversionController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)

    var args []interface{}
    
    _platform := c.Get("platform")
    if _platform != "" {
        args = append(args, models.Where{Column:"platform", Value:_platform, Compare:"like"})
    }
    _version := c.Get("version")
    if _version != "" {
        args = append(args, models.Where{Column:"version", Value:_version, Compare:"like"})
    }
    _minversion := c.Get("minversion")
    if _minversion != "" {
        args = append(args, models.Where{Column:"minversion", Value:_minversion, Compare:"like"})
    }
    _forceupdate := c.Geti("forceupdate")
    if _forceupdate != 0 {
        args = append(args, models.Where{Column:"forceupdate", Value:_forceupdate, Compare:"="})    
    }
    _updatemessage := c.Get("updatemessage")
    if _updatemessage != "" {
        args = append(args, models.Where{Column:"updatemessage", Value:_updatemessage, Compare:"like"})
    }
    _downloadurl := c.Get("downloadurl")
    if _downloadurl != "" {
        args = append(args, models.Where{Column:"downloadurl", Value:_downloadurl, Compare:"like"})
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _startreleasedate := c.Get("startreleasedate")
    _endreleasedate := c.Get("endreleasedate")
    if _startreleasedate != "" && _endreleasedate != "" {        
        var v [2]string
        v[0] = _startreleasedate
        v[1] = _endreleasedate  
        args = append(args, models.Where{Column:"releasedate", Value:v, Compare:"between"})    
    } else if  _startreleasedate != "" {          
        args = append(args, models.Where{Column:"releasedate", Value:_startreleasedate, Compare:">="})
    } else if  _endreleasedate != "" {          
        args = append(args, models.Where{Column:"releasedate", Value:_endreleasedate, Compare:"<="})            
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
                    str += ", av_" + strings.Trim(v, " ")                
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

func (c *AppversionController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)

    var args []interface{}
    
    _platform := c.Get("platform")
    if _platform != "" {
        args = append(args, models.Where{Column:"platform", Value:_platform, Compare:"like"})
        
    }
    _version := c.Get("version")
    if _version != "" {
        args = append(args, models.Where{Column:"version", Value:_version, Compare:"like"})
        
    }
    _minversion := c.Get("minversion")
    if _minversion != "" {
        args = append(args, models.Where{Column:"minversion", Value:_minversion, Compare:"like"})
        
    }
    _forceupdate := c.Geti("forceupdate")
    if _forceupdate != 0 {
        args = append(args, models.Where{Column:"forceupdate", Value:_forceupdate, Compare:"="})    
    }
    _updatemessage := c.Get("updatemessage")
    if _updatemessage != "" {
        args = append(args, models.Where{Column:"updatemessage", Value:_updatemessage, Compare:"like"})
        
    }
    _downloadurl := c.Get("downloadurl")
    if _downloadurl != "" {
        args = append(args, models.Where{Column:"downloadurl", Value:_downloadurl, Compare:"like"})
        
    }
    _status := c.Geti("status")
    if _status != 0 {
        args = append(args, models.Where{Column:"status", Value:_status, Compare:"="})    
    }
    _startreleasedate := c.Get("startreleasedate")
    _endreleasedate := c.Get("endreleasedate")

    if _startreleasedate != "" && _endreleasedate != "" {        
        var v [2]string
        v[0] = _startreleasedate
        v[1] = _endreleasedate  
        args = append(args, models.Where{Column:"releasedate", Value:v, Compare:"between"})    
    } else if  _startreleasedate != "" {          
        args = append(args, models.Where{Column:"releasedate", Value:_startreleasedate, Compare:">="})
    } else if  _endreleasedate != "" {          
        args = append(args, models.Where{Column:"releasedate", Value:_endreleasedate, Compare:"<="})            
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

func (c *AppversionController) Insert(item *models.Appversion) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewAppversionManager(conn)
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

func (c *AppversionController) Insertbatch(item *[]models.Appversion) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewAppversionManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *AppversionController) Update(item *models.Appversion) {
    
    
	conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *AppversionController) Delete(item *models.Appversion) {
    
    
    conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *AppversionController) Deletebatch(item *[]models.Appversion) {
    
    
    conn := c.NewConnection()

	manager := models.NewAppversionManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


