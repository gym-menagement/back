package models

import (
    "gym/global/config"
    "gym/models/notificationsetting"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Notificationsetting struct {
            
    Id                int64 `json:"id"`         
    User                int64 `json:"user"`         
    Enabled                int `json:"enabled"`         
    Membershipexpiry                int `json:"membershipexpiry"`         
    Membershipnear                int `json:"membershipnear"`         
    Attendanceenc                int `json:"attendanceenc"`         
    Gymannounce                int `json:"gymannounce"`         
    Systemnotice                int `json:"systemnotice"`         
    Paymentconfirm                int `json:"paymentconfirm"`         
    Pauseexpiry                int `json:"pauseexpiry"`         
    Weeklygoal                int `json:"weeklygoal"`         
    Personalrecord                int `json:"personalrecord"`         
    Quietenabled                int `json:"quietenabled"`         
    Quietstart                string `json:"quietstart"`         
    Quietend                string `json:"quietend"`         
    Createddate                string `json:"createddate"`         
    Updateddate                string `json:"updateddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type NotificationsettingManager struct {
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

func (c *Notificationsetting) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewNotificationsettingManager(conn *Connection) *NotificationsettingManager {
    var item NotificationsettingManager


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

func (p *NotificationsettingManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *NotificationsettingManager) SetIndex(index string) {
    p.Index = index
}

func (p *NotificationsettingManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *NotificationsettingManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *NotificationsettingManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *NotificationsettingManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *NotificationsettingManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select ns_id, ns_user, ns_enabled, ns_membershipexpiry, ns_membershipnear, ns_attendanceenc, ns_gymannounce, ns_systemnotice, ns_paymentconfirm, ns_pauseexpiry, ns_weeklygoal, ns_personalrecord, ns_quietenabled, ns_quietstart, ns_quietend, ns_createddate, ns_updateddate, ns_date from notificationsetting_tb")

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

func (p *NotificationsettingManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from notificationsetting_tb")

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

func (p *NotificationsettingManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select ns_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from notificationsetting_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    

    return ret.String()
}

func (p *NotificationsettingManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate notificationsetting_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *NotificationsettingManager) Insert(item *Notificationsetting) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Createddate == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Createddate = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
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
        query = "insert into notificationsetting_tb (ns_id, ns_user, ns_enabled, ns_membershipexpiry, ns_membershipnear, ns_attendanceenc, ns_gymannounce, ns_systemnotice, ns_paymentconfirm, ns_pauseexpiry, ns_weeklygoal, ns_personalrecord, ns_quietenabled, ns_quietstart, ns_quietend, ns_createddate, ns_updateddate, ns_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.User, item.Enabled, item.Membershipexpiry, item.Membershipnear, item.Attendanceenc, item.Gymannounce, item.Systemnotice, item.Paymentconfirm, item.Pauseexpiry, item.Weeklygoal, item.Personalrecord, item.Quietenabled, item.Quietstart, item.Quietend, item.Createddate, item.Updateddate, item.Date)
    } else {
        query = "insert into notificationsetting_tb (ns_user, ns_enabled, ns_membershipexpiry, ns_membershipnear, ns_attendanceenc, ns_gymannounce, ns_systemnotice, ns_paymentconfirm, ns_pauseexpiry, ns_weeklygoal, ns_personalrecord, ns_quietenabled, ns_quietstart, ns_quietend, ns_createddate, ns_updateddate, ns_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.User, item.Enabled, item.Membershipexpiry, item.Membershipnear, item.Attendanceenc, item.Gymannounce, item.Systemnotice, item.Paymentconfirm, item.Pauseexpiry, item.Weeklygoal, item.Personalrecord, item.Quietenabled, item.Quietstart, item.Quietend, item.Createddate, item.Updateddate, item.Date)
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

func (p *NotificationsettingManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notificationsetting_tb where ns_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *NotificationsettingManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notificationsetting_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NotificationsettingManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and ns_")
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

func (p *NotificationsettingManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from notificationsetting_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NotificationsettingManager) Update(item *Notificationsetting) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
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
	

	query := "update notificationsetting_tb set ns_user = ?, ns_enabled = ?, ns_membershipexpiry = ?, ns_membershipnear = ?, ns_attendanceenc = ?, ns_gymannounce = ?, ns_systemnotice = ?, ns_paymentconfirm = ?, ns_pauseexpiry = ?, ns_weeklygoal = ?, ns_personalrecord = ?, ns_quietenabled = ?, ns_quietstart = ?, ns_quietend = ?, ns_createddate = ?, ns_updateddate = ?, ns_date = ? where ns_id = ?"
	_, err := p.Exec(query, item.User, item.Enabled, item.Membershipexpiry, item.Membershipnear, item.Attendanceenc, item.Gymannounce, item.Systemnotice, item.Paymentconfirm, item.Pauseexpiry, item.Weeklygoal, item.Personalrecord, item.Quietenabled, item.Quietstart, item.Quietend, item.Createddate, item.Updateddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *NotificationsettingManager) UpdateWhere(columns []notificationsetting.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update notificationsetting_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == notificationsetting.ColumnId {
        initQuery.WriteString("ns_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnUser {
        initQuery.WriteString("ns_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnEnabled {
        initQuery.WriteString("ns_enabled = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnMembershipexpiry {
        initQuery.WriteString("ns_membershipexpiry = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnMembershipnear {
        initQuery.WriteString("ns_membershipnear = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnAttendanceenc {
        initQuery.WriteString("ns_attendanceenc = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnGymannounce {
        initQuery.WriteString("ns_gymannounce = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnSystemnotice {
        initQuery.WriteString("ns_systemnotice = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnPaymentconfirm {
        initQuery.WriteString("ns_paymentconfirm = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnPauseexpiry {
        initQuery.WriteString("ns_pauseexpiry = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnWeeklygoal {
        initQuery.WriteString("ns_weeklygoal = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnPersonalrecord {
        initQuery.WriteString("ns_personalrecord = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnQuietenabled {
        initQuery.WriteString("ns_quietenabled = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnQuietstart {
        initQuery.WriteString("ns_quietstart = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnQuietend {
        initQuery.WriteString("ns_quietend = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnCreateddate {
        initQuery.WriteString("ns_createddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnUpdateddate {
        initQuery.WriteString("ns_updateddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationsetting.ColumnDate {
        initQuery.WriteString("ns_date = ?")
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


func (p *NotificationsettingManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_user = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateEnabled(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_enabled = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateMembershipexpiry(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_membershipexpiry = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateMembershipnear(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_membershipnear = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateAttendanceenc(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_attendanceenc = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateGymannounce(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_gymannounce = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateSystemnotice(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_systemnotice = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdatePaymentconfirm(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_paymentconfirm = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdatePauseexpiry(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_pauseexpiry = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateWeeklygoal(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_weeklygoal = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdatePersonalrecord(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_personalrecord = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateQuietenabled(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_quietenabled = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateQuietstart(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_quietstart = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateQuietend(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_quietend = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateCreateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_createddate = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateUpdateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_updateddate = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationsettingManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationsetting_tb set ns_date = ? where ns_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *NotificationsettingManager) GetIdentity() int64 {
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

func (p *Notificationsetting) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *NotificationsettingManager) ReadRow(rows *sql.Rows) *Notificationsetting {
    var item Notificationsetting
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.User, &item.Enabled, &item.Membershipexpiry, &item.Membershipnear, &item.Attendanceenc, &item.Gymannounce, &item.Systemnotice, &item.Paymentconfirm, &item.Pauseexpiry, &item.Weeklygoal, &item.Personalrecord, &item.Quietenabled, &item.Quietstart, &item.Quietend, &item.Createddate, &item.Updateddate, &item.Date)
        
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
        
        return &item
    }
}

func (p *NotificationsettingManager) ReadRows(rows *sql.Rows) []Notificationsetting {
    var items []Notificationsetting

    for rows.Next() {
        var item Notificationsetting
        

        err := rows.Scan(&item.Id, &item.User, &item.Enabled, &item.Membershipexpiry, &item.Membershipnear, &item.Attendanceenc, &item.Gymannounce, &item.Systemnotice, &item.Paymentconfirm, &item.Pauseexpiry, &item.Weeklygoal, &item.Personalrecord, &item.Quietenabled, &item.Quietstart, &item.Quietend, &item.Createddate, &item.Updateddate, &item.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
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
        
        items = append(items, item)
    }


     return items
}

func (p *NotificationsettingManager) Get(id int64) *Notificationsetting {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and ns_id = ?")

    
    
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

func (p *NotificationsettingManager) GetWhere(args []interface{}) *Notificationsetting {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *NotificationsettingManager) Count(args []interface{}) int {
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

func (p *NotificationsettingManager) FindAll() []Notificationsetting {
    return p.Find(nil)
}

func (p *NotificationsettingManager) Find(args []interface{}) []Notificationsetting {
    if !p.Conn.IsConnect() {
        var items []Notificationsetting
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
                query.WriteString(" and ns_")
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
            orderby = "ns_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "ns_" + orderby
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
            orderby = "ns_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "ns_" + orderby
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
        items := make([]Notificationsetting, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *NotificationsettingManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and ns_")
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
    
    query.WriteString(" group by ns_")
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
