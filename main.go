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

			args := []string{"certonly", "--standalone","--preferred-challenges"," http"," --http-01-port","402", "-d", configs.Domain[i]}
			//output, err := RunCMD("ls", args, true)
			output, _ := RunCMD(configs.CertbotPath, args, true)
			dt := time.Now()


			fmt.Println(output)
			f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
			f.WriteString(output)
			f.WriteString("\n")

		}else{

			args := []string{"/etc/letsencrypt/live/" + configs.Domain[i] + "/fullchain.pem","/etc/letsencrypt/live/" +
				configs.Domain[i] + "/privkey.pem",">",configs.SslPath+"/"+configs.Domain[i]}
			//output, err := RunCMD("ls", args, true)
			output, _ := RunCMD("cat", args, true)
			dt := time.Now()
			fmt.Println(output)
			f.WriteString(dt.Format("01-01-2020 15:04:05 Monday") + "\n")
			f.WriteString(output)
			f.WriteString("\n")

		}

		//if err != nil {
		//	f.WriteString(dt.Format("01-02-2006 15:04:05"))
		//	f.WriteString(output)
		//} else {

		//}
		defer f.Close()
	}

}

