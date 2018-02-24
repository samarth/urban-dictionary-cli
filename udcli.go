package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const udapi = "http://api.urbandictionary.com/v0/define?term=%s"

type udresponse struct {
	Tags        []string
	Result_type string `json:"result_type"`
	List        []udresponselist
	Sounds       []string
}

type udresponselist struct {
	Definition   string
	Permalink    string
	Thumbs_up    int
	Author       string
	Word         string
	Defid        int
	Current_vote string
	Example      string
	Thumbs_down  int
}

func main() {
	var term string
	app := cli.NewApp()
	app.Name = "udcli"
	app.Usage = "udcli <term> "
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "term, t",
			Value:       "Ullu Ka Patha",
			Usage:       "Term you want to query on Urban Dictionary",
			Destination: &term,
		},
	}
	app.Action = func(c *cli.Context) error {
		queryUD(term)
		return nil
	}
	app.Run(os.Args)
}

func queryUD(term string) {
	queryuri := fmt.Sprintf(udapi, term)
	resp, err := http.Get(queryuri)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error contacting Urband Dictionary Server :(, Exiting ...")
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var udr = udresponse{}

	unmarshalererr := json.Unmarshal(bodyBytes, &udr)

	if unmarshalererr != nil {
		fmt.Println("Unable to unmarshal urban dictionary response, is it still the same?")
	}
	if len(udr.List) == 0 {
		fmt.Println()
		fmt.Println("No definitions found.")
		fmt.Println()
		os.Exit(0)
	}
	var topdef = udr.List[0]

	fmt.Println()
	fmt.Println(topdef.Definition + " ---- definition by " + topdef.Author)
	fmt.Println("Read more at -> " + topdef.Permalink)
	fmt.Println()
}
