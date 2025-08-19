package rest


import (
	"gym/controllers"
	"gym/models"

    "strings"
)

type PaymentformController struct {
	controllers.Controller
}

func (c *PaymentformController) Read(id int64) {
    
    
	conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)
	item := manager.Get(id)

    
    
    c.Set("item", item)
}

func (c *PaymentformController) Index(page int, pagesize int) {
    
    
	conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _payment := c.Geti64("payment")
    if _payment != 0 {
        args = append(args, models.Where{Column:"payment", Value:_payment, Compare:"="})    
    }
    _type := c.Geti64("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _cost := c.Geti("cost")
    if _cost != 0 {
        args = append(args, models.Where{Column:"cost", Value:_cost, Compare:"="})    
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
                    str += ", pf_" + strings.Trim(v, " ")                
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

func (c *PaymentformController) Count() {
    
    
	conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)

    var args []interface{}
    
    _gym := c.Geti64("gym")
    if _gym != 0 {
        args = append(args, models.Where{Column:"gym", Value:_gym, Compare:"="})    
    }
    _payment := c.Geti64("payment")
    if _payment != 0 {
        args = append(args, models.Where{Column:"payment", Value:_payment, Compare:"="})    
    }
    _type := c.Geti64("type")
    if _type != 0 {
        args = append(args, models.Where{Column:"type", Value:_type, Compare:"="})    
    }
    _cost := c.Geti("cost")
    if _cost != 0 {
        args = append(args, models.Where{Column:"cost", Value:_cost, Compare:"="})    
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

func (c *PaymentformController) Insert(item *models.Paymentform) {
    
    
	conn := c.NewConnection()
    
	manager := models.NewPaymentformManager(conn)
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

func (c *PaymentformController) Insertbatch(item *[]models.Paymentform) {  
    if item == nil || len(*item) == 0 {
        return
    }

    rows := len(*item)
    
    
    
	conn := c.NewConnection()
    
	manager := models.NewPaymentformManager(conn)

    for i := 0; i < rows; i++ {
	    err := manager.Insert(&((*item)[i]))
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}

func (c *PaymentformController) Update(item *models.Paymentform) {
    
    
	conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)
    err := manager.Update(item)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
        return
    }
}

func (c *PaymentformController) Delete(item *models.Paymentform) {
    
    
    conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)

    
	err := manager.Delete(item.Id)
    if err != nil {
        c.Set("code", "error")    
        c.Set("error", err)
    }
}

func (c *PaymentformController) Deletebatch(item *[]models.Paymentform) {
    
    
    conn := c.NewConnection()

	manager := models.NewPaymentformManager(conn)

    for _, v := range *item {
        
    
	    err := manager.Delete(v.Id)
        if err != nil {
            c.Set("code", "error")    
            c.Set("error", err)
            return
        }
    }
}


