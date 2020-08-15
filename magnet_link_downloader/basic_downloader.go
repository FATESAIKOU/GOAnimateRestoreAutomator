package magnet_link_downloader

import (
	"context"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BasicDownloader struct {
	StoragePath string
	ErrorFilePath string
}

func (basicDownloaderSelf *BasicDownloader) Init(downloadInfo DownloadInfo) *BasicDownloader {
	basicDownloaderSelf.StoragePath = downloadInfo.StoragePath
	basicDownloaderSelf.ErrorFilePath = downloadInfo.ErrorFilePath

	return basicDownloaderSelf
}

func (basicDownloaderSelf BasicDownloader) Download(
	magnetInfos []magnet_link_crawler.MagnetLinkInfo) magnet_link_crawler.MagnetLinkInfo {
	for _, magnetInfo := range magnetInfos {
		// TODO add timeout to downloader config
		err := DownloadMagnetLink(magnetInfo, basicDownloaderSelf.StoragePath, 2400)

		if err == nil {
			return magnetInfo
		}
	}

	return magnet_link_crawler.MagnetLinkInfo{
		Title: "",
		Episodes: []float64{},
		BtNumber: 0,
	}
}

func DownloadMagnetLink(magnetLinkInfo magnet_link_crawler.MagnetLinkInfo, storagePath string, timeout uint32) error {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
	tmpDir, err :=ioutil.TempDir("/tmp", "")
	if err != nil {
		log.Println("Fail to make temp dir")
		return nil
	}

	defer func() {
		os.RemoveAll(tmpDir)
		cancel()
	}()

	cmd := exec.CommandContext(ctxt, "webtorrent", magnetLinkInfo.MagnetLink)
	cmd.Dir = tmpDir

	log.Printf("Try to download: %s (%f)", magnetLinkInfo.Title, magnetLinkInfo.Size)

	if err := cmd.Run(); err != nil {
		if ctxt.Err() == context.DeadlineExceeded {
			log.Println("Download Timeout:", err, ":", magnetLinkInfo.Title)
		} else {
			log.Println("Download Failed:", err, ":", magnetLinkInfo.Title)
		}
		return err
	}

	files, _ := ioutil.ReadDir(tmpDir)
	for _, file := range files {
		err := os.Rename(filepath.Join(tmpDir, file.Name()), filepath.Join(storagePath, file.Name()))
		if err != nil {
			log.Println("Failed to move file:", err)
			return nil
		}
	}

	return nil
}