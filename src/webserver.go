package main

import (
	"encoding/json"
	"html/template"
	"io"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

var (
	t *template.Template
)

type SensorChannels struct {
	Label string
	Data  []float64
}

type JSONLiveData struct {
	Channels []SensorChannels
}

func startServer(port string) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	http.HandleFunc("/data", DataHandler)
	http.HandleFunc("/", IndexHandler)
	var e error

	t, e = template.ParseGlob("web/templates/*.tmpl")

	if e != nil {
		l.Fatal("Unable to parse templates: %s\n", e.Error())
	}

	e = http.ListenAndServe(port, nil)

	if e != nil {
		l.Fatal("Unable to start embeeded webserver: %s\n", e.Error())
	}
}

func DataHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		l.Fatal("Unable to parse HTTP Form: %s\n", err.Error())
	}

	l.Debug("Data Range requested: %s - %s\n", req.FormValue("start"), req.FormValue("end"))

	var startRange, endRange string

	if startRange = req.FormValue("start"); startRange == "" {
		startRange = "-10m"
	}

	if endRange = req.FormValue("end"); endRange == "" {
		endRange = "-6"
	}

	buf, err := exec.Command("/usr/local/bin/rrdtool", "fetch", "goemon.rrd", "AVERAGE", "-r6", "-s"+startRange, "-e"+endRange).Output()

	if err != nil {
		l.Fatal("Unable to get command output: %s\n", err.Error())
	}

	/*
		// Return the last 100 readings in json data
		r, e := db.Query("select * from readings order by r_date desc limit 100")

		if e != nil {
			l.Fatal("Error quering database: %s\n", e.Error())
		}
	*/

	// Split rrdtool data
	lines := strings.Split(string(buf), "\n")

	// Prepare return structure
	var data JSONLiveData

	data.Channels = make([]SensorChannels, 3)
	data.Channels[0].Label = "Phase 1"
	data.Channels[1].Label = "Phase 2"
	data.Channels[2].Label = "Phase 3"

	data.Channels[0].Data = make([]float64, len(lines)-4)
	data.Channels[1].Data = make([]float64, len(lines)-4)
	data.Channels[2].Data = make([]float64, len(lines)-4)

	i := 0

	for _, line := range lines[2 : len(lines)-2] {
		parts := strings.Split(line, " ")
		l.Debug("Parts: %s %s %s %s %s\n", parts[0], parts[1], parts[2], parts[3], parts[4])

		var c1, c2, c3 float64

		if c1, _ = strconv.ParseFloat(parts[2], 64); math.IsNaN(c1) {
			c1 = 0
		}
		if c2, _ = strconv.ParseFloat(parts[3], 64); math.IsNaN(c2) {
			c2 = 0
		}
		if c3, _ = strconv.ParseFloat(parts[4], 64); math.IsNaN(c3) {
			c3 = 0
		}

		data.Channels[0].Data[i] = c1
		data.Channels[1].Data[i] = c2
		data.Channels[2].Data[i] = c3
		i++
	}

	jsonData, e := json.Marshal(data)

	if e != nil {
		l.Fatal("Error marshalling JSON: %s\n", e.Error())
	}

	l.Debug("Marshalled JSON: %s\n", string(jsonData))

	io.WriteString(w, string(jsonData))
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	t, e := template.ParseGlob("web/templates/*.tmpl")

	if e != nil {
		l.Fatal("Unable to parse templates: %s\n", e.Error())
	}

	e = t.Execute(w, nil)

	if e != nil {
		l.Fatal("Error executing template: %s\n", e.Error())
	}
}
