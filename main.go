package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/xml"
)

type SitemapIndex struct {
  Locations []string `xml:"sitemap>loc"`
}

type News struct {
  Titles []string `xml:"url>news>title"`
  Keywords []string `xml:"url>news>keywords"`
  Locations []string `xml:"url>loc"`
  }

type NewsMap struct {
  Keyword string
  Location string
}

func main() {
    var s SitemapIndex
    var n News
    newsMap := make(map[string]NewsMap)
    resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &s)

    for _, Location := range s.Locations {
        resp, _ := http.Get(Location)
        bytes, _ := ioutil.ReadAll(resp.Body)
        xml.Unmarshal(bytes, &n)
        for i, _ :=range n.Titles {
            newsMap[n.Titles[i]] = NewsMap{n.Keywords[i], n.Locations[i]}
        }
    }

    for i, data := range newsMap {
      fmt.Println("\n", i)
      fmt.Println("\n", data.Keyword)
      fmt.Println("\n", data.Location)
    }
}
