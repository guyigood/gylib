package ftpsync

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"gylib/common"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Ftp_run struct {
	Ftplist               map[string]string
	Ftp_ser               *ftp.ServerConn
	Remote_path, Cur_path string
}

func NewFtp_Sync() *Ftp_run {
	this := new(Ftp_run)
	this.Ftplist = make(map[string]string)
	return this
}

func (this *Ftp_run) GetFilelist() []string {
	result := make([]string, 0)
	pathlist := strings.Split(this.Cur_path, ",")
	for i := 0; i < len(pathlist); i++ {
		path := pathlist[i]
		err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				result = append(result, path)
				return nil
			}
			result = append(result, path)
			return nil
		})
		if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
		}
		for _, v := range result {
			fmt.Println(v)
		}
	}
	return result
}

func (this *Ftp_run) Upload_dir_file(path_dir []string) {
	this.Connect_ftp_server()
	err1 := this.Ftp_ser.ChangeDir(this.Remote_path)
	if err1 != nil {
		this.Ftp_ser.MakeDir(this.Remote_path)
		this.Ftp_ser.ChangeDir("/")
	}
	for _, v := range path_dir {
		f, _ := os.Stat(v)
		if f.IsDir() {
			dirnew := this.Path_format(v)
			//fmt.Println("dirnew ",dirnew)
			err := this.Ftp_ser.ChangeDir(dirnew)
			if err != nil {
				fmt.Println("dirnew root", err)
				err = this.Ftp_ser.MakeDir(dirnew)
				if err != nil {
					fmt.Println("mkdir", err)
				}
				err = this.Ftp_ser.ChangeDir(dirnew)
				if err != nil {
					fmt.Println("chdir", err)
				}
			}
		} else {

			this.Ftp_ser.ChangeDir("/")
			list := strings.Split(v, string(os.PathSeparator))
			tmp_dir := ""
			for i := 0; i < (len(list) - 1); i++ {
				tmp_dir += list[i] + string(os.PathSeparator)
			}
			file_name := list[len(list)-1]
			//fmt.Println("stor ",v,file_name)
			//fmt.Println(list,tmp_dir)
			dirnew := this.Path_format(tmp_dir)
			err := this.Ftp_ser.ChangeDir(dirnew)
			if err != nil {
				fmt.Println("stor changdir", err)
			}
			err = this.Stor_Ftp(v, file_name)
			if err != nil {
				fmt.Println("stor error", err)
			}
		}
	}
	this.Loginout_ftp()
}

func (this *Ftp_run) Stor_Ftp(cur_name, file_name string) error {
	file, err := os.Open(cur_name)
	defer file.Close()
	if err != nil {
		fmt.Println("readfile eroor", err)
		return nil
	}
	ftpsize, err := this.Ftp_ser.FileSize(file_name)
	if err != nil {
		ftpsize = 0
	}
	if ftpsize == common.Get_FIle_Size(cur_name) {
		return nil
	}
	err = this.Ftp_ser.Stor(file_name, file)
	return err
}

func (this *Ftp_run) Path_format(dirstr string) string {
	dirnews := strings.Replace(dirstr, this.Cur_path, "", -1)
	dirnew := strings.Replace(dirnews, "\\", "/", -1)
	if this.Remote_path != "" {
		dirnew = "/" + this.Remote_path + "/" + dirnew
	} else {
		dirnew = "/" + dirnew
	}
	dirnew = strings.Replace(dirnew, "//", "/", -1)
	return dirnew
}

//func (this *Ftp_run) ftp_upload_file(dirPth, olddir string) (files []string, files1 []string, err error) {
//	dir, err := ioutil.ReadDir(dirPth)
//	if err != nil {
//		return nil, nil, err
//	}
//	PthSep := string(os.PathSeparator)
//	for _, fi := range dir {
//		if fi.IsDir() {
//			dirnews := strings.Replace(dirPth, olddir, "", -1)
//			dirnew := strings.Replace(dirnews, "\\", "/", -1)
//			if (this.Remote_path != "") {
//				dirnew = "/" + this.Remote_path + "/" + dirnew
//			} else {
//				dirnew = "/" + dirnew
//			}
//			dirnew = strings.Replace(dirnew, "//", "/", -1)
//			//fmt.Println(dirnew)
//			//ftpcurpath,_:=ftp_ser.CurrentDir()
//			//fmt.Println(fi.Name(),dirPth,dirnew,ftpcurpath)
//			//err=ftp_ser.ChangeDir("/")
//			////ftpcurpath,_=ftp_ser.CurrentDir()
//			////fmt.Println("1",ftpcurpath)
//			//if(err!=nil){
//			//	fmt.Println("path root",err)
//			//}
//			err = this.Ftp_ser.ChangeDir(dirnew)
//			//ftpcurpath,_=ftp_ser.CurrentDir()
//			//fmt.Println("2",dirnew,ftpcurpath)
//			if (err != nil) {
//				fmt.Println("dirnew root", err)
//				//this.Loginout_ftp()
//				this.Connect_ftp_server()
//				return this.ftp_upload_file(dirPth, olddir)
//			}
//			err = this.Ftp_ser.ChangeDir(fi.Name())
//			//ftpcurpath,_=ftp_ser.CurrentDir()
//			//fmt.Println("3",ftpcurpath)
//			if (err != nil) {
//				this.Ftp_ser.MakeDir(fi.Name())
//				this.Ftp_ser.ChangeDir(fi.Name())
//				//ftpcurpath,_=ftp_ser.CurrentDir()
//				//fmt.Println("1",ftpcurpath)
//			}
//			files1 = append(files1, dirPth+PthSep+fi.Name())
//			this.Ftp_ser.NoOp()
//			this.ftp_upload_file(dirPth+PthSep+fi.Name(), olddir)
//		} else {
//			dirnews := strings.Replace(dirPth, olddir, "", -1)
//			dirnew := strings.Replace(dirnews, "\\", "/", -1)
//			if (this.Remote_path != "") {
//				dirnew = "/" + this.Remote_path + "/" + dirnew
//			} else {
//				dirnew = "/" + dirnew
//			}
//			dirnew = strings.Replace(dirnew, "//", "/", -1)
//			err = this.Ftp_ser.ChangeDir(dirnew)
//			//ftpcurpath,_=ftp_ser.CurrentDir()
//			//fmt.Println("2",dirnew,ftpcurpath)
//			if (err != nil) {
//				fmt.Println("dirnew root", err)
//				//this.Loginout_ftp()
//				this.Connect_ftp_server()
//				return this.ftp_upload_file(dirPth, olddir)
//			}
//			//files = append(files, dirPth+PthSep+fi.Name())
//			localFile := dirPth + PthSep + fi.Name()
//			ftpsize, _ := this.Ftp_ser.FileSize(fi.Name())
//
//			//fmt.Println(localFile)
//			//ftpcurpath,_:=this.Ftp_ser.CurrentDir()
//			//fmt.Println(fi.Name(),localFile,ftpcurpath,ftpsize,common.Get_FIle_Size(localFile))
//			if (ftpsize != common.Get_FIle_Size(localFile)) {
//				file, err := os.Open(localFile)
//				defer file.Close()
//				if err != nil {
//					fmt.Println(err)
//					//this.Loginout_ftp()
//					this.Connect_ftp_server()
//					return this.ftp_upload_file(dirPth, olddir)
//				}
//				ftpcurpath, _ := this.Ftp_ser.CurrentDir()
//				fmt.Println(ftpcurpath, localFile)
//				err = this.Ftp_ser.Stor(fi.Name(), file)
//				if err != nil {
//					fmt.Println(err)
//					//this.Loginout_ftp()
//					this.Connect_ftp_server()
//					return this.ftp_upload_file(dirPth, olddir)
//				}
//			} else {
//				this.Ftp_ser.NoOp()
//			}
//		}
//	}
//	return files, files1, nil
//}

func (this *Ftp_run) Connect_ftp_server() bool {
	port := this.Ftplist["ftp_port"]
	ftpserver := this.Ftplist["ftp_addr"]
	if port != "" {
		ftpserver += ":" + port
	}
	var err error
	this.Ftp_ser, err = ftp.Connect(ftpserver)
	if err != nil {
		fmt.Println(err)
		fmt.Println(time.Now().String()[:19], this.Ftplist)
		return false
	}
	err = this.Ftp_ser.Login(this.Ftplist["ftp_user"], this.Ftplist["ftp_pass"])
	if err != nil {
		fmt.Println(err)
		return false
	}
	if this.Remote_path != "" {
		err = this.Ftp_ser.ChangeDir(this.Remote_path)
		if err != nil {
			this.Ftp_ser.MakeDir(this.Remote_path)
			this.Ftp_ser.ChangeDir(this.Remote_path)
		}
	}
	return true

}

func (this *Ftp_run) Loginout_ftp() {
	this.Ftp_ser.Quit()
}

func (this *Ftp_run) Upload_ftp_server() {
	//var err error
	if this.Connect_ftp_server() == false {
		return
	}
	list := this.GetFilelist()
	//ftp_upload_file(this.Cur_path, this.Cur_path)
	this.Upload_dir_file(list)
	this.Loginout_ftp()
}
