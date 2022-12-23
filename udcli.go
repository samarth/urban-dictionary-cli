package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/urfave/cli"
)

const udapi = "http://api.urbandictionary.com/v0/define?term=%s"

type udrSorterThumbsUp []udresponselist

func (a udrSorterThumbsUp) Len() int      { return len(a) }
func (a udrSorterThumbsUp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a udrSorterThumbsUp) Less(i, j int) bool {
	return a[i].Thumbs_up-a[i].Thumbs_down > a[j].Thumbs_up-a[j].Thumbs_down
}

type udresponse struct {
	Tags        []string
	Result_type string `json:"result_type"`
	List        []udresponselist
	Sounds      []string
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
	var high bool
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
		cli.BoolFlag{
			Name:        "highest, i",
			Usage:       "If set sort by highsted voted answer \"highest voted\"",
			Destination: &high,
		},
	}
	app.Action = func(c *cli.Context) error {
		queryUD(term, high)
		return nil
	}
	app.Run(os.Args)
}

func queryUD(term string, high bool) {
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
	if high {
		sort.Sort(udrSorterThumbsUp(udr.List))
	}

	var topdef = udr.List[0]

	fmt.Println()
	fmt.Println(topdef.Definition + " ---- definition by " + topdef.Author)
	fmt.Println("Read more at -> " + topdef.Permalink)
	fmt.Println()
}
