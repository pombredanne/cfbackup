package cfbackup

import (
	"fmt"

	"github.com/pivotalservices/beergut/persistence"
	"github.com/pivotalservices/beergut/command"
	"github.com/xchapter7x/goutil"
)

const (
	SD_PRODUCT   string = "Product"
	SD_COMPONENT string = "Component"
	SD_IDENTITY  string = "Identity"
	SD_IP        string = "Ip"
	SD_USER      string = "User"
	SD_PASS      string = "Pass"
	SD_VCAPUSER  string = "VcapUser"
	SD_VCAPPASS  string = "VcapPass"
)

type (
	stringGetterSetter interface {
		Get(string) string
		Set(string, string)
	}

	SystemDump interface {
		stringGetterSetter
		Error() error
		GetDumper() (dumper persistence.Dumper, err error)
	}

	SystemInfo struct {
		goutil.GetSet
		Product   string
		Component string
		Identity  string
		Ip        string
		User      string
		Pass      string
		VcapUser  string
		VcapPass  string
	}

	PgInfo struct {
		SystemInfo
	}

	MysqlInfo struct {
		SystemInfo
	}

	NfsInfo struct {
		SystemInfo
	}
)

func (s *SystemInfo) Get(name string) string {
	return s.GetSet.Get(s, name).(string)
}

func (s *SystemInfo) Set(name string, val string) {
	s.GetSet.Set(s, name, val)
}

func (s *NfsInfo) GetDumper() (dumper persistence.Dumper, err error) {
	return NewNFSBackup(s.Pass, s.Ip)
}

func (s *MysqlInfo) GetDumper() (dumper persistence.Dumper, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewRemoteMysqlDump(s.User, s.Pass, sshConfig)
}

func (s *PgInfo) GetDumper() (dumper persistence.Dumper, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewPgRemoteDump(2544, s.Component, s.User, s.Pass, sshConfig)
}

func (s *SystemInfo) GetDumper() (dumper persistence.Dumper, err error) {
	panic("you have to extend SystemInfo and implement GetDumper method on the child")
	return
}

func (s *SystemInfo) Error() (err error) {
	if s.Product == "" ||
		s.Component == "" ||
		s.Identity == "" ||
		s.Ip == "" ||
		s.User == "" ||
		s.Pass == "" ||
		s.VcapUser == "" ||
		s.VcapPass == "" {
		err = fmt.Errorf("invalid or incomplete system info object: %s", s)
	}
	return
}
