package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type WorkoutlogController struct {
	controllers.Controller
}

func (c *WorkoutlogController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *WorkoutlogController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _attendance := c.Geti64("attendance")
    if _attendance != 0 {
        args = append(args, models.Where{Column:"attendance", Value:_attendance, Compare:"="})    
    }
    _health := c.Geti64("health")
    if _health != 0 {
        args = append(args, models.Where{Column:"health", Value:_health, Compare:"="})    
    }
    _exercisename := c.Get("exercisename")
    if _exercisename != "" {
        args = append(args, models.Where{Column:"exercisename", Value:_exercisename, Compare:"like"})
    }
    _sets := c.Geti("sets")
    if _sets != 0 {
        args = append(args, models.Where{Column:"sets", Value:_sets, Compare:"="})    
    }
    _reps := c.Geti("reps")
    if _reps != 0 {
        args = append(args, models.Where{Column:"reps", Value:_reps, Compare:"="})    
    }
    _weight := c.Geti("weight")
    if _weight != 0 {
        args = append(args, models.Where{Column:"weight", Value:_weight, Compare:"="})    
    }
    _duration := c.Geti("duration")
    if _duration != 0 {
        args = append(args, models.Where{Column:"duration", Value:_duration, Compare:"="})    
    }
    _calories := c.Geti("calories")
    if _calories != 0 {
        args = append(args, models.Where{Column:"calories", Value:_calories, Compare:"="})    
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
                    str += ", wl_" + strings.Trim(v, " ")                
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

func (c *WorkoutlogController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
    }
    _attendance := c.Geti64("attendance")
    if _attendance != 0 {
        args = append(args, models.Where{Column:"attendance", Value:_attendance, Compare:"="})    
    }
    _health := c.Geti64("health")
    if _health != 0 {
        args = append(args, models.Where{Column:"health", Value:_health, Compare:"="})    
    }
    _exercisename := c.Get("exercisename")
    if _exercisename != "" {
        args = append(args, models.Where{Column:"exercisename", Value:_exercisename, Compare:"like"})
        
    }
    _sets := c.Geti("sets")
    if _sets != 0 {
        args = append(args, models.Where{Column:"sets", Value:_sets, Compare:"="})    
    }
    _reps := c.Geti("reps")
    if _reps != 0 {
        args = append(args, models.Where{Column:"reps", Value:_reps, Compare:"="})    
    }
    _weight := c.Geti("weight")
    if _weight != 0 {
        args = append(args, models.Where{Column:"weight", Value:_weight, Compare:"="})    
    }
    _duration := c.Geti("duration")
    if _duration != 0 {
        args = append(args, models.Where{Column:"duration", Value:_duration, Compare:"="})    
    }
    _calories := c.Geti("calories")
    if _calories != 0 {
        args = append(args, models.Where{Column:"calories", Value:_calories, Compare:"="})    
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

func (c *WorkoutlogController) Insert(item *models.Workoutlog) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewWorkoutlogManager(conn)
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

func (c *WorkoutlogController) Insertbatch(item *[]models.Workoutlog) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewWorkoutlogManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *WorkoutlogController) Update(item *models.Workoutlog) {
    
    
	conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *WorkoutlogController) Delete(item *models.Workoutlog) {
    
    
    conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *WorkoutlogController) Deletebatch(item *[]models.Workoutlog) {
    
    
    conn := c.NewConnection()

	manager := models.NewWorkoutlogManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


