package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/term"
)

type alertManAlert struct {
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

type alertManOut struct {
	Alerts            []alertManAlert `json:"alerts"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
	} `json:"commonLabels"`
	ExternalURL string `json:"externalURL"`
	GroupKey    string `json:"groupKey"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

func main() {
	ttyPath := flag.String("tty.path", "/dev/ttyUSB1", "")
	flag.Parse()

	t, err := term.Open(*ttyPath, term.Speed(9600), term.RawMode)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "Listening on 0.0.0.0:9095\n")
	http.ListenAndServe(":9095", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		amo := alertManOut{}
		err = json.Unmarshal(b, &amo)
		if err != nil {
			panic(err)
		}

		for _, alert := range amo.Alerts {
			text := fmt.Sprintf("[%s]: %s on %s: %s\r\n", strings.ToUpper(alert.Status), alert.Labels["alertname"], alert.Labels["instance"], alert.Annotations.Description)
			t.Write([]byte(text))
		}
	}))
}
