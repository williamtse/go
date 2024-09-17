package downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/williamtse/gopkg/encrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	cd  Conf
	dir string
}

type Conf struct {
	Dir string
}

func NewClient(cd Conf, logger log.Logger) *Client {
	return &Client{
		cd: cd,
	}
}

func (d *Client) makeDownloadDir() error {
	// 创建按年月日目录
	date := time.Now()
	d.dir = filepath.Join(d.cd.Dir, date.Format("2006/01/02"))
	err := os.MkdirAll(d.dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	return nil
}

// 下载文件到指定路径
func (*Client) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("获取 URL 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败: %s", resp.Status)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func (d *Client) DownloadUrls(ctx context.Context, urls []string) (string, error) {
	// 创建下载目录
	err := d.makeDownloadDir()
	if err != nil {
		return "", err
	}
	// 下载并保存文件
	var downloadedFiles []string
	for _, item := range urls {
		fileName := filepath.Base(item) // 从 URL 中提取文件名
		filePath := filepath.Join(d.dir, fileName)

		err := d.downloadFile(item, filePath)
		if err != nil {
			return "", fmt.Errorf("下载文件失败: %v", err)
		}
		downloadedFiles = append(downloadedFiles, filePath)
	}
	downloadUrl, err := d.createZip(downloadedFiles)
	if err != nil {
		return "", err
	}
	downloadUrlKey, err := encrypt.HashPassword(downloadUrl)
	if err != nil {
		return "", err
	}

	return downloadUrlKey, nil
}

// CreateZip 创建 ZIP 文件
func (d *Client) createZip(files []string) (string, error) {
	zipFilePath := filepath.Join(d.dir, time.Now().String()+".zip")
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("创建 ZIP 文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err := d.addFileToZip(zipWriter, file)
		if err != nil {
			return "", fmt.Errorf("添加文件到 ZIP 失败: %v", err)
		}
	}

	return zipFilePath, nil
}

// 将单个文件添加到 ZIP 文件中
func (d *Client) addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	zipFileWriter, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件条目失败: %v", err)
	}

	_, err = io.Copy(zipFileWriter, file)
	if err != nil {
		return fmt.Errorf("写入 ZIP 文件条目失败: %v", err)
	}

	return nil
}
