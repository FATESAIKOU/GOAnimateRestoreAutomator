package magnet_link_downloader

import (
	"context"
	"fmt"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		err := DownloadMagnetLink(magnetInfo, basicDownloaderSelf.StoragePath, 30 * uint32(math.Ceil(magnetInfo.Size)))

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
	log.Printf("Try to download: %s (%fMB)", magnetLinkInfo.Title, magnetLinkInfo.Size)
	cmd.Dir = tmpDir

	cmd.Start()

	endOfCmd := make(chan bool, 1)
	go func() {
		handleProgress(tmpDir, magnetLinkInfo.Size, endOfCmd)
	}()

	err = cmd.Wait()
	endOfCmd <- true

	if err != nil {
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

// utils
func handleProgress(tmpDir string, targetSize float64, endOfCmd chan bool) {
	preSize := 0.0
	tWidth, _ := terminal.Width()

	for {
		select {
		case <-endOfCmd:
			fmt.Printf(strings.Repeat(" ", int(tWidth)) + "\n\033[F")
			break
		default:
			nowSize, _ := dirSize(tmpDir)
			fmt.Printf("Progress: %f%% - %fMB/s\n\033[F",
				math.Min(nowSize * 100 / targetSize, 100), nowSize - preSize)
			preSize = nowSize
			time.Sleep(1 * time.Second)
			fmt.Printf(strings.Repeat(" ", int(tWidth)) + "\n\033[F")
		}
	}
}

func dirSize(path string) (float64, error) {
	var size float64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += float64(info.Size())
		}
		return err
	})
	return size / 1048576.0, err
}
