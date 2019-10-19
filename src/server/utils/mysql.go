package utils

import (
	"database/sql"
	"github.com/gihnius/gomemoize/src/memoize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/domain"
	log "github.com/sirupsen/logrus"
)

const (
	cacheTimeoutSec = 60 * 10
)

func QueryWhitelist() (whitelists []domain.Whitelist) {
	return memoize.Memoize("queryWhitelist", queryWhitelist, cacheTimeoutSec).([]domain.Whitelist)
}

func queryWhitelist() interface{} {
	ret := []domain.Whitelist{}
	db, err := sql.Open("mysql", base.MysqlDataSourceNaming)
	if err != nil {
		log.Fatal("Connect to mysql server error", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT `id`,`ip`,`create_time`,`status` FROM `nb_gen_id_whitelist`")
	defer rows.Close()
	if err != nil {
		log.Error(err)
		return ret
	}
	if rows.Err() != nil {
		log.Error(rows.Err())
		return ret
	}
	for rows.Next() {
		whitelist := domain.Whitelist{}
		err = rows.Scan(&whitelist.Id, &whitelist.Ip, &whitelist.CreateTime, &whitelist.Status)
		if err != nil {
			log.Error("Scan error", err)
		}
		ret = append(ret, whitelist)
	}
	return ret
}
