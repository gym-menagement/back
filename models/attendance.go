package models

import (
    "gym/global/config"
    "gym/models/attendance"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Attendance struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Usehealth                int64 `json:"usehealth"`         
    Gym                int64 `json:"gym"`         
    Type                attendance.Type `json:"type"`         
    Method                attendance.Method `json:"method"`         
    Checkintime                string `json:"checkintime"`         
    Checkouttime                string `json:"checkouttime"`         
    Duration                int `json:"duration"`         
    Status                attendance.Status `json:"status"`         
    Note                string `json:"note"`         
    Ip                string `json:"ip"`         
    Device                string `json:"device"`         
    Createdby                int64 `json:"createdby"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type AttendanceManager struct {
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

func (c *Attendance) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewAttendanceManager(conn *Connection) *AttendanceManager {
    var item AttendanceManager


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

func (p *AttendanceManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *AttendanceManager) SetIndex(index string) {
    p.Index = index
}

func (p *AttendanceManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *AttendanceManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *AttendanceManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *AttendanceManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *AttendanceManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select at_id, at_user, at_usehealth, at_gym, at_type, at_method, at_checkintime, at_checkouttime, at_duration, at_status, at_note, at_ip, at_device, at_createdby, at_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, uh_id, uh_order, uh_health, uh_membership, uh_user, uh_term, uh_discount, uh_startday, uh_endday, uh_gym, uh_status, uh_totalcount, uh_usedcount, uh_remainingcount, uh_qrcode, uh_lastuseddate, uh_date, g_id, g_name, g_address, g_tel, g_user, g_date from attendance_tb, user_tb, usehealth_tb, gym_tb")

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
    
    ret.WriteString("and at_user = u_id ")
    
    ret.WriteString("and at_usehealth = uh_id ")
    
    ret.WriteString("and at_gym = g_id ")
    

    return ret.String()
}

func (p *AttendanceManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from attendance_tb")

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
    
    ret.WriteString("and at_user = u_id ")
    
    ret.WriteString("and at_usehealth = uh_id ")
    
    ret.WriteString("and at_gym = g_id ")
    

    return ret.String()
}

func (p *AttendanceManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select at_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from attendance_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and at_user = u_id ")
    
    ret.WriteString("and at_usehealth = uh_id ")
    
    ret.WriteString("and at_gym = g_id ")
    

    return ret.String()
}

func (p *AttendanceManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate attendance_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *AttendanceManager) Insert(item *Attendance) error {
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
        query = "insert into attendance_tb (at_id, at_user, at_usehealth, at_gym, at_type, at_method, at_checkintime, at_checkouttime, at_duration, at_status, at_note, at_ip, at_device, at_createdby, at_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Usehealth, item.Gym, item.Type, item.Method, item.Checkintime, item.Checkouttime, item.Duration, item.Status, item.Note, item.Ip, item.Device, item.Createdby, item.Date)
    } else {
        query = "insert into attendance_tb (at_user, at_usehealth, at_gym, at_type, at_method, at_checkintime, at_checkouttime, at_duration, at_status, at_note, at_ip, at_device, at_createdby, at_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Usehealth, item.Gym, item.Type, item.Method, item.Checkintime, item.Checkouttime, item.Duration, item.Status, item.Note, item.Ip, item.Device, item.Createdby, item.Date)
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

func (p *AttendanceManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from attendance_tb where at_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *AttendanceManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from attendance_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AttendanceManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and at_")
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

func (p *AttendanceManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from attendance_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AttendanceManager) Update(item *Attendance) error {
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
	

	query := "update attendance_tb set at_user = ?, at_usehealth = ?, at_gym = ?, at_type = ?, at_method = ?, at_checkintime = ?, at_checkouttime = ?, at_duration = ?, at_status = ?, at_note = ?, at_ip = ?, at_device = ?, at_createdby = ?, at_date = ? where at_id = ?"
	_, err := p.Exec(query, item.User, item.Usehealth, item.Gym, item.Type, item.Method, item.Checkintime, item.Checkouttime, item.Duration, item.Status, item.Note, item.Ip, item.Device, item.Createdby, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *AttendanceManager) UpdateWhere(columns []attendance.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update attendance_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == attendance.ColumnId {
        initQuery.WriteString("at_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnUser {
        initQuery.WriteString("at_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnUsehealth {
        initQuery.WriteString("at_usehealth = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnGym {
        initQuery.WriteString("at_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnType {
        initQuery.WriteString("at_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnMethod {
        initQuery.WriteString("at_method = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnCheckintime {
        initQuery.WriteString("at_checkintime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnCheckouttime {
        initQuery.WriteString("at_checkouttime = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnDuration {
        initQuery.WriteString("at_duration = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnStatus {
        initQuery.WriteString("at_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnNote {
        initQuery.WriteString("at_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnIp {
        initQuery.WriteString("at_ip = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnDevice {
        initQuery.WriteString("at_device = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnCreatedby {
        initQuery.WriteString("at_createdby = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == attendance.ColumnDate {
        initQuery.WriteString("at_date = ?")
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


func (p *AttendanceManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_user = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateUsehealth(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_usehealth = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_gym = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateType(value attendance.Type, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_type = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateMethod(value attendance.Method, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_method = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateCheckintime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_checkintime = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateCheckouttime(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_checkouttime = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateDuration(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_duration = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateStatus(value attendance.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_status = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_note = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateIp(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_ip = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateDevice(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_device = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateCreatedby(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_createdby = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AttendanceManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update attendance_tb set at_date = ? where at_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *AttendanceManager) GetIdentity() int64 {
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

func (p *Attendance) InitExtra() {
    p.Extra = map[string]interface{}{
            "type":     attendance.GetType(p.Type),
            "method":     attendance.GetMethod(p.Method),
            "status":     attendance.GetStatus(p.Status),

    }
}

func (p *AttendanceManager) ReadRow(rows *sql.Rows) *Attendance {
    var item Attendance
    var err error

    var _user User
    var _usehealth Usehealth
    var _gym Gym
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Usehealth, &item.Gym, &item.Type, &item.Method, &item.Checkintime, &item.Checkouttime, &item.Duration, &item.Status, &item.Note, &item.Ip, &item.Device, &item.Createdby, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_usehealth.Id, &_usehealth.Order, &_usehealth.Health, &_usehealth.Membership, &_usehealth.User, &_usehealth.Term, &_usehealth.Discount, &_usehealth.Startday, &_usehealth.Endday, &_usehealth.Gym, &_usehealth.Status, &_usehealth.Totalcount, &_usehealth.Usedcount, &_usehealth.Remainingcount, &_usehealth.Qrcode, &_usehealth.Lastuseddate, &_usehealth.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date)
        
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
        _user.InitExtra()
        item.AddExtra("user",  _user)
_usehealth.InitExtra()
        item.AddExtra("usehealth",  _usehealth)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)

        return &item
    }
}

func (p *AttendanceManager) ReadRows(rows *sql.Rows) []Attendance {
    var items []Attendance

    for rows.Next() {
        var item Attendance
        var _user User
        var _usehealth Usehealth
        var _gym Gym
        

        err := rows.Scan(&item.Id, &item.User, &item.Usehealth, &item.Gym, &item.Type, &item.Method, &item.Checkintime, &item.Checkouttime, &item.Duration, &item.Status, &item.Note, &item.Ip, &item.Device, &item.Createdby, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date, &_usehealth.Id, &_usehealth.Order, &_usehealth.Health, &_usehealth.Membership, &_usehealth.User, &_usehealth.Term, &_usehealth.Discount, &_usehealth.Startday, &_usehealth.Endday, &_usehealth.Gym, &_usehealth.Status, &_usehealth.Totalcount, &_usehealth.Usedcount, &_usehealth.Remainingcount, &_usehealth.Qrcode, &_usehealth.Lastuseddate, &_usehealth.Date, &_gym.Id, &_gym.Name, &_gym.Address, &_gym.Tel, &_gym.User, &_gym.Date)
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
        _user.InitExtra()
        item.AddExtra("user",  _user)
_usehealth.InitExtra()
        item.AddExtra("usehealth",  _usehealth)
_gym.InitExtra()
        item.AddExtra("gym",  _gym)

        items = append(items, item)
    }


     return items
}

func (p *AttendanceManager) Get(id int64) *Attendance {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and at_id = ?")

    
    query.WriteString(" and at_user = u_id")
    
    query.WriteString(" and at_usehealth = uh_id")
    
    query.WriteString(" and at_gym = g_id")
    
    
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

func (p *AttendanceManager) GetWhere(args []interface{}) *Attendance {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *AttendanceManager) Count(args []interface{}) int {
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

func (p *AttendanceManager) FindAll() []Attendance {
    return p.Find(nil)
}

func (p *AttendanceManager) Find(args []interface{}) []Attendance {
    if !p.Conn.IsConnect() {
        var items []Attendance
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
                query.WriteString(" and at_")
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
            orderby = "at_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "at_" + orderby
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
            orderby = "at_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "at_" + orderby
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
        items := make([]Attendance, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *AttendanceManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and at_")
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
    
    query.WriteString(" group by at_")
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
