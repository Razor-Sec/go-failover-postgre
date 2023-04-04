package main

import (
	// "database/sql"

	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	// "reflect"

	// "os/exec"

	yaml "gopkg.in/yaml.v3"

	_ "github.com/lib/pq"
)

type Configuration struct {
	Local struct {
		Host     string `yaml:"host,omitempty"`
		Port     int    `yaml:"port,omitempty"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password"`
		PgCtl    string `yaml:"pg-ctl,omitempty"`
		PgData   string `yaml:"pg-data,omitempty"`
	} `yaml:"local,omitempty"`
	Remote []struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port,omitempty"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password"`
	} `yaml:"remote"`
}

const (
	defaultUser        = "postgres"
	defaultPort        = 5432
	defaultLocalHost   = "localhost"
	defaultLocalPgCtl  = "/usr/bin/pg_ctl"
	defaultLocalPgData = "/var/lib/pgsql/data"
	defaultConfDir     = "conf.yaml"
)

func parseYaml(file string) (*Configuration, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	conf := &Configuration{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return conf, err
}

func defaultingStr(variable string, value string) string {
	if variable == "" {
		variable = value
	}
	return variable
}

func defaultingInt(variable int, value int) int {
	if variable == 0 {
		variable = value
	}
	return variable
}

func misConfRemote(kind string, variable string, remoteNumber int) {
	if variable == "" {
		exitCode := 2
		log.Printf("[ERROR] Need %s for remote server number #%d. Exit status %d\n", kind, remoteNumber+1, exitCode)
		fmt.Printf("%s [ERROR] Need %s for remote server number #%d. Exit status %d\n", time.Now().Format(time.RFC3339), kind, remoteNumber+1, exitCode)
		os.Exit(exitCode)
	}
}

func main() {
	filePath := flag.String("config-file", "", "Configuration file")
	flag.Parse()
	if *filePath == "" {
		*filePath = "conf.yaml"
	}
	*filePath = defaultingStr(*filePath, defaultConfDir)
	conf, err := parseYaml(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	conf.Local.Host = defaultingStr(conf.Local.Host, defaultLocalHost)
	conf.Local.Port = defaultingInt(conf.Local.Port, defaultPort)
	conf.Local.User = defaultingStr(conf.Local.User, defaultUser)
	conf.Local.PgCtl = defaultingStr(conf.Local.PgCtl, defaultLocalPgCtl)
	conf.Local.PgData = defaultingStr(conf.Local.PgData, defaultLocalPgData)

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	logPath := path + "/log/" + time.Now().Format("01-02-2006") + ".log"
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if conf.Local.Password == "" {
		exitCode := 2
		log.Printf("[ERROR] Need password for local host. Exit status %d\n", exitCode)
		fmt.Printf("%s [ERROR] Need password for local host. Exit status %d\n", time.Now().Format(time.RFC3339), exitCode)
		os.Exit(exitCode)
	}

	var failedTry, attempt int = 0, 3
	for i := range conf.Remote {
		// if conf.Remote[i].Port == 0 {
		// 	conf.Remote[i].Port = 5432
		// }
		conf.Remote[i].Port = defaultingInt(conf.Remote[i].Port, defaultPort)
		misConfRemote("host", conf.Remote[i].Host, i)
		misConfRemote("password", conf.Remote[i].Password, i)
		conf.Remote[i].User = defaultingStr(conf.Remote[i].User, defaultUser)
	}
	for {
		localCheck := checkDB("local", conf.Local.Host, conf.Local.Port, conf.Local.User, conf.Local.Password)
		if localCheck == false {
			exitCode := 3
			fmt.Printf("%s [ERROR] Your local Database DOWN. Exit status %d\n", time.Now().Format(time.RFC3339), exitCode)
			log.Printf("[ERROR] Your local Database DOWN. Exit status %d\n", exitCode)
			os.Exit(exitCode)
		}
		// var misConf bool
		var remoteCheck []bool
		for i := range conf.Remote {
			remoteCheck = append(remoteCheck, checkDB("remote", conf.Remote[i].Host, conf.Remote[i].Port, conf.Remote[i].User, conf.Remote[i].Password))
		}
		var upServer int
		for i := range remoteCheck {
			if remoteCheck[i] == true {
				upServer++
			}
		}
		if upServer == 0 {
			failedTry++
			attempt--
			// println(failedTry, attempt)
			if attempt <= 0 {
				fmt.Printf("%s [WARN] %d server(s) of %d is UP. Need promotion.\n", time.Now().Format(time.RFC3339), upServer, len(remoteCheck))
				log.Printf("[WARN] %d server(s) of %d is UP. Need promotion.\n", upServer, len(remoteCheck))
				exitCode := 0
				fmt.Printf("%s [WARN] Local database has been Promoted. Exit status %d\n", time.Now().Format(time.RFC3339), exitCode)
				log.Printf("[WARN] Local database has been Promoted. Exit status %d\n", exitCode)
				promote(conf.Local.PgCtl, conf.Local.Port, conf.Local.PgData) // (localpg string, port int, user string, password string)
				// failedTry, attempt = 0, 3
				os.Exit(exitCode)
			} else {
				fmt.Printf("%s [INFO] %d server(s) of %d is UP. Need promotion in %d attempt(s).\n", time.Now().Format(time.RFC3339), upServer, len(remoteCheck), attempt)
				log.Printf("[INFO] %d server(s) of %d is UP. Need promotion in %d attempt(s).\n", upServer, len(remoteCheck), attempt)
			}
			// fmt.Println(failedTry, attempt)
		}
		// fmt.Println(upServer, remoteCheck)
		time.Sleep(5 * time.Second)
		// fmt.Println(remoteCheck)
	}
}

// const (
// 	//dbname   = "postgres"
// 	interval = 6
// 	fail     = 3
// 	//logfile  = "/tmp/logtest"
// )

// func main() {

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
// 		var check1 bool = mainDB(*host1, *port1, *user1, *password1)                 // host, port, user, password
// 		var check2 bool = mainDB(*host2, *port2, *user2, *password2)                 // db 2
// 		var checklocal bool = mainDB(*localhost, *localport, *localuser, *localpass) // db local
// 		//println(check1, check2, checklocal)

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
// 				fmt.Println(time.Now().Format(time.RFC3339), "[ERROR] TIME OUT GO PROMOTE....")
// 				log.Println("[ERROR] TIME OUT GO PROMOTE....")
// 				promote(*localpg, *localport, *localdata) // (localpg string, port int, user string, password string)
// 				// function promote
// 				break
// 			}
// 		}
// 		time.Sleep(interval * time.Second)
// 	}
// }

//func mainGO(host1 string)
// func checkingDB1(err error) bool { // ini ga guna
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

func checkDB(dbType string, host string, port int, user string, password string) bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable connect_timeout=1", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	// defer cancel()
	err = db.Ping()
	if err != nil {
		// fmt.Println(time.Now().Format(time.RFC3339), "[WARN] Database", dbType, "at", host, ":", port, "DOWN")
		fmt.Printf("%s [WARN] Database %s at %s:%d is DOWN \n", time.Now().Format(time.RFC3339), dbType, host, port)
		log.Printf("[WARN] Database %s at %s:%d is DOWN \n", dbType, host, port)
		// log.Println("[WARN] Database", dbType, "at", host, ":", port, "DOWN")
		return false
	}

	fmt.Printf("%s [INFO] Database %s at %s:%d is UP \n", time.Now().Format(time.RFC3339), dbType, host, port)
	log.Printf("[INFO] Database %s at %s:%d is UP \n", dbType, host, port)
	// log.Println("[INFO] Database", dbType, "at", host, ":", port, "UP")
	if dbType == "local" {
		localStatus := checkStatus(db)
		if localStatus == "t" {
			fmt.Println(time.Now().Format(time.RFC3339), "[INFO] Status LOCAL Database on STANDBY MODE")
			log.Println("[INFO] Status LOCAL Database on STANDBY MODE")
		} else {
			fmt.Println(time.Now().Format(time.RFC3339), "[WARN] Status LOCAL Database on MASTER MODE")
			log.Println("[WARN] Status LOCAL Database on MASTER MODE")
		}
	}
	return true
}

// func mainDB(host string, port int, user string, password string) bool {
//

// 	localHostAddr, err := net.InterfaceAddrs()

// 	for i := range localHostAddr {

// 	}
// 	fmt.Println(time.Now().Format(time.RFC3339), "[INFO] Database ", host, ":", port, "UP")
// 	log.Println("[INFO] Database ", host, ":", port, "UP")
// 	return true
// }

func promote(pgCtl string, port int, PgData string) {
	//fmt.Println(localpg, "promote", "-D", localdata)
	run, err := exec.Command(pgCtl, "promote", "-D", PgData).Output() //check this for promote
	//run, err := exec.Command("su", "-", "postgres", "-c", "whoami").Output() //check this for promote
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(run))
}

func checkStatus(db *sql.DB) string {
	rows, err := db.Query("SELECT pg_is_in_recovery from pg_is_in_recovery();")
	//rows, err := db.Query("SELECT * from test;")
	var res string
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&res); err != nil {
			log.Fatalln(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}
	return res
}
