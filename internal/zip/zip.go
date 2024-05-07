package zip

import (
	"archive/zip"
	"fmt"
	"strings"

	"github.com/1boombacks1/zipViewer/internal/model"
)

type ZipReader struct{}

func New() ZipReader {
	return ZipReader{}
}

func (z ZipReader) GetNamesByExt(path string, ext string) ([]model.File, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("ZipReader - GetNames - zip.OpenReader: %w", err)
	}
	defer r.Close()

	files := make([]model.File, 0, len(r.File))
	for _, f := range r.File {
		if strings.Compare(ext, "*") == 0 {
			file := model.File{
				Name:         f.Name,
				ModifiedDate: f.Modified.Format("2006-01-02 15:04:05"),
				Size:         f.FileInfo().Size() / 1024,
			}
			files = append(files, file)
		} else if strings.Contains(f.Name, ".") {
			fields := strings.Split(f.Name, ".")
			fileExt := fields[len(fields)-1]

			if strings.Compare(fileExt, ext) == 0 {
				file := model.File{
					Name:         f.Name,
					ModifiedDate: f.Modified.Format("2006-01-02 15:04:05"),
					Size:         f.FileInfo().Size() / 1024,
				}
				files = append(files, file)
			}
		}
	}

	return files, nil
}
