package main

import (
	"log"
	"time"

	"zero/core/stores/clickhouse"
	"zero/core/stores/sqlx"
)

func main() {
	conn := clickhouse.New("tcp://127.0.0.1:9000")
	_, err := conn.Exec(`
			CREATE TABLE IF NOT EXISTS example (
				country_code FixedString(2),
				os_id        UInt8,
				browser_id   UInt8,
				categories   Array(Int16),
				action_day   Date,
				action_time  DateTime
			) engine=Memory
		`)
	if err != nil {
		log.Fatal(err)
	}

	conn.Transact(func(session sqlx.Session) error {
		stmt, err := session.Prepare("INSERT INTO example (country_code, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for i := 0; i < 10; i++ {
			_, err := stmt.Exec("RU", 10+i, 100+i, []int16{1, 2, 3}, time.Now(), time.Now())
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	var items []struct {
		CountryCode string    `db:"country_code"`
		OsID        uint8     `db:"os_id"`
		BrowserID   uint8     `db:"browser_id"`
		Categories  []int16   `db:"categories"`
		ActionTime  time.Time `db:"action_time"`
	}

	err = conn.QueryRows(&items, "SELECT country_code, os_id, browser_id, categories, action_time FROM example")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_time: %s",
			item.CountryCode, item.OsID, item.BrowserID, item.Categories, item.ActionTime)
	}

	if _, err := conn.Exec("DROP TABLE example"); err != nil {
		log.Fatal(err)
	}
}
