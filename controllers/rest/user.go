package rest


import (
	"gym/controllers"
	"gym/models"

	"gym/models/user"

    "strings"
)

type UserController struct {
	controllers.Controller
}

func (c *UserController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *UserController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    var args []interface{}
    
    _loginid := c.Get("loginid")
    if _loginid != "" {
        args = append(args, models.Where{Column:"loginid", Value:_loginid, Compare:"like"})
    }
    _passwd := c.Get("passwd")
    if _passwd != "" {
        args = append(args, models.Where{Column:"passwd", Value:_passwd, Compare:"like"})
    }
    _email := c.Get("email")
    if _email != "" {
        args = append(args, models.Where{Column:"email", Value:_email, Compare:"like"})
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _tel := c.Get("tel")
    if _tel != "" {
        args = append(args, models.Where{Column:"tel", Value:_tel, Compare:"like"})
    }
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
    }
    _image := c.Get("image")
    if _image != "" {
        args = append(args, models.Where{Column:"image", Value:_image, Compare:"like"})
    }
    _sex := c.Geti("sex")
    if _sex != 0 {
        args = append(args, models.Where{Column:"sex", Value:_sex, Compare:"="})    
    }
    _startbirth := c.Get("startbirth")
    _endbirth := c.Get("endbirth")
    if _startbirth != "" && _endbirth != "" {        
        var v [2]string
        v[0] = _startbirth
        v[1] = _endbirth  
        args = append(args, models.Where{Column:"birth", Value:v, Compare:"between"})    
    } else if  _startbirth != "" {          
        args = append(args, models.Where{Column:"birth", Value:_startbirth, Compare:">="})
    } else if  _endbirth != "" {          
        args = append(args, models.Where{Column:"birth", Value:_endbirth, Compare:"<="})            
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _connectid := c.Get("connectid")
    if _connectid != "" {
        args = append(args, models.Where{Column:"connectid", Value:_connectid, Compare:"like"})
    }
    _level := c.Geti("level")
    if _level != 0 {
        args = append(args, models.Where{Column:"level", Value:_level, Compare:"="})    
    }
    _role := c.Geti("role")
    if _role != 0 {
        args = append(args, models.Where{Column:"role", Value:_role, Compare:"="})    
    }
    _use := c.Geti("use")
    if _use != 0 {
        args = append(args, models.Where{Column:"use", Value:_use, Compare:"="})    
    }
    _startlogindate := c.Get("startlogindate")
    _endlogindate := c.Get("endlogindate")
    if _startlogindate != "" && _endlogindate != "" {        
        var v [2]string
        v[0] = _startlogindate
        v[1] = _endlogindate  
        args = append(args, models.Where{Column:"logindate", Value:v, Compare:"between"})    
    } else if  _startlogindate != "" {          
        args = append(args, models.Where{Column:"logindate", Value:_startlogindate, Compare:">="})
    } else if  _endlogindate != "" {          
        args = append(args, models.Where{Column:"logindate", Value:_endlogindate, Compare:"<="})            
    }
    _startlastchangepasswddate := c.Get("startlastchangepasswddate")
    _endlastchangepasswddate := c.Get("endlastchangepasswddate")
    if _startlastchangepasswddate != "" && _endlastchangepasswddate != "" {        
        var v [2]string
        v[0] = _startlastchangepasswddate
        v[1] = _endlastchangepasswddate  
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:v, Compare:"between"})    
    } else if  _startlastchangepasswddate != "" {          
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:_startlastchangepasswddate, Compare:">="})
    } else if  _endlastchangepasswddate != "" {          
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:_endlastchangepasswddate, Compare:"<="})            
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
                    str += ", u_" + strings.Trim(v, " ")                
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

func (c *UserController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    var args []interface{}
    
    _loginid := c.Get("loginid")
    if _loginid != "" {
        args = append(args, models.Where{Column:"loginid", Value:_loginid, Compare:"like"})
        
    }
    _passwd := c.Get("passwd")
    if _passwd != "" {
        args = append(args, models.Where{Column:"passwd", Value:_passwd, Compare:"like"})
        
    }
    _email := c.Get("email")
    if _email != "" {
        args = append(args, models.Where{Column:"email", Value:_email, Compare:"like"})
        
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
        
    }
    _tel := c.Get("tel")
    if _tel != "" {
        args = append(args, models.Where{Column:"tel", Value:_tel, Compare:"like"})
        
    }
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
        
    }
    _image := c.Get("image")
    if _image != "" {
        args = append(args, models.Where{Column:"image", Value:_image, Compare:"like"})
        
    }
    _sex := c.Geti("sex")
    if _sex != 0 {
        args = append(args, models.Where{Column:"sex", Value:_sex, Compare:"="})    
    }
    _startbirth := c.Get("startbirth")
    _endbirth := c.Get("endbirth")

    if _startbirth != "" && _endbirth != "" {        
        var v [2]string
        v[0] = _startbirth
        v[1] = _endbirth  
        args = append(args, models.Where{Column:"birth", Value:v, Compare:"between"})    
    } else if  _startbirth != "" {          
        args = append(args, models.Where{Column:"birth", Value:_startbirth, Compare:">="})
    } else if  _endbirth != "" {          
        args = append(args, models.Where{Column:"birth", Value:_endbirth, Compare:"<="})            
    }
    _type := c.Geti("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _connectid := c.Get("connectid")
    if _connectid != "" {
        args = append(args, models.Where{Column:"connectid", Value:_connectid, Compare:"like"})
        
    }
    _level := c.Geti("level")
    if _level != 0 {
        args = append(args, models.Where{Column:"level", Value:_level, Compare:"="})    
    }
    _role := c.Geti("role")
    if _role != 0 {
        args = append(args, models.Where{Column:"role", Value:_role, Compare:"="})    
    }
    _use := c.Geti("use")
    if _use != 0 {
        args = append(args, models.Where{Column:"use", Value:_use, Compare:"="})    
    }
    _startlogindate := c.Get("startlogindate")
    _endlogindate := c.Get("endlogindate")

    if _startlogindate != "" && _endlogindate != "" {        
        var v [2]string
        v[0] = _startlogindate
        v[1] = _endlogindate  
        args = append(args, models.Where{Column:"logindate", Value:v, Compare:"between"})    
    } else if  _startlogindate != "" {          
        args = append(args, models.Where{Column:"logindate", Value:_startlogindate, Compare:">="})
    } else if  _endlogindate != "" {          
        args = append(args, models.Where{Column:"logindate", Value:_endlogindate, Compare:"<="})            
    }
    _startlastchangepasswddate := c.Get("startlastchangepasswddate")
    _endlastchangepasswddate := c.Get("endlastchangepasswddate")

    if _startlastchangepasswddate != "" && _endlastchangepasswddate != "" {        
        var v [2]string
        v[0] = _startlastchangepasswddate
        v[1] = _endlastchangepasswddate  
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:v, Compare:"between"})    
    } else if  _startlastchangepasswddate != "" {          
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:_startlastchangepasswddate, Compare:">="})
    } else if  _endlastchangepasswddate != "" {          
        args = append(args, models.Where{Column:"lastchangepasswddate", Value:_endlastchangepasswddate, Compare:"<="})            
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

func (c *UserController) Insert(item *models.UserUpdate) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewUserManager(conn)
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

func (c *UserController) Insertbatch(item *[]models.UserUpdate) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewUserManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *UserController) Update(item *models.UserUpdate) {
    
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *UserController) Delete(item *models.User) {
    
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *UserController) Deletebatch(item *[]models.User) {
    
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}



func (c *UserController) GetByLoginid(loginid string) *models.User {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.GetByLoginid(loginid)
    
    c.Set("item", item)
    
    
    
    return item
    
}


func (c *UserController) GetByConnectid(connectid string) *models.User {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.GetByConnectid(connectid)
    
    c.Set("item", item)
    
    
    
    return item
    
}


func (c *UserController) CountByLoginid(loginid string) int {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.CountByLoginid(loginid)
    
    
    
    c.Set("count", item)
    
    return item
    
}

// @Put()
func (c *UserController) UpdateLogindateById(logindate string ,id int64) {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    manager.UpdateLogindateById(logindate, id)
    
}


func (c *UserController) FindByLevel(level user.Level) []models.User {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.FindByLevel(level)
    
    
    c.Set("items", item)
    
    
    return item
    
}


func (c *UserController) FindByEmail(email string) []models.User {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.FindByEmail(email)
    
    
    c.Set("items", item)
    
    
    return item
    
}


func (c *UserController) FindByTel(tel string) []models.User {
    
    conn := c.NewConnection()

	manager := models.NewUserManager(conn)
    
    item := manager.FindByTel(tel)
    
    
    c.Set("items", item)
    
    
    return item
    
}

