package zipfile

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)

// 压缩文件
// src 可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
type Gyzip struct {
}

func NewGyzip() *Gyzip {
	that := new(Gyzip)
	return that
}

func (that *Gyzip) Compress(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	files := []*os.File{f}
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := that.compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (that *Gyzip) compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		// 增加对空目录的判断
		if len(fileInfos) <= 0 {
			header, err := zip.FileInfoHeader(info)
			header.Name = prefix
			if err != nil {
				fmt.Println("error is:" + err.Error())
				return err
			}
			_, err = zw.CreateHeader(header)
			if err != nil {
				fmt.Println("create error is:" + err.Error())
				return err
			}
			file.Close()
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = that.compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 解压
func (that *Gyzip) DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(that.getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func (that *Gyzip) getDir(path string) string {
	return that.subString(path, 0, strings.LastIndex(path, "/"))
}

func (that *Gyzip) subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
