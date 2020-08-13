package magnet_link_crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
)

type AnimateMagnetInfo map[string]map[float64][]MagnetLinkInfo

type MagnetLinkInfo struct {
	Title string
	MagnetLink string
	Episodes []float64
	BtNumber int
}

// Public
func GetAnimateMagnetInfo(pageUrl string, cfg *AnimateRequestInfo) AnimateMagnetInfo {
	animateMagnetInfo := make(AnimateMagnetInfo)

	for animateKey, animateStatus := range cfg.AnimateStatus {
		// initialization
		animateMagnetInfo[animateKey] = make(map[float64][]MagnetLinkInfo)
		teamIds := animateStatus.PreferTeamIds
		if len(teamIds) == 0 {
			teamIds = append(teamIds, "")
		}

		// Crawl and collect magnet link info
		for _, teamId := range teamIds {
			resp, _ := getPage(pageUrl + "?keyword=" + url.PathEscape(animateKey) + "&team_id=" + teamId)
			magnetLinkInfos := extractDmhyMagnetLinkInfo(resp, *animateStatus)
			episodeMagnetMap := genEpisodeMagnetMap(magnetLinkInfos, *animateStatus)

			for episode, magnetLinkInfos := range episodeMagnetMap {
				animateMagnetInfo[animateKey][episode] = append(
					animateMagnetInfo[animateKey][episode], magnetLinkInfos...)
			}
		}

		// Sort result with BtNumber
		for episode, _ := range animateMagnetInfo[animateKey] {
			sort.Slice(animateMagnetInfo[animateKey][episode], func(i int, j int) bool {
				return animateMagnetInfo[animateKey][episode][i].BtNumber >
					animateMagnetInfo[animateKey][episode][j].BtNumber
			})
		}
	}

	return animateMagnetInfo
}

func DumpAnimateMagnetInfo(animateMagnetInfo AnimateMagnetInfo) {
	for animateKey, episodeMagnetLinkInfos := range animateMagnetInfo {
		fmt.Println("=========================")
		fmt.Println("Name: " + animateKey)
		for episodeId, episodeMagnetLinkInfo := range episodeMagnetLinkInfos {
			fmt.Println("Episode: ", episodeId)
			for _, magnetLinkInfo := range episodeMagnetLinkInfo {
				fmt.Println("MagnetLink/btNums: ",
					magnetLinkInfo.MagnetLink[:50]+ ".../" + strconv.Itoa(magnetLinkInfo.BtNumber))
			}
		}
		fmt.Println("=========================")
	}
}

// Private
func getPage(pageUrl string) (resp *http.Response, err error) {
	// Setup cookie
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name: "cf_clearance",
		Value: "f1d8bd94c77e76a423871e94de9ec2ce64cdb366-1597199043-0-1zc9740c68z47380d19z8273ec9-150",
		Path: "/",
		Domain: ".dmhy.org",
	})
	u, _ := url.Parse(pageUrl)
	jar.SetCookies(u, cookies)

	// Set cookie
	client := &http.Client{ Jar: jar }
	req, _ := http.NewRequest("GET", u.String(), nil)

	// Setup header
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")

	// Do request
	resp, err = client.Do(req)

	return
}

func extractDmhyMagnetLinkInfo(resp *http.Response, animateStatus AnimateStatus) []MagnetLinkInfo {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Failed to parse dmhy response: ", err)
	}

	var magnetLinkInfos []MagnetLinkInfo
	doc.Find(".tablesorter tbody tr").Each(func(_ int, s *goquery.Selection) {
		title := s.Find(".title a[target=_blank]").Text()
		magnetLink, _ := s.Find(".download-arrow").Attr("href")
		btNumber, _ := strconv.Atoi(s.Find(".btl_1").Text())
		episodes := parseEpisode(title, animateStatus.PreferParser)

		if len(episodes) > 0 {
			magnetLinkInfos = append(magnetLinkInfos, MagnetLinkInfo{
				Title: title,
				MagnetLink: magnetLink,
				Episodes: episodes,
				BtNumber: btNumber,
			})
		}
	})

	return magnetLinkInfos
}

func parseEpisode(title string, preferParser string) []float64 {
	var episodeList []float64

	floatPattern := regexp.MustCompile(`\d+\.?\d*`)
	var singlePattern *regexp.Regexp

	if len(preferParser) > 0 {
		singlePattern = regexp.MustCompile(preferParser)
	} else {
		singlePattern = regexp.MustCompile(
			`(\[\d+\.?\d*\])|(【\d+\.?\d*】)|(「\d+\.?\d*」)|(第\d+\.?\d*集)|(第\d+\.?\d*話)|(第\d+\.?\d*话)`)
	}

	singleM := singlePattern.FindAllString(title, -1)
	if len(singleM) > 0 {
		for _, episodeStr := range singleM {
			episodeId, _ := strconv.ParseFloat(floatPattern.FindString(episodeStr), 10)
			if 0 < episodeId && episodeId < 1000 {
				episodeList = append(episodeList, episodeId)
			}
		}

		return episodeList
	}

	multiPattern := regexp.MustCompile(
		`(\[\d+\.?\d*\-\d+\.?\d*\])|(【\d+\.?\d*\-\d+\.?\d*】)|(「\d+\.?\d*\-\d+\.?\d*」)|` +
			`(第\d+\.?\d*\-\d+\.?\d*集)|(第\d+\.?\d*\-\d+\.?\d*話)|(第\d+\.?\d*\-\d+\.?\d*话)`)

	multiM := multiPattern.FindAllString(title, -1)
	if len(multiM) > 0 {
		for _, episodeStr := range multiM {
			episodePureStrM := floatPattern.FindAllString(episodeStr, -1)

			if checkEpisodeLimit(episodePureStrM) {
				start, _ := strconv.ParseFloat(episodePureStrM[0], 10)
				end, _ := strconv.ParseFloat(episodePureStrM[1], 10)

				var episodes []float64

				episodes = append(episodes, start)
				for i := int(math.Floor(start + 1)); i <= int(math.Ceil(end - 1)); i ++ {
					if i < 10 {
						episodes = append(episodes, float64(i))
					} else {
						episodes = append(episodes, float64(i))
					}
				}
				episodes = append(episodes, end)

				episodeList = append(episodeList, episodes...)
			}
		}

		return episodeList
	}

	return []float64{}
}

func checkEpisodeLimit(episodes []string) bool {
	if len(episodes) < 1 {
		return false
	}

	for _, episodeStr := range episodes {
		episodeInt, _ := strconv.Atoi(episodeStr)
		if episodeInt < 0 || episodeInt > 1024 {
			return false
		}
	}

	return true
}

func genEpisodeMagnetMap(magnetLinkInfos []MagnetLinkInfo, animateStatus AnimateStatus) map[float64][]MagnetLinkInfo {
	episodeMagnetMap := make(map[float64][]MagnetLinkInfo)

	for _, magnetLinkInfo := range magnetLinkInfos {
		for _, episode := range magnetLinkInfo.Episodes {
			episodeMagnetMap[episode] = append(
				episodeMagnetMap[episode], magnetLinkInfo)
		}
	}

	for _, completedEpisode := range animateStatus.CompletedEpisodes {
		delete(episodeMagnetMap, completedEpisode)
	}

	return episodeMagnetMap
}