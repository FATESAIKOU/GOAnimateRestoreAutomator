package magnet_link_downloader

// The settings for magnet link downloading
type DownloadInfo struct {
	StoragePath string
	ErrorFilePath string
}

func (downloadInfoSelf *DownloadInfo) Load(storagePath string, errorFilePath string) *DownloadInfo {
	downloadInfoSelf.StoragePath = storagePath
	downloadInfoSelf.ErrorFilePath = errorFilePath

	return downloadInfoSelf
}