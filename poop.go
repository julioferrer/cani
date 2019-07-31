package main

import (
	"fmt"
	"github.com/tarm/serial"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const longPoop = int64(10)

var lastPoop int64

func main() {
	go serialRead()
	http.HandleFunc("/", gotPoop)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serialRead() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	lastPoop = time.Now().Unix()
	poopStart := lastPoop

	for true {
		buf := make([]byte, 4)
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", string(buf[:n]))
		// Wait because this might be a long pooper who sits in  the darkness for _longPoop_ seconds
		if canIPoop() {
			writeTimestamp("poop.log", fmt.Sprintf("%d %d\n", poopStart, lastPoop))
			poopStart = time.Now().Unix()
		}
		lastPoop = time.Now().Unix()
		// write lastPoop
	}
}

func writeTimestamp(fileName, text string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func gotPoop(w http.ResponseWriter, r *http.Request) {
	if "last" == r.URL.Path[1:] {
		fmt.Fprintf(w, "Last Poop was %d seconds ago", time.Now().Unix()-lastPoop)
		return
	}

	respTpl, err := template.New("resp").Parse(htmlTmpl)
	if err != nil {
		log.Fatal("Template f up ", err)
	}
	imgSrc := "https://i.imgur.com/QhWI4Mg.gif"
	decision := "Nope"

	if canIPoop() {
		imgSrc = "https://i.imgur.com/l3Bs46c.jpg"
		decision = "Yes you can!"
	}
	respTpl.Execute(w, struct {
		ImgSrc   string
		Decision string
	}{imgSrc, decision})
}

func canIPoop() bool {
	return time.Now().Unix()-lastPoop > longPoop
}

const htmlTmpl = `
<head>
	<title>Can I Poop?</title>
	<meta property="og:title" content="Can I Poop?" />
	<meta property="og:description" content="{{.Decision}}" />
	<meta property="og:image" content="{{.ImgSrc}}" />
</head>
<body>
	<img src='{{.ImgSrc}}' />
</body>
`
