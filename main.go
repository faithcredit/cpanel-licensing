package main

import (

	// "bytes"
	// "time"

	"fmt"
	"lisa/utils"
	"os"
	"strings"

	php "github.com/syyongx/php2go"
	// "log"
	// "os/exec"
)

func MainFunc(input string) {

	utils.NewLine()

	postData := "key=" + utils.SoftwareKey()
	resultData := utils.CallAPI("license/info", postData)

	var expiry, brand, website, currentIP, version, todayDate string

	if !php.Empty(resultData) {

		if resultData.Status == 1 {
			if resultData.Expiry != "" {
				expiry = utils.FormatDateTime(resultData.Expiry)
			} else {
				expiry = ""
			}
			if resultData.Brand != "" {
				brand = resultData.Brand

			} else {
				brand = ""
			}
			if resultData.Url != "" {
				website = resultData.Url

			} else {
				website = ""
			}
			if resultData.Ip != "" {
				currentIP = resultData.Ip

			} else {
				currentIP = ""
			}
			if resultData.Version != "" {
				version = resultData.Version

			} else {
				version = ""
			}
			if resultData.CurrentDate != "" {
				todayDate = utils.FormatDateTime(resultData.CurrentDate)
			} else {
				todayDate = ""
			}

			var hostname, cpanelVersion, kernel, totalAccounts string
			hostname = utils.Exec_php("hostname")
			cpanelVersion = utils.Exec_php("cat /usr/local/cpanel/version")
			kernel = utils.Exec_php("uname -r")
			totalAccounts = utils.Exec_php(`find "/var/cpanel/users" -maxdepth 1 -type f -print | wc -l`)
			serverType := "cPanel"
			utils.ColorTab(" || Welcome to "+brand+" Licensing System ", "s", "0")
			utils.Color("(v"+version+")", "i", "1")
			utils.ColorTab(" || Thank you for using our licensing System", "s")
			utils.NewLine()
			utils.Line()
			utils.ColorTab(" || Our Website :        "+website, "w")
			utils.ColorTab(" || cPanel Version :     "+cpanelVersion, "w", "0")

			utils.ColorTab("|| Kernel Version :     "+kernel, "w", "0")
			utils.ColorTab("|| Total Accounts :     "+totalAccounts, "w", "0")
			utils.ColorTab("|| License Type :       "+serverType, "w")
			utils.ColorTab(" || Server IPv4 :        "+currentIP, "w")
			utils.ColorTab(" || Hostname :           "+hostname, "w", "0")
			utils.ColorTab("|| Today Date:         "+todayDate, "w")
			utils.ColorTab(" || Expiry Date:        "+expiry, "w")
			utils.Line()
			utils.NewLine()
			if input == "--cplicense" || input == "--skip-stable" {
				utils.ColorTab(" || Installing cPanel License ...", "nc")
				utils.NewLine()
				//skipStableVersionCheck = (input == "--skip-stable") ? 1 : 0;
				skipStableVersionCheck := 1

				utils.InstallcPanelLicense(skipStableVersionCheck)

			} else if input == "--installssl" || input == "--installssl2" {

				utils.ColorTab(" || Installing SSL on cPanel Services ...", "nc")
				utils.NewLine()
				if input == "--installssl" {
					utils.InstallSSL()
				} else {
					utils.InstallSSL2()
				}
			} else if input == "--update" {
				utils.ColorTab(" || Updating cPanel ...", "nc")
				utils.NewLine()
				utils.Exec_no_output("/scripts/upcp --force")
			} else if input == "--fleetssl" {
				utils.ColorTab(" || Installing FleetSSL ...", "nc")
				utils.NewLine()
				utils.InstallFleetSSL()
			} else if input == "--wordpress-toolkit" {
				utils.ColorTab(" || Installing Wordpress Toolkit ...", "nc")
				utils.NewLine()
				utils.InstallWordpressToolkit()
			} else if input == "--letsencrypt" {
				utils.ColorTab(" || Installing Letsencrypt SSL ...", "nc")
				utils.NewLine()
				utils.InstallLetsencryptSSL()
			} else if input == "--uninstall" {
				utils.ColorTab(" || Uninstalling cPanel License ...", "nc")
				utils.NewLine()
				utils.UninstallcPanelLicense()
			} else if input == "--help" {
				utils.ColorTab(" || --installssl", "s", "0")
				utils.Color("       :   Install SSL on cPanel Services", "nc", "1")
				utils.ColorTab(" || --update", "s", "0")
				utils.Color("           :   Update cPanel", "nc", "1")
				utils.ColorTab(" || --fleetssl", "s", "0")
				utils.Color("         :   Install FleetSSL on cPanel", "nc", "1")
				utils.ColorTab(" || --wordpress-toolkit", "s", "0")
				utils.Color(":   Install Wordpress Toolkit on cPanel", "nc", "1")
				utils.ColorTab(" || --letsencrypt", "s", "0")
				utils.Color("      :   Install Letsencrypt SSL on cPanel", "nc", "1")
				utils.ColorTab(" || --uninstall", "s", "0")
				utils.Color("        :   Uninstall the cPanel License", "nc", "1")
				utils.NewLine()
				os.Exit(0)
			}
			os.Exit(0)
		} else if php.InArray(resultData.Status, []int{2, 3}) {
			utils.UninstallcPanelLicense()
		} else {
			str := resultData.Expiry
			utils.ColorTabExit(" || Error: X501. "+str, "e")
		}
	} else {
		utils.ColorTabExit("", " || Error: X4. Null response from server", "e")
	}

}

func main() {
	// Enter commmand
	// scanner := bufio.NewScanner(os.Stdin)
	// utils.Color(" Default Command:--cPlicense", "s")
	// utils.Line()
	// utils.Color(" Enter command:", "s", "0")
	// scanner.Scan()
	// var input string = scanner.Text()

	// if input == "" {
	// 	input = "--cPlicense"
	// }

	// utils.Color(" Command:\t"+input, "s")
	// utils.Line()
	var input string = "--cPlicense"
	if len(os.Args) > 1 {
		input = os.Args[1]
	}
	input = strings.ToLower(input)

	if !php.InArray(input, []string{"--cplicense", "--installssl", "--installssl2", "--uninstall", "--fleetssl", "--letsencrypt", "--wordpress-toolkit", "--update", "--help", "--skip-stable", "--checklicense", "--force"}) {
		utils.NewLine()
		utils.ColorTabExit(" || Invalid Argument", "e")
	}
	if !utils.File_exists("/usr/local/BLBIN") {
		fmt.Print("/usr: ")

		utils.Exec_no_output("cd /usr/local/ &> /dev/null")
		fmt.Print("/wget: ")

		utils.Exec_no_output("wget -q -O clbinary.tar.gz --no-check-certificate https://apiv2.cpanelseller.com/download/" + utils.SoftwareKey() + "/blbin &> /dev/null")
		fmt.Print("/tar: ")

		utils.Exec_no_output("tar -xpf clbinary.tar.gz --directory /usr/local &> /dev/null")

	}

	if input == "--checklicense" {
		utils.CheckLicense()
	} else {
		if input == "--force" {
			utils.Kill_licensecp()
			input = "--cplicense"
		}
		MainFunc(input)
	}
}
