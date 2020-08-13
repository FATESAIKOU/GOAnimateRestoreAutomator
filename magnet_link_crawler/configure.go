package magnet_link_crawler

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// The log and settings for magnet link crawling
type AnimateRequestInfo struct {
	AnimateStatus map[string]*AnimateStatus
}

type AnimateStatus struct {
	CompletedEpisodes []float64
	PreferTeamIds []string
	PreferParser string
}

func (animateRequestInfoSelf *AnimateRequestInfo) LoadJson(jsonFilePath string) *AnimateRequestInfo {
	// Read json raw
	rawJson, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatal("Fail to read animate request json file: ", err)
	}

	// Initialize map
	animateRequestInfoSelf.AnimateStatus = make(map[string]*AnimateStatus)

	// Parse json
	var halfParseJson map[string]json.RawMessage
	err = json.Unmarshal(rawJson, &halfParseJson)

	var animateStatusTmp *AnimateStatus
	for animateKeyword, statusJson := range halfParseJson {
		animateStatusTmp = new(AnimateStatus)
		err = json.Unmarshal(statusJson, animateStatusTmp)

		animateRequestInfoSelf.AnimateStatus[animateKeyword] = animateStatusTmp
	}

	if err != nil {
		log.Fatal("Fail to parse animate request json file: ", err)
	}

	return animateRequestInfoSelf
}

func (animateRequestInfoSelf *AnimateRequestInfo) SaveJson(jsonFilePath string) {
	rawCfgJson, _ := json.MarshalIndent(animateRequestInfoSelf.AnimateStatus, "", "\t")
	ioutil.WriteFile(jsonFilePath, rawCfgJson, 0644)
}

func (animateStatusSelf *AnimateStatus) CommitEpisode(episodes ... float64) {
	animateStatusSelf.CompletedEpisodes = append(animateStatusSelf.CompletedEpisodes, episodes ...)
}

func (animateStatusSelf AnimateStatus) IsComplete(episode float64) bool {
	for _, completedEpisode := range animateStatusSelf.CompletedEpisodes {
		if completedEpisode == episode {
			return true
		}
	}

	return false
}