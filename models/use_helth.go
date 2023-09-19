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


type UseHelth struct {
	Id					int64 `json:"id"`
	Order				int64 `json:"order"`
	Helth				int64 `json:"helth"`
	User				int64 `json:"user"`
	Rocker				int64 `json:"rocker"`
	Term				int64 `json:"term"`
	Discount			int64 `json:"discount"`
	Startday			string `json:"startday"`
	Endday				string `'json:"endday"`
	Date				string `json:"date"`

	Extra				map[string]interface{} `json:"extra"`
}

type UseHelthManager struct {
	Conn	*sql.DB
	Tx		*sql.Tx
	Result	*sql.Result
	Index	string
}

func (c *UseHelth) AddExtra(key string, value interface{}) {
	c.Extra[key] = value
}

func NewUseHelthManager(conn interface{}) *UseHelthManager {
	var item UseHelthManager

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

func (p *UseHelthManager) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *UseHelthManager) SetIndex(index string) {
	p.Index = index
}

func (p *UseHelthManager) Exec(query string, params ...interface{}) (sql.Result, error) {
	if p.Conn != nil {
		return p.Conn.Exec(query, params...)
	} else {
		return p.Tx.Exec(query, params...)
	}
}

func (p *UseHelthManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if p.Conn != nil {
		return p.Conn.Query(query, params...)
	} else {
		return p.Tx.Query(query + " FOR UPDATE", params...)
	}
}

func (p *UseHelthManager) GetQeury() string {
	ret := ""

	str := "select uh_id, uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date from use_helth_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ")"
	}

	ret += "where 1=1 "

	return ret;
}

func (p *UseHelthManager) GetQeurySelect() string {
	ret := ""

	str := "select count(*) from use_helth_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret;
}

func (p *UseHelthManager) Truncate() error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "truncate use_helth_tb "
	p.Exec(query)

	return nil
}

func (p *UseHelthManager) Insert(item *UseHelth) error {
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
		query = "insert into use_helth_tb (uh_id, uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Id, item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
	} else {
		query = "insert into use_helth_tb (uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
	}

	if err == nil {
		p.Result = &res
	} else {
		log.Println(err)
		p.Result = nil
	}

	return err
}

func (p *UseHelthManager) Delete(id int64) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "delete from use_helth_tb where uh_id = ?"
	_, err := p.Exec(query, id)

	return err
}

func (p *UseHelthManager) Update(item *UseHelth) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "update use_helth_tb set uh_order = ?, uh_helth = ?, uh_user = ?, uh_rocker = ?, uh_term = ?, uh_discount = ?, uh_startday = ?, uh_endday = ?, uh_date = ? where uh_id = ?"
	_, err := p.Exec(query, item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date, item.Id)

	return err
}

func (p *UseHelthManager) GetIdentity() int64 {
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

func (p *UseHelth) InitExtra() {
	p.Extra = map[string]interface{}{

	}
}

func (p *UseHelthManager) ReadRow(rows *sql.Rows) *UseHelth {
	var item UseHelth
	var err error

	if rows.Next() {
		err = rows.Scan(&item.Id, &item.Order, &item.Helth, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date)
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

func (p *UseHelthManager) ReadRows(rows *sql.Rows) *[]UseHelth {
	var items []UseHelth

	for rows.Next() {
		var item UseHelth

		err := rows.Scan(&item.Id, &item.Order, &item.Helth, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date)

		if err != nil {
			log.Printf("ReadRows error : %v\n", err)
			break
		}

		item.InitExtra()

		items = append(items, item)
	}
	return &items
}

func (p *UseHelthManager) Get(id int64) *UseHelth {
	if p.Conn == nil && p.Tx == nil {
		return nil
	}

	query := p.GetQeury() + " and uh_id = ?"

	rows, err := p.Query(query, id)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	return p.ReadRow(rows)
}

func (p *UseHelthManager) Count(args []interface{}) int {
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
				query += " and uh_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and uh_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and uh_" + item.Column + " " + item.Compare + " ?"
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

func (p *UseHelthManager) Find(args []interface{}) *[]UseHelth {
	if p.Conn == nil && p.Tx == nil {
		var items []UseHelth
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
				query += " and uh_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and uh_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and uh_" + item.Column + " " + item.Compare + " ?"
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
			orderby = "uh_id"
		} else {
			orderby = "uh_" + orderby
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
			orderby = "uh_id"
		} else {
			orderby = "uh_" + orderby
		}
		query += " order by " + orderby
	}

	rows, err := p.Query(query, params...)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		var items []UseHelth
		return &items
	}

	defer rows.Close()

	return p.ReadRows(rows)
}

func (p *UseHelthManager) GetByName(loginid string, args ...interface{}) *UseHelth {
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