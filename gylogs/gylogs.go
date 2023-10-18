package gylogs

import (
	"fmt"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	"os"
	"time"
)

type Gylogs struct {
	DirName, Filename string
}

func NewGylogs() *Gylogs {
	that := new(Gylogs)
	return that
}

func (that *Gylogs) SetFixFile(dirname, filename string) *Gylogs {
	dirpath := dirname
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	//dirpath += common.Int2Date_str(time.Now().Unix()) + "/"
	//if !common.PathExists(dirpath) {
	//	os.Mkdir(dirpath, os.ModePerm)
	//}
	logname := filename
	if logname == "" {
		logname = "1"
	}
	that.Filename = dirpath + logname
	return that
}

func (that *Gylogs) SetDirFile(dirname, filename string) *Gylogs {
	dirpath := dirname
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	dirpath += common.Int2Date_str(time.Now().Unix()) + "/"
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	logname := filename
	if logname == "" {
		logname = "1"
	}
	that.Filename = dirpath + logname + "-" + datatype.Type2str(time.Now().Hour()) + ".log"
	return that
}

func (that *Gylogs) SaveToFile(msg string) {
	f, err := os.OpenFile(that.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write([]byte(msg + "\n"))

}

func (that *Gylogs) StringToSaveFile(filename, msg string) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write([]byte(msg + "\n"))

}
