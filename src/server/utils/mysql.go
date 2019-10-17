package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/domain"
	log "github.com/sirupsen/logrus"
)

func QueryWhitelist() (whitelists []domain.Whitelist) {
	db, err := sql.Open("mysql", base.MysqlDataSourceNaming)
	if err != nil {
		log.Fatal("Connect to mysql server error", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT `id`,`ip`,`create_time`,`status` FROM `nb_gen_id_whitelist`")
	defer rows.Close()
	if err != nil {
		log.Error(err)
		return
	}
	if rows.Err() != nil {
		log.Error(rows.Err())
		return
	}
	for rows.Next() {
		whitelist := domain.Whitelist{}
		err = rows.Scan(&whitelist.Id, &whitelist.Ip, &whitelist.CreateTime, &whitelist.Status)
		if err != nil {
			log.Error("Scan error", err)
		}
		whitelists = append(whitelists, whitelist)
	}
	return
}
