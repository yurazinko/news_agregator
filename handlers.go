package main

import (
    "net/http"
    "html/template"
    "io/ioutil"
    "encoding/xml"
    "sync"
)

var wg sync.WaitGroup

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

type NewsAggPage struct {
    Title string
    News map[string]NewsMap
}

func newsRoutine(c chan News, Location string) {
    defer wg.Done()
    var n News
    resp, _ := http.Get(Location)
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &n)
    resp.Body.Close()
    c <- n
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
    var s SitemapIndex
    newsMap := make(map[string]NewsMap)
    resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &s)
    resp.Body.Close()
    queue := make(chan News, 30)

    for _, Location := range s.Locations {
        wg.Add(1)
        go newsRoutine(queue, Location)
    }

    wg.Wait()
    close(queue)

    for elem := range queue {
        for i, _ := range elem.Keywords {
            newsMap[elem.Titles[i]] = NewsMap{elem.Keywords[i], elem.Locations[i]}
        }
    }

    p := NewsAggPage{Title: "Amazing News Aggregator", News: newsMap }
    t, _ := template.ParseFiles("feed.html")
    t.Execute(w, p)
}

func main() {
    http.HandleFunc("/", newsAggHandler)
    http.ListenAndServe(":8000", nil)
}
