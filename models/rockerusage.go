package models

import (
    "gym/global/config"
    "gym/models/rockerusage"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Rockerusage struct {
            
    Id                int64 `json:"id"`         
    Rocker                int64 `json:"rocker"`         
    User                int64 `json:"user"`         
    Membership                int64 `json:"membership"`         
    Startdate                string `json:"startdate"`         
    Enddate                string `json:"enddate"`         
    Status                rockerusage.Status `json:"status"`         
    Deposit                int `json:"deposit"`         
    Monthlyfee                int `json:"monthlyfee"`         
    Note                string `json:"note"`         
    Assignedby                int64 `json:"assignedby"`         
    Assigneddate                string `json:"assigneddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type RockerusageManager struct {
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

func (c *Rockerusage) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewRockerusageManager(conn *Connection) *RockerusageManager {
    var item RockerusageManager


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

func (p *RockerusageManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *RockerusageManager) SetIndex(index string) {
    p.Index = index
}

func (p *RockerusageManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *RockerusageManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *RockerusageManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *RockerusageManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *RockerusageManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select ru_id, ru_rocker, ru_user, ru_membership, ru_startdate, ru_enddate, ru_status, ru_deposit, ru_monthlyfee, ru_note, ru_assignedby, ru_assigneddate, ru_date, r_id, r_group, r_name, r_available, r_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date, m_id, m_gym, m_user, m_name, m_sex, m_birth, m_phonenum, m_address, m_image, m_date, u_id, u_loginid, u_passwd, u_email, u_name, u_tel, u_address, u_image, u_sex, u_birth, u_type, u_connectid, u_level, u_role, u_use, u_logindate, u_lastchangepasswddate, u_date from rockerusage_tb, rocker_tb, user_tb, membership_tb, user_tb")

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
    
    ret.WriteString("and ru_rocker = r_id ")
    
    ret.WriteString("and ru_user = u_id ")
    
    ret.WriteString("and ru_membership = m_id ")
    
    ret.WriteString("and ru_assignedby = u_id ")
    

    return ret.String()
}

func (p *RockerusageManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from rockerusage_tb")

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
    
    ret.WriteString("and ru_rocker = r_id ")
    
    ret.WriteString("and ru_user = u_id ")
    
    ret.WriteString("and ru_membership = m_id ")
    
    ret.WriteString("and ru_assignedby = u_id ")
    

    return ret.String()
}

func (p *RockerusageManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select ru_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from rockerusage_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and ru_rocker = r_id ")
    
    ret.WriteString("and ru_user = u_id ")
    
    ret.WriteString("and ru_membership = m_id ")
    
    ret.WriteString("and ru_assignedby = u_id ")
    

    return ret.String()
}

func (p *RockerusageManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate rockerusage_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *RockerusageManager) Insert(item *Rockerusage) error {
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
	
    if item.Assigneddate == "" {
       item.Assigneddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into rockerusage_tb (ru_id, ru_rocker, ru_user, ru_membership, ru_startdate, ru_enddate, ru_status, ru_deposit, ru_monthlyfee, ru_note, ru_assignedby, ru_assigneddate, ru_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Rocker, item.User, item.Membership, item.Startdate, item.Enddate, item.Status, item.Deposit, item.Monthlyfee, item.Note, item.Assignedby, item.Assigneddate, item.Date)
    } else {
        query = "insert into rockerusage_tb (ru_rocker, ru_user, ru_membership, ru_startdate, ru_enddate, ru_status, ru_deposit, ru_monthlyfee, ru_note, ru_assignedby, ru_assigneddate, ru_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Rocker, item.User, item.Membership, item.Startdate, item.Enddate, item.Status, item.Deposit, item.Monthlyfee, item.Note, item.Assignedby, item.Assigneddate, item.Date)
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

func (p *RockerusageManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from rockerusage_tb where ru_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *RockerusageManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from rockerusage_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *RockerusageManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and ru_")
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

func (p *RockerusageManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from rockerusage_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *RockerusageManager) Update(item *Rockerusage) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Startdate == "" {
       item.Startdate = "1000-01-01 00:00:00"
    }
	
    if item.Enddate == "" {
       item.Enddate = "1000-01-01 00:00:00"
    }
	
    if item.Assigneddate == "" {
       item.Assigneddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update rockerusage_tb set ru_rocker = ?, ru_user = ?, ru_membership = ?, ru_startdate = ?, ru_enddate = ?, ru_status = ?, ru_deposit = ?, ru_monthlyfee = ?, ru_note = ?, ru_assignedby = ?, ru_assigneddate = ?, ru_date = ? where ru_id = ?"
	_, err := p.Exec(query, item.Rocker, item.User, item.Membership, item.Startdate, item.Enddate, item.Status, item.Deposit, item.Monthlyfee, item.Note, item.Assignedby, item.Assigneddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *RockerusageManager) UpdateWhere(columns []rockerusage.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update rockerusage_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == rockerusage.ColumnId {
        initQuery.WriteString("ru_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnRocker {
        initQuery.WriteString("ru_rocker = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnUser {
        initQuery.WriteString("ru_user = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnMembership {
        initQuery.WriteString("ru_membership = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnStartdate {
        initQuery.WriteString("ru_startdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnEnddate {
        initQuery.WriteString("ru_enddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnStatus {
        initQuery.WriteString("ru_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnDeposit {
        initQuery.WriteString("ru_deposit = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnMonthlyfee {
        initQuery.WriteString("ru_monthlyfee = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnNote {
        initQuery.WriteString("ru_note = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnAssignedby {
        initQuery.WriteString("ru_assignedby = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnAssigneddate {
        initQuery.WriteString("ru_assigneddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == rockerusage.ColumnDate {
        initQuery.WriteString("ru_date = ?")
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


func (p *RockerusageManager) UpdateRocker(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_rocker = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateUser(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_user = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateMembership(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_membership = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateStartdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_startdate = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateEnddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_enddate = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateStatus(value rockerusage.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_status = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateDeposit(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_deposit = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateMonthlyfee(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_monthlyfee = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateNote(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_note = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateAssignedby(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_assignedby = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateAssigneddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_assigneddate = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *RockerusageManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update rockerusage_tb set ru_date = ? where ru_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *RockerusageManager) GetIdentity() int64 {
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

func (p *Rockerusage) InitExtra() {
    p.Extra = map[string]interface{}{
            "status":     rockerusage.GetStatus(p.Status),

    }
}

func (p *RockerusageManager) ReadRow(rows *sql.Rows) *Rockerusage {
    var item Rockerusage
    var err error

    var _rocker Rocker
    var _memberuser User
    var _membership Membership
    var _assignedbyuser User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Rocker, &item.User, &item.Membership, &item.Startdate, &item.Enddate, &item.Status, &item.Deposit, &item.Monthlyfee, &item.Note, &item.Assignedby, &item.Assigneddate, &item.Date, &_rocker.Id, &_rocker.Group, &_rocker.Name, &_rocker.Available, &_rocker.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_membership.Id, &_membership.Gym, &_membership.User, &_membership.Name, &_membership.Sex, &_membership.Birth, &_membership.Phonenum, &_membership.Address, &_membership.Image, &_membership.Date, &_assignedbyuser.Id, &_assignedbyuser.Loginid, &_assignedbyuser.Passwd, &_assignedbyuser.Email, &_assignedbyuser.Name, &_assignedbyuser.Tel, &_assignedbyuser.Address, &_assignedbyuser.Image, &_assignedbyuser.Sex, &_assignedbyuser.Birth, &_assignedbyuser.Type, &_assignedbyuser.Connectid, &_assignedbyuser.Level, &_assignedbyuser.Role, &_assignedbyuser.Use, &_assignedbyuser.Logindate, &_assignedbyuser.Lastchangepasswddate, &_assignedbyuser.Date)
        
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
		
        if item.Assigneddate == "0000-00-00 00:00:00" || item.Assigneddate == "1000-01-01 00:00:00" || item.Assigneddate == "9999-01-01 00:00:00" {
            item.Assigneddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Assigneddate = strings.ReplaceAll(strings.ReplaceAll(item.Assigneddate, "T", " "), "Z", "")
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
        _rocker.InitExtra()
        item.AddExtra("rocker",  _rocker)
_memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_membership.InitExtra()
        item.AddExtra("membership",  _membership)
_assignedbyuser.InitExtra()
        item.AddExtra("assignedbyuser",  _assignedbyuser)

        return &item
    }
}

func (p *RockerusageManager) ReadRows(rows *sql.Rows) []Rockerusage {
    var items []Rockerusage

    for rows.Next() {
        var item Rockerusage
        var _rocker Rocker
        var _memberuser User
        var _membership Membership
        var _assignedbyuser User
        

        err := rows.Scan(&item.Id, &item.Rocker, &item.User, &item.Membership, &item.Startdate, &item.Enddate, &item.Status, &item.Deposit, &item.Monthlyfee, &item.Note, &item.Assignedby, &item.Assigneddate, &item.Date, &_rocker.Id, &_rocker.Group, &_rocker.Name, &_rocker.Available, &_rocker.Date, &_memberuser.Id, &_memberuser.Loginid, &_memberuser.Passwd, &_memberuser.Email, &_memberuser.Name, &_memberuser.Tel, &_memberuser.Address, &_memberuser.Image, &_memberuser.Sex, &_memberuser.Birth, &_memberuser.Type, &_memberuser.Connectid, &_memberuser.Level, &_memberuser.Role, &_memberuser.Use, &_memberuser.Logindate, &_memberuser.Lastchangepasswddate, &_memberuser.Date, &_membership.Id, &_membership.Gym, &_membership.User, &_membership.Name, &_membership.Sex, &_membership.Birth, &_membership.Phonenum, &_membership.Address, &_membership.Image, &_membership.Date, &_assignedbyuser.Id, &_assignedbyuser.Loginid, &_assignedbyuser.Passwd, &_assignedbyuser.Email, &_assignedbyuser.Name, &_assignedbyuser.Tel, &_assignedbyuser.Address, &_assignedbyuser.Image, &_assignedbyuser.Sex, &_assignedbyuser.Birth, &_assignedbyuser.Type, &_assignedbyuser.Connectid, &_assignedbyuser.Level, &_assignedbyuser.Role, &_assignedbyuser.Use, &_assignedbyuser.Logindate, &_assignedbyuser.Lastchangepasswddate, &_assignedbyuser.Date)
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
		
        if item.Assigneddate == "0000-00-00 00:00:00" || item.Assigneddate == "1000-01-01 00:00:00" || item.Assigneddate == "9999-01-01 00:00:00" {
            item.Assigneddate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Assigneddate = strings.ReplaceAll(strings.ReplaceAll(item.Assigneddate, "T", " "), "Z", "")
        }
		
        if item.Date == "0000-00-00 00:00:00" || item.Date == "1000-01-01 00:00:00" || item.Date == "9999-01-01 00:00:00" {
            item.Date = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Date = strings.ReplaceAll(strings.ReplaceAll(item.Date, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _rocker.InitExtra()
        item.AddExtra("rocker",  _rocker)
_memberuser.InitExtra()
        item.AddExtra("memberuser",  _memberuser)
_membership.InitExtra()
        item.AddExtra("membership",  _membership)
_assignedbyuser.InitExtra()
        item.AddExtra("assignedbyuser",  _assignedbyuser)

        items = append(items, item)
    }


     return items
}

func (p *RockerusageManager) Get(id int64) *Rockerusage {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and ru_id = ?")

    
    query.WriteString(" and ru_rocker = r_id")
    
    query.WriteString(" and ru_user = u_id")
    
    query.WriteString(" and ru_membership = m_id")
    
    query.WriteString(" and ru_assignedby = u_id")
    
    
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

func (p *RockerusageManager) GetWhere(args []interface{}) *Rockerusage {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *RockerusageManager) Count(args []interface{}) int {
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

func (p *RockerusageManager) FindAll() []Rockerusage {
    return p.Find(nil)
}

func (p *RockerusageManager) Find(args []interface{}) []Rockerusage {
    if !p.Conn.IsConnect() {
        var items []Rockerusage
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
                query.WriteString(" and ru_")
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
            orderby = "ru_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "ru_" + orderby
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
            orderby = "ru_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "ru_" + orderby
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
        items := make([]Rockerusage, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *RockerusageManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and ru_")
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
    
    query.WriteString(" group by ru_")
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
