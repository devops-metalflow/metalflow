package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"path"
	"strings"
	"time"
)

type SshConfig struct {
	Address  string
	Port     int
	Username string
	Password string
	Protocol string
	Timeout  int // 默认ssh连接超时时间
}

// Option 尝试一下go编程范式Functional Options: https://coolshell.cn/articles/21146.html#Functional_Options
type Option func(*SshConfig)

func Timeout(timeout int) Option {
	return func(config *SshConfig) {
		config.Timeout = timeout
	}
}

func Protocol(p string) Option {
	return func(config *SshConfig) {
		config.Protocol = p
	}
}

func NewSshConfig(addr string, port int, username, password string, options ...Option) *SshConfig {
	sshConfig := SshConfig{
		Address:  addr,
		Port:     port,
		Username: username,
		Password: password,
		Protocol: "tcp",
		Timeout:  5,
	}
	for _, option := range options {
		option(&sshConfig)
	}
	return &sshConfig
}

// GetSshClient 获取ssh连接
func GetSshClient(config *SshConfig) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(config.Password))

	clientConfig = &ssh.ClientConfig{
		User:    config.Username,
		Auth:    auth,
		Timeout: time.Second * time.Duration(config.Timeout),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", config.Address, config.Port)
	client, err = ssh.Dial(config.Protocol, addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("无法连接ssh, 地址%s, 错误信息%v", addr, err)
	}
	return client, nil
}

// IsSafetyCmd 判断命令是否运行的安全命令
func IsSafetyCmd(cmd string) error {
	// 避免rm * 或 rm /*等命令直接出现, 删除命令指定全路径
	c := strings.ToLower(cmd)
	// List of dangerous command patterns to exclude
	dangerousPatterns := []string{
		"rm /",
		"rm -rf /",
	}
	for _, pattern := range dangerousPatterns {
		if strings.HasPrefix(c, pattern) && len(strings.Split(c, "/")) <= 2 {
			return fmt.Errorf("rm命令%s不能删除小于2级目录的文件", cmd)
		}
	}
	return nil
}
