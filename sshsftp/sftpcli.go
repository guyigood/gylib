package sshsftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type SFTPClient struct {
	User, Password, Host string
	Is_ready, Port       int
	//sshClient    *ssh.Client
	sftpCli *sftp.Client
}

func NewSFTPClient() *SFTPClient {
	this := new(SFTPClient)
	this.Is_ready = 0
	return this
}

func (this *SFTPClient) SetParam(user, password, host string, port int) *SFTPClient {
	this.User = user
	this.Password = password
	this.Host = host
	this.Port = port
	return this
}

func (this *SFTPClient) Connect() {
	var err error
	this.sftpCli, err = this.InitCon()
	if err != nil {
		fmt.Println(err)
	} else {
		this.Is_ready = 1
	}
}

func (this *SFTPClient) InitCon() (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(this.Password))

	clientConfig = &ssh.ClientConfig{
		User:            this.User,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", this.Host, this.Port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		fmt.Println("ionitcon dial", err)
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		fmt.Println("ionitcon newclient", err)
		return nil, err
	}
	return sftpClient, nil
}

func (this *SFTPClient) Close() {
	this.sftpCli.Close()
	this.Is_ready = 0
}

func (this *SFTPClient) UploadFile(localFilePath string, remotePath string) bool {
	if this.Is_ready == 0 {
		fmt.Println("uploadfile ", this.Is_ready)
		return false
	}

	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		return false
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := this.sftpCli.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		//log.Fatal(err)
		return false
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return false
		//log.Fatal(err)

	}
	_, err2 := dstFile.Write(ff)
	if err2 != nil {
		fmt.Println("uploadfile write", err2)
		return false
	}
	return true
	//fmt.Println(localFilePath + "  copy file to remote server finished!")
}

func (this *SFTPClient) UploadDirectory(localPath string, remotePath string) {
	if this.Is_ready == 0 {
		fmt.Println("upload not is ready")
		return
	}
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		//log.Fatal("read dir list fail ", err)
		fmt.Println("read dir list fail ", err)
		return
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			this.sftpCli.Mkdir(remoteFilePath)
			this.UploadDirectory(localFilePath, remoteFilePath)
		} else {
			this.UploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	fmt.Println(localPath + "  copy directory to remote server finished!")
}
