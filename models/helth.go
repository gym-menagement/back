package models

import (
	"gym/config"

	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)


type Helth struct {
	Id					int64 `json:"id"`
	Category			int64 `json:"category"`
	Term				int64 `json:"term"`
	Name				string `json:"name"`
	Count				int `json:"count"`
	Cost				int `json:"cost"`
	Discount			int64 `json:"discount"`
	CostDiscount		int `json:"costdiscount"`
	Content				string `json:"content"`
	Date				string `json:"date"`

	Extra				map[string]interface{} `json:"extra"`
}

type HelthManager struct {
	Conn	*sql.DB
	Tx		*sql.Tx
	Result	*sql.Result
	Index	string
}

func (c *Helth) AddExtra(key string, value interface{}) {
	c.Extra[key] = value
}

func NewHelthManager(conn interface{}) *HelthManager {
	var item HelthManager

	if conn == nil {
		item.Conn = NewConnection()
	} else {
		if v, ok := conn.(*sql.DB); ok {
			item.Conn = v
			item.Tx = nil
		} else {
			item.Tx = conn.(*sql.Tx)
			item.Conn = nil
		}
	}
	
	item.Index = ""
	return &item
}

func (p *HelthManager) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *HelthManager) SetIndex(index string) {
	p.Index = index
}

func (p *HelthManager) Exec(query string, params ...interface{}) (sql.Result, error) {
	if p.Conn != nil {
		return p.Conn.Exec(query, params...)
	} else {
		return p.Tx.Exec(query, params...)
	}
}

func (p *HelthManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if p.Conn != nil {
		return p.Conn.Query(query, params...)
	} else {
		return p.Tx.Query(query + " FOR UPDATE", params...)
	}
}

func (p *HelthManager) GetQeury() string {
	ret := ""

	str := "select h_id, h_category, h_term, h_name, h_count, h_cost, h_discount, h_costdiscount, h_content, h_date from helth_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ")"
	}

	ret += "where 1=1 "

	return ret;
}

func (p *HelthManager) GetQeurySelect() string {
	ret := ""

	str := "select count(*) from helth_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret;
}

func (p *HelthManager) Truncate() error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "truncate helth_tb "
	p.Exec(query)

	return nil
}

func (p *HelthManager) Insert(item *Helth) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	if item.Date == "" {
		t := time.Now()
		item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}

	query := ""
	var res sql.Result
	var err error
	if item.Id > 0 {
		query = "insert into helth_tb (h_id, h_category, h_term, h_name, h_count, h_cost, h_discount, h_costdiscount, h_content, h_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Id, item.Category, item.Term, item.Name, item.Count, item.Cost, item.Discount, item.CostDiscount, item.Content, item.Date)
	} else {
		query = "insert into helth_tb (h_category, h_term, h_name, h_count, h_cost, h_discount, h_costdiscount, h_content, h_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Category, item.Term, item.Name, item.Count, item.Cost, item.Discount, item.CostDiscount, item.Content, item.Date)
	}

	if err == nil {
		p.Result = &res
	} else {
		log.Println(err)
		p.Result = nil
	}

	return err
}

func (p *HelthManager) Delete(id int64) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "delete from helth_tb where h_id = ?"
	_, err := p.Exec(query, id)

	return err
}

func (p *HelthManager) Update(item *Helth) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "update helth_tb set h_category = ?, h_term = ?, h_name = ?, h_count = ?, h_cost = ?, h_discount = ?, h_costdiscount = ?, h_content = ?, h_date = ? where h_id = ?"
	_, err := p.Exec(query, item.Category, item.Term, item.Name, item.Count, item.Cost, item.Discount, item.CostDiscount, item.Content, item.Date, item.Id)

	return err
}

func (p *HelthManager) GetIdentity() int64 {
	if p.Result == nil && p.Tx == nil {
		return 0
	}

	id, err := (*p.Result).LastInsertId()

	if err != nil {
		return 0
	} else {
		return id
	}
}

func (p *Helth) InitExtra() {
	p.Extra = map[string]interface{}{

	}
}

func (p *HelthManager) ReadRow(rows *sql.Rows) *Helth {
	var item Helth
	var err error

	if rows.Next() {
		err = rows.Scan(&item.Id, &item.Category, &item.Term, &item.Name, &item.Count, &item.Cost, &item.Discount, &item.CostDiscount, &item.Content, &item.Date)
	} else {
		return nil
	}
	if err != nil {
		return nil
	} else {
		item.InitExtra()
		return &item
	}
}

func (p *HelthManager) ReadRows(rows *sql.Rows) *[]Helth {
	var items []Helth

	for rows.Next() {
		var item Helth

		err := rows.Scan(&item.Id, &item.Category, &item.Term, &item.Name, &item.Count, &item.Cost, &item.Discount, &item.CostDiscount, &item.Content, &item.Date)

		if err != nil {
			log.Printf("ReadRows error : %v\n", err)
			break
		}

		item.InitExtra()

		items = append(items, item)
	}
	return &items
}

func (p *HelthManager) Get(id int64) *Helth {
	if p.Conn == nil && p.Tx == nil {
		return nil
	}

	query := p.GetQeury() + " and h_id = ?"

	rows, err := p.Query(query, id)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	return p.ReadRow(rows)
}

func (p *HelthManager) Count(args []interface{}) int {
	if p.Conn == nil && p.Tx == nil {
		return 0
	}

	var params []interface{}
	query := p.GetQeurySelect() + " where 1=1 "

	for _, arg := range args {
		switch v := arg.(type) {
		case Where:
			item := v

			if item.Compare == "in" {
				query += " and h_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and h_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and h_" + item.Column + " " + item.Compare + " ?"
				if item.Compare == "like" {
					params = append(params, "%" + item.Value.(string) + "%")
				} else {
					params = append(params, item.Value)
				}
			}
		}
	}

	rows, err := p.Query(query, params...)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return 0
	}

	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	cnt := 0
	err = rows.Scan(&cnt)

	if err != nil {
		return 0
	} else {
		return cnt
	}
}

func (p *HelthManager) Find(args []interface{}) *[]Helth {
	if p.Conn == nil && p.Tx == nil {
		var items []Helth
		return &items
	}

	var params []interface{}
	query := p.GetQeury()

	page := 0
	pagesize := 0
	orderby := ""

	for _, arg := range args {
		switch v := arg.(type) {
		case PagingType:
			item := v
			page = item.Page
			pagesize = item.Pagesize
			break
		case OrderingType:
			item := v
			orderby = item.Order
			break
		case LimitType:
			item := v
			page = 1
			pagesize = item.Limit
			break
		case OptionType:
			item := v
			if item.Limit > 0 {
				page = 1
				pagesize = item.Limit
			} else {
				page = item.Page
				pagesize = item.Pagesize
			}
			orderby = item.Order
			break
		case Where:
			item := v

			if item.Compare == "in" {
				query += " and h_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and h_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and h_" + item.Column + " " + item.Compare + " ?"
				if item.Compare == "like" {
					params = append(params, "%" + item.Value.(string) + "%")
				} else {
					params = append(params, item.Value)
				}
			}
		}
	}

	startpage := (page -1) * pagesize

	if page > 0 && pagesize > 0 {
		if orderby == "" {
			orderby = "h_id"
		} else {
			orderby = "h_" + orderby
		}
		query += " order by " + orderby
		if config.Database == "mysql" {
			query += " limit ? offset ?"
			params = append(params, pagesize)
			params = append(params, startpage)
		} else if config.Database == "mssql" || config.Database == "sqlserver" {
			query += "OFFSET ? ROWS FITCH NEXT ? ROWS ONLY"
			params = append(params, startpage)
			params = append(params, pagesize)
		}
	} else {
		if orderby == "" {
			orderby = "h_id"
		} else {
			orderby = "h_" + orderby
		}
		query += " order by " + orderby
	}

	rows, err := p.Query(query, params...)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		var items []Helth
		return &items
	}

	defer rows.Close()

	return p.ReadRows(rows)
}

func (p *HelthManager) GetByName(loginid string, args ...interface{}) *Helth {
    if loginid != "" {
        args = append(args, Where{Column:"name", Value:loginid, Compare:"="})        
    }
    
    items := p.Find(args)

    if items != nil && len(*items) > 0 {
        return &(*items)[0]
    } else {
        return nil
    }
}