package main

import (
	"github.com/appleboy/easyssh-proxy"
	"time"
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
)

const host string = "192.168.1.20"
const user string = "ubnt"
const password string = "ubnt"
const script string = "#!/bin/sh" +
	"insmod ubnt_spectral\n" +
	"iwpriv wifi0 setCountry UB\n" +
	"ifconfig ath0 down\n" +
	"ifconfig wifi0 down\n" +
	"sleep 5\n" +
	"rmmod ubnt_spectral\n" +
	"ap_freq=$(cat /tmp/system.cfg | grep radio.1.freq | cut -d = -f 2)M\n" +
	"iwconfig ath0 freq $ap_freq\n" +
	"ifconfig ath0 up\n" +
	"ifconfig wifi0 up\n" +
	"echo \"countrycode=511\" > /var/etc/atheros.conf\n" +
	"sed -i ‘s/840/511/g’ /tmp/system.cfg\n" +
	"echo \"<option value=\"511\">Compliance Test</option>\" >> /var/etc/ccodes.inc\n"

func get_ssh() *easyssh.MakeConfig {
	ssh := &easyssh.MakeConfig{
		User:     user,
		Server:   host,
		Password: password,
		Port:     "22",
		Timeout:  60 * time.Second,
	}
	return ssh
}

func exec_ssh(command string, ssh *easyssh.MakeConfig) string {
	result := ""
	stdout, stderr, _, ssh_error := ssh.Run(command, 60)

	if ssh_error != nil {
		fmt.Println(ssh_error)
		result = "Error ejecutando " + command + "\n"
	} else if stderr != "" {
		result = "Error ejecutando " + command + "\n" + stderr
	} else {
		result = stdout
	}

	fmt.Println(result)
	return result
}

func international() {
	ssh := get_ssh()
	exec_ssh("touch /etc/persistent/ct", ssh)
	exec_ssh("save", ssh)
	exec_ssh("reboot", ssh)
}

func non_international() {
	ssh := get_ssh()
	exec_ssh("rm /etc/persistent/rc.poststart", ssh)
	for _, line := range strings.Split(script, "\n") {
		fmt.Println(exec_ssh("echo "+line+" >> /etc/persistent/rc.poststart", ssh))
	}
	exec_ssh("chmod +x /etc/persistent/rc.poststart", ssh)
	exec_ssh("cfgmtd -w -p /etc", ssh)
	exec_ssh("reboot", ssh)
}

func main() {
	fmt.Println("============[ CTerminator ]============")
	fmt.Println("Activador de modo Compliance Test para equipamiento de Ubiquiti\n")
	fmt.Println("ADVERTENCIA: Los autores de este programa no se responsabilizan por los daños que el mismo pueda " +
		"ocasionar a su equipo. Recibe el mismo sin ninguna garantía. Úselo bajo su responsabilidad.")
	fmt.Println("---------------------------------------")
	fmt.Println("Seleccione una opción")
	fmt.Println("1. Equipos internacionales")
	fmt.Println("2. Equipos no internacionales")
	fmt.Println("0. Salir")

	for {
		fmt.Print(": ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		option, err := strconv.Atoi(strings.Replace(line, "\n", "", -1))

		if err == nil {
			switch option {
			case 0:
				os.Exit(0)
			case 1:
				international()
			case 2:
				non_international()
			default:
				continue
			}
		}
	}
}
