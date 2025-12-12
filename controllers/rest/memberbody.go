package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type MemberbodyController struct {
	controllers.Controller
}

func (c *MemberbodyController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *MemberbodyController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _height := c.Geti("height")
    if _height != 0 {
        args = append(args, models.Where{Column:"height", Value:_height, Compare:"="})    
    }
    _weight := c.Geti("weight")
    if _weight != 0 {
        args = append(args, models.Where{Column:"weight", Value:_weight, Compare:"="})    
    }
    _bodyfat := c.Geti("bodyfat")
    if _bodyfat != 0 {
        args = append(args, models.Where{Column:"bodyfat", Value:_bodyfat, Compare:"="})    
    }
    _musclemass := c.Geti("musclemass")
    if _musclemass != 0 {
        args = append(args, models.Where{Column:"musclemass", Value:_musclemass, Compare:"="})    
    }
    _bmi := c.Geti("bmi")
    if _bmi != 0 {
        args = append(args, models.Where{Column:"bmi", Value:_bmi, Compare:"="})    
    }
    _skeletalmuscle := c.Geti("skeletalmuscle")
    if _skeletalmuscle != 0 {
        args = append(args, models.Where{Column:"skeletalmuscle", Value:_skeletalmuscle, Compare:"="})    
    }
    _bodywater := c.Geti("bodywater")
    if _bodywater != 0 {
        args = append(args, models.Where{Column:"bodywater", Value:_bodywater, Compare:"="})    
    }
    _chest := c.Geti("chest")
    if _chest != 0 {
        args = append(args, models.Where{Column:"chest", Value:_chest, Compare:"="})    
    }
    _waist := c.Geti("waist")
    if _waist != 0 {
        args = append(args, models.Where{Column:"waist", Value:_waist, Compare:"="})    
    }
    _hip := c.Geti("hip")
    if _hip != 0 {
        args = append(args, models.Where{Column:"hip", Value:_hip, Compare:"="})    
    }
    _arm := c.Geti("arm")
    if _arm != 0 {
        args = append(args, models.Where{Column:"arm", Value:_arm, Compare:"="})    
    }
    _thigh := c.Geti("thigh")
    if _thigh != 0 {
        args = append(args, models.Where{Column:"thigh", Value:_thigh, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
    }
    _startmeasureddate := c.Get("startmeasureddate")
    _endmeasureddate := c.Get("endmeasureddate")
    if _startmeasureddate != "" && _endmeasureddate != "" {        
        var v [2]string
        v[0] = _startmeasureddate
        v[1] = _endmeasureddate  
        args = append(args, models.Where{Column:"measureddate", Value:v, Compare:"between"})    
    } else if  _startmeasureddate != "" {          
        args = append(args, models.Where{Column:"measureddate", Value:_startmeasureddate, Compare:">="})
    } else if  _endmeasureddate != "" {          
        args = append(args, models.Where{Column:"measureddate", Value:_endmeasureddate, Compare:"<="})            
    }
    _measuredby := c.Geti64("measuredby")
    if _measuredby != 0 {
        args = append(args, models.Where{Column:"measuredby", Value:_measuredby, Compare:"="})    
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
                    str += ", mb_" + strings.Trim(v, " ")                
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

func (c *MemberbodyController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _height := c.Geti("height")
    if _height != 0 {
        args = append(args, models.Where{Column:"height", Value:_height, Compare:"="})    
    }
    _weight := c.Geti("weight")
    if _weight != 0 {
        args = append(args, models.Where{Column:"weight", Value:_weight, Compare:"="})    
    }
    _bodyfat := c.Geti("bodyfat")
    if _bodyfat != 0 {
        args = append(args, models.Where{Column:"bodyfat", Value:_bodyfat, Compare:"="})    
    }
    _musclemass := c.Geti("musclemass")
    if _musclemass != 0 {
        args = append(args, models.Where{Column:"musclemass", Value:_musclemass, Compare:"="})    
    }
    _bmi := c.Geti("bmi")
    if _bmi != 0 {
        args = append(args, models.Where{Column:"bmi", Value:_bmi, Compare:"="})    
    }
    _skeletalmuscle := c.Geti("skeletalmuscle")
    if _skeletalmuscle != 0 {
        args = append(args, models.Where{Column:"skeletalmuscle", Value:_skeletalmuscle, Compare:"="})    
    }
    _bodywater := c.Geti("bodywater")
    if _bodywater != 0 {
        args = append(args, models.Where{Column:"bodywater", Value:_bodywater, Compare:"="})    
    }
    _chest := c.Geti("chest")
    if _chest != 0 {
        args = append(args, models.Where{Column:"chest", Value:_chest, Compare:"="})    
    }
    _waist := c.Geti("waist")
    if _waist != 0 {
        args = append(args, models.Where{Column:"waist", Value:_waist, Compare:"="})    
    }
    _hip := c.Geti("hip")
    if _hip != 0 {
        args = append(args, models.Where{Column:"hip", Value:_hip, Compare:"="})    
    }
    _arm := c.Geti("arm")
    if _arm != 0 {
        args = append(args, models.Where{Column:"arm", Value:_arm, Compare:"="})    
    }
    _thigh := c.Geti("thigh")
    if _thigh != 0 {
        args = append(args, models.Where{Column:"thigh", Value:_thigh, Compare:"="})    
    }
    _note := c.Get("note")
    if _note != "" {
        args = append(args, models.Where{Column:"note", Value:_note, Compare:"like"})
        
    }
    _startmeasureddate := c.Get("startmeasureddate")
    _endmeasureddate := c.Get("endmeasureddate")

    if _startmeasureddate != "" && _endmeasureddate != "" {        
        var v [2]string
        v[0] = _startmeasureddate
        v[1] = _endmeasureddate  
        args = append(args, models.Where{Column:"measureddate", Value:v, Compare:"between"})    
    } else if  _startmeasureddate != "" {          
        args = append(args, models.Where{Column:"measureddate", Value:_startmeasureddate, Compare:">="})
    } else if  _endmeasureddate != "" {          
        args = append(args, models.Where{Column:"measureddate", Value:_endmeasureddate, Compare:"<="})            
    }
    _measuredby := c.Geti64("measuredby")
    if _measuredby != 0 {
        args = append(args, models.Where{Column:"measuredby", Value:_measuredby, Compare:"="})    
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

func (c *MemberbodyController) Insert(item *models.Memberbody) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewMemberbodyManager(conn)
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

func (c *MemberbodyController) Insertbatch(item *[]models.Memberbody) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewMemberbodyManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *MemberbodyController) Update(item *models.Memberbody) {
    
    
	conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *MemberbodyController) Delete(item *models.Memberbody) {
    
    
    conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *MemberbodyController) Deletebatch(item *[]models.Memberbody) {
    
    
    conn := c.NewConnection()

	manager := models.NewMemberbodyManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


