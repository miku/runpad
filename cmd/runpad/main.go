package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	etherpadlite "github.com/FabianWe/etherpadlite-golang"
	_ "github.com/joho/godotenv/autoload"
	"github.com/miku/runpad/padutil"
)

var (
	baseURL      = flag.String("u", os.Getenv("RUNPAD_BASE_URL"), "etherpad base URL")
	apiKey       = flag.String("a", os.Getenv("RUNPAD_APIKEY"), "etherpad api key")
	identifier   = flag.String("p", "runpad", "pad name to watch")
	listPads     = flag.Bool("l", false, "list pads")
	showPad      = flag.Bool("c", false, "show pad contents and info")
	listSnippets = flag.Bool("s", false, "list snippets")
	runSnippet   = flag.Int("r", -1, "run snippet with given id, zero-indexed")
)

func main() {
	flag.Parse()
	var (
		ctx = context.Background()
		e   = etherpadlite.NewEtherpadLite(*apiKey)
	)
	e.BaseURL = *baseURL
	ether := padutil.Etherpad{*e}
	switch {
	case *listPads:
		resp, err := ether.ListAllPads(ctx)
		if err != nil {
			log.Fatal(err)
		}
		padIDs, ok := resp.Data["padIDs"]
		if !ok {
			log.Fatal("unexpected api response")
		}
		for i, v := range padIDs.([]interface{}) {
			log.Printf("%v %v", i, v)
		}
		os.Exit(0)
	case *showPad:
		ok, err := ether.Exists(ctx, *identifier)
		if err != nil {
			log.Fatal("cannot reach pad: %v", err)
		}
		if !ok {
			log.Fatal("pad not found")
		}
		revs, err := ether.GetPadSavedRevisions(ctx, *identifier)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("url: %s/p/%s\n", strings.Replace(*baseURL, "/api", "", 1), *identifier)
		fmt.Printf("revs: %v\n\n", revs.Data.SavedRevisions)
		fmt.Println("----\n")
		// At this point we know the pad exists and we can access it.
		padText, err := ether.GetPadText(ctx, *identifier)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(padText.Data.Text)
	case *listSnippets:
		ok, err := ether.Exists(ctx, *identifier)
		if err != nil {
			log.Fatal("cannot reach pad: %v", err)
		}
		if !ok {
			log.Fatal("pad not found")
		}
		// At this point we know the pad exists and we can access it.
		padText, err := ether.GetPadText(ctx, *identifier)
		if err != nil {
			log.Fatal(err)
		}
		text := padutil.Text{
			Content: padText.Data.Text,
		}
		for i, v := range text.Snippets() {
			log.Printf("snippet #%d contains %d lines", i, v.NumLines())
			b, err := json.Marshal(v)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(b))
		}
	case *runSnippet >= 0:
		ok, err := ether.Exists(ctx, *identifier)
		if err != nil {
			log.Fatal("cannot reach pad: %v", err)
		}
		if !ok {
			log.Fatal("pad not found")
		}
		// At this point we know the pad exists and we can access it.
		padText, err := ether.GetPadText(ctx, *identifier)
		if err != nil {
			log.Fatal(err)
		}
		text := padutil.Text{
			Content: padText.Data.Text,
		}
		snippets := text.Snippets()
		if *runSnippet >= len(snippets) {
			log.Fatalf("invalid snippet id: %v (%v)", *runSnippet, len(snippets))
		}
		var (
			snippet = snippets[*runSnippet]
			runner  padutil.Runner
		)
		switch snippet.Tag {
		case "python":
			runner = &padutil.SimpleFileRunner{
				Prefix: []string{"python"},
			}
		case "go":
			runner = &padutil.SimpleFileRunner{
				Prefix: []string{"go", "run"},
			}
		}
		if err := runner.Run(os.Stdout, snippet); err != nil {
			log.Fatal(err)
		}
	}
}
