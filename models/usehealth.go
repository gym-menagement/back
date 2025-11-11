package models

import (
    "gym/global/config"
    "gym/models/usehealth"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Usehealth struct {
            
    Id                int64 `json:"id"`         
    Order                int64 `json:"order"`         
    Health                int64 `json:"health"`         
    User                int64 `json:"user"`         
    Rocker                int64 `json:"rocker"`         
    Term                int64 `json:"term"`         
    Discount                int64 `json:"discount"`         
    Startday                string `json:"startday"`         
    Endday                string `json:"endday"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type UsehealthManager struct {
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

func (c *Usehealth) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewUsehealthManager(conn *Connection) *UsehealthManager {
    var item UsehealthManager


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

func (p *UsehealthManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *UsehealthManager) SetIndex(index string) {
    p.Index = index
}

func (p *UsehealthManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *UsehealthManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *UsehealthManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *UsehealthManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
    if p.Isolation {
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

func (p *UsehealthManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select uh_id, uh_order, uh_health, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date, o_id, o_membership, o_date, h_id, h_category, h_term, h_name, h_count, h_cost, h_discount, h_costdiscount, h_content, h_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, r_id, r_group, r_name, r_available, r_date, t_id, t_gym, t_daytype, t_name, t_term, t_date, d_id, d_name, d_discount, d_date from usehealth_tb, order_tb, health_tb, user_tb, rocker_tb, term_tb, discount_tb")

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
    
    ret.WriteString("and uh_order = o_id ")
    
    ret.WriteString("and uh_health = h_id ")
    
    ret.WriteString("and uh_user = u_id ")
    
    ret.WriteString("and uh_rocker = r_id ")
    
    ret.WriteString("and uh_term = t_id ")
    
    ret.WriteString("and uh_discount = d_id ")
    

    return ret.String()
}

func (p *UsehealthManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from usehealth_tb")

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
    
    ret.WriteString("and uh_order = o_id ")
    
    ret.WriteString("and uh_health = h_id ")
    
    ret.WriteString("and uh_user = u_id ")
    
    ret.WriteString("and uh_rocker = r_id ")
    
    ret.WriteString("and uh_term = t_id ")
    
    ret.WriteString("and uh_discount = d_id ")
    

    return ret.String()
}

func (p *UsehealthManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select uh_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from usehealth_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and uh_order = o_id ")
    
    ret.WriteString("and uh_health = h_id ")
    
    ret.WriteString("and uh_user = u_id ")
    
    ret.WriteString("and uh_rocker = r_id ")
    
    ret.WriteString("and uh_term = t_id ")
    
    ret.WriteString("and uh_discount = d_id ")
    

    return ret.String()
}

func (p *UsehealthManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate usehealth_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *UsehealthManager) Insert(item *Usehealth) error {
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
        query = "insert into usehealth_tb (uh_id, uh_order, uh_health, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Order, item.Health, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
    } else {
        query = "insert into usehealth_tb (uh_order, uh_health, uh_user, uh_rocker, uh_term, uh_discount, uh_startday, uh_endday, uh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Order, item.Health, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date)
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

func (p *UsehealthManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from usehealth_tb where uh_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *UsehealthManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from usehealth_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *UsehealthManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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

func (p *UsehealthManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from usehealth_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *UsehealthManager) Update(item *Usehealth) error {
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
	

	query := "update usehealth_tb set uh_order = ?, uh_health = ?, uh_user = ?, uh_rocker = ?, uh_term = ?, uh_discount = ?, uh_startday = ?, uh_endday = ?, uh_date = ? where uh_id = ?"
	_, err := p.Exec(query, item.Order, item.Health, item.User, item.Rocker, item.Term, item.Discount, item.Startday, item.Endday, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *UsehealthManager) UpdateWhere(columns []usehealth.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update usehealth_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == usehealth.ColumnId {
        initQuery.WriteString("uh_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnOrder {
        initQuery.WriteString("uh_order = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnHealth {
        initQuery.WriteString("uh_health = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnUser {
        initQuery.WriteString("uh_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnRocker {
        initQuery.WriteString("uh_rocker = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnTerm {
        initQuery.WriteString("uh_term = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnDiscount {
        initQuery.WriteString("uh_discount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnStartday {
        initQuery.WriteString("uh_startday = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnEndday {
        initQuery.WriteString("uh_endday = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealth.ColumnDate {
        initQuery.WriteString("uh_date = ?")
        initParams = append(initParams, v.Value)
        } else {
        
        }
    }

    initQuery.WriteString(" where 1=1 ")

    query, params := p.MakeQuery(initQuery.String(), "", initParams, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

/*


func (p *UsehealthManager) UpdateOrder(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_order = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateHealth(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_health = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_user = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateRocker(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_rocker = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateTerm(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_term = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateDiscount(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_discount = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateStartday(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_startday = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateEndday(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_endday = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealth_tb set uh_date = ? where uh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *UsehealthManager) GetIdentity() int64 {
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

func (p *Usehealth) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *UsehealthManager) ReadRow(rows *sql.Rows) *Usehealth {
    var item Usehealth
    var err error

    var _order Order
    var _health Health
    var _user User
    var _rocker Rocker
    var _term Term
    var _discount Discount
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Order, &item.Health, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date, &_order.Id, &_order.Membership, &_order.Date, &_health.Id, &_health.Category, &_health.Term, &_health.Name, &_health.Count, &_health.Cost, &_health.Discount, &_health.Costdiscount, &_health.Content, &_health.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_rocker.Id, &_rocker.Group, &_rocker.Name, &_rocker.Available, &_rocker.Date, &_term.Id, &_term.Gym, &_term.Daytype, &_term.Name, &_term.Term, &_term.Date, &_discount.Id, &_discount.Name, &_discount.Discount, &_discount.Date)
        
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
        _order.InitExtra()
        item.AddExtra("order",  _order)
_health.InitExtra()
        item.AddExtra("health",  _health)
_user.InitExtra()
        item.AddExtra("user",  _user)
_rocker.InitExtra()
        item.AddExtra("rocker",  _rocker)
_term.InitExtra()
        item.AddExtra("term",  _term)
_discount.InitExtra()
        item.AddExtra("discount",  _discount)

        return &item
    }
}

func (p *UsehealthManager) ReadRows(rows *sql.Rows) []Usehealth {
    var items []Usehealth

    for rows.Next() {
        var item Usehealth
        var _order Order
        var _health Health
        var _user User
        var _rocker Rocker
        var _term Term
        var _discount Discount
        

        err := rows.Scan(&item.Id, &item.Order, &item.Health, &item.User, &item.Rocker, &item.Term, &item.Discount, &item.Startday, &item.Endday, &item.Date, &_order.Id, &_order.Membership, &_order.Date, &_health.Id, &_health.Category, &_health.Term, &_health.Name, &_health.Count, &_health.Cost, &_health.Discount, &_health.Costdiscount, &_health.Content, &_health.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_rocker.Id, &_rocker.Group, &_rocker.Name, &_rocker.Available, &_rocker.Date, &_term.Id, &_term.Gym, &_term.Daytype, &_term.Name, &_term.Term, &_term.Date, &_discount.Id, &_discount.Name, &_discount.Discount, &_discount.Date)
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
        _order.InitExtra()
        item.AddExtra("order",  _order)
_health.InitExtra()
        item.AddExtra("health",  _health)
_user.InitExtra()
        item.AddExtra("user",  _user)
_rocker.InitExtra()
        item.AddExtra("rocker",  _rocker)
_term.InitExtra()
        item.AddExtra("term",  _term)
_discount.InitExtra()
        item.AddExtra("discount",  _discount)

        items = append(items, item)
    }


     return items
}

func (p *UsehealthManager) Get(id int64) *Usehealth {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and uh_id = ?")

    
    query.WriteString(" and uh_order = o_id")
    
    query.WriteString(" and uh_health = h_id")
    
    query.WriteString(" and uh_user = u_id")
    
    query.WriteString(" and uh_rocker = r_id")
    
    query.WriteString(" and uh_term = t_id")
    
    query.WriteString(" and uh_discount = d_id")
    
    
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

func (p *UsehealthManager) GetWhere(args []interface{}) *Usehealth {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *UsehealthManager) Count(args []interface{}) int {
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

func (p *UsehealthManager) FindAll() []Usehealth {
    return p.Find(nil)
}

func (p *UsehealthManager) Find(args []interface{}) []Usehealth {
    if !p.Conn.IsConnect() {
        var items []Usehealth
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
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "uh_" + orderby
                }
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
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "uh_" + orderby
                }
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
        items := make([]Usehealth, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *UsehealthManager) GroupBy(name string, args []interface{}) []Groupby {
    if !p.Conn.IsConnect() {
        var items []Groupby
        return items
    }

    var params []interface{}
    baseQuery := p.GetQueryGroup(name)
    var query strings.Builder
    pos := 1

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
        case Base:
             item := v

             baseQuery = item.Query
        }
    }
    
    query.WriteString(" group by uh_")
    query.WriteString(name)

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        var items []Groupby
        return items
    }

    defer rows.Close()

    var items []Groupby

    for rows.Next() {
        var item Groupby
        err := rows.Scan(&item.Value, &item.Count)
        if err != nil {
           if p.Log {
                log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        items = append(items, item)
    }

    return items
}
