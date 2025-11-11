package models

import (
    "gym/global/config"
    "gym/models/appversion"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "gym/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Appversion struct {
            
    Id                int64 `json:"id"`         
    Platform                string `json:"platform"`         
    Version                string `json:"version"`         
    Minversion                string `json:"minversion"`         
    Forceupdate                appversion.Forceupdate `json:"forceupdate"`         
    Updatemessage                string `json:"updatemessage"`         
    Downloadurl                string `json:"downloadurl"`         
    Status                appversion.Status `json:"status"`         
    Releasedate                string `json:"releasedate"`         
    Createddate                string `json:"createddate"`         
    Date                string `json:"date"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type AppversionManager struct {
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

func (c *Appversion) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewAppversionManager(conn *Connection) *AppversionManager {
    var item AppversionManager


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

func (p *AppversionManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *AppversionManager) SetIndex(index string) {
    p.Index = index
}

func (p *AppversionManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *AppversionManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *AppversionManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *AppversionManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
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

func (p *AppversionManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select av_id, av_platform, av_version, av_minversion, av_forceupdate, av_updatemessage, av_downloadurl, av_status, av_releasedate, av_createddate, av_date from appversion_tb")

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

func (p *AppversionManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from appversion_tb")

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

func (p *AppversionManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select av_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from appversion_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    

    return ret.String()
}

func (p *AppversionManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate appversion_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *AppversionManager) Insert(item *Appversion) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Releasedate == "" {
       item.Releasedate = "1000-01-01 00:00:00"
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
        query = "insert into appversion_tb (av_id, av_platform, av_version, av_minversion, av_forceupdate, av_updatemessage, av_downloadurl, av_status, av_releasedate, av_createddate, av_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Platform, item.Version, item.Minversion, item.Forceupdate, item.Updatemessage, item.Downloadurl, item.Status, item.Releasedate, item.Createddate, item.Date)
    } else {
        query = "insert into appversion_tb (av_platform, av_version, av_minversion, av_forceupdate, av_updatemessage, av_downloadurl, av_status, av_releasedate, av_createddate, av_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Platform, item.Version, item.Minversion, item.Forceupdate, item.Updatemessage, item.Downloadurl, item.Status, item.Releasedate, item.Createddate, item.Date)
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

func (p *AppversionManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from appversion_tb where av_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *AppversionManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from appversion_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AppversionManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
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
                query.WriteString(" and av_")
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

func (p *AppversionManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from appversion_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AppversionManager) Update(item *Appversion) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Releasedate == "" {
       item.Releasedate = "1000-01-01 00:00:00"
    }
	
    if item.Createddate == "" {
       item.Createddate = "1000-01-01 00:00:00"
    }
	
    if item.Date == "" {
       item.Date = "1000-01-01 00:00:00"
    }
	

	query := "update appversion_tb set av_platform = ?, av_version = ?, av_minversion = ?, av_forceupdate = ?, av_updatemessage = ?, av_downloadurl = ?, av_status = ?, av_releasedate = ?, av_createddate = ?, av_date = ? where av_id = ?"
	_, err := p.Exec(query, item.Platform, item.Version, item.Minversion, item.Forceupdate, item.Updatemessage, item.Downloadurl, item.Status, item.Releasedate, item.Createddate, item.Date, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *AppversionManager) UpdateWhere(columns []appversion.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update appversion_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == appversion.ColumnId {
        initQuery.WriteString("av_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnPlatform {
        initQuery.WriteString("av_platform = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnVersion {
        initQuery.WriteString("av_version = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnMinversion {
        initQuery.WriteString("av_minversion = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnForceupdate {
        initQuery.WriteString("av_forceupdate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnUpdatemessage {
        initQuery.WriteString("av_updatemessage = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnDownloadurl {
        initQuery.WriteString("av_downloadurl = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnStatus {
        initQuery.WriteString("av_status = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnReleasedate {
        initQuery.WriteString("av_releasedate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnCreateddate {
        initQuery.WriteString("av_createddate = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == appversion.ColumnDate {
        initQuery.WriteString("av_date = ?")
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


func (p *AppversionManager) UpdatePlatform(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_platform = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateVersion(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_version = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateMinversion(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_minversion = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateForceupdate(value appversion.Forceupdate, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_forceupdate = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateUpdatemessage(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_updatemessage = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateDownloadurl(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_downloadurl = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateStatus(value appversion.Status, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_status = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateReleasedate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_releasedate = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateCreateddate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_createddate = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AppversionManager) UpdateDate(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update appversion_tb set av_date = ? where av_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *AppversionManager) GetIdentity() int64 {
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

func (p *Appversion) InitExtra() {
    p.Extra = map[string]interface{}{
            "forceupdate":     appversion.GetForceupdate(p.Forceupdate),
            "status":     appversion.GetStatus(p.Status),

    }
}

func (p *AppversionManager) ReadRow(rows *sql.Rows) *Appversion {
    var item Appversion
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Platform, &item.Version, &item.Minversion, &item.Forceupdate, &item.Updatemessage, &item.Downloadurl, &item.Status, &item.Releasedate, &item.Createddate, &item.Date)
        
        if item.Releasedate == "0000-00-00 00:00:00" || item.Releasedate == "1000-01-01 00:00:00" || item.Releasedate == "9999-01-01 00:00:00" {
            item.Releasedate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Releasedate = strings.ReplaceAll(strings.ReplaceAll(item.Releasedate, "T", " "), "Z", "")
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
        
        return &item
    }
}

func (p *AppversionManager) ReadRows(rows *sql.Rows) []Appversion {
    var items []Appversion

    for rows.Next() {
        var item Appversion
        

        err := rows.Scan(&item.Id, &item.Platform, &item.Version, &item.Minversion, &item.Forceupdate, &item.Updatemessage, &item.Downloadurl, &item.Status, &item.Releasedate, &item.Createddate, &item.Date)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Releasedate == "0000-00-00 00:00:00" || item.Releasedate == "1000-01-01 00:00:00" || item.Releasedate == "9999-01-01 00:00:00" {
            item.Releasedate = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Releasedate = strings.ReplaceAll(strings.ReplaceAll(item.Releasedate, "T", " "), "Z", "")
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
        
        items = append(items, item)
    }


     return items
}

func (p *AppversionManager) Get(id int64) *Appversion {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and av_id = ?")

    
    
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

func (p *AppversionManager) GetWhere(args []interface{}) *Appversion {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *AppversionManager) Count(args []interface{}) int {
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

func (p *AppversionManager) FindAll() []Appversion {
    return p.Find(nil)
}

func (p *AppversionManager) Find(args []interface{}) []Appversion {
    if !p.Conn.IsConnect() {
        var items []Appversion
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
                query.WriteString(" and av_")
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
            orderby = "av_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "av_" + orderby
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
            orderby = "av_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "av_" + orderby
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
        items := make([]Appversion, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}





func (p *AppversionManager) GroupBy(name string, args []interface{}) []Groupby {
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
                query.WriteString(" and av_")
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
    
    query.WriteString(" group by av_")
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
