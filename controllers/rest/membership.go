package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type MembershipController struct {
	controllers.Controller
}

func (c *MembershipController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *MembershipController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
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
    _phonenum := c.Get("phonenum")
    if _phonenum != "" {
        args = append(args, models.Where{Column:"phonenum", Value:_phonenum, Compare:"like"})
    }
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
    }
    _image := c.Get("image")
    if _image != "" {
        args = append(args, models.Where{Column:"image", Value:_image, Compare:"like"})
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
                    str += ", m_" + strings.Trim(v, " ")                
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

func (c *MembershipController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
        
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
    _phonenum := c.Get("phonenum")
    if _phonenum != "" {
        args = append(args, models.Where{Column:"phonenum", Value:_phonenum, Compare:"like"})
        
    }
    _address := c.Get("address")
    if _address != "" {
        args = append(args, models.Where{Column:"address", Value:_address, Compare:"like"})
        
    }
    _image := c.Get("image")
    if _image != "" {
        args = append(args, models.Where{Column:"image", Value:_image, Compare:"like"})
        
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

func (c *MembershipController) Insert(item *models.Membership) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewMembershipManager(conn)
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

func (c *MembershipController) Insertbatch(item *[]models.Membership) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewMembershipManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *MembershipController) Update(item *models.Membership) {
    
    
	conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *MembershipController) Delete(item *models.Membership) {
    
    
    conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *MembershipController) Deletebatch(item *[]models.Membership) {
    
    
    conn := c.NewConnection()

	manager := models.NewMembershipManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


