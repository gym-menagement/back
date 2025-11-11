package models

import (
    "gym/global/config"
    "gym/models/membershipusage"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Membershipusage struct {
            
    Id                int64 `json:"id"`         
    Membership                int64 `json:"membership"`         
    User                int64 `json:"user"`         
    Type                membershipusage.Type `json:"type"`         
    Totaldays                int `json:"totaldays"`         
    Useddays                int `json:"useddays"`         
    Remainingdays                int `json:"remainingdays"`         
    Totalcount                int `json:"totalcount"`         
    Usedcount                int `json:"usedcount"`         
    Remainingcount                int `json:"remainingcount"`         
    Startdate                string `json:"startdate"`         
    Enddate                string `json:"enddate"`         
    Status                membershipusage.Status `json:"status"`         
    Pausedays                int `json:"pausedays"`         
    Lastuseddate                string `json:"lastuseddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type MembershipusageManager struct {
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

func (c *Membershipusage) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewMembershipusageManager(conn *Connection) *MembershipusageManager {
    var item MembershipusageManager


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

func (p *MembershipusageManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *MembershipusageManager) SetIndex(index string) {
    p.Index = index
}

func (p *MembershipusageManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *MembershipusageManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *MembershipusageManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *MembershipusageManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *MembershipusageManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select mu_id, mu_membership, mu_user, mu_type, mu_totaldays, mu_useddays, mu_remainingdays, mu_totalcount, mu_usedcount, mu_remainingcount, mu_startdate, mu_enddate, mu_status, mu_pausedays, mu_lastuseddate, mu_date, m_id, m_gym, m_user, m_name, m_sex, m_birth, m_phonenum, m_address, m_image, m_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from membershipusage_tb, membership_tb, user_tb")

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
    
    ret.WriteString("and mu_membership = m_id ")
    
    ret.WriteString("and mu_user = u_id ")
    

    return ret.String()
}

func (p *MembershipusageManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from membershipusage_tb")

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
    
    ret.WriteString("and mu_membership = m_id ")
    
    ret.WriteString("and mu_user = u_id ")
    

    return ret.String()
}

func (p *MembershipusageManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select mu_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from membershipusage_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and mu_membership = m_id ")
    
    ret.WriteString("and mu_user = u_id ")
    

    return ret.String()
}

func (p *MembershipusageManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate membershipusage_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *MembershipusageManager) Insert(item *Membershipusage) error {
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
        query = "insert into membershipusage_tb (mu_id, mu_membership, mu_user, mu_type, mu_totaldays, mu_useddays, mu_remainingdays, mu_totalcount, mu_usedcount, mu_remainingcount, mu_startdate, mu_enddate, mu_status, mu_pausedays, mu_lastuseddate, mu_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Membership, item.User, item.Type, item.Totaldays, item.Useddays, item.Remainingdays, item.Totalcount, item.Usedcount, item.Remainingcount, item.Startdate, item.Enddate, item.Status, item.Pausedays, item.Lastuseddate, item.Date)
    } else {
        query = "insert into membershipusage_tb (mu_membership, mu_user, mu_type, mu_totaldays, mu_useddays, mu_remainingdays, mu_totalcount, mu_usedcount, mu_remainingcount, mu_startdate, mu_enddate, mu_status, mu_pausedays, mu_lastuseddate, mu_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Membership, item.User, item.Type, item.Totaldays, item.Useddays, item.Remainingdays, item.Totalcount, item.Usedcount, item.Remainingcount, item.Startdate, item.Enddate, item.Status, item.Pausedays, item.Lastuseddate, item.Date)
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

func (p *MembershipusageManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from membershipusage_tb where mu_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *MembershipusageManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from membershipusage_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MembershipusageManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and mu_")
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

func (p *MembershipusageManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from membershipusage_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *MembershipusageManager) Update(item *Membershipusage) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Lastuseddate == "" {
       item.Lastuseddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update membershipusage_tb set mu_membership = ?, mu_user = ?, mu_type = ?, mu_totaldays = ?, mu_useddays = ?, mu_remainingdays = ?, mu_totalcount = ?, mu_usedcount = ?, mu_remainingcount = ?, mu_startdate = ?, mu_enddate = ?, mu_status = ?, mu_pausedays = ?, mu_lastuseddate = ?, mu_date = ? where mu_id = ?"
	_, err := p.Exec(query, item.Membership, item.User, item.Type, item.Totaldays, item.Useddays, item.Remainingdays, item.Totalcount, item.Usedcount, item.Remainingcount, item.Startdate, item.Enddate, item.Status, item.Pausedays, item.Lastuseddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *MembershipusageManager) UpdateWhere(columns []membershipusage.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update membershipusage_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == membershipusage.ColumnId {
        initQuery.WriteString("mu_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnMembership {
        initQuery.WriteString("mu_membership = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnUser {
        initQuery.WriteString("mu_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnType {
        initQuery.WriteString("mu_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnTotaldays {
        initQuery.WriteString("mu_totaldays = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnUseddays {
        initQuery.WriteString("mu_useddays = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnRemainingdays {
        initQuery.WriteString("mu_remainingdays = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnTotalcount {
        initQuery.WriteString("mu_totalcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnUsedcount {
        initQuery.WriteString("mu_usedcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnRemainingcount {
        initQuery.WriteString("mu_remainingcount = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnStartdate {
        initQuery.WriteString("mu_startdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnEnddate {
        initQuery.WriteString("mu_enddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnStatus {
        initQuery.WriteString("mu_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnPausedays {
        initQuery.WriteString("mu_pausedays = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnLastuseddate {
        initQuery.WriteString("mu_lastuseddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == membershipusage.ColumnDate {
        initQuery.WriteString("mu_date = ?")
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


func (p *MembershipusageManager) UpdateMembership(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_membership = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_user = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateType(value membershipusage.Type, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_type = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateTotaldays(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_totaldays = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateUseddays(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_useddays = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateRemainingdays(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_remainingdays = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateTotalcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_totalcount = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateUsedcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_usedcount = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateRemainingcount(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_remainingcount = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateStartdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_startdate = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateEnddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_enddate = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateStatus(value membershipusage.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_status = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdatePausedays(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_pausedays = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateLastuseddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_lastuseddate = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *MembershipusageManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update membershipusage_tb set mu_date = ? where mu_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *MembershipusageManager) GetIdentity() int64 {
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

func (p *Membershipusage) InitExtra() {
    p.Extra = map[string]interface{}{
            "type":     membershipusage.GetType(p.Type),
            "status":     membershipusage.GetStatus(p.Status),

    }
}

func (p *MembershipusageManager) ReadRow(rows *sql.Rows) *Membershipusage {
    var item Membershipusage
    var err error

    var _membership Membership
    var _user User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Membership, &item.User, &item.Type, &item.Totaldays, &item.Useddays, &item.Remainingdays, &item.Totalcount, &item.Usedcount, &item.Remainingcount, &item.Startdate, &item.Enddate, &item.Status, &item.Pausedays, &item.Lastuseddate, &item.Date, &_membership.Id, &_membership.Gym, &_membership.User, &_membership.Name, &_membership.Sex, &_membership.Birth, &_membership.Phonenum, &_membership.Address, &_membership.Image, &_membership.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
        
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
        _membership.InitExtra()
        item.AddExtra("membership",  _membership)
_user.InitExtra()
        item.AddExtra("user",  _user)

        return &item
    }
}

func (p *MembershipusageManager) ReadRows(rows *sql.Rows) []Membershipusage {
    var items []Membershipusage

    for rows.Next() {
        var item Membershipusage
        var _membership Membership
        var _user User
        

        err := rows.Scan(&item.Id, &item.Membership, &item.User, &item.Type, &item.Totaldays, &item.Useddays, &item.Remainingdays, &item.Totalcount, &item.Usedcount, &item.Remainingcount, &item.Startdate, &item.Enddate, &item.Status, &item.Pausedays, &item.Lastuseddate, &item.Date, &_membership.Id, &_membership.Gym, &_membership.User, &_membership.Name, &_membership.Sex, &_membership.Birth, &_membership.Phonenum, &_membership.Address, &_membership.Image, &_membership.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Email, &_user.Name, &_user.Tel, &_user.Address, &_user.Image, &_user.Sex, &_user.Birth, &_user.Type, &_user.Connectid, &_user.Level, &_user.Role, &_user.Use, &_user.Logindate, &_user.Lastchangepasswddate, &_user.Date)
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
        _membership.InitExtra()
        item.AddExtra("membership",  _membership)
_user.InitExtra()
        item.AddExtra("user",  _user)

        items = append(items, item)
    }


     return items
}

func (p *MembershipusageManager) Get(id int64) *Membershipusage {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and mu_id = ?")

    
    query.WriteString(" and mu_membership = m_id")
    
    query.WriteString(" and mu_user = u_id")
    
    
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

func (p *MembershipusageManager) GetWhere(args []interface{}) *Membershipusage {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *MembershipusageManager) Count(args []interface{}) int {
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

func (p *MembershipusageManager) FindAll() []Membershipusage {
    return p.Find(nil)
}

func (p *MembershipusageManager) Find(args []interface{}) []Membershipusage {
    if !p.Conn.IsConnect() {
        var items []Membershipusage
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
                query.WriteString(" and mu_")
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
            orderby = "mu_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mu_" + orderby
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
            orderby = "mu_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "mu_" + orderby
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
        items := make([]Membershipusage, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *MembershipusageManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and mu_")
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
    
    query.WriteString(" group by mu_")
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
