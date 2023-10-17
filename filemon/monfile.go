package filemon

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type FileNotify struct {
	FilePath    string
	Host        string
	Remote_Path string
	Watcher     *fsnotify.Watcher
	Mutex       sync.Mutex
}

func NewFileNotifySever() *FileNotify {
	this := new(FileNotify)
	return this
}

func (this *FileNotify) Print(args ...interface{}) {
	fmt.Println(time.Now(), args)
}
func (this *FileNotify) isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		this.Print("error:", err.Error())
		return false
	}
	return fileInfo.IsDir()
}
func (this *FileNotify) watchPath(filePath string) {
	this.Print("watchPath:", filePath)
	err := this.Watcher.Watch(filePath)
	if err != nil {
		this.Print(err.Error())
		return
	}
}
func (this *FileNotify) broweDir(path string) {
	this.Print("broweDir:", path)
	dir, err := os.Open(path)
	if err != nil {
		this.Print("error:", err.Error())
		return
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		this.Print("error:", err.Error())
		return
	}
	for _, name := range names {
		dirPath := path + "/" + name
		if !this.isDir(dirPath) {
			continue
		}
		this.watchPath(dirPath)
		this.broweDir(dirPath)
	}
}

func (this *FileNotify) Run_file_sync() {
	var err error
	this.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer this.Watcher.Close()
	this.broweDir(this.FilePath)
	this.watchPath(this.FilePath)
	this.dealWatch()
}
func (this *FileNotify) copy(event *fsnotify.FileEvent) *exec.Cmd {
	return exec.Command(
		"scp",
		"-r",
		"-P 23456",
		event.Name,
		this.Host+":"+this.Remote_Path+strings.TrimPrefix(event.Name, this.FilePath))
}
func (this *FileNotify) remove(event *fsnotify.FileEvent) *exec.Cmd {
	return exec.Command(
		"ssh",
		"-p 23456",
		this.Host,
		`rm -r `+this.Remote_Path+strings.TrimPrefix(event.Name, this.FilePath)+``)
}
func (this *FileNotify) dealWatch() {
	for {
		func() {
			//mutex.Lock()
			//defer mutex.Unlock()
			select {
			case event := <-this.Watcher.Event:
				this.Print("event: ", event)
				var cmd *exec.Cmd
				if event.IsCreate() || event.IsModify() {
					cmd = this.copy(event)
				}
				if event.IsDelete() || event.IsRename() {
					cmd = this.remove(event)
				}
				this.Print("cmd:", cmd.Args)
				stderr, err := cmd.StderrPipe()
				if err != nil {
					this.Print(err.Error())
					return
				}
				defer stderr.Close()
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					this.Print(err.Error())
					return
				}
				defer stdout.Close()
				if err = cmd.Start(); err != nil {
					this.Print(err.Error())
					return
				}
				errBytes, err := ioutil.ReadAll(stderr)
				if err != nil {
					this.Print(err.Error())
					return
				}
				outBytes, err := ioutil.ReadAll(stdout)
				if err != nil {
					this.Print(err.Error())
					return
				}
				if len(errBytes) != 0 {
					this.Print("errors:", string(errBytes))
				}
				if len(outBytes) != 0 {
					this.Print("output:", string(outBytes))
				}
				if err = cmd.Wait(); err != nil {
					this.Print(err.Error())
				}
			case err := <-this.Watcher.Error:
				this.Print("error: ", err.Error())
			}
		}()
	}
}
