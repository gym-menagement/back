package models

import (
    "gym/global/config"
    "gym/models/usehealthusage"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Usehealthusage struct {
            
    Id                int64 `json:"id"`         
    Gym                int64 `json:"gym"`         
    Usehealth                int64 `json:"usehealth"`         
    Membership                int64 `json:"membership"`         
    User                int64 `json:"user"`         
    Attendance                int64 `json:"attendance"`         
    Type                usehealthusage.Type `json:"type"`         
    Usedcount                int `json:"usedcount"`         
    Remainingcount                int `json:"remainingcount"`         
    Checkintime                string `json:"checkintime"`         
    Checkouttime                string `json:"checkouttime"`         
    Duration                int `json:"duration"`         
    Note                string `json:"note"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type UsehealthusageManager struct {
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

func (c *Usehealthusage) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewUsehealthusageManager(conn *Connection) *UsehealthusageManager {
    var item UsehealthusageManager


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

func (p *UsehealthusageManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *UsehealthusageManager) SetIndex(index string) {
    p.Index = index
}

func (p *UsehealthusageManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *UsehealthusageManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *UsehealthusageManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *UsehealthusageManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *UsehealthusageManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select uhu_id, uhu_gym, uhu_usehealth, uhu_membership, uhu_user, uhu_attendance, uhu_type, uhu_usedcount, uhu_remainingcount, uhu_checkintime, uhu_checkouttime, uhu_duration, uhu_note, uhu_date, g_id, g_name, g_address, g_tel, g_user, g_date, uh_id, uh_order, uh_health, uh_membership, uh_user, uh_term, uh_discount, uh_startday, uh_endday, uh_gym, uh_status, uh_totalcount, uh_usedcount, uh_remainingcount, uh_qrcode, uh_lastuseddate, uh_date, m_id, m_user, m_gym, m_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, at_id, at_user, at_usehealth, at_gym, at_type, at_method, at_checkintime, at_checkouttime, at_duration, at_status, at_note, at_ip, at_device, at_createdby, at_date from usehealthusage_tb, gym_tb, usehealth_tb, membership_tb, user_tb, attendance_tb")

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
    
    ret.WriteString("and uhu_gym = g_id ")
    
    ret.WriteString("and uhu_usehealth = uh_id ")
    
    ret.WriteString("and uhu_membership = m_id ")
    
    ret.WriteString("and uhu_user = u_id ")
    
    ret.WriteString("and uhu_attendance = at_id ")
    

    return ret.String()
}

func (p *UsehealthusageManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from usehealthusage_tb")

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
    
    ret.WriteString("and uhu_gym = g_id ")
    
    ret.WriteString("and uhu_usehealth = uh_id ")
    
    ret.WriteString("and uhu_membership = m_id ")
    
    ret.WriteString("and uhu_user = u_id ")
    
    ret.WriteString("and uhu_attendance = at_id ")
    

    return ret.String()
}

func (p *UsehealthusageManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select uhu_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from usehealthusage_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and uhu_gym = g_id ")
    
    ret.WriteString("and uhu_usehealth = uh_id ")
    
    ret.WriteString("and uhu_membership = m_id ")
    
    ret.WriteString("and uhu_user = u_id ")
    
    ret.WriteString("and uhu_attendance = at_id ")
    

    return ret.String()
}

func (p *UsehealthusageManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate usehealthusage_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *UsehealthusageManager) Insert(item *Usehealthusage) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Checkintime == "" {
       item.Checkintime = "1000-01-01 00:00:00"
    }
	
    if item.Checkouttime == "" {
       item.Checkouttime = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into usehealthusage_tb (uhu_id, uhu_gym, uhu_usehealth, uhu_membership, uhu_user, uhu_attendance, uhu_type, uhu_usedcount, uhu_remainingcount, uhu_checkintime, uhu_checkouttime, uhu_duration, uhu_note, uhu_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Gym, item.Usehealth, item.Membership, item.User, item.Attendance, item.Type, item.Usedcount, item.Remainingcount, item.Checkintime, item.Checkouttime, item.Duration, item.Note, item.Date)
    } else {
        query = "insert into usehealthusage_tb (uhu_gym, uhu_usehealth, uhu_membership, uhu_user, uhu_attendance, uhu_type, uhu_usedcount, uhu_remainingcount, uhu_checkintime, uhu_checkouttime, uhu_duration, uhu_note, uhu_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Gym, item.Usehealth, item.Membership, item.User, item.Attendance, item.Type, item.Usedcount, item.Remainingcount, item.Checkintime, item.Checkouttime, item.Duration, item.Note, item.Date)
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

func (p *UsehealthusageManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from usehealthusage_tb where uhu_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *UsehealthusageManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from usehealthusage_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *UsehealthusageManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and uhu_")
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

func (p *UsehealthusageManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from usehealthusage_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *UsehealthusageManager) Update(item *Usehealthusage) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Checkintime == "" {
       item.Checkintime = "1000-01-01 00:00:00"
    }
	
    if item.Checkouttime == "" {
       item.Checkouttime = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update usehealthusage_tb set uhu_gym = ?, uhu_usehealth = ?, uhu_membership = ?, uhu_user = ?, uhu_attendance = ?, uhu_type = ?, uhu_usedcount = ?, uhu_remainingcount = ?, uhu_checkintime = ?, uhu_checkouttime = ?, uhu_duration = ?, uhu_note = ?, uhu_date = ? where uhu_id = ?"
	_, err := p.Exec(query, item.Gym, item.Usehealth, item.Membership, item.User, item.Attendance, item.Type, item.Usedcount, item.Remainingcount, item.Checkintime, item.Checkouttime, item.Duration, item.Note, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *UsehealthusageManager) UpdateWhere(columns []usehealthusage.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update usehealthusage_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == usehealthusage.ColumnId {
        initQuery.WriteString("uhu_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnGym {
        initQuery.WriteString("uhu_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnUsehealth {
        initQuery.WriteString("uhu_usehealth = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnMembership {
        initQuery.WriteString("uhu_membership = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnUser {
        initQuery.WriteString("uhu_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnAttendance {
        initQuery.WriteString("uhu_attendance = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnType {
        initQuery.WriteString("uhu_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnUsedcount {
        initQuery.WriteString("uhu_usedcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnRemainingcount {
        initQuery.WriteString("uhu_remainingcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnCheckintime {
        initQuery.WriteString("uhu_checkintime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnCheckouttime {
        initQuery.WriteString("uhu_checkouttime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnDuration {
        initQuery.WriteString("uhu_duration = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnNote {
        initQuery.WriteString("uhu_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == usehealthusage.ColumnDate {
        initQuery.WriteString("uhu_date = ?")
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


func (p *UsehealthusageManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_gym = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateUsehealth(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_usehealth = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateMembership(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_membership = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_user = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateAttendance(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_attendance = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateType(value usehealthusage.Type, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_type = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateUsedcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_usedcount = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateRemainingcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_remainingcount = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateCheckintime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_checkintime = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateCheckouttime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_checkouttime = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateDuration(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_duration = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_note = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *UsehealthusageManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update usehealthusage_tb set uhu_date = ? where uhu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *UsehealthusageManager) GetIdentity() int64 {
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

func (p *Usehealthusage) InitExtra() {
    p.Extra = map[string]interface{}{
            "type":     usehealthusage.GetType(p.Type),

    }
}

func (p *UsehealthusageManager) ReadRow(rows *sql.Rows) *Usehealthusage {
    var item Usehealthusage
    var err error

    var _gym Gym
    var _usehealth Usehealth
    var _membership Membership
    var _user User
    var _attendance Attendance
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Gym, &item.Usehealth, &item.Membership, &item.User, &item.Attendance, &item.Type, &item.Usedcount, &item.Remainingcount, &item.Checkintime, &item.Checkouttime, &item.Duration, &item.Note, &item.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_usehealth.Id, &_usehealth.Order, &_usehealth.Health, &_usehealth.Membership, &_usehealth.User, &_usehealth.Term, &_usehealth.Discount, &_usehealth.Startday, &_usehealth.Endday, &_usehealth.Gym, &_usehealth.Status, &_usehealth.Totalcount, &_usehealth.Usedcount, &_usehealth.Remainingcount, &_usehealth.Qrcode, &_usehealth.Lastuseddate, &_usehealth.Date, &_membership.Id, &_membership.User, &_membership.Gym, &_membership.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_attendance.Id, &_attendance.User, &_attendance.Usehealth, &_attendance.Gym, &_attendance.Type, &_attendance.Method, &_attendance.Checkintime, &_attendance.Checkouttime, &_attendance.Duration, &_attendance.Status, &_attendance.Note, &_attendance.Ip, &_attendance.Device, &_attendance.Createdby, &_attendance.Date)
        
        if item.Checkintime == "0000-00-00 00:00:00" || item.Checkintime == "1000-01-01 00:00:00" || item.Checkintime == "9999-01-01 00:00:00" {
            item.Checkintime = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Checkintime = strings.ReplaceAll(strings.ReplaceAll(item.Checkintime, "T", " "), "Z", "")
        }
		
        if item.Checkouttime == "0000-00-00 00:00:00" || item.Checkouttime == "1000-01-01 00:00:00" || item.Checkouttime == "9999-01-01 00:00:00" {
            item.Checkouttime = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Checkouttime = strings.ReplaceAll(strings.ReplaceAll(item.Checkouttime, "T", " "), "Z", "")
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
_usehealth.InitExtra()
        item.AddExtra("usehealth",  _usehealth)
_membership.InitExtra()
        item.AddExtra("membership",  _membership)
_user.InitExtra()
        item.AddExtra("user",  _user)
_attendance.InitExtra()
        item.AddExtra("attendance",  _attendance)

        return &item
    }
}

func (p *UsehealthusageManager) ReadRows(rows *sql.Rows) []Usehealthusage {
    var items []Usehealthusage

    for rows.Next() {
        var item Usehealthusage
        var _gym Gym
        var _usehealth Usehealth
        var _membership Membership
        var _user User
        var _attendance Attendance
        

        err := rows.Scan(&item.Id, &item.Gym, &item.Usehealth, &item.Membership, &item.User, &item.Attendance, &item.Type, &item.Usedcount, &item.Remainingcount, &item.Checkintime, &item.Checkouttime, &item.Duration, &item.Note, &item.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date, &_usehealth.Id, &_usehealth.Order, &_usehealth.Health, &_usehealth.Membership, &_usehealth.User, &_usehealth.Term, &_usehealth.Discount, &_usehealth.Startday, &_usehealth.Endday, &_usehealth.Gym, &_usehealth.Status, &_usehealth.Totalcount, &_usehealth.Usedcount, &_usehealth.Remainingcount, &_usehealth.Qrcode, &_usehealth.Lastuseddate, &_usehealth.Date, &_membership.Id, &_membership.User, &_membership.Gym, &_membership.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_attendance.Id, &_attendance.User, &_attendance.Usehealth, &_attendance.Gym, &_attendance.Type, &_attendance.Method, &_attendance.Checkintime, &_attendance.Checkouttime, &_attendance.Duration, &_attendance.Status, &_attendance.Note, &_attendance.Ip, &_attendance.Device, &_attendance.Createdby, &_attendance.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Checkintime == "0000-00-00 00:00:00" || item.Checkintime == "1000-01-01 00:00:00" || item.Checkintime == "9999-01-01 00:00:00" {
            item.Checkintime = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Checkintime = strings.ReplaceAll(strings.ReplaceAll(item.Checkintime, "T", " "), "Z", "")
        }
		
        if item.Checkouttime == "0000-00-00 00:00:00" || item.Checkouttime == "1000-01-01 00:00:00" || item.Checkouttime == "9999-01-01 00:00:00" {
            item.Checkouttime = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Checkouttime = strings.ReplaceAll(strings.ReplaceAll(item.Checkouttime, "T", " "), "Z", "")
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
_usehealth.InitExtra()
        item.AddExtra("usehealth",  _usehealth)
_membership.InitExtra()
        item.AddExtra("membership",  _membership)
_user.InitExtra()
        item.AddExtra("user",  _user)
_attendance.InitExtra()
        item.AddExtra("attendance",  _attendance)

        items = append(items, item)
    }


     return items
}

func (p *UsehealthusageManager) Get(id int64) *Usehealthusage {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and uhu_id = ?")

    
    query.WriteString(" and uhu_gym = g_id")
    
    query.WriteString(" and uhu_usehealth = uh_id")
    
    query.WriteString(" and uhu_membership = m_id")
    
    query.WriteString(" and uhu_user = u_id")
    
    query.WriteString(" and uhu_attendance = at_id")
    
    
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

func (p *UsehealthusageManager) GetWhere(args []interface{}) *Usehealthusage {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *UsehealthusageManager) Count(args []interface{}) int {
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

func (p *UsehealthusageManager) FindAll() []Usehealthusage {
    return p.Find(nil)
}

func (p *UsehealthusageManager) Find(args []interface{}) []Usehealthusage {
    if !p.Conn.IsConnect() {
        var items []Usehealthusage
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
                query.WriteString(" and uhu_")
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
            orderby = "uhu_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "uhu_" + orderby
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
            orderby = "uhu_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "uhu_" + orderby
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
        items := make([]Usehealthusage, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *UsehealthusageManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and uhu_")
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
    
    query.WriteString(" group by uhu_")
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
