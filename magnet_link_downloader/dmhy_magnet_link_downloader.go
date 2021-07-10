package magnet_link_downloader

import (
	"os"
	"fmt"
	"path/filepath"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
)

// Public
func DownloadMagnet(animateMagnetInfo magnet_link_crawler.AnimateMagnetInfo, downloadInfo DownloadInfo,
	requestInfo magnet_link_crawler.AnimateRequestInfo) map[string]*magnet_link_crawler.AnimateStatus {
	fmt.Println(downloadInfo)

	// Initialize a basic magnet downloader
	downloader := new(BasicDownloader).Init(downloadInfo)
	newDownloadeds := make(map[string]*magnet_link_crawler.AnimateStatus)

	for animateKey, episodeMagnetMaps := range animateMagnetInfo {
		realStoragePath := filepath.Join(downloadInfo.StoragePath, animateKey)
		if _, err := os.Stat(realStoragePath); os.IsNotExist(err) {
			os.MkdirAll(realStoragePath, 0755)
		}
		downloader.StoragePath = realStoragePath

		fmt.Println("================================")
		fmt.Println("AnimateKeyword:", animateKey)

		newDownloadeds[animateKey] = new(magnet_link_crawler.AnimateStatus)
		for episode, magnetLinkInfos := range episodeMagnetMaps {
			if requestInfo.AnimateStatusMap[animateKey].IsComplete(episode) {
				continue
			}

			fmt.Println("Episode:", episode)

			// Use downloader to download animate with magnet link
			downloadedMagnetLinkInfo := downloader.Download(magnetLinkInfos)

			// Update downloaded info
			requestInfo.AnimateStatusMap[animateKey].CommitEpisode(downloadedMagnetLinkInfo.Episodes...)
			newDownloadeds[animateKey].CommitEpisode(downloadedMagnetLinkInfo.Episodes...)
		}
	}

	return newDownloadeds
}
