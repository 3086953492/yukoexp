package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

// 文件类型判断
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	fmt.Println(ext)
	for _, ext := range imageExtensions {
		if ext == filepath.Ext(filename) {
			return true
		}
	}
	return false
}

// 获取文件夹内所有文件并返回
func GetFiles(reportID int) ([]map[string]interface{}, error) {
	// 生成文件夹路径
    folderPath := fmt.Sprintf("attachments/%d", reportID)

    // 获取该文件夹中的所有文件
    files, err := os.ReadDir(folderPath)
    if err != nil {
        return nil, err
    }

    var fileData []map[string]interface{}
    // 遍历文件夹中的文件
    for _, file := range files {
        fileName := file.Name()
        filePath := filepath.Join(folderPath, fileName)

        fileInfo := map[string]interface{}{
            "file_name": fileName,
            "file_path": filePath,
            "is_image":  isImageFile(fileName),
        }
        fileData = append(fileData, fileInfo)
    }

    return fileData, nil
}