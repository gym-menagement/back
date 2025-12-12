package models

import (
    "gym/global/config"
    "gym/models/inquiry"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Inquiry struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Gym                int64 `json:"gym"`         
    Type                inquiry.Type `json:"type"`         
    Title                string `json:"title"`         
    Content                string `json:"content"`         
    Status                inquiry.Status `json:"status"`         
    Answer                string `json:"answer"`         
    Answeredby                int64 `json:"answeredby"`         
    Answereddate                string `json:"answereddate"`         
    Createddate                string `json:"createddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type InquiryManager struct {
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

func (c *Inquiry) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewInquiryManager(conn *Connection) *InquiryManager {
    var item InquiryManager


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

func (p *InquiryManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *InquiryManager) SetIndex(index string) {
    p.Index = index
}

func (p *InquiryManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *InquiryManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *InquiryManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *InquiryManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *InquiryManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select iq_id, iq_user, iq_gym, iq_type, iq_title, iq_content, iq_status, iq_answer, iq_answeredby, iq_answereddate, iq_createddate, iq_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, g_id, g_name, g_address, g_tel, g_user, g_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from inquiry_tb, user_tb, gym_tb, user_tb")

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
    
    ret.WriteString("and iq_user = u_id ")
    
    ret.WriteString("and iq_gym = g_id ")
    
    ret.WriteString("and iq_answeredby = u_id ")
    

    return ret.String()
}

func (p *InquiryManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from inquiry_tb")

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
    
    ret.WriteString("and iq_user = u_id ")
    
    ret.WriteString("and iq_gym = g_id ")
    
    ret.WriteString("and iq_answeredby = u_id ")
    

    return ret.String()
}

func (p *InquiryManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select iq_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from inquiry_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and iq_user = u_id ")
    
    ret.WriteString("and iq_gym = g_id ")
    
    ret.WriteString("and iq_answeredby = u_id ")
    

    return ret.String()
}

func (p *InquiryManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate inquiry_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *InquiryManager) Insert(item *Inquiry) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Answereddate == "" {
       item.Answereddate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into inquiry_tb (iq_id, iq_user, iq_gym, iq_type, iq_title, iq_content, iq_status, iq_answer, iq_answeredby, iq_answereddate, iq_createddate, iq_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Gym, item.Type, item.Title, item.Content, item.Status, item.Answer, item.Answeredby, item.Answereddate, item.Createddate, item.Date)
    } else {
        query = "insert into inquiry_tb (iq_user, iq_gym, iq_type, iq_title, iq_content, iq_status, iq_answer, iq_answeredby, iq_answereddate, iq_createddate, iq_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Gym, item.Type, item.Title, item.Content, item.Status, item.Answer, item.Answeredby, item.Answereddate, item.Createddate, item.Date)
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

func (p *InquiryManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from inquiry_tb where iq_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *InquiryManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from inquiry_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *InquiryManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and iq_")
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

func (p *InquiryManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from inquiry_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *InquiryManager) Update(item *Inquiry) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Answereddate == "" {
       item.Answereddate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update inquiry_tb set iq_user = ?, iq_gym = ?, iq_type = ?, iq_title = ?, iq_content = ?, iq_status = ?, iq_answer = ?, iq_answeredby = ?, iq_answereddate = ?, iq_createddate = ?, iq_date = ? where iq_id = ?"
	_, err := p.Exec(query, item.User, item.Gym, item.Type, item.Title, item.Content, item.Status, item.Answer, item.Answeredby, item.Answereddate, item.Createddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *InquiryManager) UpdateWhere(columns []inquiry.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update inquiry_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == inquiry.ColumnId {
        initQuery.WriteString("iq_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnUser {
        initQuery.WriteString("iq_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnGym {
        initQuery.WriteString("iq_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnType {
        initQuery.WriteString("iq_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnTitle {
        initQuery.WriteString("iq_title = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnContent {
        initQuery.WriteString("iq_content = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnStatus {
        initQuery.WriteString("iq_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnAnswer {
        initQuery.WriteString("iq_answer = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnAnsweredby {
        initQuery.WriteString("iq_answeredby = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnAnswereddate {
        initQuery.WriteString("iq_answereddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnCreateddate {
        initQuery.WriteString("iq_createddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == inquiry.ColumnDate {
        initQuery.WriteString("iq_date = ?")
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


func (p *InquiryManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_user = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_gym = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateType(value inquiry.Type, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_type = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateTitle(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_title = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateContent(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_content = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateStatus(value inquiry.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_status = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateAnswer(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_answer = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateAnsweredby(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_answeredby = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateAnswereddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_answereddate = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateCreateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_createddate = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *InquiryManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update inquiry_tb set iq_date = ? where iq_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *InquiryManager) GetIdentity() int64 {
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

func (p *Inquiry) InitExtra() {
    p.Extra = map[string]interface{}{
            "type":     inquiry.GetType(p.Type),
            "status":     inquiry.GetStatus(p.Status),

    }
}

func (p *InquiryManager) ReadRow(rows *sql.Rows) *Inquiry {
    var item Inquiry
    var err error

    var _inquireruser User
    var _gym Gym
    var _answeredbyuser User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Gym, &item.Type, &item.Title, &item.Content, &item.Status, &item.Answer, &item.Answeredby, &item.Answereddate, &item.Createddate, &item.Date, &_inquireruser.Id, &_inquireruser.Loginid, &_inquireruser.Passwd, &_inquireruser.Email, &_inquireruser.Name, &_inquireruser.Tel, &_inquireruser.Address, &_inquireruser.Image, &_inquireruser.Sex, &_inquireruser.Birth, &_inquireruser.Type, &_inquireruser.Connectid, &_inquireruser.Level, &_inquireruser.Role, &_inquireruser.Use, &_inquireruser.Logindate, &_inquireruser.Lastchangepasswddate, &_inquireruser.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_answeredbyuser.Id, &_answeredbyuser.Loginid, &_answeredbyuser.Passwd, &_answeredbyuser.Email, &_answeredbyuser.Name, &_answeredbyuser.Tel, &_answeredbyuser.Address, &_answeredbyuser.Image, &_answeredbyuser.Sex, &_answeredbyuser.Birth, &_answeredbyuser.Type, &_answeredbyuser.Connectid, &_answeredbyuser.Level, &_answeredbyuser.Role, &_answeredbyuser.Use, &_answeredbyuser.Logindate, &_answeredbyuser.Lastchangepasswddate, &_answeredbyuser.Date)
        
        if item.Answereddate == "0000-00-00 00:00:00" || item.Answereddate == "1000-01-01 00:00:00" || item.Answereddate == "9999-01-01 00:00:00" {
            item.Answereddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Answereddate = strings.ReplaceAll(strings.ReplaceAll(item.Answereddate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
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
        _inquireruser.InitExtra()
        item.AddExtra("inquireruser",  _inquireruser)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)
_answeredbyuser.InitExtra()
        item.AddExtra("answeredbyuser",  _answeredbyuser)

        return &item
    }
}

func (p *InquiryManager) ReadRows(rows *sql.Rows) []Inquiry {
    var items []Inquiry

    for rows.Next() {
        var item Inquiry
        var _inquireruser User
        var _gym Gym
        var _answeredbyuser User
        

        err := rows.Scan(&item.Id, &item.User, &item.Gym, &item.Type, &item.Title, &item.Content, &item.Status, &item.Answer, &item.Answeredby, &item.Answereddate, &item.Createddate, &item.Date, &_inquireruser.Id, &_inquireruser.Loginid, &_inquireruser.Passwd, &_inquireruser.Email, &_inquireruser.Name, &_inquireruser.Tel, &_inquireruser.Address, &_inquireruser.Image, &_inquireruser.Sex, &_inquireruser.Birth, &_inquireruser.Type, &_inquireruser.Connectid, &_inquireruser.Level, &_inquireruser.Role, &_inquireruser.Use, &_inquireruser.Logindate, &_inquireruser.Lastchangepasswddate, &_inquireruser.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_answeredbyuser.Id, &_answeredbyuser.Loginid, &_answeredbyuser.Passwd, &_answeredbyuser.Email, &_answeredbyuser.Name, &_answeredbyuser.Tel, &_answeredbyuser.Address, &_answeredbyuser.Image, &_answeredbyuser.Sex, &_answeredbyuser.Birth, &_answeredbyuser.Type, &_answeredbyuser.Connectid, &_answeredbyuser.Level, &_answeredbyuser.Role, &_answeredbyuser.Use, &_answeredbyuser.Logindate, &_answeredbyuser.Lastchangepasswddate, &_answeredbyuser.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Answereddate == "0000-00-00 00:00:00" || item.Answereddate == "1000-01-01 00:00:00" || item.Answereddate == "9999-01-01 00:00:00" {
            item.Answereddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Answereddate = strings.ReplaceAll(strings.ReplaceAll(item.Answereddate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _inquireruser.InitExtra()
        item.AddExtra("inquireruser",  _inquireruser)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)
_answeredbyuser.InitExtra()
        item.AddExtra("answeredbyuser",  _answeredbyuser)

        items = append(items, item)
    }


     return items
}

func (p *InquiryManager) Get(id int64) *Inquiry {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and iq_id = ?")

    
    query.WriteString(" and iq_user = u_id")
    
    query.WriteString(" and iq_gym = g_id")
    
    query.WriteString(" and iq_answeredby = u_id")
    
    
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

func (p *InquiryManager) GetWhere(args []interface{}) *Inquiry {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *InquiryManager) Count(args []interface{}) int {
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

func (p *InquiryManager) FindAll() []Inquiry {
    return p.Find(nil)
}

func (p *InquiryManager) Find(args []interface{}) []Inquiry {
    if !p.Conn.IsConnect() {
        var items []Inquiry
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
                query.WriteString(" and iq_")
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
            orderby = "iq_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "iq_" + orderby
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
            orderby = "iq_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "iq_" + orderby
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
        items := make([]Inquiry, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *InquiryManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and iq_")
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
    
    query.WriteString(" group by iq_")
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
