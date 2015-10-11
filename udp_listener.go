package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	// _ "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"log"
	"net"
	"net/http"
	"time"
)

const maxPacketSize = 1024 * 1024

var db *sql.DB

type Stat struct {
	Key   string         `json:"key"`
	Value types.JsonText `json:"value"`
}

func (s *Stat) persist() {
	var incomingId int
	err := db.QueryRow(`INSERT INTO incoming(received_at, data) VALUES($1, $2) RETURNING id`, time.Now(), s.Value).Scan(&incomingId)
	if err != nil {
		log.Fatal("persist()", err)
	}
	fmt.Printf("class: %s id: %d\n", s.Key, incomingId)
}

func countRows() int {
	var rows int
	err := db.QueryRow("select count(*) from incoming").Scan(&rows)
	if err != nil {
		log.Fatal("rowcount: ", err)
	}

	return rows
}

func listenToUDP(conn *net.UDPConn) {
	buffer := make([]byte, maxPacketSize)
	for {
		bytes, err := conn.Read(buffer)
		if err != nil {
			log.Println("UDP read error: ", err.Error())
			continue
		}

		msg := make([]byte, bytes)
		copy(msg, buffer)

		var stat Stat
		err = json.Unmarshal(msg, &stat)
		if err != nil {
			log.Println(err)
			log.Printf("%+v\n", stat)
		}

		stat.persist()
	}
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	count := countRows()
	jmap := map[string]int{
		"status": 200,
		"count":  count,
	}

	output, err := json.MarshalIndent(jmap, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		log.Fatal("ResolveUDPAddr", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("ListenUDP", err)
	}

	db, err = sql.Open("postgres", "user=markmulder dbname=notifier sslmode=disable")
	if err != nil {
		log.Fatal("DB Open()", err)
	}
	defer db.Close()

	rows := countRows()
	fmt.Println("Total rows: ", rows)

	go listenToUDP(conn)
	http.HandleFunc("/", handleCount)

	fmt.Println("Will start listening on port 8000")
	http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}

}
