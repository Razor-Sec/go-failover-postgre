package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	yaml "gopkg.in/yaml.v3"

	_ "github.com/lib/pq"
)

	
type Configuration struct {
	Local struct {
		Host string `yaml:"host,omitempty"`
		Port int    `yaml:"port,omitempty"`
		User string `yaml:"user,omitempty"`
	} `yaml:"local"`
	Remote []struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port,omitempty"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password",omitempty`
	} `yaml:"remote"`
}

func parseYaml(file String) (*Configuration, error) {
	yamlConf, err := os.ReadFile(file)
    if err != nil {
        log.Printf("yamlConf.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlConf, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return yamlConf, err
}

func main(){
	test, err := parseYaml("test.yaml")
	fmt.Printf("%#v",test)
}

// const (
// 	//dbname   = "postgres"
// 	interval = 6
// 	fail     = 3
// 	//logfile  = "/tmp/logtest"
// )

// func main() {
// 	path, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	LOG_FILE := path + "/log/" + time.Now().Format("01-02-2006") + ".log"
// 	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer logFile.Close()
// 	log.SetOutput(logFile)
// 	currentTime := time.Now().Format(time.RFC3339)
// 	//interval := flag.Int("interval", interval, "Time for check")
// 	//fail := flag.Int("fail", fail, "Time for fail")
// 	host1 := flag.String("host1", "localhost", "")
// 	host2 := flag.String("host2", "localhost", "")
// 	user1 := flag.String("user1", "postgres", "")
// 	user2 := flag.String("user2", "postgres", "")
// 	port1 := flag.Int("port1", 5432, "a int var")
// 	port2 := flag.Int("port2", 5432, "a int var")
// 	password1 := flag.String("password1", "", "")
// 	password2 := flag.String("password2", "", "")
// 	localdata := flag.String("localdata", "/var/lib/pgsql/data", "")
// 	localpg := flag.String("localpg", "/usr/bin/pg_ctl", "")
// 	localhost := flag.String("localhost", "localhost", "")
// 	localport := flag.Int("localport", 5432, "a int var")
// 	localuser := flag.String("localuser", "postgres", "")
// 	localpass := flag.String("localpass", "", "")
// 	flag.Parse()
// 	fmt.Println(*host1, *port1, *user1, *password1, *localdata)
// 	fmt.Println(*host2, *port2, *user2, *password2, *localdata)
// 	var faill int = 0
// 	for {
// 		var check1 bool = mainDB(*host1, *port1, *user1, *password1)                 // url_host, port, user, password
// 		var check2 bool = mainDB(*host2, *port2, *user2, *password2)                 // db 2
// 		var checklocal bool = mainDB(*localhost, *localport, *localuser, *localpass) // db local
// 		//println(check1, check2, checklocal)
// 		if !checklocal {
// 			log.Panic("[FAIL] Your local Database DOWN")
// 			panic("[FAIL] Your local Database DOWN")
// 		}
// 		if check1 == true || check2 == true {
// 			fmt.Println(currentTime, "[INFO] No need promote")
// 			log.Println("[INFO] No need promote")
// 			faill = 0
// 		} else {
// 			faill++
// 			fmt.Println(time.Now().Format(time.RFC3339), "[WARN] Need promote ,TIMES : ", faill)
// 			log.Println("[WARN] Need promote ,", faill)
// 			if faill >= fail {
// 				fmt.Println(time.Now().Format(time.RFC3339), "[FAIL] TIME OUT GO PROMOTE....")
// 				log.Println("[FAIL] TIME OUT GO PROMOTE....")
// 				promote(*localpg, *localport, *localdata) // (localpg string, port int, user string, password string)
// 				// function promote
// 				break
// 			}
// 		}
// 		time.Sleep(interval * time.Second)
// 	}
// }

// //func mainGO(host1 string)
// // func checkingDB1(err error) bool { // ini ga guna
// // 	if err != nil {
// // 		return false
// // 	}
// // 	return true
// // }

// func mainDB(url_host string, port int, user string, password string) bool {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s sslmode=disable", url_host, port, user, password)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer db.Close()
// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println(time.Now().Format(time.RFC3339), "[WARN] Database ", url_host, ":", port, "DOWN")
// 		log.Println("[WARN] Database ", url_host, ":", port, "DOWN")
// 		return false
// 	}
// 	if url_host == "localhost" {
// 		var statusdb string = checkStatus(db)
// 		if statusdb == "t" || statusdb == "true" {
// 			fmt.Println(time.Now().Format(time.RFC3339), "[INFO] Status LOCAL Database on STANDBY MODE")
// 			log.Println("[INFO] Status LOCAL Database on STANDBY MODE")
// 		} else {
// 			fmt.Println(time.Now().Format(time.RFC3339), "[WARN] Status LOCAL Database on MASTER MODE")
// 			log.Println("[WARN] Status LOCAL Database on MASTER MODE")
// 			//fmt.Println("Please change status database to standby mode")
// 		}
// 	} else {
// 		fmt.Println(time.Now().Format(time.RFC3339), "[INFO] Database ", url_host, ":", port, "UP")
// 		log.Println("[INFO] Database ", url_host, ":", port, "UP")
// 	}
// 	return true
// }

// func promote(localpg string, port int, localdata string) {
// 	//fmt.Println(localpg, "promote", "-D", localdata)
// 	run, err := exec.Command(localpg, "promote", "-D", localdata).Output() //check this for promote
// 	//run, err := exec.Command("su", "-", "postgres", "-c", "whoami").Output() //check this for promote
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(string(run))
// }

// func checkStatus(db *sql.DB) string {
// 	rows, err := db.Query("SELECT pg_is_in_recovery from pg_is_in_recovery();")
// 	//rows, err := db.Query("SELECT * from test;")
// 	var res string
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//fmt.Println(rows) //GG NEXT DEK
// 	defer rows.Close()
// 	for rows.Next() {
// 		if err := rows.Scan(&res); err != nil {
// 			log.Fatalln(err)
// 		}
// 	}
// 	if err := rows.Err(); err != nil {
// 		log.Fatalln(err)
// 	}
// 	return res
// }
