package scan

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/log"
)

type SFTPClient struct {
	IP       string
	Port     string
	Username string
	Password string
	taskID   int
}

func NewSFTPClient(ip, port, username, password string) *SFTPClient {
	return &SFTPClient{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SFTPClient) ScanSshFiles(remotePath string, taskID int) error {
	s.taskID = taskID
	// Set up the SSH client configuration
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Connect to the remote server
	client, err := ssh.Dial("tcp", s.IP+":"+s.Port, config)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return fmt.Errorf("连接超时: %v", err)
		}
		return fmt.Errorf("无法连接到SSH服务器，可能是用户名或密码错误: %v", err)
	}
	defer client.Close()

	// Create a new SFTP session
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("创建SFTP会话失败：%s", err.Error()))
		return fmt.Errorf("创建SFTP会话失败: %v", err)
	}
	defer sftp.Close()

	// Process the remote files or directories
	return s.processRemotePath(sftp, remotePath)
}

func (s *SFTPClient) processRemotePath(sftp *sftp.Client, remotePath string) error {
	// Check if the remote path is a file or a directory
	fileInfo, err := sftp.Stat(remotePath)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("ssh扫描错误：%s", err.Error()))
		return err
	}

	if fileInfo.IsDir() {
		// Process each file in the remote directory
		files, err := sftp.ReadDir(remotePath)
		if err != nil {
			return err
		}

		for _, file := range files {
			filePath := remotePath + "/" + file.Name()
			if file.IsDir() {
				// Recursively process subdirectories
				err := s.processRemotePath(sftp, filePath)
				if err != nil {
					return err
				}
			} else {
				// Download the file
				err := s.downloadFile(sftp, filePath, file)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Download the single file
		err := s.downloadFile(sftp, remotePath, fileInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SFTPClient) ReadFile(filePath string) ([]byte, error) {
	// SSH client config
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the SSH server

	client, err := ssh.Dial("tcp", s.IP+":"+s.Port, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}
	defer client.Close()

	// Create a new SFTP session
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Open the remote file
	file, err := sftpClient.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open remote file: %v", err)
	}
	defer file.Close()

	// Read the file content
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("failed to read remote file content: %v", err)
	}

	return buf.Bytes(), nil
}

func (s *SFTPClient) downloadFile(sftp *sftp.Client, remotePath string, fileInfo os.FileInfo) error {
	// Open the remote file
	srcFile, err := sftp.Open(remotePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Read the contents of the remote file
	contents, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}

	// TODO: 处理获取到的文件内容
	//fmt.Println("File contents of:", remotePath)
	// TODO 下面是检测模块
	modificationTime := fileInfo.ModTime().Unix()
	size := fileInfo.Size()

	sshScanner := NewSshScanner(remotePath, true)
	results, err := sshScanner.SshScan(contents)
	if err != nil {
		return err
	} else {
		taskCheckList := dbpos.TaskCheckList{
			TaskBaseID:       int64(s.taskID),
			Path:             remotePath,
			Results:          results,
			ModificationTime: modificationTime,
			Size:             strconv.FormatInt(size, 10),
		}
		db.Model(&dbpos.TaskCheckList{}).Create(&taskCheckList)
	}
	return nil
}
