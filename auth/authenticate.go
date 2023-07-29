package league_auth

import (
	"context"
	"log"
	"regexp"
	"runtime"
	"time"

	"github.com/abdfnx/gosh"
)

var DEFAULT_PROCESS_NAME = "LeagueClientUx"

type AuthenticationOptions struct {
	AwaitConnection bool
	Timeout         time.Duration
}

type Credentials struct {
	Port        string
	Password    string
	PID         string
	Certificate string
}

type LeagueAuth struct {
	AuthenticationOptions
	Credentials *Credentials
}

func Init(options AuthenticationOptions) *LeagueAuth {
	return &LeagueAuth{
		AuthenticationOptions: options,
	}
}

func (l *LeagueAuth) Authenticate(ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, l.Timeout)
	defer cancel()

	if l.Credentials == nil {
		l.Credentials = &Credentials{}
	}

	portRegex, _ := regexp.Compile(`--app-port=([0-9]+)(?:\s|"|$)`)
	passwordRegex, _ := regexp.Compile(`--remoting-auth-token=(.+?)(?:\s|"|$)`)
	pidRegex, _ := regexp.Compile(`--app-pid=([0-9]+)(?:\s|"|$)`)

	isWindows := runtime.GOOS == "windows"

	var getArgsCommand string

	if isWindows {
		getArgsCommand = "Get-CimInstance -Query \"SELECT * from Win32_Process WHERE name LIKE '" + DEFAULT_PROCESS_NAME + ".exe'\" | Select-Object -ExpandProperty CommandLine"
	} else {
		getArgsCommand = "ps x -o args | grep '" + DEFAULT_PROCESS_NAME + "'"
	}

	err, rawOutput, _ := gosh.RunOutput(getArgsCommand)
	if err != nil {
		log.Println("Error getting League Client credentials")
		log.Fatalln(err)
	}

	commandOutput := regexp.MustCompile("/\n|\r/g").ReplaceAllString(rawOutput, "")

	if len(commandOutput) == 0 {
		log.Fatalln("LeagueClient not found")
	}

	l.Credentials.Certificate = RiotCertificate
	l.Credentials.Port = portRegex.FindStringSubmatch(commandOutput)[1]
	l.Credentials.Password = passwordRegex.FindStringSubmatch(commandOutput)[1]
	l.Credentials.PID = pidRegex.FindStringSubmatch(commandOutput)[1]
}
