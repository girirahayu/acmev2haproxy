package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"fmt"
	"time"
)

//type Configs struct {
//	Configs []Config `json:"configs"`
//}

type Config struct {
	Domain      []string `json:"domain"`
	SslPath     string   `json:"ssl_path"`
	CertbotPath string   `json:"certbot_path"`
	Email		string	 `json:"email"`
}

func RunCMD(path string, args []string, debug bool) (out string, err error) {
	cmd := exec.Command(path, args...)
	var b []byte
	b, err = cmd.CombinedOutput()
	out = string(b)

	if debug {
		fmt.Println(strings.Join(cmd.Args[:], " "))
		if err != nil {
				fmt.Println(out)
		}
	}
	return
}


func main(){
	jsonFile, e := os.Open("config.json")
	if e != nil {
		fmt.Println(e)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var configs Config

	json.Unmarshal(byteValue, &configs)

	for i := 0; i < len(configs.Domain); i++ {

		f, _ := os.OpenFile("/tmp/ssl.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if _, e := os.Stat("/etc/letsencrypt/live/" + configs.Domain[i] + "/fullchain.pem"); os.IsNotExist(e) {

			args := []string{"certonly", "--standalone", "--agree-tos", "--email", configs.Email, "--http-01-port=402", "-d", configs.Domain[i]}
			//output, err := RunCMD("ls", args, true)
			output, ee := RunCMD(configs.CertbotPath, args, true)
			dt := time.Now()

			if ee != nil {
				f.WriteString(dt.Format("January 02, 2006 15:04:05 Monday") + "\n")
				f.WriteString(output+"\n")
			}else{

				f.WriteString(dt.Format("January 02, 2006 15:04:05 Monday") + "\n")
				f.WriteString(output+"\n")

				fcert, _ := os.Create(configs.SslPath + "/" +configs.Domain[i])
				cert, _ := ioutil.ReadFile("/etc/letsencrypt/live/" + configs.Domain[i] + "/fullchain.pem")
				key, _ := ioutil.ReadFile("/etc/letsencrypt/live/" + configs.Domain[i] + "/privkey.pem")
				fcert.WriteString(string(cert) + string(key))
				dt := time.Now()
				if _, e := os.Stat(configs.SslPath + "/" + configs.Domain[i]); !os.IsNotExist(e) {

					f.WriteString(dt.Format("January 02, 2006 15:04:05 Monday") + "\n")
					f.WriteString(configs.Domain[i]+" allready combined\n")
				}

			}


		}else{

			fcert, _ := os.Create(configs.SslPath + "/" +configs.Domain[i])
			cert, _ := ioutil.ReadFile("/etc/letsencrypt/live/" + configs.Domain[i] + "/fullchain.pem")
			key, _ := ioutil.ReadFile("/etc/letsencrypt/live/" + configs.Domain[i] + "/privkey.pem")
			fcert.WriteString(string(cert) + string(key))
			dt := time.Now()
			if _, e := os.Stat(configs.SslPath + "/" + configs.Domain[i]); !os.IsNotExist(e) {

				f.WriteString(dt.Format("January 02, 2006 15:04:05 Monday") + "\n")
					f.WriteString(configs.Domain[i]+" allready combined\n")
				}


		}

		defer f.Close()
	}

}

