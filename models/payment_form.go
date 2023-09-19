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


type PaymentForm struct {
	Id					int64 `json:"id"`
	Gym					int64 `json:"gym"`
	Payment				int64 `json:"payment"`
	Type				int64 `json:"type"`
	Cost				int `json:"cost"`
	Date				string `json:"date"`

	Extra				map[string]interface{} `json:"extra"`
}

type PaymentFormManager struct {
	Conn	*sql.DB
	Tx		*sql.Tx
	Result	*sql.Result
	Index	string
}

func (c *PaymentForm) AddExtra(key string, value interface{}) {
	c.Extra[key] = value
}

func NewPaymentFormManager(conn interface{}) *PaymentFormManager {
	var item PaymentFormManager

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

func (p *PaymentFormManager) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *PaymentFormManager) SetIndex(index string) {
	p.Index = index
}

func (p *PaymentFormManager) Exec(query string, params ...interface{}) (sql.Result, error) {
	if p.Conn != nil {
		return p.Conn.Exec(query, params...)
	} else {
		return p.Tx.Exec(query, params...)
	}
}

func (p *PaymentFormManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if p.Conn != nil {
		return p.Conn.Query(query, params...)
	} else {
		return p.Tx.Query(query + " FOR UPDATE", params...)
	}
}

func (p *PaymentFormManager) GetQeury() string {
	ret := ""

	str := "select pf_id, pf_gym, pf_payment, pf_type, pf_cost, pf_date from payment_form_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ")"
	}

	ret += "where 1=1 "

	return ret;
}

func (p *PaymentFormManager) GetQeurySelect() string {
	ret := ""

	str := "select count(*) from payment_form_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret;
}

func (p *PaymentFormManager) Truncate() error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "truncate payment_form_tb "
	p.Exec(query)

	return nil
}

func (p *PaymentFormManager) Insert(item *PaymentForm) error {
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
		query = "insert into payment_form_tb (pf_id, pf_gym, pf_payment, pf_type, pf_cost, pf_date) values (?, ?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Id, item.Gym, item.Payment, item.Type, item.Cost, item.Date)
	} else {
		query = "insert into payment_form_tb (pf_gym, pf_payment, pf_type, pf_cost, pf_date) values (?, ?, ?, ?, ?)"
		res, err = p.Exec(query , item.Gym, item.Payment, item.Type, item.Cost, item.Date)
	}

	if err == nil {
		p.Result = &res
	} else {
		log.Println(err)
		p.Result = nil
	}

	return err
}

func (p *PaymentFormManager) Delete(id int64) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "delete from payment_form_tb where pf_id = ?"
	_, err := p.Exec(query, id)

	return err
}

func (p *PaymentFormManager) Update(item *PaymentForm) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "update payment_form_tb set pf_gym = ?, pf_payment = ?, pf_type = ?, pf_cost = ?, pf_date = ? where pf_id = ?"
	_, err := p.Exec(query, item.Gym, item.Payment, item.Date, item.Type, item.Cost, item.Id)

	return err
}

func (p *PaymentFormManager) GetIdentity() int64 {
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

func (p *PaymentForm) InitExtra() {
	p.Extra = map[string]interface{}{

	}
}

func (p *PaymentFormManager) ReadRow(rows *sql.Rows) *PaymentForm {
	var item PaymentForm
	var err error

	if rows.Next() {
		err = rows.Scan(&item.Id, &item.Gym, &item.Payment, &item.Type, &item.Cost, &item.Date)
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

func (p *PaymentFormManager) ReadRows(rows *sql.Rows) *[]PaymentForm {
	var items []PaymentForm

	for rows.Next() {
		var item PaymentForm

		err := rows.Scan(&item.Id, &item.Gym, &item.Payment, &item.Type, &item.Cost, &item.Date)

		if err != nil {
			log.Printf("ReadRows error : %v\n", err)
			break
		}

		item.InitExtra()

		items = append(items, item)
	}
	return &items
}

func (p *PaymentFormManager) Get(id int64) *PaymentForm {
	if p.Conn == nil && p.Tx == nil {
		return nil
	}

	query := p.GetQeury() + " and pf_id = ?"

	rows, err := p.Query(query, id)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	return p.ReadRow(rows)
}

func (p *PaymentFormManager) Count(args []interface{}) int {
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
				query += " and pf_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and pf_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and pf_" + item.Column + " " + item.Compare + " ?"
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

func (p *PaymentFormManager) Find(args []interface{}) *[]PaymentForm {
	if p.Conn == nil && p.Tx == nil {
		var items []PaymentForm
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
				query += " and pf_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and pf_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and pf_" + item.Column + " " + item.Compare + " ?"
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
			orderby = "pf_id"
		} else {
			orderby = "pf_" + orderby
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
			orderby = "pf_id"
		} else {
			orderby = "pf_" + orderby
		}
		query += " order by " + orderby
	}

	rows, err := p.Query(query, params...)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		var items []PaymentForm
		return &items
	}

	defer rows.Close()

	return p.ReadRows(rows)
}

func (p *PaymentFormManager) GetByPayment(loginid string, args ...interface{}) *PaymentForm {
    if loginid != "" {
        args = append(args, Where{Column:"payment", Value:loginid, Compare:"="})        
    }
    
    items := p.Find(args)

    if items != nil && len(*items) > 0 {
        return &(*items)[0]
    } else {
        return nil
    }
}