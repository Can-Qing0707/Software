package service

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"training_eval_system/config"
	"training_eval_system/internal/model"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) SaveFile(file io.Reader, filename string) (string, error) {
	uploadDir := config.AppConfig.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	saveName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(filename))
	savePath := filepath.Join(uploadDir, saveName)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	return "/uploads/" + saveName, nil
}

func (s *FileService) GetFilePath(relativePath string) string {
	return filepath.Join(config.AppConfig.Upload.Dir, filepath.Base(relativePath))
}

func (s *FileService) DeleteFile(fileURL string) error {
	path := s.GetFilePath(fileURL)
	return os.Remove(path)
}

func (s *FileService) ExtractContent(files model.FileList) string {
	var parts []string

	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name))
		absPath := s.GetFilePath(f.URL)

		var text string
		switch ext {
		case ".docx":
			text = extractDocxText(absPath)
		case ".pdf":
			text = extractPdfText(absPath)
		case ".txt", ".md", ".vue", ".js", ".ts", ".go", ".py", ".java", ".cs",
			".html", ".css", ".json", ".xml", ".yaml", ".yml", ".sql", ".sh":
			text = extractPlainText(absPath)
		case ".zip":
			text = extractZipText(absPath)
		default:
			continue
		}

		if text != "" {
			parts = append(parts, fmt.Sprintf("=== %s ===\n%s", f.Name, text))
		}
	}

	return strings.Join(parts, "\n\n")
}

func extractPlainText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	text := string(data)
	if len(text) > 100000 {
		text = text[:100000] + "\n... (内容过长，已截断)"
	}
	return text
}

func extractDocxText(path string) string {
	r, err := zip.OpenReader(path)
	if err != nil {
		return ""
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return ""
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return ""
			}
			return parseDocxXML(data)
		}
	}
	return ""
}

type docxBody struct {
	Paragraphs []docxParagraph `xml:"p"`
}

type docxParagraph struct {
	Runs []docxRun `xml:"r"`
}

type docxRun struct {
	Text string `xml:"t"`
}

func parseDocxXML(data []byte) string {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var lines []string
	var currentLine strings.Builder
	inParagraph := false

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch el := tok.(type) {
		case xml.StartElement:
			if el.Name.Local == "p" {
				inParagraph = true
				currentLine.Reset()
			} else if el.Name.Local == "t" && inParagraph {
				for {
					innerTok, innerErr := decoder.Token()
					if innerErr != nil {
						break
					}
					if cd, ok := innerTok.(xml.CharData); ok {
						currentLine.Write(cd)
					} else if _, ok := innerTok.(xml.EndElement); ok {
						break
					}
				}
			}
		case xml.EndElement:
			if el.Name.Local == "p" && inParagraph {
				inParagraph = false
				line := strings.TrimSpace(currentLine.String())
				if line != "" {
					lines = append(lines, line)
				}
			}
		}
	}

	result := strings.Join(lines, "\n")
	if len(result) > 100000 {
		result = result[:100000] + "\n... (内容过长，已截断)"
	}
	return result
}

func extractPdfText(path string) string {
	return extractPlainTextMetadata(path)
}

func extractPlainTextMetadata(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	text := string(data)
	readable := 0
	for _, r := range text {
		if r >= 32 && r < 127 || r == '\n' || r == '\r' || r == '\t' || r > 127 {
			readable++
		}
	}

	if len(text) > 0 && float64(readable)/float64(len(text)) < 0.3 {
		return ""
	}

	const maxSize = 50000
	if len(text) > maxSize {
		text = text[:maxSize] + "\n... (内容过长，已截断)"
	}
	return text
}

func extractZipText(path string) string {
	r, err := zip.OpenReader(path)
	if err != nil {
		return ""
	}
	defer r.Close()

	var result strings.Builder
	textExts := []string{".txt", ".md", ".vue", ".js", ".ts", ".go", ".py", ".java", ".cs",
		".html", ".css", ".json", ".xml", ".yaml", ".yml", ".sql", ".sh", ".c", ".h", ".cpp", ".hpp"}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		isText := false
		for _, te := range textExts {
			if ext == te {
				isText = true
				break
			}
		}
		if !isText {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		text := string(data)
		if len(text) > 5000 {
			text = text[:5000] + "\n... (文件较长，已截断)"
		}
		result.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", f.Name, text))

		if result.Len() > 100000 {
			result.WriteString("\n... (ZIP内容过多，已截断)")
			break
		}
	}
	return result.String()
}
