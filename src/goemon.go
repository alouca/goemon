package main

import (
	"bufio"
	"config"
	//"database/sql"
	"encoding/xml"
	"flag"
	//_ "github.com/mattn/go-sqlite3"
	"fmt"
	"logger"
	"os/exec"
	"rs232"
)

// Sample Message
/* 
 <msg>
	 <src>CC128-v1.29</src>
	 <dsb>00003</dsb>
	 <time>23:15:48</time>
	 <tmpr>31.2</tmpr>
	 <sensor>0</sensor>
	 <id>02884</id>
	 <type>1</type>
	 <ch1>
	 	<watts>00194</watts>
	 </ch1>
	 <ch2>
	 	<watts>01024</watts>
	 </ch2>
	 <ch3>
	 	<watts>00240</watts>
	 </ch3>
	 </msg>
*/

type CCMsg struct {
	XMLName         xml.Name `xml:"msg"`
	Source          string   `xml:"src"`
	Dsb             string
	Time            string
	Temperature     float32 `xml:"tmpr"`
	Sensor          int
	Id              int
	Type            int
	Channel1Reading int `xml:"ch1>watts"`
	Channel2Reading int `xml:"ch2>watts"`
	Channel3Reading int `xml:"ch3>watts"`
}

var (
	configConfigFile string
	l                *logger.Logger
	c                *config.Config
	//db               *sql.DB
)

func init() {
	flag.StringVar(&configConfigFile, "config", "goemon.json", "Specify the GoEmon Server Configuration File")
	flag.Parse()

	c = config.LoadConfigFromFile(configConfigFile)

	verbose := c.GetBool("logger.verbose")
	debug := c.GetBool("logger.debug")
	logOut := c.GetString("logger.output")

	if logOut == "stdout" {
		l = logger.CreateLogger(verbose, debug)
	} else {
		l = logger.CreateLoggerWithFile(verbose, debug, logOut)
	}

}

func main() {
	var err error

	go startServer(c.GetString("webserver.address"))

	// Open RS232 Port
	port, err := rs232.OpenPort(c.GetString("serial.path"), c.GetInt("serial.baud"), rs232.S_8N1)
	if err != nil {
		l.Fatal("Error opening port: %s", err)
	}
	defer port.Close()

	r := bufio.NewReader(&port)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			l.Fatal("Error reading:  %s", err)
		}
		var msg CCMsg

		e := xml.Unmarshal(line, &msg)

		if e != nil {
			l.Fatal("Error parsing data: %s\n", e.Error())
		}

		l.Debug("<: %s", line)
		total := msg.Channel1Reading + msg.Channel2Reading + msg.Channel3Reading
		l.Debug("Parsed: Temperature: %f - CH1 %d - CH2 %d - CH3 %d (Total: %d)\n", msg.Temperature, msg.Channel1Reading, msg.Channel2Reading, msg.Channel3Reading, total)

		cmd := exec.Command("/usr/local/bin/rrdupdate", "goemon.rrd", fmt.Sprintf("N:%f:%d:%d:%d", msg.Temperature, msg.Channel1Reading, msg.Channel2Reading, msg.Channel3Reading))
		cmd.Run()
		//stmt.Exec(msg.Temperature, msg.Channel1Reading, msg.Channel2Reading, msg.Channel3Reading)

		// Update the RRD File

	}
}
