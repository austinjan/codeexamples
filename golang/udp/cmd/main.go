package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// generateSyslogMessage generates a syslog message
func generateSyslogMessage(servity int, msg string) string {
	now := time.Now()
	// get current host url
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	// create a syslog message

	message := fmt.Sprintf("<%d>1 %s %s %s", servity, now.Format(time.RFC3339), host, msg)

	return message
}

func sendUDPMessage(url, message string) error {
	// create udp connection
	conn, err := net.Dial("udp", url)
	if err != nil {
		return err
	}
	defer conn.Close()

	// send message
	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func sendSyslog(url string, servity int, message string) error {
	// create a syslog message
	syslogMessage := generateSyslogMessage(servity, message)
	return sendUDPMessage(url, syslogMessage)
}

func main() {
	// udpctl localhost:8080 {message}
	var servity int
	// sevirity flag 0-7
	flag.IntVar(&servity, "s", 6, "syslog servity (0-7)")
	// type flag
	var typ string
	flag.StringVar(&typ, "t", "syslog", "message type (syslog, plain)")
	flag.Parse()

	// check args
	if len(os.Args) != 3 {
		// print usage
		fmt.Println("Usage: udpctl host:port message")
		os.Exit(1)
	}
	logrus.Debug("args: ", os.Args)
	url := os.Args[1]
	message := os.Args[2]
	// send upd message to  url

	switch typ {
	case "syslog":
		err := sendSyslog(url, servity, message)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	default:
		err := sendUDPMessage(url, message)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}

	fmt.Println("Message sent successfully")
}
