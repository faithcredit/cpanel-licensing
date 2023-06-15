package utils

import (
	"net/url"

	// "os"
	"regexp"
	"strings"
	"time"
)

func InstallSSL() {

	InstallSSL2()
	Exec_no_output("rm -rf /root/.acme.sh/ca.cer &> /dev/null")

	Exec_no_output("rm -rf /root/.acme.sh/cert.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/private.key &> /dev/null")

	Exec_no_output("rm -rf /root/.acme.sh/" + Exec_php("hostname") + " &> /dev/null")

	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/ca.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/cert.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/private.key &> /dev/null")
	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	Exec_no_output("cd /root &> /dev/null")
	Exec_no_output("git clone https://github.com/acmesh-official/acme.sh.git &> /dev/null")
	Exec_no_output("cd ./acme.sh &> /dev/null")
	Exec_no_output("./acme.sh --install -m support@" + Exec_php("hostname") + " &> /dev/null")
	Exec_no_output("cd ~ &> /dev/null")
	Exec_no_output("acme.sh --issue -d " + Exec_php("hostname") + " -w /var/www/html &> /dev/null")
	Exec_no_output("timedatectl set-timezone GMT &> /dev/null")

	folder := Exec_php("find /root/.acme.sh -type d -name " + Exec_php("hostname") + " -print")

	domain := folder[strings.LastIndex(folder, "/")+1:]

	cert := File_get_contents("/root/.acme.sh/" + domain + "/" + Exec_php("hostname") + ".cer")

	key := File_get_contents("/root/.acme.sh/" + domain + "/" + Exec_php("hostname") + ".key")

	//ca = File_get_contents("/root/.acme.sh/thunder.paihost.com/ca.cer");
	Exec_no_output("chmod +x /usr/local/cpanel/cpsrvd &> /dev/null")
	cert1 := url.QueryEscape(cert)
	key1 := url.QueryEscape(key)

	//ca1 = urlencode(ca);
	ColorTab(" || Installing SSL on FTP... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=ftp crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_ftpd &> /dev/null")
	Exec_no_output("/scripts/restartsrv_ftpserver &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on Exim... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=exim crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_exim &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on dovecot... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=dovecot crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_dovecot &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on cPanel... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=cpanel crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Color("OK", "s")
	time.Sleep(1 * time.Second)
	Exec_no_output("service cpanel restart &> /dev/null")
	time.Sleep(3 * time.Second)
	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	LogMessage("SSL installed on cpanel services")
	ColorTabExit(" || Done. Thank you", "w")
}

func InstallSSL2() {

	Exec_no_output("rm -rf /root/.acme.sh/ca.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/cert.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/private.key &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/" + Exec_php("hostname") + " &> /dev/null")
	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/ca.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/cert.cer &> /dev/null")
	Exec_no_output("rm -rf /root/.acme.sh/private.key &> /dev/null")
	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	Exec_no_output("cd /root &> /dev/null")
	Exec_no_output("git clone https://github.com/Neilpang/acme.sh.git &> /dev/null")
	Exec_no_output("cd ./acme.sh &> /dev/null")
	Exec_no_output("./acme.sh --install -m support@" + Exec_php("hostname") + " &> /dev/null")
	Exec_no_output("cd ~ &> /dev/null")
	Exec_no_output("acme.sh --issue -d " + Exec_php("hostname") + " -w /var/www/html &> /dev/null")
	Exec_no_output("timedatectl set-timezone GMT &> /dev/null")

	cert := File_get_contents("/root/.acme.sh/" + Exec_php("hostname") + "_ecc/" + Exec_php("hostname") + ".cer")

	key := File_get_contents("/root/.acme.sh/" + Exec_php("hostname") + "_ecc/" + Exec_php("hostname") + ".key")

	//ca = file_get_contents("/root/.acme.sh/thunder.paihost.com/ca.cer");
	Exec_no_output("chmod +x /usr/local/cpanel/cpsrvd &> /dev/null")
	cert1 := url.QueryEscape(cert)
	key1 := url.QueryEscape(key)
	//ca1 = urlencode(ca);
	ColorTab(" || Installing SSL on FTP... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=ftp crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_ftpd &> /dev/null")
	Exec_no_output("/scripts/restartsrv_ftpserver &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on Exim... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=exim crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_exim &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on dovecot... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=dovecot crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Exec_no_output("/scripts/restartsrv_dovecot &> /dev/null")
	Color("OK", "s")
	ColorTab(" || Installing SSL on cPanel... ", "nc", "0")
	Exec_no_output("/usr/sbin/whmapi1 install_service_ssl_certificate service=cpanel crt=" + cert1 + " key=" + key1 + " &> /dev/null")
	Color("OK", "s")
	time.Sleep(1 * time.Second)
	Exec_no_output("service cpanel restart &> /dev/null")
	time.Sleep(3 * time.Second)
	Exec_no_output("rm -rf /root/acme.sh &> /dev/null")
	LogMessage("SSL installed on cpanel services")
	ColorTabExit(" || Done. Thank you", "w")
}
func InstallLetsencryptSSL() {
	Exec_no_output("/usr/local/cpanel/scripts/install_lets_encrypt_autossl_provider")
	LogMessage("Letsencrypt autossl installed")
	ColorTabExit(" || Done. Thank you", "w")
}

func InstallFleetSSL() {
	Exec_no_output("rm -rf  /etc/letsencrypt-cpanel.licence  &> /dev/null")
	Exec_no_output("wget -q -O /etc/letsencrypt-cpanel.licence --no-check-certificate https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/letsencryptcpanellicense &> /dev/null")
	Exec_no_output("rpm -i --no-check-certificate https://apiv2.cpanelseller.com/download/" + SoftwareKey() + "/gitssl &> /dev/null")
	output_check_license := Exec_output("le-cp self-test")
	if matched, _ := regexp.MatchString("Licensed on", output_check_license); matched {
		LogMessage("Fleetssl installed")
		ColorTabExit(" || Done. Thank you", "w")
	}
	LogMessage("Fleetssl not installed")
	ColorTabExit(" || FleetSSL not installed", "e")
}
func InstallWordpressToolkit() {
	Exec_php("wget -O /root/wordpresstoolkit.sh https://wp-toolkit.plesk.com/cPanel/installer.sh > /dev/null 2>&1")
	Exec_php("chmod +x /root/wordpresstoolkit.sh > /dev/null 2>&1")
	Exec_php("/root/wordpresstoolkit.sh > /dev/null 2>&1")
	Exec_php("rm -rf /root/wordpresstoolkit.sh > /dev/null 2>&1")
	LogMessage("Wordpress toolkit installed")
	ColorTabExit(" || Done. Thank you", "w")
}
