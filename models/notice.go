package models

import (
    "gym/global/config"
    "gym/models/notice"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Notice struct {
            
    Id                int64 `json:"id"`         
    Gym                int64 `json:"gym"`         
    Title                string `json:"title"`         
    Content                string `json:"content"`         
    Type                notice.Type `json:"type"`         
    Ispopup                notice.Ispopup `json:"ispopup"`         
    Ispush                notice.Ispush `json:"ispush"`         
    Target                notice.Target `json:"target"`         
    Viewcount                int `json:"viewcount"`         
    Startdate                string `json:"startdate"`         
    Enddate                string `json:"enddate"`         
    Status                notice.Status `json:"status"`         
    Createdby                int64 `json:"createdby"`         
    Createddate                string `json:"createddate"`         
    Updateddate                string `json:"updateddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type NoticeManager struct {
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

func (c *Notice) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewNoticeManager(conn *Connection) *NoticeManager {
    var item NoticeManager


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

func (p *NoticeManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *NoticeManager) SetIndex(index string) {
    p.Index = index
}

func (p *NoticeManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *NoticeManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *NoticeManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *NoticeManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *NoticeManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select nt_id, nt_gym, nt_title, nt_content, nt_type, nt_ispopup, nt_ispush, nt_target, nt_viewcount, nt_startdate, nt_enddate, nt_status, nt_createdby, nt_createddate, nt_updateddate, nt_date, g_id, g_name, g_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from notice_tb, gym_tb, user_tb")

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
    
    ret.WriteString("and nt_gym = g_id ")
    
    ret.WriteString("and nt_createdby = u_id ")
    

    return ret.String()
}

func (p *NoticeManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from notice_tb")

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
    
    ret.WriteString("and nt_gym = g_id ")
    
    ret.WriteString("and nt_createdby = u_id ")
    

    return ret.String()
}

func (p *NoticeManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select nt_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from notice_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and nt_gym = g_id ")
    
    ret.WriteString("and nt_createdby = u_id ")
    

    return ret.String()
}

func (p *NoticeManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate notice_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *NoticeManager) Insert(item *Notice) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Updateddate == "" {
       item.Updateddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into notice_tb (nt_id, nt_gym, nt_title, nt_content, nt_type, nt_ispopup, nt_ispush, nt_target, nt_viewcount, nt_startdate, nt_enddate, nt_status, nt_createdby, nt_createddate, nt_updateddate, nt_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Gym, item.Title, item.Content, item.Type, item.Ispopup, item.Ispush, item.Target, item.Viewcount, item.Startdate, item.Enddate, item.Status, item.Createdby, item.Createddate, item.Updateddate, item.Date)
    } else {
        query = "insert into notice_tb (nt_gym, nt_title, nt_content, nt_type, nt_ispopup, nt_ispush, nt_target, nt_viewcount, nt_startdate, nt_enddate, nt_status, nt_createdby, nt_createddate, nt_updateddate, nt_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Gym, item.Title, item.Content, item.Type, item.Ispopup, item.Ispush, item.Target, item.Viewcount, item.Startdate, item.Enddate, item.Status, item.Createdby, item.Createddate, item.Updateddate, item.Date)
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

func (p *NoticeManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notice_tb where nt_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *NoticeManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notice_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NoticeManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and nt_")
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

func (p *NoticeManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from notice_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NoticeManager) Update(item *Notice) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Updateddate == "" {
       item.Updateddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update notice_tb set nt_gym = ?, nt_title = ?, nt_content = ?, nt_type = ?, nt_ispopup = ?, nt_ispush = ?, nt_target = ?, nt_viewcount = ?, nt_startdate = ?, nt_enddate = ?, nt_status = ?, nt_createdby = ?, nt_createddate = ?, nt_updateddate = ?, nt_date = ? where nt_id = ?"
	_, err := p.Exec(query, item.Gym, item.Title, item.Content, item.Type, item.Ispopup, item.Ispush, item.Target, item.Viewcount, item.Startdate, item.Enddate, item.Status, item.Createdby, item.Createddate, item.Updateddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *NoticeManager) UpdateWhere(columns []notice.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update notice_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == notice.ColumnId {
        initQuery.WriteString("nt_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnGym {
        initQuery.WriteString("nt_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnTitle {
        initQuery.WriteString("nt_title = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnContent {
        initQuery.WriteString("nt_content = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnType {
        initQuery.WriteString("nt_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnIspopup {
        initQuery.WriteString("nt_ispopup = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnIspush {
        initQuery.WriteString("nt_ispush = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnTarget {
        initQuery.WriteString("nt_target = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnViewcount {
        initQuery.WriteString("nt_viewcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnStartdate {
        initQuery.WriteString("nt_startdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnEnddate {
        initQuery.WriteString("nt_enddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnStatus {
        initQuery.WriteString("nt_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnCreatedby {
        initQuery.WriteString("nt_createdby = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnCreateddate {
        initQuery.WriteString("nt_createddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnUpdateddate {
        initQuery.WriteString("nt_updateddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notice.ColumnDate {
        initQuery.WriteString("nt_date = ?")
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


func (p *NoticeManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_gym = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateTitle(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_title = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateContent(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_content = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateType(value notice.Type, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_type = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateIspopup(value notice.Ispopup, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_ispopup = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateIspush(value notice.Ispush, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_ispush = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateTarget(value notice.Target, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_target = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateViewcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_viewcount = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateStartdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_startdate = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateEnddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_enddate = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateStatus(value notice.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_status = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateCreatedby(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_createdby = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateCreateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_createddate = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateUpdateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_updateddate = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NoticeManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notice_tb set nt_date = ? where nt_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *NoticeManager) GetIdentity() int64 {
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

func (p *Notice) InitExtra() {
    p.Extra = map[string]interface{}{
            "type":     notice.GetType(p.Type),
            "ispopup":     notice.GetIspopup(p.Ispopup),
            "ispush":     notice.GetIspush(p.Ispush),
            "target":     notice.GetTarget(p.Target),
            "status":     notice.GetStatus(p.Status),

    }
}

func (p *NoticeManager) ReadRow(rows *sql.Rows) *Notice {
    var item Notice
    var err error

    var _gym Gym
    var _user User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Gym, &item.Title, &item.Content, &item.Type, &item.Ispopup, &item.Ispush, &item.Target, &item.Viewcount, &item.Startdate, &item.Enddate, &item.Status, &item.Createdby, &item.Createddate, &item.Updateddate, &item.Date, &_gym.Id, &_gym.Name, &_gym.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
        
        if item.Startdate == "0000-00-00 00:00:00" || item.Startdate == "1000-01-01 00:00:00" || item.Startdate == "9999-01-01 00:00:00" {
            item.Startdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startdate = strings.ReplaceAll(strings.ReplaceAll(item.Startdate, "T", " "), "Z", "")
        }
		
        if item.Enddate == "0000-00-00 00:00:00" || item.Enddate == "1000-01-01 00:00:00" || item.Enddate == "9999-01-01 00:00:00" {
            item.Enddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Enddate = strings.ReplaceAll(strings.ReplaceAll(item.Enddate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
        }
		
        if item.Updateddate == "0000-00-00 00:00:00" || item.Updateddate == "1000-01-01 00:00:00" || item.Updateddate == "9999-01-01 00:00:00" {
            item.Updateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updateddate = strings.ReplaceAll(strings.ReplaceAll(item.Updateddate, "T", " "), "Z", "")
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
        _gym.InitExtra()
        item.AddExtra("gym",  _gym)
_user.InitExtra()
        item.AddExtra("user",  _user)

        return &item
    }
}

func (p *NoticeManager) ReadRows(rows *sql.Rows) []Notice {
    var items []Notice

    for rows.Next() {
        var item Notice
        var _gym Gym
        var _user User
        

        err := rows.Scan(&item.Id, &item.Gym, &item.Title, &item.Content, &item.Type, &item.Ispopup, &item.Ispush, &item.Target, &item.Viewcount, &item.Startdate, &item.Enddate, &item.Status, &item.Createdby, &item.Createddate, &item.Updateddate, &item.Date, &_gym.Id, &_gym.Name, &_gym.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Startdate == "0000-00-00 00:00:00" || item.Startdate == "1000-01-01 00:00:00" || item.Startdate == "9999-01-01 00:00:00" {
            item.Startdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Startdate = strings.ReplaceAll(strings.ReplaceAll(item.Startdate, "T", " "), "Z", "")
        }
		
        if item.Enddate == "0000-00-00 00:00:00" || item.Enddate == "1000-01-01 00:00:00" || item.Enddate == "9999-01-01 00:00:00" {
            item.Enddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Enddate = strings.ReplaceAll(strings.ReplaceAll(item.Enddate, "T", " "), "Z", "")
        }
		
        if item.Createddate == "0000-00-00 00:00:00" || item.Createddate == "1000-01-01 00:00:00" || item.Createddate == "9999-01-01 00:00:00" {
            item.Createddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createddate = strings.ReplaceAll(strings.ReplaceAll(item.Createddate, "T", " "), "Z", "")
        }
		
        if item.Updateddate == "0000-00-00 00:00:00" || item.Updateddate == "1000-01-01 00:00:00" || item.Updateddate == "9999-01-01 00:00:00" {
            item.Updateddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updateddate = strings.ReplaceAll(strings.ReplaceAll(item.Updateddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _gym.InitExtra()
        item.AddExtra("gym",  _gym)
_user.InitExtra()
        item.AddExtra("user",  _user)

        items = append(items, item)
    }


     return items
}

func (p *NoticeManager) Get(id int64) *Notice {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and nt_id = ?")

    
    query.WriteString(" and nt_gym = g_id")
    
    query.WriteString(" and nt_createdby = u_id")
    
    
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

func (p *NoticeManager) GetWhere(args []interface{}) *Notice {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *NoticeManager) Count(args []interface{}) int {
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

func (p *NoticeManager) FindAll() []Notice {
    return p.Find(nil)
}

func (p *NoticeManager) Find(args []interface{}) []Notice {
    if !p.Conn.IsConnect() {
        var items []Notice
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
                query.WriteString(" and nt_")
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
            orderby = "nt_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "nt_" + orderby
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
            orderby = "nt_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "nt_" + orderby
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
        items := make([]Notice, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *NoticeManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and nt_")
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
    
    query.WriteString(" group by nt_")
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
