package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const destPath = "../content/post/old"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const form = "Mon, 02 Jan 2006 15:04:05 MST"
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(form, v)
	if err != nil {
		return err
	}
	*t = Time{parse}
	return nil
}

type Entry struct {
	Title           string
	Description     string
	BlogEntryNumber int
	PubDate         Time
}
type BlogDay struct {
	Entries []Entry `xml:"Entries>BlogEntry"`
}

type FrontMatter struct {
	Title   string    `json:"title"`
	Section string    `json:"section"`
	Date    time.Time `json:"date"`
	Archive []string  `json:"archive"`
}

func main() {
	log.SetFlags(0)
	log.Printf("Removing %s", destPath)
	os.RemoveAll(destPath)

	log.Printf("Creating %s", destPath)
	os.Mkdir(destPath, os.ModePerm)

	log.Printf("Reading old blog posts")
	var entries []Entry

	err := filepath.Walk(".", func(path string, fi os.FileInfo, _ error) error {
		if fi.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}

		log.Printf("  Loading %s", path)
		xml_data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		bd := BlogDay{}
		err = xml.Unmarshal(xml_data, &bd)
		if err != nil {
			return err
		}

		for _, e := range bd.Entries {
			if e.Title == "" {
				if e.Description == "" {
					continue
				}
				e.Title = fmt.Sprintf("Untitled #%d", e.BlogEntryNumber)
			}
			log.Printf("    Parsed post %d: %s", e.BlogEntryNumber, e.Title)
			entries = append(entries, e)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed %d entries", len(entries))

	loc, _ := time.LoadLocation("US/Pacific")

	log.Printf("Creating new posts")
	for _, e := range entries {
		log.Printf("  Creating post %d: %s", e.BlogEntryNumber, e.Title)
		f, _ := os.Create(fmt.Sprintf("%s/%05d.md", destPath, e.BlogEntryNumber))
		fm := FrontMatter{
			Title:   e.Title,
			Section: "post",
			Date:    e.PubDate.In(loc),
			Archive: []string{e.PubDate.In(loc).Format("2006/01/02")},
		}
		fms, _ := json.MarshalIndent(fm, "", "  ")
		f.Write(fms)
		f.WriteString("\n{{< verbatim >}}\n")
		f.WriteString(e.Description)
		f.WriteString("\n{{< /verbatim >}}\n")
	}
}
