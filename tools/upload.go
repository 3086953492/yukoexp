package tools

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

// 上传附件并保存到指定目录
func UploadAttachments(reportID int, attachments []*multipart.FileHeader, attachmentDir string) error {
	// 创建报销单对应的目录
	reportDir := filepath.Join(".",attachmentDir, strconv.Itoa(reportID))
	err := os.MkdirAll(reportDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建附件目录失败: %v", err)
	}
	

	// 遍历附件并上传
	for _, file := range attachments {
		// 打开文件
		filePath := filepath.Join(reportDir, file.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建文件失败: %v", err)
		}
		defer dst.Close()

		// 将文件内容从上传的文件写入目标路径
		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开文件失败: %v", err)
		}
		defer src.Close()

		_, err = dst.ReadFrom(src)
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}
	return nil
}