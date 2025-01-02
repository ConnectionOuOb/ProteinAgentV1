//##### package & import
package main


import (
	"os"
	"fmt"
	"log"
	"time"
	"os/exec"
	"net/http"
	"io/ioutil"
	"encoding/hex"
	"crypto/sha256"
	"html/template"
	"path/filepath"
)



//##### Struct & Const
var numParrell int
var dirData, dirAFPG, dirAFDB, endDate string

const (
	PROGRAM = "Alphafold Web Server        "
	VERSION = "1.0                         "
	PRGDATE = "2022.02                     "
	AUTHORS = "LI,YAN-JIE                  "
)



//##### Functions
//### Error Report
func ExportError(funcName string, err error) {
	if err != nil {
		log.Fatal("\n# Error - ", funcName," :\n< ", err, " >\n")
	} //if
} //func ExportError


//### Router
func routePage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(16)

	switch r.Method {
		case "GET" :
			t, _ := template.ParseFiles("user.html")
			t.Execute(w, nil)

		case "POST" :
			if len(r.Form["UserSequ"]) == 0 { return }

			userTime := time.Now().Format("2006-01-02_15-04-05")
			userSequ := r.Form["UserSequ"][0]

			sessionID := Sha256Encoding(userTime + userSequ)

			dirSession := filepath.Join(dirData, sessionID)
			CreateDirIfNotExist(dirSession)

			http.Redirect(w, r, "http://120.126.17.200:8082/" + sessionID, 302)

			pathFasta := filepath.Join(dirSession, "querySeq.fasta")
			fileConts := ">" + userTime + "\n" + userSequ + "\n"

			ioutil.WriteFile(pathFasta, []byte(fileConts), 0777)

			RunAlphaFold(dirSession)

	} //switch r.Method
} //func routePage


//### System
func CreateDirIfNotExist(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0777)
		ExportError("os.MkdirAll", err)
	} //if
} //func CreateDirIfNotExist

func Sha256Encoding(src string) string {
    m := sha256.New()
    m.Write([]byte(src))

    return hex.EncodeToString(m.Sum(nil))
} //func Sha256Encoding

func InitialProgram(Argvs []string, lenArgs int) {
	if len(Argvs) < lenArgs {
		log.Fatal("\n# Usage :\n " + Argvs[0] + " dir_data")
	} //if

	CreateDirIfNotExist(Argvs[1])

	dirData = Argvs[1]
	dirAFPG = "/SSD_Intel/dir_ssd/AFserver/alphafold/docker/run_docker.py"
	dirAFDB = "/SSD_Intel/dir_ssd/AFserver/AFDB"
	endDate = "--max_template_date=2021-12-31"
} //func IniPrg

func RunAlphaFold(dirSession string) bool {
	pathFasta  := "--fasta_paths=" + filepath.Join(dirSession, "querySeq.fasta")
	pathData   := "--data_dir=" + dirAFDB
	pathOutput := "--output_dir=" + dirSession

	for true {
		if numParrell == 0 {
			fmt.Println(dirSession, "Start")

			numParrell = 1

			cmd := exec.Command("python3", dirAFPG, pathFasta, pathData, pathOutput, endDate)
			cmd.Start()
			cmd.Wait()

			numParrell = 0

			fmt.Println(dirSession, "End")

			break
		} else {
			time.Sleep(30 * time.Second) 
		}
	}

	return true
}



//##### Main
func main() {
	//### Start Servers
	fmt.Println("┌---------- Program Information ------------┐")
	fmt.Println("| Name     :  ", PROGRAM,                  "|")
	fmt.Println("| Version  :  ", VERSION,                  "|")
	fmt.Println("| Date     :  ", PRGDATE,                  "|")
	fmt.Println("| Authors  :  ", AUTHORS,                  "|")
	fmt.Println("└-------------------------------------------┘")

	//### Setup
	InitialProgram(os.Args, 2)

	//### Routers
	http.HandleFunc("/", routePage)

	//### Listen at TCL Port
	err := http.ListenAndServe(":8081", nil)
	ExportError("http.ListenAndServe", err)
} //func main()


