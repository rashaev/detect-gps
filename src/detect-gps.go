package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/adrianmo/go-nmea"
	flag "github.com/spf13/pflag"
	"go.bug.st/serial.v1"
)

var baudrate int
var timeOut time.Duration

func init() {
	flag.IntVarP(&baudrate, "baudrate", "b", 115200, "baud rate")
	flag.DurationVarP(&timeOut, "timeout", "t", 5s*time.Second, "max time in seconds to read data from serial port")
}

func readPort(ttyport string, timeout time.Duration) string {
	mode := &serial.Mode{
		BaudRate: baudrate,
	}

	chResult := make(chan string, 1)

	port, err := serial.Open(ttyport, mode)
	if err != nil {
		return ""
	}

	defer port.Close()

	regexGps := regexp.MustCompile(`GNRMC`)

	go func() {
		for {
			reader := bufio.NewReader(port)
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				nextString := scanner.Text()

				newstr := regexGps.FindString(nextString)
				if newstr != "" {
					sentence, err := nmea.Parse(nextString)
					if err == nil {
						if sentence.DataType() == nmea.TypeRMC {
							m := sentence.(nmea.RMC)
							chResult <- fmt.Sprintf("Latitude:%s;Longitude:%s;%s", nmea.FormatGPS(m.Latitude), nmea.FormatGPS(m.Longitude), ttyport)
							break
						}
					}
				}
			}
		}
	}()

	select {
	case data := <-chResult:
		return data
	case <-time.After(timeout):
		return ""
	}
}

func main() {
	flag.Parse()

	var ttySerial []string

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	re := regexp.MustCompile(`^/dev/ttyS[0-9]+$`)

	for _, port := range ports {
		if re.MatchString(port) == true {
			ttySerial = append(ttySerial, port)
		}
	}

	for _, ttyport := range ttySerial {
		strRes := readPort(ttyport, timeOut)
		if strRes != "" {
			fmt.Println(strRes)
			break
		}
	}
}
