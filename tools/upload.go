package tools

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxAttachmentBytes = 32 << 20 // 单附件上限 32 MiB

// sanitizeAttachmentSegment 去掉前导路径分隔符，使与 main 中 http.Dir("."+attachmentDir) 的物理根一致。
func sanitizeAttachmentSegment(attachmentDir string) string {
	s := strings.TrimSpace(attachmentDir)
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimPrefix(s, `\`)
	return s
}

// AttachmentReportDir 返回某报销单附件目录的相对路径（与 StaticFS 使用的 ./attachments 对齐）。
func AttachmentReportDir(reportID int, attachmentDir string) string {
	seg := sanitizeAttachmentSegment(attachmentDir)
	if seg == "" {
		seg = "attachments"
	}
	return filepath.Join(".", seg, strconv.Itoa(reportID))
}

func splitNameExt(base string) (stem, ext string) {
	i := strings.LastIndex(base, ".")
	if i <= 0 || i == len(base)-1 {
		return base, ""
	}
	return base[:i], base[i:]
}

func pickUniqueFilename(dir, base string) (string, error) {
	if base == "" || base == "." || base == ".." {
		base = "file"
	}
	candidate := base
	for n := 0; n < 10000; n++ {
		full := filepath.Join(dir, candidate)
		_, err := os.Stat(full)
		if os.IsNotExist(err) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
		stem, ext := splitNameExt(base)
		candidate = fmt.Sprintf("%s_%d%s", stem, n+1, ext)
	}
	return "", fmt.Errorf("无法为 %s 生成唯一文件名", base)
}

func saveMultipartFile(reportDir string, fh *multipart.FileHeader) error {
	finalName, err := pickUniqueFilename(reportDir, filepath.Base(fh.Filename))
	if err != nil {
		return err
	}

	src, err := fh.Open()
	if err != nil {
		return fmt.Errorf("打开上传流失败 %s: %w", finalName, err)
	}
	defer src.Close()

	tmp, err := os.CreateTemp(reportDir, ".upload-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := tmp.Name()

	cleanupTmp := func() {
		tmp.Close()
		_ = os.Remove(tmpPath)
	}

	limited := io.LimitReader(src, fh.Size+1)
	n, err := io.Copy(tmp, limited)
	if err != nil {
		cleanupTmp()
		return fmt.Errorf("写入临时文件失败 %s: %w", finalName, err)
	}
	if n != fh.Size {
		cleanupTmp()
		return fmt.Errorf("附件 %s 大小校验失败: 期望 %d 字节, 实际 %d 字节", finalName, fh.Size, n)
	}

	if err := tmp.Sync(); err != nil {
		cleanupTmp()
		return fmt.Errorf("同步临时文件失败 %s: %w", finalName, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("关闭临时文件失败 %s: %w", finalName, err)
	}

	finalPath := filepath.Join(reportDir, finalName)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("落盘失败 %s: %w", finalName, err)
	}
	return nil
}

// UploadAttachments 将 multipart 附件原子写入 report 子目录；任一步失败会删除整目录。
func UploadAttachments(reportID int, attachments []*multipart.FileHeader, attachmentDir string) error {
	reportDir := AttachmentReportDir(reportID, attachmentDir)
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return fmt.Errorf("创建附件目录失败: %w", err)
	}

	ok := false
	defer func() {
		if !ok {
			_ = os.RemoveAll(reportDir)
		}
	}()

	anySaved := false
	for _, fh := range attachments {
		if fh == nil {
			continue
		}
		if fh.Size == 0 {
			continue
		}
		if fh.Size < 0 {
			return fmt.Errorf("非法附件大小: %s", fh.Filename)
		}
		if fh.Size > maxAttachmentBytes {
			return fmt.Errorf("附件超过大小上限 (%d MiB): %s", maxAttachmentBytes>>20, fh.Filename)
		}
		if err := saveMultipartFile(reportDir, fh); err != nil {
			return err
		}
		anySaved = true
	}

	if !anySaved {
		return fmt.Errorf("没有可保存的有效附件（大小均为 0）")
	}

	ok = true
	return nil
}
