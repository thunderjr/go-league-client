package league_auth

import (
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"time"

	"github.com/abdfnx/gosh"
)

var DEFAULT_PROCESS_NAME = "LeagueClientUx"

type Credentials struct {
	Port        string
	Password    string
	PID         string
	Certificate string
}

type AuthenticationOptions struct {
	AwaitConnection bool
	PoolInterval    *time.Duration
}

type LeagueAuth struct {
	AuthenticationOptions
	AuthSuccess chan bool
	Credentials *Credentials
}

func Init(options AuthenticationOptions) *LeagueAuth {
	if options.PoolInterval == nil {
		options.PoolInterval = new(time.Duration)
		*options.PoolInterval = 5 * time.Second
	}

	return &LeagueAuth{
		AuthSuccess:           make(chan bool),
		AuthenticationOptions: options,
	}
}

func (l *LeagueAuth) Authenticate() {
	portRegex, _ := regexp.Compile(`--app-port=([0-9]+)(?:\s|"|$)`)
	passwordRegex, _ := regexp.Compile(`--remoting-auth-token=(.+?)(?:\s|"|$)`)
	pidRegex, _ := regexp.Compile(`--app-pid=([0-9]+)(?:\s|"|$)`)

	isWindows := runtime.GOOS == "windows"

	var getArgsCommand string

	if isWindows {
		getArgsCommand = "Get-CimInstance -Query \"SELECT * from Win32_Process WHERE name LIKE '" + DEFAULT_PROCESS_NAME + ".exe'\" | Select-Object -ExpandProperty CommandLine"
	} else {
		getArgsCommand = "ps x -o args | grep '" + DEFAULT_PROCESS_NAME + "' | grep -v grep"
	}

	var commandOutput string
	for {
		err, rawOutput, _ := gosh.RunOutput(getArgsCommand)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				continue
			}

			log.Println("Error getting League Client credentials")
			log.Fatalln(err)
		}

		commandOutput = regexp.MustCompile("/\n|\r/g").ReplaceAllString(rawOutput, "")
		if len(commandOutput) != 0 || !l.AwaitConnection {
			break
		}

		time.Sleep(*l.AuthenticationOptions.PoolInterval)
	}

	if len(commandOutput) == 0 {
		log.Fatalln("LeagueClient not found")
		return
	}

	if l.Credentials == nil {
		l.Credentials = &Credentials{}
	}

	l.Credentials.Certificate = RiotCertificate
	l.Credentials.Port = portRegex.FindStringSubmatch(commandOutput)[1]
	l.Credentials.Password = passwordRegex.FindStringSubmatch(commandOutput)[1]
	l.Credentials.PID = pidRegex.FindStringSubmatch(commandOutput)[1]

	l.AuthSuccess <- true
}
