package main

import (
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "regexp"

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


func ConvToFloat32s(num_strs []string) []float32 {
    res := []float32{}

    for i := range num_strs {
        f, _ := strconv.ParseFloat(num_strs[i], 32)
        res = append(res, float32(f))
    }

    return res
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
        re := regexp.MustCompile(`(\[.*?\])|(【.*?】)`)
        matchs := re.FindAllString(rows[i].title, -1)
        if len(matchs) < 3 {
            continue
        }

        re2 := regexp.MustCompile(`\d+(\.\d+)*`)
        num_strs := re2.FindAllString(matchs[2], -1)
        row_epis := ConvToFloat32s(num_strs)
        if len(row_epis) == 0 {
            continue
        }

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
            episodes: row_epis,
            magnet: rows[i].magnet})
    }

    return tmp_cands
}


func GetCands(keyword string, team_ids []int, episodes []float32) []Candidate {
    cands := []Candidate{}
    path := fmt.Sprintf("https://share.dmhy.org/topics/list?keyword=%s", keyword)

    paths := []string{}
    for i := range team_ids {
        paths = append(paths, fmt.Sprintf(
            path + "&team_id=%d", team_ids[i]))
    }

    if len(paths) == 0 {
        paths = append(paths, path)
    }

    for i := range paths {
        resp, _ := GetContent(paths[i])
        rows := RowReader(resp)

        tmp_cands := ExtractCands(rows, episodes, keyword)

        cands = append(cands, tmp_cands...)
    }

    return cands
}
