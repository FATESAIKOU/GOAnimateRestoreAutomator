package main

import (
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "regexp"
    "net/http/cookiejar"

    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
    "net/url"
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


func GetContent(ur string) (resp *http.Response, err error) {
    // Setup cookie
    jar, _ := cookiejar.New(nil)
    var cookies []*http.Cookie
    cookies = append(cookies, &http.Cookie{
        Name: "cf_clearance",
        Value: "af596069241d985207dc25c5f55a8742a494959f-1596954475-0-1ze41e1b5bzaccd7c85zeb3f3455-150",
        Path: "/",
        Domain: ".dmhy.org",
    })
    u, _ := url.Parse(ur)
    jar.SetCookies(u, cookies)

    // Set cookie
    client := &http.Client{ Jar: jar }
    req, _ := http.NewRequest("GET", ur, nil)

    // Setup header
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")

    // Do request
    resp, err = client.Do(req)

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

// TODO Add Range Support!!
func ConvToFloat32Range(num_strs []string) []float32 {
    res := []float32{}

    if len(num_strs) == 1 || (len(num_strs) == 2 && num_strs[0][0] != 'v') {
        num, _ := strconv.ParseFloat(num_strs[0], 32)
        res = append(res, float32(num))
    }

    return res
}


func ExtractCands(rows []Row, episodes []float32, keyword string, new_episodes map[float32]bool) []Candidate {
    tmp_cands := []Candidate{}

    // Load Episodes
    has_episodes := make(map[float32]bool)
    LoadEpisodesToMap(has_episodes, episodes)

    // Extract Cands
    for i := range rows {
        // regexp
        re := regexp.MustCompile(`(\[.*?\])|(【.*?】)|(「.*?」)`)
        matchs := re.FindAllString(rows[i].title, -1)
        if len(matchs) < 5 {
            continue
        }

        re2 := regexp.MustCompile(`(\d+)|(v\d+)`)
        num_strs := re2.FindAllString(matchs[2], -1)
        row_epis := ConvToFloat32Range(num_strs)
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

    new_episodes := make(map[float32]bool)
    for i := range paths {
        resp, _ := GetContent(paths[i])
        rows := RowReader(resp)

        tmp_cands := ExtractCands(rows, episodes, keyword, new_episodes)

        cands = append(cands, tmp_cands...)
    }

    return cands
}
