package utils

import (
	"encoding/json"
	"fmt"

	// "os"
	// "bytes"
	"os/exec"
	"strings"

	curl "github.com/andelf/go-curl"
	php "github.com/syyongx/php2go"
)

const SysShellToUse = "bash"

func Exec_php(command string) string {
	cmd, _ := exec.Command("bash", "-c", command).Output()
	NewLine()
	return string(cmd)
}

func Exec_no_output(command string) {
	cmd, _ := exec.Command("bash", "-c", command).Output()
	fmt.Printf("%s", cmd)
}

func CallAPI(params_optional ...string) Response {
	urlPost := ""
	postData := ""

	if len(params_optional) > 0 {
		urlPost = params_optional[0]
	}
	if len(params_optional) > 1 {
		postData = params_optional[1]
	}

	url := GetURL() + urlPost

	easy := curl.EasyInit()
	defer easy.Cleanup()
	easy.Setopt(curl.OPT_URL, url)

	if !php.Empty(postData) {
		easy.Setopt(curl.OPT_POSTFIELDS, postData)
		easy.Setopt(curl.OPT_CUSTOMREQUEST, "POST")
	} else {
		easy.Setopt(curl.OPT_CUSTOMREQUEST, "GET")

	}

	easy.Setopt(curl.OPT_USERAGENT, "QConnect")
	easy.Setopt(curl.OPT_TIMEOUT, 30)
	easy.Setopt(curl.OPT_SSL_VERIFYHOST, false)
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	// easy.Setopt(curl.CURLOPT_RETURNTRANSFER, true)
	easy.Setopt(curl.OPT_HEADER, true)
	// easy.Setopt(curl.OPT_RETURNTRANSFER, 1) //
	// make a callback function
	var bodyText []byte
	fooTest := func(buf []byte, userdata interface{}) bool {
		// println("DEBUG: size=>", len(buf))
		// println("DEBUG: content=>", string(buf))
		bodyText = buf
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	var respData Response
	err := json.Unmarshal(bodyText, &respData)

	if err != nil {
		fmt.Println("err:", err)
	}

	// get the response code
	respCode, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	// respData.StatusCode = fmt.Sprintf("%v", respCode)
	respData.StatusCode = Str2int(Interface2Str(respCode))

	return respData
}

func LogMessage(message string) int {
	postData := "key==" + SoftwareKey() + "&message=" + message
	CallAPI("license/info", postData)
	return 1
}

func CheckStableVersion(skipStableVersionCheck int) {
	if skipStableVersionCheck == 0 {
		var stableVersion = ""
		var currentversion = ""
		var cv = ""
		var sv = ""
		var ncv = 0
		var nsv = 0
		easy := curl.EasyInit()
		easy.Setopt(curl.OPT_URL, "https://apiv2.cpanelseller.com/download/"+SoftwareKey()+"/stableversion")
		// easy.Setopt(curl.CURLOPT_RETURNTRANSFER, true)
		///////////////////////
		var server_output string
		fooTest := func(buf []byte, userdata interface{}) bool {
			// println("DEBUG: size=>", len(buf))
			// println("DEBUG: content=>", string(buf))
			server_output = string(buf)
			return true
		}

		easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

		if err := easy.Perform(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
		///////////////////////////
		status, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		if err != nil {
			fmt.Println(err)
		}
		http_status := status

		if http_status == 200 {
			stableVersion = server_output
		}
		defer easy.Cleanup()

		currentversion = File_get_contents("/usr/local/cpanel/version")
		currentversion = Str_replace("\n", "", currentversion)
		if !php.Empty(stableVersion) && !php.Empty(currentversion) {
			cv = Str_replace(".", "", currentversion)
			sv = Str_replace(".", "", stableVersion)
			ncv = Str2int(cv)
			nsv = Str2int(sv)

			if ncv != nsv {
				easy := curl.EasyInit()
				easy.Setopt(curl.OPT_URL, "https://apiv2.cpanelseller.com/download/"+SoftwareKey()+"/update")
				// easy.Setopt(curl.CURLOPT_RETURNTRANSFER, true)

				var server_output string
				fooTest := func(buf []byte, userdata interface{}) bool {
					// println("DEBUG: size=>", len(buf))
					// println("DEBUG: content=>", string(buf))
					server_output = string(buf)
					return true
				}

				easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

				if err := easy.Perform(); err != nil {
					fmt.Printf("ERROR: %v\n", err)
				}

				status, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
				if err != nil {
					fmt.Println(err)
				}
				http_status := status
				defer easy.Cleanup()

				if http_status == 200 {
					File_put_contents("/etc/cpupdate.conf", fmt.Sprintf("%v", server_output))
				}
				ColorTab("|| Updating cPanel... It may takes 5 to 10 minutes...", "w")
				NewLine()
				Exec_no_output("/scripts/upcp --force")
			}
		}
	}
}
func RemoveTrialBanners() {
	Exec_no_output(`echo "" > /usr/local/cpanel/whostmgr/docroot/templates/menu/_trial.tmpl &> /dev/null`)
	Exec_no_output(`sed -i -s "s/_is_trial/_is_tria1/g" "/usr/local/cpanel/base/show_template.stor" &> /dev/null`)
	Exec_no_output(`sed -i -s "s/IS_TRIAL/IS_TRIA1/g" "/usr/local/cpanel/base/resetpass.cgi" &> /dev/null`)
	Exec_no_output(`sed -i -e "s/CPANEL.CPFLAGS.item(\"trial\")/False/g" "/usr/local/cpanel/base/frontend/paper_lantern/_assets/master_glass/master_content.html.tt" &> /dev/null`)
	Exec_no_output(`sed -i -e "s/CPANEL.CPFLAGS.item(\"trial\")/False/g" "/usr/local/cpanel/base/frontend/paper_lantern/_assets/master_retro/master_content.html.tt" &> /dev/null`)
	Exec_no_output(`sed -i -e "s/CPANEL.CPFLAGS.item(\"trial\")/False/g" "/usr/local/cpanel/base/frontend/paper_lantern/_assets/master_content.html.tt" &> /dev/null`)
	Exec_no_output(`sed -i -e "s/CPANEL.CPFLAGS.item(\"trial\")/False/g" "/usr/local/cpanel/base/frontend/jupiter/_assets/master_retro/master_content.html.tt" &> /dev/null`)
	Exec_no_output(`sed -i -e "s/CPANEL.CPFLAGS.item(\"trial\")/False/g" "/usr/local/cpanel/base/frontend/jupiter/_assets/master_content.html.tt" &> /dev/null`)
}

func Exec_output(cmd string) string {
	out, _ := exec.Command("bash", "-c", cmd).Output()
	output := strings.Split(string(out), "\n")
	if len(output) > 1 {
		return strings.Join(output, "\r\n")
	}
	return output[0]
}
