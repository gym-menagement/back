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


type RockerGroup struct {
	Id					int64 `json:"id"`
	Gym					int64 `json:"gym"`
	Name				string `json:"name"`
	Date				string `json:"date"`

	Extra				map[string]interface{} `json:"extra"`
}

type RockerGroupManager struct {
	Conn	*sql.DB
	Tx		*sql.Tx
	Result	*sql.Result
	Index	string
}

func (c *RockerGroup) AddExtra(key string, value interface{}) {
	c.Extra[key] = value
}

func NewRockerGroupManager(conn interface{}) *RockerGroupManager {
	var item RockerGroupManager

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

func (p *RockerGroupManager) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *RockerGroupManager) SetIndex(index string) {
	p.Index = index
}

func (p *RockerGroupManager) Exec(query string, params ...interface{}) (sql.Result, error) {
	if p.Conn != nil {
		return p.Conn.Exec(query, params...)
	} else {
		return p.Tx.Exec(query, params...)
	}
}

func (p *RockerGroupManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if p.Conn != nil {
		return p.Conn.Query(query, params...)
	} else {
		return p.Tx.Query(query + " FOR UPDATE", params...)
	}
}

func (p *RockerGroupManager) GetQeury() string {
	ret := ""

	str := "select rg_id, rg_gym, rg_name, rg_date from rocker_group_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ")"
	}

	ret += "where 1=1 "

	return ret;
}

func (p *RockerGroupManager) GetQeurySelect() string {
	ret := ""

	str := "select count(*) from rocker_group_tb "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret;
}

func (p *RockerGroupManager) Truncate() error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "truncate rocker_group_tb "
	p.Exec(query)

	return nil
}

func (p *RockerGroupManager) Insert(item *RockerGroup) error {
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
		query = "insert into rocker_group_tb (rg_id, rg_gym, rg_name, rg_date) values (?, ?, ?, ?)"
		res, err = p.Exec(query , item.Id, item.Gym, item.Name, item.Date)
	} else {
		query = "insert into rocker_group_tb (rg_gym, rg_name, rg_date) values (?, ?, ?)"
		res, err = p.Exec(query , item.Gym, item.Name, item.Date)
	}

	if err == nil {
		p.Result = &res
	} else {
		log.Println(err)
		p.Result = nil
	}

	return err
}

func (p *RockerGroupManager) Delete(id int64) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "delete from rocker_group_tb where rg_id = ?"
	_, err := p.Exec(query, id)

	return err
}

func (p *RockerGroupManager) Update(item *RockerGroup) error {
	if p.Conn == nil && p.Tx == nil {
		return errors.New("Connection Error")
	}

	query := "update rocker_group_tb set rg_gym = ?, rg_name = ?, rg_date = ? where rg_id = ?"
	_, err := p.Exec(query, item.Gym, item.Name, item.Date, item.Id)

	return err
}

func (p *RockerGroupManager) GetIdentity() int64 {
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

func (p *RockerGroup) InitExtra() {
	p.Extra = map[string]interface{}{

	}
}

func (p *RockerGroupManager) ReadRow(rows *sql.Rows) *RockerGroup {
	var item RockerGroup
	var err error

	if rows.Next() {
		err = rows.Scan(&item.Id, &item.Gym, &item.Name, &item.Date)
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

func (p *RockerGroupManager) ReadRows(rows *sql.Rows) *[]RockerGroup {
	var items []RockerGroup

	for rows.Next() {
		var item RockerGroup

		err := rows.Scan(&item.Id, &item.Gym, &item.Name, &item.Date)

		if err != nil {
			log.Printf("ReadRows error : %v\n", err)
			break
		}

		item.InitExtra()

		items = append(items, item)
	}
	return &items
}

func (p *RockerGroupManager) Get(id int64) *RockerGroup {
	if p.Conn == nil && p.Tx == nil {
		return nil
	}

	query := p.GetQeury() + " and rg_id = ?"

	rows, err := p.Query(query, id)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	return p.ReadRow(rows)
}

func (p *RockerGroupManager) Count(args []interface{}) int {
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
				query += " and rg_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and rg_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and rg_" + item.Column + " " + item.Compare + " ?"
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

func (p *RockerGroupManager) Find(args []interface{}) *[]RockerGroup {
	if p.Conn == nil && p.Tx == nil {
		var items []RockerGroup
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
				query += " and rg_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
			} else if item.Compare == "between" {
				query += " and rg_" + item.Column + " between ? and ?"

				s := item.Value.([2]string)
				params = append(params, s[0])
				params = append(params, s[1])
			} else {
				query += " and rg_" + item.Column + " " + item.Compare + " ?"
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
			orderby = "rg_id"
		} else {
			orderby = "rg_" + orderby
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
			orderby = "rg_id"
		} else {
			orderby = "rg_" + orderby
		}
		query += " order by " + orderby
	}

	rows, err := p.Query(query, params...)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		var items []RockerGroup
		return &items
	}

	defer rows.Close()

	return p.ReadRows(rows)
}

func (p *RockerGroupManager) GetByName(loginid string, args ...interface{}) *RockerGroup {
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