package rss

import (
	"database/sql"
	"fmt"
	"log"
	"tgbot/app"
	"time"
)

var (
	ErrRssIsDeleted  = fmt.Errorf("rss is deleted")
	ErrRssIsNil      = fmt.Errorf("rss is nil")
	ErrRssURLInvalid = fmt.Errorf("rss url is invalid")
)

const (
	tableRssResourceDDL = `
CREATE TABLE IF NOT EXISTS t_rss_resource (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	url TEXT NOT NULL,
	title TEXT NOT NULL,
	link TEXT UNIQUE NOT NULL,
	last_item_title TEXT NOT NULL,
	last_item_link TEXT NOT NULL,
	last_item_publish_time INTEGER NOT NULL,
	is_deleted INTEGER DEFAULT 0 NOT NULL,
	create_time INTEGER DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time INTEGER NOT NULL
);`

	tableRssUserDDL = `
CREATE TABLE IF NOT EXISTS t_rss_user (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	rss_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	is_deleted INTEGER DEFAULT 0 NOT NULL,
	create_time INTEGER DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time INTEGER NOT NULL
);`
)

func createTables() error {
	db := app.DB

	_, err := db.Exec(tableRssResourceDDL)
	if err != nil {
		return err
	}

	_, err = db.Exec(tableRssUserDDL)
	if err != nil {
		return err
	}

	return nil
}

type RssResource struct {
	ID                  int64  `db:"id"`
	URL                 string `db:"url"`
	Title               string `db:"title"`
	Link                string `db:"link"`
	LastItemTitle       string `db:"last_item_title"`
	LastItemLink        string `db:"last_item_link"`
	LastItemPublishTime int64  `db:"last_item_publish_time"`
	IsDeleted           int64  `db:"is_deleted"`
	CreateTime          int64  `db:"create_time"`
	UpdateTime          int64  `db:"update_time"`
}

func GetRssResources() (*[]RssResource, error) {
	db := app.DB

	res := []RssResource{}
	err := db.Select(&res, "SELECT * FROM t_rss_resource WHERE is_deleted = 0")
	if err != nil {
		log.Println("Error: db.Select failed,", err)
		return nil, err
	}

	return &res, nil
}

func GetRssResourceByID(id int64) (*RssResource, error) {
	db := app.DB

	res := RssResource{}
	err := db.Get(&res, "SELECT * FROM t_rss_resource WHERE id = ?", id)
	if err != nil {
		log.Println("Error: db.Get failed,", err)
		return nil, err
	}

	return &res, nil
}

func GetRssResourceByURL(url string) (*RssResource, error) {
	db := app.DB

	res := RssResource{}
	err := db.Get(&res, "SELECT * FROM t_rss_resource WHERE url = ?", url)
	if err != nil {
		log.Println("Error: db.Get failed,", err)
		return nil, err
	}

	return &res, nil
}

func GetRssResourceByUserID(userID int) (*[]RssResource, error) {
	db := app.DB

	res := []RssResource{}

	err := db.Select(&res, "SELECT r.id, r.url, r.title, r.link, r.last_item_title, r.last_item_link, r.last_item_publish_time, r.is_deleted, r.create_time, r.update_time FROM t_rss_resource r LEFT JOIN t_rss_user u ON (u.rss_id = r.id) WHERE r.is_deleted = 0 AND u.is_deleted = 0 AND u.user_id = ?", userID)
	if err != nil {
		log.Println("Error: db.Select failed,", err)
		return nil, err
	}

	return &res, nil
}

func GetUserIDsByRssID(rssID int64) ([]int, error) {
	db := app.DB

	res := []int{}

	err := db.Select(&res, "SELECT u.user_id FROM t_rss_user u LEFT JOIN t_rss_resource r ON (r.id = u.rss_id) WHERE u.is_deleted = 0 AND r.is_deleted = 0 AND u.rss_id = ?", rssID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func AddRssResource(r *RssResource) (int64, error) {
	if r == nil {
		return 0, ErrRssIsNil
	}

	if len(r.URL) == 0 {
		return 0, ErrRssURLInvalid
	}

	l, err := GetRssResourceByURL(r.URL)
	if err == nil {
		if l.IsDeleted != 0 {
			err = UpdateRssResource(l.ID, map[string]interface{}{
				"title":                  r.Title,
				"link":                   r.Link,
				"last_item_title":        r.LastItemTitle,
				"last_item_link":         r.LastItemLink,
				"last_item_publish_time": r.LastItemPublishTime,
				"is_deleted":             0,
			})
			if err != nil {
				return 0, err
			}
			return l.ID, nil
		} else {
			return l.ID, nil
		}
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	db := app.DB

	res, err := db.Exec("INSERT INTO t_rss_resource (url, title, link, last_item_title, last_item_link, last_item_publish_time, is_deleted, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", r.URL, r.Title, r.Link, r.LastItemTitle, r.LastItemLink, r.LastItemPublishTime, 0, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		log.Println("Error: db.Exec failed,", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error: LastInsertId failed,", err)
		return 0, err
	}

	return id, nil
}

func AddRssUser(rssID int64, userID int) error {
	db := app.DB

	r, err := GetRssResourceByID(rssID)
	if err != nil {
		return err
	}

	if r.IsDeleted != 0 {
		return ErrRssIsDeleted
	}

	isDeleted := 1
	err = db.Get(&isDeleted, "SELECT u.is_deleted FROM t_rss_user u LEFT JOIN t_rss_resource r ON (r.id = u.rss_id) WHERE r.is_deleted = 0 AND u.user_id = ? AND u.rss_id = ?", userID, rssID)
	if err == nil {
		if isDeleted == 0 {
			return nil
		} else {
			_, err = UpdateRssUser(rssID, userID, map[string]interface{}{
				"is_deleted": 0,
			})
			if err != nil {
				return err
			}
			return nil
		}
	} else if err != sql.ErrNoRows {
		log.Println("Error: db.Get failed,", err)
		return err
	}

	_, err = db.Exec("INSERT INTO t_rss_user (rss_id, user_id, create_time, update_time) VALUES (?, ?, ?, ?)", rssID, userID, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		log.Println("Error: db.Exec failed,", err)
		return err
	}

	return nil
}

func UpdateRssUser(rssID int64, userID int, fields map[string]interface{}) (int64, error) {
	if fields == nil || len(fields) == 0 {
		return 0, nil
	}

	db := app.DB

	sql := `
UPDATE t_rss_user
SET
`
	for k := range fields {
		if k == "update_time" || k == "id" || k == "rss_id" || k == "user_id" {
			continue
		}
		sql += fmt.Sprintf("  %s = :%s,\n", k, k)
	}

	sql += "  update_time = :update_time\n"
	sql += "WHERE rss_id = :rss_id AND user_id = :user_id"

	fields["rss_id"] = rssID
	fields["user_id"] = userID
	fields["update_time"] = time.Now().Unix()

	res, err := db.NamedExec(sql, fields)
	if err != nil {
		log.Println("Error: NamedExec failed,", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error: RowsAffected failed,", err)
		return 0, err
	}

	return rowsAffected, nil
}

func UpdateRssResource(id int64, fields map[string]interface{}) error {
	if fields == nil || len(fields) == 0 {
		return nil
	}

	db := app.DB

	sql := `
UPDATE t_rss_resource
SET
`
	for k := range fields {
		if k == "update_time" || k == "id" {
			continue
		}
		sql += fmt.Sprintf("  %s = :%s,\n", k, k)
	}

	sql += "  update_time = :update_time\n"
	sql += "WHERE id = :id"

	fields["id"] = id
	fields["update_time"] = time.Now().Unix()

	_, err := db.NamedExec(sql, fields)
	if err != nil {
		log.Println("Error: NamedExec failed,", err)
		return err
	}

	return nil
}
