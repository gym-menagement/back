package models

import (
    "gym/global/config"
    "gym/models/memberqr"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Memberqr struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Code                string `json:"code"`         
    Imageurl                string `json:"imageurl"`         
    Isactive                memberqr.Isactive `json:"isactive"`         
    Expiredate                string `json:"expiredate"`         
    Generateddate                string `json:"generateddate"`         
    Lastuseddate                string `json:"lastuseddate"`         
    Usecount                int `json:"usecount"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type MemberqrManager struct {
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

func (c *Memberqr) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewMemberqrManager(conn *Connection) *MemberqrManager {
    var item MemberqrManager


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

func (p *MemberqrManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *MemberqrManager) SetIndex(index string) {
    p.Index = index
}

func (p *MemberqrManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *MemberqrManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *MemberqrManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *MemberqrManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *MemberqrManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select mq_id, mq_user, mq_code, mq_imageurl, mq_isactive, mq_expiredate, mq_generateddate, mq_lastuseddate, mq_usecount, mq_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from memberqr_tb, user_tb")

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
    
    ret.WriteString("and mq_user = u_id ")
    

    return ret.String()
}

func (p *MemberqrManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from memberqr_tb")

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
    
    ret.WriteString("and mq_user = u_id ")
    

    return ret.String()
}

func (p *MemberqrManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select mq_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from memberqr_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and mq_user = u_id ")
    

    return ret.String()
}

func (p *MemberqrManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate memberqr_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *MemberqrManager) Insert(item *Memberqr) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Expiredate == "" {
       item.Expiredate = "1000-01-01 00:00:00"
    }
	
    if item.Generateddate == "" {
       item.Generateddate = "1000-01-01 00:00:00"
    }
	
    if item.Lastuseddate == "" {
       item.Lastuseddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into memberqr_tb (mq_id, mq_user, mq_code, mq_imageurl, mq_isactive, mq_expiredate, mq_generateddate, mq_lastuseddate, mq_usecount, mq_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Code, item.Imageurl, item.Isactive, item.Expiredate, item.Generateddate, item.Lastuseddate, item.Usecount, item.Date)
    } else {
        query = "insert into memberqr_tb (mq_user, mq_code, mq_imageurl, mq_isactive, mq_expiredate, mq_generateddate, mq_lastuseddate, mq_usecount, mq_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Code, item.Imageurl, item.Isactive, item.Expiredate, item.Generateddate, item.Lastuseddate, item.Usecount, item.Date)
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

func (p *MemberqrManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from memberqr_tb where mq_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *MemberqrManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from memberqr_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MemberqrManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and mq_")
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

func (p *MemberqrManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from memberqr_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MemberqrManager) Update(item *Memberqr) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Expiredate == "" {
       item.Expiredate = "1000-01-01 00:00:00"
    }
	
    if item.Generateddate == "" {
       item.Generateddate = "1000-01-01 00:00:00"
    }
	
    if item.Lastuseddate == "" {
       item.Lastuseddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update memberqr_tb set mq_user = ?, mq_code = ?, mq_imageurl = ?, mq_isactive = ?, mq_expiredate = ?, mq_generateddate = ?, mq_lastuseddate = ?, mq_usecount = ?, mq_date = ? where mq_id = ?"
	_, err := p.Exec(query, item.User, item.Code, item.Imageurl, item.Isactive, item.Expiredate, item.Generateddate, item.Lastuseddate, item.Usecount, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *MemberqrManager) UpdateWhere(columns []memberqr.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update memberqr_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == memberqr.ColumnId {
        initQuery.WriteString("mq_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnUser {
        initQuery.WriteString("mq_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnCode {
        initQuery.WriteString("mq_code = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnImageurl {
        initQuery.WriteString("mq_imageurl = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnIsactive {
        initQuery.WriteString("mq_isactive = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnExpiredate {
        initQuery.WriteString("mq_expiredate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnGenerateddate {
        initQuery.WriteString("mq_generateddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnLastuseddate {
        initQuery.WriteString("mq_lastuseddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnUsecount {
        initQuery.WriteString("mq_usecount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == memberqr.ColumnDate {
        initQuery.WriteString("mq_date = ?")
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


func (p *MemberqrManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_user = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateCode(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_code = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateImageurl(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_imageurl = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateIsactive(value memberqr.Isactive, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_isactive = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateExpiredate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_expiredate = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateGenerateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_generateddate = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateLastuseddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_lastuseddate = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateUsecount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_usecount = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MemberqrManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update memberqr_tb set mq_date = ? where mq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *MemberqrManager) GetIdentity() int64 {
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

func (p *Memberqr) InitExtra() {
    p.Extra = map[string]interface{}{
            "isactive":     memberqr.GetIsactive(p.Isactive),

    }
}

func (p *MemberqrManager) ReadRow(rows *sql.Rows) *Memberqr {
    var item Memberqr
    var err error

    var _user User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Code, &item.Imageurl, &item.Isactive, &item.Expiredate, &item.Generateddate, &item.Lastuseddate, &item.Usecount, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
        
        if item.Expiredate == "0000-00-00 00:00:00" || item.Expiredate == "1000-01-01 00:00:00" || item.Expiredate == "9999-01-01 00:00:00" {
            item.Expiredate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Expiredate = strings.ReplaceAll(strings.ReplaceAll(item.Expiredate, "T", " "), "Z", "")
        }
		
        if item.Generateddate == "0000-00-00 00:00:00" || item.Generateddate == "1000-01-01 00:00:00" || item.Generateddate == "9999-01-01 00:00:00" {
            item.Generateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Generateddate = strings.ReplaceAll(strings.ReplaceAll(item.Generateddate, "T", " "), "Z", "")
        }
		
        if item.Lastuseddate == "0000-00-00 00:00:00" || item.Lastuseddate == "1000-01-01 00:00:00" || item.Lastuseddate == "9999-01-01 00:00:00" {
            item.Lastuseddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Lastuseddate = strings.ReplaceAll(strings.ReplaceAll(item.Lastuseddate, "T", " "), "Z", "")
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
        _user.InitExtra()
        item.AddExtra("user",  _user)

        return &item
    }
}

func (p *MemberqrManager) ReadRows(rows *sql.Rows) []Memberqr {
    var items []Memberqr

    for rows.Next() {
        var item Memberqr
        var _user User
        

        err := rows.Scan(&item.Id, &item.User, &item.Code, &item.Imageurl, &item.Isactive, &item.Expiredate, &item.Generateddate, &item.Lastuseddate, &item.Usecount, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Expiredate == "0000-00-00 00:00:00" || item.Expiredate == "1000-01-01 00:00:00" || item.Expiredate == "9999-01-01 00:00:00" {
            item.Expiredate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Expiredate = strings.ReplaceAll(strings.ReplaceAll(item.Expiredate, "T", " "), "Z", "")
        }
		
        if item.Generateddate == "0000-00-00 00:00:00" || item.Generateddate == "1000-01-01 00:00:00" || item.Generateddate == "9999-01-01 00:00:00" {
            item.Generateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Generateddate = strings.ReplaceAll(strings.ReplaceAll(item.Generateddate, "T", " "), "Z", "")
        }
		
        if item.Lastuseddate == "0000-00-00 00:00:00" || item.Lastuseddate == "1000-01-01 00:00:00" || item.Lastuseddate == "9999-01-01 00:00:00" {
            item.Lastuseddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Lastuseddate = strings.ReplaceAll(strings.ReplaceAll(item.Lastuseddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _user.InitExtra()
        item.AddExtra("user",  _user)

        items = append(items, item)
    }


     return items
}

func (p *MemberqrManager) Get(id int64) *Memberqr {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and mq_id = ?")

    
    query.WriteString(" and mq_user = u_id")
    
    
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

func (p *MemberqrManager) GetWhere(args []interface{}) *Memberqr {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *MemberqrManager) Count(args []interface{}) int {
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

func (p *MemberqrManager) FindAll() []Memberqr {
    return p.Find(nil)
}

func (p *MemberqrManager) Find(args []interface{}) []Memberqr {
    if !p.Conn.IsConnect() {
        var items []Memberqr
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
                query.WriteString(" and mq_")
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
            orderby = "mq_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mq_" + orderby
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
            orderby = "mq_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mq_" + orderby
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
        items := make([]Memberqr, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *MemberqrManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and mq_")
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
    
    query.WriteString(" group by mq_")
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
