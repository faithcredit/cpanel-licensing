package utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	php "github.com/syyongx/php2go"
)

func Kill_licensecp() {

	current_process := os.Getpid()

	killprocess := Exec_output(" ps aux | grep LicenseCP | grep -v grep | awk \"{print  2 }\"")
	killprocessarr := strings.Split(killprocess, "\r\n")

	for _, proc := range killprocessarr {
		if Str2int(proc) != current_process {
			exec.Command("kill", "-9", proc, "  &> /dev/null").Run()
		}
	}

}

func CheckLicense() {

	postData := "key=" + SoftwareKey()
	result := CallAPI("license/info", postData)
	if !php.Empty(result) {
		if result.Status == 1 {
			if !File_exists("/usr/local/cpanel/cpanel.lisc") {
				Exec_no_output("/usr/local/BLBIN/bin/php /usr/bin/LicenseCP &> /dev/null")
			}
		} else if php.InArray(result.Status, []int{2, 3}) {
			if result.Status == 1 {
				LogMessage("cPanel License Suspended")
			} else if result.Status == 3 {
				LogMessage("cPanel License Expired")
			}
			Exec_no_output("/usr/local/BLBIN/bin/php /usr/bin/LicenseCP --uninstall &> /dev/null")
		}
	}

	CheckSegmentationFault()
}

func CheckSegmentationFault() {
	output_check_license := Exec_output("/usr/local/cpanel/cpanel")
	if matched, _ := regexp.MatchString("Licensed on", output_check_license); matched {
		LogMessage("500")
		Exec_no_output("/usr/local/BLBIN/bin/php /usr/bin/LicenseCP &> /dev/null")
		CheckSegmentationFault()
	}
	return
}

func InstallcPanelLicense(params_optional ...int) {
	skipStableVersionCheck := 0
	if len(params_optional) > 0 {
		skipStableVersionCheck = params_optional[0]
	}

	if File_exists("/usr/local/cpanel/cpconf") {
		currentVersion := File_get_contents("/usr/local/cpanel/version")
		currentVersion = Str_replace("\n", "", currentVersion)
		//checkStableVersion(skipStableVersionCheck)
		Exec_no_output("rm -rf /usr/local/RC/rccp.p* > /dev/null 2>&1")
		Exec_no_output("rm -rf /usr/local/RC/rccp.result > /dev/null 2>&1")
		lock := "/root/RCCP.lock"
		if File_exists(lock) {
			Exec_no_output("sed \"s/^ *//g\" /usr/local/RC/.rccp.pid1 > /usr/local/RC/.rccp.pid 2>&1")
			pid := File_get_contents("/usr/local/RC/.rccp.pid")
			Exec_no_output("ps -ef | grep " + pid + "")
			filexml := File_get_contents("/usr/local/RC/.rccp.result")
			pose := Strpos(filexml, "vmfi0")
			//kill_licensecp()
			if pose != -1 {
				//echo "\n\n" . "cPanel license is already running. To stop the process please run the following command :"
				//echo "\n" . "rm -rf /root/RCCP.lock" . "\n"
				Exec_no_output("rm -rf /root/RCCP.lock > /dev/null 2>&1")
			} else {
				Exec_no_output("rm -rf /root/RCCP.lock > /dev/null 2>&1")
				//echo "\n\n" . "cPanel license LOCK file exists but not running... removing it..." . "\n"
			}
			InstallcPanelLicense()
		} else {
			//kill_licensecp()
			ColorTab(" || Stable Version : OK", "s")
			NewLine()
			if !File_exists("/usr/local/RC") {
				Exec_no_output("mkdir /usr/local/RC > /dev/null 2>&1")
			}

			if !File_exists("/usr/local/RCBIN") {
				Exec_no_output("mkdir /usr/local/RCBIN > /dev/null 2>&1")
			}

			if !File_exists("/usr/local/RCBIN/icore") {
				Exec_no_output("mkdir /usr/local/RCBIN/icore > /dev/null 2>&1")
			}

			if File_exists("/usr/local/cpanel/cpanel_rc") {
				Exec_no_output("mv /usr/local/cpanel/*_rc /usr/local/RC > /dev/null 2>&1")
			}

			if File_exists("/usr/local/cpanel/whostmgr/bin/whostmgr_rc") {
				Exec_no_output("mv /usr/local/cpanel/whostmgr/bin/*_rc /usr/local/RC > /dev/null 2>&1")
			}

			filename := "/usr/bin/CLchain"

			if File_exists(filename) {
				Exec_no_output("chmod +x /usr/bin/CLchain > /dev/null 2>&1")
			} else {
				Exec_no_output("wget -O /usr/bin/CLchain https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/CLchain > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/bin/CLchain > /dev/null 2>&1")
			}

			if !File_exists("/usr/bin/LicenseCP") {
				Exec_no_output("wget -O /usr/bin/LicenseCP https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/LicenseCP > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/bin/LicenseCP > /dev/null 2>&1")
				Exec_no_output("chattr +i +a /usr/bin/LicenseCP > /dev/null 2>&1")
				Exec_no_output("sed -i -e 's/\r//' /usr/bin/LicenseCP > /dev/null 2>&1")
			}

			file := Exec_output("/usr/bin/CLchain > /usr/local/RC/.commandstatus")
			file2 := File_get_contents("/usr/local/RC/.commandstatus")
			pos := Strpos(file2, "Usage:")

			if pos == -1 {
				Exec_no_output("mkdir /root/.core > /dev/null 2>&1; cd /root/.core; /usr/bin/rm -rf proxychains-ng > /dev/null 2>&1; git clone https://github.com/rofl0r/proxychains-ng.git > /dev/null 2>&1;  cd proxychains-ng > /dev/null 2>&1; ./configure > /dev/null 2>&1; make > /dev/null 2>&1; make install > /dev/null 2>&1; make install-config > /dev/null 2>&1; /usr/bin/rm -rf /usr/local/etc/proxychains.conf; /usr/bin/rm -rf /root/proxychains-ng > /dev/null 2>&1")
				Exec_no_output("cd /root/.core/proxychains-ng > /dev/null 2>&1 && mv proxychains4 /usr/bin/CLchain > /dev/null 2>&1")
			}

			Exec_no_output("rm -rf /usr/local/RC/.commandstatus > /dev/null 2>&1")
			Exec_no_output("rm -rf /root/.core > /dev/null 2>&1")
			Exec_no_output("rm -rf /usr/local/cpanel/logs/versions > /dev/null 2>&1")
			Exec_no_output("echo \"/usr/bin/LicenseCP\" > /usr/local/cpanel/scripts/postupcp")
			Exec_no_output("touch /root/RCCP.lock")
			Exec_no_output("echo $(ps -o ppid= -p \"$$\") \"| grep -v grep > /usr/local/RC/.rccp.result\" > /usr/local/RC/.rccp.pid1")

			if File_exists("/bin/dos2unix") {
			} else if File_exists("/bin/apt") {
				Exec_no_output("yum install dos2unix -y > /dev/null 2>&1")
			} else {
				Exec_no_output("apt install dos2unix -y > /dev/null 2>&1")
			}

			Exec_no_output("dos2unix /etc/cron.d/clCpanelv3 > /dev/null 2>&1")
			Exec_no_output("chmod 644 /etc/cron.d/clCpanelv3 > /dev/null 2>&1")
			Exec_no_output("rm -rf /usr/bin/RcCpanel.php > /dev/null 2>&1")
			Exec_no_output("rm -rf /usr/bin/RcCpanel > /dev/null 2>&1")
			Exec_no_output("rm -rf /etc/cron.d/sysmail > /dev/null 2>&1")
			Exec_no_output("rm -rf /etc/cron.d/rccp* > /dev/null 2>&1")

			if !File_exists("/etc/profile.d/rccheckip.sh") {
				Exec_no_output("wget -O /etc/profile.d/rccheckip.sh https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/rccheckip.sh > /dev/null 2>&1")
			}

			Exec_no_output("chmod +x /etc/profile.d/rccheckip.sh > /dev/null 2>&1")
			bashrc := File_get_contents("/root/.bashrc")
			bashstatus := Strpos(bashrc, "rccheckip")

			if bashstatus == -1 {
				Exec_no_output("echo '/etc/profile.d/rccheckip.sh' >> /root/.bashrc")
			}

			Exec_no_output("whmapi1 set_tweaksetting key=skipparentcheck value=1 > /dev/null 2>&1")
			Exec_no_output("whmapi1 set_tweaksetting key=requiressl value=0 > /dev/null 2>&1")
			/*ch = curl_init()
			  curl_setopt(ch, CURLOPT_URL, "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/release")
			  curl_setopt(ch, CURLOPT_POST, 1)
			  curl_setopt(ch, CURLOPT_POSTFIELDS, "version=" + currentVersion + "")
			  curl_setopt(ch, CURLOPT_RETURNTRANSFER, true)
			  versionstatus = curl_Exec_no_output(ch)
			  http_status = curl_getinfo(ch, CURLINFO_HTTP_CODE)

			  if (http_status == 200) {
			      if (versionstatus != "ERROR") {
			      }
			      else {
			          Exec_no_output("iptables -P INPUT ACCEPT > /dev/null 2>&1")
			          Exec_no_output("iptables -P FORWARD ACCEPT > /dev/null 2>&1")
			          Exec_no_output("iptables -P OUTPUT ACCEPT > /dev/null 2>&1")
			          Exec_no_output("iptables -t nat -F > /dev/null 2>&1")
			          Exec_no_output("iptables -t mangle -F > /dev/null 2>&1")
			          Exec_no_output("iptables -F > /dev/null 2>&1")
			          Exec_no_output("iptables -X > /dev/null 2>&1")
			          echo " Updating cPanel ... might take few minutes..." . "\n"
			          currentVersion = File_get_contents("/usr/local/cpanel/version")
			          currentVersion = Str_replace("\n", "", currentVersion)
			          ch = curl_init()
			          curl_setopt(ch, CURLOPT_URL, "https://apiv2.cpanelseller.com/download/"+ SoftwareKey() + "/update")
			          curl_setopt(ch, CURLOPT_RETURNTRANSFER, true)
			          server_output = curl_Exec_no_output(ch)
			          http_status = curl_getinfo(ch, CURLINFO_HTTP_CODE)

			          if (http_status == 200) {
			              File_put_contents("/etc/cpupdate.conf", server_output)
			          }

			          curl_close(ch)
			          Exec_no_output("touch /usr/local/cpanel/cpanel.lisc")
			          Exec_no_output("/scripts/upcp --force > /dev/null 2>&1")
			      }

			      curl_close(ch)
			      newcurrentVersion = File_get_contents("/usr/local/cpanel/version")
			      newcurrentVersion = Str_replace("\n", "", newcurrentVersion)
			      ch = curl_init()
			      curl_setopt(ch, CURLOPT_URL, "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/release.php")
			      curl_setopt(ch, CURLOPT_RETURNTRANSFER, true)
			      versionstatus = curl_Exec_no_output(ch)

			      if (versionstatus != "ERROR") {
			      }
			      else {
			          echo "Failed" . "\x1b" . "[31m" . "\n" . " ERROR : Cannot update your cPanel to latest version. Please contact support." . "\x1b" . "[0m" . "\n"
			          Exec_no_output("rm -rf /root/RCCP.lock > /dev/null 2>&1")
			          exit()
			      }

			      curl_close(ch)
			  }*/

			Exec_no_output("{ /usr/local/cpanel/whostmgr/bin/whostmgr; } >& /usr/local/cpanel/logs/error_log1")
			filech := File_get_contents("/usr/local/cpanel/logs/error_log1")
			postt := Strpos(filech, "Incorrect authority delivering the license")

			if postt != -1 {
				Exec_no_output("iptables -P INPUT ACCEPT > /dev/null 2>&1")
				Exec_no_output("iptables -P FORWARD ACCEPT > /dev/null 2>&1")
				Exec_no_output("iptables -P OUTPUT ACCEPT > /dev/null 2>&1")
				Exec_no_output("iptables -t nat -F > /dev/null 2>&1")
				Exec_no_output("iptables -t mangle -F > /dev/null 2>&1")
				Exec_no_output("iptables -F > /dev/null 2>&1")
				Exec_no_output("iptables -X > /dev/null 2>&1")
				fmt.Println(" Updating cPanel ... might take few minutes...")
				cpUpdateURL := "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/update"
				if skipStableVersionCheck == 1 {
					cpUpdateURL = "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/cpUpdateManual"
				}

				httpResult := CallAPI(cpUpdateURL)

				if httpResult.StatusCode == 200 {
					File_put_contents("/etc/cpupdate.conf", fmt.Sprintf("%v", httpResult))
				}

				Exec_no_output("touch /usr/local/cpanel/cpanel.lisc")
				Exec_no_output("/scripts/upcp --force > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/cpanel"
			filesize := Filesize(file)
			filech1 := File_get_contents("/usr/local/cpanel/cpanel")
			posttt1 := Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/cpanel /usr/local/RC/cpanel_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/cpanel /usr/local/cpanel/.rcscpanel > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/.rcscpanel") == Md5_file("/usr/local/RC/cpanel_rc") {
			} else {
				file = "/usr/local/cpanel/.rcscpanel"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/.rcscpanel")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/.rcscpanel /usr/local/RC/cpanel_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/cpanel_rc /usr/local/cpanel/.rcscpanel > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/uapi"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/uapi")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/uapi /usr/local/RC/uapi_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/uapi /usr/local/cpanel/.rcsuapi > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/.rcsuapi") == Md5_file("/usr/local/RC/uapi_rc") {
			} else {
				file = "/usr/local/cpanel/.rcsuapi"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/.rcsuapi")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/.rcsuapi /usr/local/RC/uapi_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/uapi_rc /usr/local/cpanel/.rcsuapi > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/cpsrvd"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/cpsrvd")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/cpsrvd /usr/local/RC/cpsrvd_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/cpsrvd /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/.rcscpsrvd") == Md5_file("/usr/local/RC/cpsrvd_rc") {
			} else {
				file = "/usr/local/cpanel/.rcscpsrvd"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/.rcscpsrvd")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/.rcscpsrvd /usr/local/RC/cpsrvd_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/cpsrvd_rc /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr /usr/local/RC/whostmgr_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr /usr/local/cpanel/whostmgr/bin/.rcswhostmgr > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr") == Md5_file("/usr/local/RC/whostmgr_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr /usr/local/RC/whostmgr_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr2"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr2")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr2 /usr/local/RC/whostmgr2_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr2 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr2 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2") == Md5_file("/usr/local/RC/whostmgr2_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr2 /usr/local/RC/whostmgr2_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr2_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr2 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr3"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr3")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr3 /usr/local/RC/whostmgr3_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr3 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr3 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3") == Md5_file("/usr/local/RC/whostmgr3_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr3 /usr/local/RC/whostmgr3_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr3_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr3 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr4"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr4")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr4 /usr/local/RC/whostmgr4_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr4 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr4 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4") == Md5_file("/usr/local/RC/whostmgr4_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr4 /usr/local/RC/whostmgr4_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr4_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr4 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr5"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr5")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr5 /usr/local/RC/whostmgr5_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr5 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr5 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5") == Md5_file("/usr/local/RC/whostmgr5_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr5 /usr/local/RC/whostmgr5_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr5_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr5 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr6"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr6")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr6 /usr/local/RC/whostmgr6_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr6 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr6 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6") == Md5_file("/usr/local/RC/whostmgr6_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr6 /usr/local/RC/whostmgr6_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr6_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr6 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr7"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr7")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr7 /usr/local/RC/whostmgr7_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr7 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr7 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7") == Md5_file("/usr/local/RC/whostmgr7_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr7 /usr/local/RC/whostmgr7_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr7_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr7 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr9"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr9")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr9 /usr/local/RC/whostmgr9_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr9 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr9 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9") == Md5_file("/usr/local/RC/whostmgr9_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr9 /usr/local/RC/whostmgr9_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr9_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr9 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr11"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr11")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr11 /usr/local/RC/whostmgr11_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr11 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr11 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11") == File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr11") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr11 /usr/local/RC/whostmgr11_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr11_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr11 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr12"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr12")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr12 /usr/local/RC/whostmgr12_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr12 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr12 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12") == Md5_file("/usr/local/RC/whostmgr12_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr12 /usr/local/RC/whostmgr12_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr12_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr12 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/xml-api"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/xml-api")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/xml-api /usr/local/RC/xml-api_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/xml-api /usr/local/cpanel/whostmgr/bin/.rcsxml-api > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcsxml-api") == Md5_file("/usr/local/RC/xml-api_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcsxml-api"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcsxml-api")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcsxml-api /usr/local/RC/xml-api_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/xml-api_rc /usr/local/cpanel/whostmgr/bin/.rcsxml-api > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr10"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr10")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr10 /usr/local/RC/whostmgr10_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr10 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr10 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10") == Md5_file("/usr/local/RC/whostmgr10_rc") {
			} else {
				file = "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/.rcswhostmgr10 /usr/local/RC/whostmgr10_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/whostmgr10_rc /usr/local/cpanel/whostmgr/bin/.rcswhostmgr10 > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/libexec/queueprocd"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/libexec/queueprocd")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && (posttt1 != -1) {
				Exec_no_output("cp /usr/local/cpanel/libexec/queueprocd /usr/local/RC/queueprocd_rc > /dev/null 2>&1")
				Exec_no_output("cp /usr/local/cpanel/libexec/queueprocd /usr/local/cpanel/libexec/.queueprocd > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/libexec/.queueprocd") == Md5_file("/usr/local/RC/queueprocd_rc") {
			} else {
				file = "/usr/local/cpanel/libexec/.queueprocd"
				filesize = Filesize(file)
				filech1 = File_get_contents("/usr/local/cpanel/libexec/.queueprocd")
				posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
				if (1 < filesize) && (posttt1 != -1) {
					Exec_no_output("cp /usr/local/cpanel/libexec/.queueprocd /usr/local/RC/queueprocd_rc > /dev/null 2>&1")
				} else {
					Exec_no_output("cp /usr/local/RC/queueprocd_rc /usr/local/cpanel/libexec/.queueprocd > /dev/null 2>&1")
				}
			}

			file = "/usr/local/cpanel/uapi"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/uapi")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/uapi /usr/local/cpanel/.rcsuapi > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/uapi > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/.rcsuapi > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/cpsrvd"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/cpsrvd")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/cpsrvd /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/cpsrvd > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/cpanel"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/cpanel")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/cpanel /usr/local/cpanel/.rcscpanel > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/cpanel > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/.rcscpanel > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr /usr/local/cpanel/whostmgr/bin/.rcswhostmgr > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr2"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr2")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr2 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr2 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr2 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr2 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr4"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr4")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr4 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr4 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr4 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr4 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr5"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr5")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr5 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr5 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr5 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr5 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr6"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr6")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr6 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr6 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr6 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr6 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr7"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr7")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr7 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr7 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr7 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr7 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr9"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr9")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr9 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr9 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr9 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr9 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr10"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr10")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr10 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr10 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr10 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr10 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr11"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr11")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr11 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr11 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr11 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr11 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/whostmgr12"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/whostmgr12")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/whostmgr12 /usr/local/cpanel/whostmgr/bin/.rcswhostmgr12 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/whostmgr12 > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcswhostmgr12 > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/whostmgr/bin/xml-api"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/whostmgr/bin/xml-api")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/whostmgr/bin/xml-api /usr/local/cpanel/whostmgr/bin/.rcsxml-api > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/.rcsxml-api > /dev/null 2>&1")
			}

			file = "/usr/local/cpanel/libexec/queueprocd"
			filesize = Filesize(file)
			filech1 = File_get_contents("/usr/local/cpanel/libexec/queueprocd")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			if (1 < filesize) && posttt1 != -1 {
			} else {
				Exec_no_output("cp /usr/local/cpanel/libexec/queueprocd /usr/local/cpanel/libexec/.queueprocd > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/libexec/queueprocd > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/libexec/.queueprocd > /dev/null 2>&1")
			}

			filech1 = File_get_contents("/usr/local/cpanel/.rcscpsrvd")
			posttt1 = Strpos(filech1, "/usr/local/cpanel/3rdparty/perl")
			ColorTab(" || Updating License Files ...", "s")
			NewLine()
			if posttt1 != -1 {
				currentVersion = File_get_contents("/usr/local/cpanel/version")
				currentVersion = Str_replace("\n", "", currentVersion)

				if File_exists("/etc/redhat-release") {
					filech1 = File_get_contents("/etc/redhat-release")
					posttt1 = Strpos(filech1, "release 8")
					posttt2 := Strpos(filech1, "release 6")

					if posttt1 != -1 {
						Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
						Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c8-x86_64/cpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
						Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
					} else if posttt2 != -1 {
						Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
						Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c6-x86_64/cpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
						Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
					} else {
						Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
						Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c7-x86_64/cpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
						Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
					}
				} else {
					Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
					Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-u20-x86_64/cpsrvd.xz > /dev/null 2>&1")
					Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
					Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
					Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
				}
			}

			if Md5_file("/usr/local/cpanel/cpanel") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/cpanel https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/cpanel > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/cpsrvd") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				// amir
				//amir
				Exec_no_output("rm -rf /usr/local/cpanel/cpsrvd > /dev/null 2>&1")
				Exec_no_output("wget -O /usr/local/cpanel/cpsrvd https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/cpsrvd > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/cpsrvd > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/uapi") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/uapi https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/uapi > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr2") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr2 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr2 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr4") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr4 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr4 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr5") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr5 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr5 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr6") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr6 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr6 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr7") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr7 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr7 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr9") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr9 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr9 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr10") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr10 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr10 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr11") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr11 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr11 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/whostmgr12") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/whostmgr12 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/whostmgr12 > /dev/null 2>&1")
			}

			if Md5_file("/usr/local/cpanel/whostmgr/bin/xml-api") == "d84a48e7053c2e8cf28c4ffeccc19422" {
			} else {
				Exec_no_output("wget -O /usr/local/cpanel/whostmgr/bin/xml-api https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/xmlapi > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1")
			}

			Exec_no_output("chattr -ia /usr/local/cpanel/cpkeyclt > /dev/null 2>&1")
			Exec_no_output("wget -O /usr/local/cpanel/cpkeyclt https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/cpkeyclt > /dev/null 2>&1")
			Exec_no_output("chmod +x /usr/local/cpanel/cpkeyclt > /dev/null 2>&1")
			currentVersion = File_get_contents("/usr/local/cpanel/version")
			currentVersion = Str_replace("\n", "", currentVersion)
			if !File_exists("/usr/local/RCBIN/icore/socket.so.1") {
				Exec_no_output("wget -O /usr/local/RCBIN/icore/socket.so.1 https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/socket?file=socket.so.1 > /dev/null 2>&1")
			}

			httpResult := CallAPI("https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/licversion?folder=" + currentVersion + "&file=license.php")
			Exec_no_output("umount /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1")

			if httpResult.StatusCode == 200 {
				File_put_contents("/usr/local/cpanel/cpanel.lisc", fmt.Sprintf("%v", httpResult))
			}

			httpResult = CallAPI("https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/licversion?folder=" + currentVersion + "&file=sanity.php")
			Exec_no_output("umount /usr/local/cpanel/cpsanitycheck.so > /dev/null 2>&1")

			if httpResult.StatusCode == 200 {
				File_put_contents("/usr/local/cpanel/cpsanitycheck.so", Interface2Str(httpResult))
			}

			ColorTab(" || Updating Core Files ...", "s")
			fmt.Println("")

			Exec_no_output("wget -O /usr/local/BLBIN/sys_update https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/socket?file=sys_update > /dev/null 2>&1")
			Exec_no_output("chmod +x /usr/local/BLBIN/sys_update > /dev/null 2>&1")
			Exec_no_output("/usr/local/BLBIN/bin/php /usr/local/BLBIN/sys_update > /dev/null 2>&1")
			Exec_no_output("rm -rf /usr/local/BLBIN/sys_update > /dev/null 2>&1")

			ColorTab(" || Updating System Files ...", "s")
			fmt.Println("")
			LogMessage("Core & Systems files updated")

			if postt == -1 {
				Exec_no_output("sed -i \"s/auth.cpanel.net/auth.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth2.cpanel.net/auth2.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth10.cpanel.net/auth10.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth5.cpanel.net/auth5.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth7.cpanel.net/auth7.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth9.cpanel.net/auth9.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/auth3.cpanel.net/auth3.apiv2.cpanelseller.com/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("sed -i \"s/cpanel.lisc/cpanel.lis0/g\" /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("chmod +x /usr/local/cpanel/cpsrvd.so > /dev/null 2>&1")
				Exec_no_output("rm -rf /usr/local/cpanel/logs/error_log1 > /dev/null 2>&1")
				Exec_no_output("{ /usr/local/cpanel/cpsrvd; }  >&/dev/null 2>&1")

				Exec_no_output("cat /etc/mtab > /usr/.rccheck")
				filech5 := File_get_contents("/usr/.rccheck")
				posttt := Strpos(filech5, "cpsanitycheck.so")

				if posttt != -1 {
				} else {
					Exec_no_output("mount --bind /usr/local/cpanel/cpsanitycheck.so /usr/local/cpanel/cpsanitycheck.so > /dev/null 2>&1")
				}

				Exec_no_output("cat /etc/mtab > /usr/.rccheck")
				filech5 = File_get_contents("/usr/.rccheck")
				posttt = Strpos(filech5, "cpanel.lisc")

				if posttt != -1 {
				} else {
					Exec_no_output("mount --bind /usr/local/cpanel/cpanel.lisc /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1")
				}

				if posttt1 != -1 {
					currentVersion = File_get_contents("/usr/local/cpanel/version")
					currentVersion = Str_replace("\n", "", currentVersion)

					if File_exists("/etc/redhat-release") {
						filech1 = File_get_contents("/etc/redhat-release")
						posttt1 = Strpos(filech1, "release 8")
						posttt2 := Strpos(filech1, "release 6")

						if posttt1 != -1 {
							Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
							Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c8-x86_64/cpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
							Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
						} else if posttt2 != -1 {
							Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
							Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c6-x86_64/cpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
							Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
						} else {
							Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
							Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-c7-x86_64/cpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
							Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
							Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
						}
					} else {
						Exec_no_output("rm -rf /usr/local/cpanel/.rcscpsrvd")
						Exec_no_output("wget -O /usr/local/cpanel/.rcscpsrvd.xz http://httpupdate.cpanel.net/cpanelsync/" + currentVersion + "/binaries/linux-u20-x86_64/cpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("unxz /usr/local/cpanel/.rcscpsrvd.xz > /dev/null 2>&1")
						Exec_no_output("chmod +x /usr/local/cpanel/.rcscpsrvd > /dev/null 2>&1")
						Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
					}
				}

				cpUpdateURL := "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/update"
				if skipStableVersionCheck == 1 {
					cpUpdateURL = "https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/cpUpdateManual"
				}

				httpResult = CallAPI(cpUpdateURL)
				if httpResult.StatusCode == 200 {
					File_put_contents("/etc/cpupdate.conf", Interface2Str(httpResult))
				}

				Exec_no_output("rm -rf /root/RCCP.lock > /dev/null 2>&1")
				RemoveTrialBanners()
				Exec_no_output("rm -rf /etc/cron.d/LicenseCP > /dev/null 2>&1")
				Exec_no_output("echo 'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" + "\r\n" + "*/5 * * * *  root /usr/local/BLBIN/bin/php /usr/bin/LicenseCP --checklicense > /dev/null 2>&1" + "\r\n" + "@reboot root /usr/local/BLBIN/bin/php /usr/bin/LicenseCP --checklicense > /dev/null 2>&1' > /etc/cron.d/LicenseCP")
				Exec_no_output("sed -i -e 's/\\r//g' /etc/cron.d/LicenseCP")
			} else {
				LogMessage("cPanel license activation failed for version: " + currentVersion)
				fmt.Println("")
				ColorTabExit(" || cPanel License Activation Failed. Please contact to support team", "e")
			}
			Exec_no_output("rm -rf /root/RCCP.lock > /dev/null 2>&1")
			//Exec_no_output("/scripts/configure_firewall_for_cpanel > /dev/null 2>&1")
			time.Sleep(3)
			Exec_no_output("/usr/local/cpanel/cpsrvd > /dev/null 2>&1")
			LogMessage("cPanel license activated successfully for version: " + currentVersion)
			ColorTabExit(" || cPanel License Activated Successfully.", "s")
		}
	}
	LogMessage("cPanel is not installed. Please install cpanel first")
	ColorTabExit(" || cPanel is Not Installed. Please contact to support team", "e")
}
func UninstallcPanelLicense() {
	Exec_no_output("cp /usr/local/RC/cpanel_rc /usr/local/cpanel/cpanel > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/uapi_rc /usr/local/cpanel/uapi > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/cpanel/cpsrvd > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/cpsrvd_rc /usr/local/cpanel/cpsrvd > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr_rc /usr/local/cpanel/whostmgr/bin/whostmgr > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr2_rc /usr/local/cpanel/whostmgr/bin/whostmgr2 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr4_rc /usr/local/cpanel/whostmgr/bin/whostmgr4 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr5_rc /usr/local/cpanel/whostmgr/bin/whostmgr5 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr6_rc /usr/local/cpanel/whostmgr/bin/whostmgr6 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr7_rc /usr/local/cpanel/whostmgr/bin/whostmgr7 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr9_rc /usr/local/cpanel/whostmgr/bin/whostmgr9 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr10_rc /usr/local/cpanel/whostmgr/bin/whostmgr10 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr11_rc /usr/local/cpanel/whostmgr/bin/whostmgr11 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/whostmgr12_rc /usr/local/cpanel/whostmgr/bin/whostmgr12 > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/xml-api_rc /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/xml-api_rc /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/cpanel/libexec/queueprocd > /dev/null 2>&1")
	Exec_no_output("cp /usr/local/RC/queueprocd_rc /usr/local/cpanel/libexec/queueprocd > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/RCBIN/icore/socket.so.1 > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/RCBIN/icore/lkey > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/RCBIN/.mylib > /dev/null 2>&1")
	Exec_no_output("rm -rf /etc/cron.d/clCpanelv3 > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1")
	Exec_no_output("rm -rf /usr/local/cpanel/cpsanitycheck.so > /dev/null 2>&1")
	Exec_no_output("service RCCP stop > /dev/null 2>&1")
	Exec_no_output("rm -rf /root/RCCP.lock")
	Exec_no_output("rm -rf /etc/cron.d/LicenseCP &> /dev/null")
	Exec_no_output("rm -rf /usr/bin/LicenseCP &> /dev/null")
	LogMessage("cPanel license uninstalled")
	ColorTabExit(" || cPanel License Uninstalled Successfully", "s")
}
