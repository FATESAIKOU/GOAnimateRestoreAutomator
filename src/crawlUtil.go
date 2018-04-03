package main

import (
    "fmt"
    "net/http"
    "strings"

    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
)

type Row struct {
    title string
    magnet string
}

type HttpError struct {
    original string
}


func GetContent(url string) (resp *http.Response, err error) {
    resp, err = http.Get(url)

    return
}


func RowReader(resp *http.Response) []Row {
    page := html.NewTokenizer(resp.Body)
    var title string
    rows := []Row{}

    inside_title := false
    for {
        tt := page.Next()
        t  := page.Token()

        if tt == html.ErrorToken {
            break
        }

        if tt == html.TextToken {
            if inside_title {
                title = fmt.Sprintf("%s%s", title, t.Data)
            }
        }

        if t.DataAtom == atom.A {
            if tt == html.StartTagToken && getAttr(t, "title") == "磁力下載" {
                title = strings.TrimSpace(title)
                rows = append(rows, Row{title: title, magnet: getAttr(t, "href")})
                title = ""
            }

            if inside_title && tt == html.EndTagToken {
                inside_title = false
            }

            if tt == html.StartTagToken && getAttr(t, "target") == "_blank"  &&
                strings.Contains(getAttr(t, "href"), "/topics/view") {
                inside_title = true
            }
        }
    }

    return rows
}


func getAttr(tag html.Token, attr string) (value string) {
    for i := range tag.Attr {
        if tag.Attr[i].Key == attr {
            value = tag.Attr[i].Val
            return
        }
    }

    return
}


func getRows(keyword string, team_id int) []Row {
    path := fmt.Sprintf(
        "https://share.dmhy.org/topics/list?team_id=%d&keyword=%s",
        team_id, keyword)
    resp, _ := GetContent(path)
    rows := RowReader(resp)

    return rows
}
