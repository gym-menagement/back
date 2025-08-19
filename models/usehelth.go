package models

import (
    "gym/global/config"
    
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Usehelth struct {
            
    Id                int64 `json:"id"`         
    Order                int64 `json:"order"`         
    Helth                int64 `json:"helth"`         
    User                int64 `json:"user"`         
    Rocker                int64 `json:"rocker"`         
    Term                int64 `json:"term"`         
    Discount                int64 `json:"discount"`         
    Startday                string `json:"startday"`         
    Endday                string `json:"endday"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type UsehelthManager struct {
    Conn    *Connection
    Result  *sql.Result
    Index   string
    Isolation   bool
    SelectQuery  string
    JoinQuery string
    CountQuery   string
    GroupQuery string
    SelectLog bool
    Log bool
}

func (c *Usehelth) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewUsehelthManager(conn *Connection) *UsehelthManager {
    var item UsehelthManager


    if conn == nil {
        item.Conn = NewConnection()
        item.Isolation = false
    } else {
        item.Conn = conn 
        item.Isolation = conn.Isolation
    }

    item.Index = ""
    item.SelectLog = config.Log.Database
    item.Log = config.Log.Database

    return &item
}

func (p *UsehelthManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *UsehelthManager) SetIndex(index string) {
    p.Index = index
}

func (p *UsehelthManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *UsehelthManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *UsehelthManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *UsehelthManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
    if p.Isolation == true {
        query += " for update"
    }

    if p.SelectLog {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Query(query, params...)
}

func (p *UsehelthManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select uh_id, uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date from usehelth_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    

    return ret.String()
}

func (p *UsehelthManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from usehelth_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    

    return ret.String()
}

func (p *UsehelthManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate usehelth_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *UsehelthManager) Insert(item *Usehelth) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Startday == "" {
       item.Startday = "1000-01-01 00:00:00"
    }
	
    if item.Endday == "" {
       item.Endday = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into usehelth_tb (uh_id, uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
    } else {
        query = "insert into usehelth_tb (uh_order, uh_helth, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
    }
    
    if err == nil {
        p.Result = &res
        
    } else {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        p.Result = nil
    }

    return err
}
func (p *UsehelthManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from usehelth_tb where uh_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}
func (p *UsehelthManager) Update(item *Usehelth) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Startday == "" {
       item.Startday = "1000-01-01 00:00:00"
    }
	
    if item.Endday == "" {
       item.Endday = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update usehelth_tb set uh_order = ?, uh_helth = ?, uh_user = ?, uh_rocker = ?, uh_term = ?, uh_discount = ?, uh_startday = ?, uh_endday = ?, uh_date = ? where uh_id = ?"
	_, err := p.Exec(query, item.Order, item.Helth, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *UsehelthManager) GetIdentity() int64 {
    if !p.Conn.IsConnect() {
        return 0
    }

    id, err := (*p.Result).LastInsertId()

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return 0
    } else {
        return id
    }
}

func (p *Usehelth) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *UsehelthManager) ReadRow(rows *sql.Rows) *Usehelth {
    var item Usehelth
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Order, &item.Helth, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date)
        
        if item.Startday == "0000-00-00 00:00:00" || item.Startday == "1000-01-01 00:00:00" || item.Startday == "9999-01-01 00:00:00" {
            item.Startday = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startday = strings.ReplaceAll(strings.ReplaceAll(item.Startday, "T", " "), "Z", "")
        }
		
        if item.Endday == "0000-00-00 00:00:00" || item.Endday == "1000-01-01 00:00:00" || item.Endday == "9999-01-01 00:00:00" {
            item.Endday = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Endday = strings.ReplaceAll(strings.ReplaceAll(item.Endday, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

    } else {
        return nil
    }

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return nil
    } else {

        item.InitExtra()
        
        return &item
    }
}

func (p *UsehelthManager) ReadRows(rows *sql.Rows) []Usehelth {
    var items []Usehelth

    for rows.Next() {
        var item Usehelth
        
    
        err := rows.Scan(&item.Id, &item.Order, &item.Helth, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Startday == "0000-00-00 00:00:00" || item.Startday == "1000-01-01 00:00:00" || item.Startday == "9999-01-01 00:00:00" {
            item.Startday = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startday = strings.ReplaceAll(strings.ReplaceAll(item.Startday, "T", " "), "Z", "")
        }
		
        if item.Endday == "0000-00-00 00:00:00" || item.Endday == "1000-01-01 00:00:00" || item.Endday == "9999-01-01 00:00:00" {
            item.Endday = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Endday = strings.ReplaceAll(strings.ReplaceAll(item.Endday, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		
        
        item.InitExtra()        
        
        items = append(items, item)
    }


     return items
}

func (p *UsehelthManager) Get(id int64) *Usehelth {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and uh_id = ?")

    
    
    rows, err := p.Query(query.String(), id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
       return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *UsehelthManager) GetWhere(args []interface{}) *Usehelth {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *UsehelthManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
    var params []interface{}
    if initParams != nil {
        params = append(params, initParams...)
    }

    pos := 1

    var query strings.Builder
	query.WriteString(initQuery)

    for _, arg := range args {
        switch v := arg.(type) {        
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and uh_")
            }
            query.WriteString(item.Column)

            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        }        
    }

	query.WriteString(postQuery)

    return query.String(), params
}

func (p *UsehelthManager) Count(args []interface{}) int {
    if !p.Conn.IsConnect() {
        return 0
    }

    query, params := p.MakeQuery(p.GetQuerySelect(), p.GroupQuery, nil, args)
    rows, err := p.Query(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
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

func (p *UsehelthManager) FindAll() []Usehelth {
    return p.Find(nil)
}

func (p *UsehelthManager) Find(args []interface{}) []Usehelth {
    if !p.Conn.IsConnect() {
        var items []Usehelth
        return items
    }

    var params []interface{}
    baseQuery := p.GetQuery()

    var query strings.Builder

    page := 0
    pagesize := 0
    orderby := ""

    pos := 1
    
    for _, arg := range args {
        switch v := arg.(type) {
        case PagingType:
            item := v
            page = item.Page
            pagesize = item.Pagesize            
        case OrderingType:
            item := v
            orderby = item.Order
        case LimitType:
            item := v
            page = 1
            pagesize = item.Limit
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
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and uh_")
            }
            query.WriteString(item.Column)
            
            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        case Base:
             item := v

             baseQuery = item.Query
        }
    }

    query.WriteString(p.GroupQuery)
    
    startpage := (page - 1) * pagesize
    
    if page > 0 && pagesize > 0 {
        if orderby == "" {
            orderby = "uh_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                orderby = "uh_" + orderby
            }
            
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
        if config.Database.Type == config.Postgresql {
            query.WriteString(fmt.Sprintf(" limit $%v offset $%v", pos, pos + 1))
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Mysql {
            query.WriteString(" limit ? offset ?")
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Sqlserver {
            query.WriteString("OFFSET ? ROWS FETCH NEXT ? ROWS ONLY")
            params = append(params, startpage)
            params = append(params, pagesize)
        }
    } else {
        if orderby == "" {
            orderby = "uh_id"
        } else {
            if !strings.Contains(orderby, "_") {
                orderby = "uh_" + orderby
            }
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
    }

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        var items []Usehelth
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}




