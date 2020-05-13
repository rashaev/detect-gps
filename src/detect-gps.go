package main

import (
	"go.bug.st/serial.v1"
	"github.com/adrianmo/go-nmea"
	"fmt"
	"bufio"
	"log"
	"regexp"
	"time"
)


func readPort(ttyport string, timeout time.Duration) string {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	chResult := make(chan string, 1)
	
	port, err := serial.Open(ttyport, mode)
		if err != nil {
			return ""
		}
	defer port.Close()

	regexGps := regexp.MustCompile(`GNRMC`)

	go func(){
		for {
			s, _ :=  bufio.NewReader(port).ReadString('\n')
			newstr := regexGps.FindString(s)
			if newstr != "" {
				sentence, err := nmea.Parse(s)
				if err == nil {
					if sentence.DataType() == nmea.TypeRMC {
						m := sentence.(nmea.RMC)
						chResult <- fmt.Sprintf("Latitude:%s;Longitude:%s;%s", nmea.FormatGPS(m.Latitude), nmea.FormatGPS(m.Longitude), ttyport)
						break
						}
				}			
			}	
		}
	}()

	select {
	case data := <- chResult:
		return data
	case <-time.After(timeout):
		return ""
	}
		
}

func main(){
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
		strRes := readPort(ttyport, 3 * time.Second)
		if strRes != "" {
			fmt.Println(strRes)
			break
		}
	}
}