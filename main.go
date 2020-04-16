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

			args := []string{"certonly", "--standalone", "--agree-tos", "--email","admin@codigo.id", "--http-01-port=402", "-d", configs.Domain[i]}
			//output, err := RunCMD("ls", args, true)
			output, ee := RunCMD(configs.CertbotPath, args, true)
			dt := time.Now()


			if ee != nil {
				f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
				f.WriteString(output+"\n")
			}else{

				f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
				f.WriteString(output+"\n")

				args := []string{configs.SslPath+" "+configs.Domain[i]}
				output, er := RunCMD("./combine.sh", args, true)
				dt := time.Now()

				if er != nil {
					f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
					f.WriteString(output+"\n")
				}else{
					if _, e := os.Stat(configs.SslPath + "/" + configs.Domain[i]); !os.IsNotExist(e) {
						f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
						f.WriteString(configs.Domain[i] + " Combine Aman.\n")
					}
				}
			}



		}else{

			args := []string{configs.SslPath+" "+configs.Domain[i]}
			output, er := RunCMD("./combine.sh", args, true)
			dt := time.Now()

			if er != nil {
				f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
				f.WriteString(output+"\n")
			}else{
				if _, e := os.Stat(configs.SslPath + "/" + configs.Domain[i]); !os.IsNotExist(e) {
					f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
					f.WriteString(output+"\n")
				}
			}

		}


		//}
		defer f.Close()
	}

}

