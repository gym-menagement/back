package models

import (
    "gym/global/config"
    
    "gym/models/user"
    
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type User struct {
            
    Id                int64 `json:"id"`         
    Loginid                string `json:"loginid"`         
    Passwd                string `json:"passwd"`         
    Email                string `json:"email"`         
    Name                string `json:"name"`         
    Tel                string `json:"tel"`         
    Address                string `json:"address"`         
    Image                string `json:"image"`         
    Sex                int `json:"sex"`         
    Birth                string `json:"birth"`         
    Type                user.Type `json:"type"`         
    Connectid                string `json:"connectid"`         
    Level                user.Level `json:"level"`         
    Role                user.Role `json:"role"`         
    Use                user.Use `json:"use"`         
    Logindate                string `json:"logindate"`         
    Lastchangepasswddate                string `json:"lastchangepasswddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type UserManager struct {
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

func (c *User) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewUserManager(conn *Connection) *UserManager {
    var item UserManager


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

func (p *UserManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *UserManager) SetIndex(index string) {
    p.Index = index
}

func (p *UserManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *UserManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *UserManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *UserManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *UserManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from user_tb")

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

func (p *UserManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from user_tb")

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

func (p *UserManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate user_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *UserManager) Insert(item *User) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Birth == "" {
       item.Birth = "1000-01-01 00:00:00"
    }
	
    if item.Logindate == "" {
       item.Logindate = "1000-01-01 00:00:00"
    }
	
    if item.Lastchangepasswddate == "" {
       item.Lastchangepasswddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into user_tb (u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Loginid, item.Passwd, item.Email, item.Name, item.Tel, item.Address, item.Image, item.Sex, item.Birth, item.Type, item.Connectid, item.Level, item.Role, item.Use, item.Logindate, item.Lastchangepasswddate, item.Date)
    } else {
        query = "insert into user_tb (u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Loginid, item.Passwd, item.Email, item.Name, item.Tel, item.Address, item.Image, item.Sex, item.Birth, item.Type, item.Connectid, item.Level, item.Role, item.Use, item.Logindate, item.Lastchangepasswddate, item.Date)
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
func (p *UserManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from user_tb where u_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}
func (p *UserManager) Update(item *User) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Birth == "" {
       item.Birth = "1000-01-01 00:00:00"
    }
	
    if item.Logindate == "" {
       item.Logindate = "1000-01-01 00:00:00"
    }
	
    if item.Lastchangepasswddate == "" {
       item.Lastchangepasswddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update user_tb set u_loginid = ?, u_passwd = ?, u_email = ?, u_name = ?, u_tel = ?, u_address = ?, u_image = ?, u_sex = ?, u_birth = ?, u_type = ?, u_connectid = ?, u_level = ?, u_role = ?, u_use = ?, u_logindate = ?, u_lastchangepasswddate = ?, u_date = ? where u_id = ?"
	_, err := p.Exec(query, item.Loginid, item.Passwd, item.Email, item.Name, item.Tel, item.Address, item.Image, item.Sex, item.Birth, item.Type, item.Connectid, item.Level, item.Role, item.Use, item.Logindate, item.Lastchangepasswddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *UserManager) GetIdentity() int64 {
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

func (p *User) InitExtra() {
    p.Extra = map[string]interface{}{
            "level":     user.GetLevel(p.Level),
            "use":     user.GetUse(p.Use),
            "type":     user.GetType(p.Type),
            "role":     user.GetRole(p.Role),

    }
}

func (p *UserManager) ReadRow(rows *sql.Rows) *User {
    var item User
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Email, &item.Name, &item.Tel, &item.Address, &item.Image, &item.Sex, &item.Birth, &item.Type, &item.Connectid, &item.Level, &item.Role, &item.Use, &item.Logindate, &item.Lastchangepasswddate, &item.Date)
        
        if item.Birth == "0000-00-00 00:00:00" || item.Birth == "1000-01-01 00:00:00" || item.Birth == "9999-01-01 00:00:00" {
            item.Birth = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Birth = strings.ReplaceAll(strings.ReplaceAll(item.Birth, "T", " "), "Z", "")
        }
		
        if item.Logindate == "0000-00-00 00:00:00" || item.Logindate == "1000-01-01 00:00:00" || item.Logindate == "9999-01-01 00:00:00" {
            item.Logindate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Logindate = strings.ReplaceAll(strings.ReplaceAll(item.Logindate, "T", " "), "Z", "")
        }
		
        if item.Lastchangepasswddate == "0000-00-00 00:00:00" || item.Lastchangepasswddate == "1000-01-01 00:00:00" || item.Lastchangepasswddate == "9999-01-01 00:00:00" {
            item.Lastchangepasswddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Lastchangepasswddate = strings.ReplaceAll(strings.ReplaceAll(item.Lastchangepasswddate, "T", " "), "Z", "")
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

func (p *UserManager) ReadRows(rows *sql.Rows) []User {
    var items []User

    for rows.Next() {
        var item User
        
    
        err := rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Email, &item.Name, &item.Tel, &item.Address, &item.Image, &item.Sex, &item.Birth, &item.Type, &item.Connectid, &item.Level, &item.Role, &item.Use, &item.Logindate, &item.Lastchangepasswddate, &item.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Birth == "0000-00-00 00:00:00" || item.Birth == "1000-01-01 00:00:00" || item.Birth == "9999-01-01 00:00:00" {
            item.Birth = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Birth = strings.ReplaceAll(strings.ReplaceAll(item.Birth, "T", " "), "Z", "")
        }
		
        if item.Logindate == "0000-00-00 00:00:00" || item.Logindate == "1000-01-01 00:00:00" || item.Logindate == "9999-01-01 00:00:00" {
            item.Logindate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Logindate = strings.ReplaceAll(strings.ReplaceAll(item.Logindate, "T", " "), "Z", "")
        }
		
        if item.Lastchangepasswddate == "0000-00-00 00:00:00" || item.Lastchangepasswddate == "1000-01-01 00:00:00" || item.Lastchangepasswddate == "9999-01-01 00:00:00" {
            item.Lastchangepasswddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Lastchangepasswddate = strings.ReplaceAll(strings.ReplaceAll(item.Lastchangepasswddate, "T", " "), "Z", "")
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

func (p *UserManager) Get(id int64) *User {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and u_id = ?")

    
    
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

func (p *UserManager) GetWhere(args []interface{}) *User {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *UserManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and u_")
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

func (p *UserManager) Count(args []interface{}) int {
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

func (p *UserManager) FindAll() []User {
    return p.Find(nil)
}

func (p *UserManager) Find(args []interface{}) []User {
    if !p.Conn.IsConnect() {
        var items []User
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
                query.WriteString(" and u_")
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
            orderby = "u_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                orderby = "u_" + orderby
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
            orderby = "u_id"
        } else {
            if !strings.Contains(orderby, "_") {
                orderby = "u_" + orderby
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
        var items []User
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}


func (p *UserManager) GetByLoginid(loginid string, args ...interface{}) *User {
    if loginid != "" {
        args = append(args, Where{Column:"loginid", Value:loginid, Compare:"="})        
    }
    
    items := p.Find(args)

    if len(items) > 0 {
        return &items[0]
    } else {
        return nil
    }
}

func (p *UserManager) GetByConnectid(connectid string, args ...interface{}) *User {
    if connectid != "" {
        args = append(args, Where{Column:"connectid", Value:connectid, Compare:"="})        
    }
    
    items := p.Find(args)

    if len(items) > 0 {
        return &items[0]
    } else {
        return nil
    }
}

func (p *UserManager) CountByLoginid(loginid string, args ...interface{}) int {
    rets := make([]interface{}, 0)
    rets = append(rets, args...)
    
    if loginid != "" { 
        rets = append(rets, Where{Column:"loginid", Value:loginid, Compare:"="})
     }
    
    return p.Count(rets)
}

func (p *UserManager) UpdateLogindateById(logindate string, id int64) error {
     if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "update user_tb set u_logindate = ? where 1=1 and u_id = ?"
	_, err := p.Exec(query, logindate, id)

    return err    
}

func (p *UserManager) FindByLevel(level user.Level, args ...interface{}) []User {
    rets := make([]interface{}, 0)
    rets = append(rets, args...)

    if level != 0 { 
        rets = append(rets, Where{Column:"level", Value:level, Compare:"="})
     }
    
    return p.Find(rets)
}



