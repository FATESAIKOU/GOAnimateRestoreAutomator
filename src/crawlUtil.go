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

type Candidate struct {
    keyword string
    episodes []float32
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
            if tt == html.StartTagToken && GetAttr(t, "title") == "磁力下載" {
                title = strings.TrimSpace(title)
                rows = append(rows, Row{title: title, magnet: GetAttr(t, "href")})
                title = ""
            }

            if inside_title && tt == html.EndTagToken {
                inside_title = false
            }

            if tt == html.StartTagToken && GetAttr(t, "target") == "_blank"  &&
                strings.Contains(GetAttr(t, "href"), "/topics/view") {
                inside_title = true
            }
        }
    }

    return rows
}


func GetAttr(tag html.Token, attr string) (value string) {
    for i := range tag.Attr {
        if tag.Attr[i].Key == attr {
            value = tag.Attr[i].Val
            return
        }
    }

    return
}


func LoadEpisodesToMap(m map[float32]bool, episodes []float32) {
    for i := range episodes {
        m[episodes[i]] = true
    }
}


func ExtractCands(rows []Row, episodes []float32, keyword string) []Candidate {
    tmp_cands := []Candidate{}

    // Load Episodes
    has_episodes := make(map[float32]bool)
    LoadEpisodesToMap(has_episodes, episodes)

    // Extract Cands
    new_episodes := make(map[float32]bool)
    for i := range rows {
        // regexp
        // row_epis := reg(rows[i].title)
        row_epis := []float32{1, 2, 3}

        // is in episodes?
        is_downloaded := false
        for j := range row_epis {
            if has_episodes[row_epis[j]] == true ||
                new_episodes[row_epis[j]] == true {
                is_downloaded = true
                break
            }
        }

        if is_downloaded {
            is_downloaded = false
            continue
        }

        // update new episodes
        LoadEpisodesToMap(new_episodes, row_epis)

        // append
        tmp_cands = append(tmp_cands, Candidate{
            keyword: keyword,
            episodes: []float32{0.0},
            magnet: rows[i].magnet})
    }

    return tmp_cands
}


func GetCands(keyword string, team_ids []int, episodes []float32) []Candidate {
    cands := []Candidate{}

    for i := range team_ids {
        path := fmt.Sprintf(
            "https://share.dmhy.org/topics/list?team_id=%d&keyword=%s",
            team_ids[i], keyword)
        resp, _ := GetContent(path)
        rows := RowReader(resp)

        tmp_cands := ExtractCands(rows, episodes, keyword)

        cands = append(cands, tmp_cands...)
    }

    return cands
}
