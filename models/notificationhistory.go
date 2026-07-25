package models

import (
    "gym/global/config"
    "gym/models/notificationhistory"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Notificationhistory struct {
            
    Id                int64 `json:"id"`         
    Sender                int64 `json:"sender"`         
    Receiver                int64 `json:"receiver"`         
    Gym                int64 `json:"gym"`         
    Type                int `json:"type"`         
    Title                string `json:"title"`         
    Body                string `json:"body"`         
    Data                string `json:"data"`         
    Status                int `json:"status"`         
    Errormessage                string `json:"errormessage"`         
    Sentdate                string `json:"sentdate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type NotificationhistoryManager struct {
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

func (c *Notificationhistory) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewNotificationhistoryManager(conn *Connection) *NotificationhistoryManager {
    var item NotificationhistoryManager


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

func (p *NotificationhistoryManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *NotificationhistoryManager) SetIndex(index string) {
    p.Index = index
}

func (p *NotificationhistoryManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *NotificationhistoryManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *NotificationhistoryManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *NotificationhistoryManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *NotificationhistoryManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select nh_id, nh_sender, nh_receiver, nh_gym, nh_type, nh_title, nh_body, nh_data, nh_status, nh_errormessage, nh_sentdate, nh_date from notificationhistory_tb")

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

func (p *NotificationhistoryManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from notificationhistory_tb")

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

func (p *NotificationhistoryManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select nh_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from notificationhistory_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    

    return ret.String()
}

func (p *NotificationhistoryManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate notificationhistory_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *NotificationhistoryManager) Insert(item *Notificationhistory) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Sentdate == "" {
       item.Sentdate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into notificationhistory_tb (nh_id, nh_sender, nh_receiver, nh_gym, nh_type, nh_title, nh_body, nh_data, nh_status, nh_errormessage, nh_sentdate, nh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Sender, item.Receiver, item.Gym, item.Type, item.Title, item.Body, item.Data, item.Status, item.Errormessage, item.Sentdate, item.Date)
    } else {
        query = "insert into notificationhistory_tb (nh_sender, nh_receiver, nh_gym, nh_type, nh_title, nh_body, nh_data, nh_status, nh_errormessage, nh_sentdate, nh_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Sender, item.Receiver, item.Gym, item.Type, item.Title, item.Body, item.Data, item.Status, item.Errormessage, item.Sentdate, item.Date)
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

func (p *NotificationhistoryManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notificationhistory_tb where nh_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *NotificationhistoryManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from notificationhistory_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NotificationhistoryManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and nh_")
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

func (p *NotificationhistoryManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from notificationhistory_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *NotificationhistoryManager) Update(item *Notificationhistory) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Sentdate == "" {
       item.Sentdate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update notificationhistory_tb set nh_sender = ?, nh_receiver = ?, nh_gym = ?, nh_type = ?, nh_title = ?, nh_body = ?, nh_data = ?, nh_status = ?, nh_errormessage = ?, nh_sentdate = ?, nh_date = ? where nh_id = ?"
	_, err := p.Exec(query, item.Sender, item.Receiver, item.Gym, item.Type, item.Title, item.Body, item.Data, item.Status, item.Errormessage, item.Sentdate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *NotificationhistoryManager) UpdateWhere(columns []notificationhistory.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update notificationhistory_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == notificationhistory.ColumnId {
        initQuery.WriteString("nh_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnSender {
        initQuery.WriteString("nh_sender = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnReceiver {
        initQuery.WriteString("nh_receiver = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnGym {
        initQuery.WriteString("nh_gym = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnType {
        initQuery.WriteString("nh_type = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnTitle {
        initQuery.WriteString("nh_title = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnBody {
        initQuery.WriteString("nh_body = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnData {
        initQuery.WriteString("nh_data = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnStatus {
        initQuery.WriteString("nh_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnErrormessage {
        initQuery.WriteString("nh_errormessage = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnSentdate {
        initQuery.WriteString("nh_sentdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == notificationhistory.ColumnDate {
        initQuery.WriteString("nh_date = ?")
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


func (p *NotificationhistoryManager) UpdateSender(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_sender = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateReceiver(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_receiver = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateGym(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_gym = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateType(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_type = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateTitle(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_title = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateBody(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_body = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateData(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_data = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateStatus(value int, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_status = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateErrormessage(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_errormessage = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateSentdate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_sentdate = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *NotificationhistoryManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update notificationhistory_tb set nh_date = ? where nh_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *NotificationhistoryManager) GetIdentity() int64 {
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

func (p *Notificationhistory) InitExtra() {
    p.Extra = map[string]interface{}{

    }
}

func (p *NotificationhistoryManager) ReadRow(rows *sql.Rows) *Notificationhistory {
    var item Notificationhistory
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Sender, &item.Receiver, &item.Gym, &item.Type, &item.Title, &item.Body, &item.Data, &item.Status, &item.Errormessage, &item.Sentdate, &item.Date)
        
        if item.Sentdate == "0000-00-00 00:00:00" || item.Sentdate == "1000-01-01 00:00:00" || item.Sentdate == "9999-01-01 00:00:00" {
            item.Sentdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Sentdate = strings.ReplaceAll(strings.ReplaceAll(item.Sentdate, "T", " "), "Z", "")
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

func (p *NotificationhistoryManager) ReadRows(rows *sql.Rows) []Notificationhistory {
    var items []Notificationhistory

    for rows.Next() {
        var item Notificationhistory
        

        err := rows.Scan(&item.Id, &item.Sender, &item.Receiver, &item.Gym, &item.Type, &item.Title, &item.Body, &item.Data, &item.Status, &item.Errormessage, &item.Sentdate, &item.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Sentdate == "0000-00-00 00:00:00" || item.Sentdate == "1000-01-01 00:00:00" || item.Sentdate == "9999-01-01 00:00:00" {
            item.Sentdate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Sentdate = strings.ReplaceAll(strings.ReplaceAll(item.Sentdate, "T", " "), "Z", "")
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

func (p *NotificationhistoryManager) Get(id int64) *Notificationhistory {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and nh_id = ?")

    
    
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

func (p *NotificationhistoryManager) GetWhere(args []interface{}) *Notificationhistory {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *NotificationhistoryManager) Count(args []interface{}) int {
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

func (p *NotificationhistoryManager) FindAll() []Notificationhistory {
    return p.Find(nil)
}

func (p *NotificationhistoryManager) Find(args []interface{}) []Notificationhistory {
    if !p.Conn.IsConnect() {
        var items []Notificationhistory
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
                query.WriteString(" and nh_")
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
            orderby = "nh_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "nh_" + orderby
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
            orderby = "nh_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "nh_" + orderby
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
        items := make([]Notificationhistory, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *NotificationhistoryManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and nh_")
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
    
    query.WriteString(" group by nh_")
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
