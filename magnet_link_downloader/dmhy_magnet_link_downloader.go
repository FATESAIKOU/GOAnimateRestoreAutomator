package magnet_link_downloader

import (
	"fmt"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
)

// Public
func DownloadMagnet(animateMagnetInfo magnet_link_crawler.AnimateMagnetInfo, downloadInfo DownloadInfo,
	requestInfo magnet_link_crawler.AnimateRequestInfo) magnet_link_crawler.AnimateRequestInfo{
	fmt.Println(downloadInfo)

	// TODO initialize a basic magnet downloader
	// downloader := BasicDownloader.Init(downloaderInfo)

	for animateKey, episodeMagnetMaps := range animateMagnetInfo {
		fmt.Println("AnimateKeyword: ", animateKey)

		for episode, magnetLinkInfos := range episodeMagnetMaps {
			if requestInfo.AnimateStatus[animateKey].IsComplete(episode) {
				continue
			}

			fmt.Println("Episode: ", episode)
			fmt.Println(len(magnetLinkInfos))

			// TODO use downloader to download animate with magnet link
			// downloadedMagnetLinkInfo := downloader.Download(magnetLinkInfos, downloadInfo)

			// TODO update downloaded info
			// requestInfo.AnimateStatus[animateKey].CommitEpisode(downloadedMagnetLinkInfo.Episodes...)
		}
	}

	return requestInfo
}