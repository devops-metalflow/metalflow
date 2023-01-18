package vncproxy

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"metalflow/pkg/global"
	"net"
	"reflect"
	"strconv"
	"time"
)

const (
	VersionLength  = 12
	RfbVersion     = 3.8
	AuthStatusFail = '\x00'
	AuthStatusPass = '\x02'
)

type AuthType = int

var (
	None       = 1
	VenEncrypt = 19
	X509None   = 260
)

// nolint:gocyclo
func Connect(addr string, source, target net.Conn) (net.Conn, error) {
	isVenEncrypt, err := checkIsEncrypt(addr)
	if err != nil {
		return nil, err
	}
	if !isVenEncrypt {
		return target, nil
	}
	targetVersion, err := readCon(target, VersionLength)
	if err != nil {
		return nil, err
	}
	tv := parseVersion(targetVersion)
	if tv != RfbVersion {
		return nil, errors.New("Security proxying requires RFB protocol version 3.8 , but tenant asked for " + string(targetVersion))
	}
	_, err = target.Write(targetVersion)
	if err != nil {
		return nil, err
	}
	_, err = source.Write(targetVersion)
	if err != nil {
		return nil, err
	}
	sourceVersion, err := readCon(source, VersionLength)
	if err != nil {
		return nil, err
	}
	v := parseVersion(sourceVersion)
	if v != RfbVersion {
		return nil, errors.New("Security proxying requires RFB protocol version 3.8 , but tenant asked for " + string(sourceVersion))
	}
	authType, err := readCon(target, 1)
	if err != nil {
		return nil, err
	}
	if byte2int(authType) == 0 {
		return nil, errors.New("negotiation failed: " + string(authType))
	}
	f, err := readCon(target, byte2int(authType))
	if err != nil {
		return nil, err
	}
	permittedAuthType := make([]int, 0)
	for _, t := range f {
		permittedAuthType = append(permittedAuthType, int(t))
	}
	data := []byte("\x01\x01")
	_, err = source.Write(data)
	if err != nil {
		return nil, err
	}
	clientAuth, _ := readCon(source, 1)
	if byte2int(clientAuth) != None {
		return nil, errors.New("negotiation failed: " + string(clientAuth))
	}
	if permittedAuthType[0] != VenEncrypt {
		return nil, errors.New("is not VenEncrypt conn")
	}
	_, err = target.Write(f)
	if err != nil {
		return nil, err
	}
	return securityHandshake(target)
}

// parseVersion parse rfb proto version.
func parseVersion(version []byte) float64 {
	versionStr := string(version)
	result, _ := strconv.ParseFloat(fmt.Sprintf("%v.%v", str2int(versionStr[4:7]), str2int(versionStr[8:11])), 64) // nolint: gomnd
	return result
}

func checkIsEncrypt(addr string) (bool, error) {
	target, err := net.DialTimeout("tcp", addr, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	if err != nil {
		return false, err
	}
	targetVersion, err := readCon(target, VersionLength)
	if err != nil {
		return false, err
	}
	tv := parseVersion(targetVersion)
	if tv != RfbVersion {
		return false, errors.New("Security proxying requires RFB protocol version 3.8 , but tenant asked for " + string(targetVersion))
	}
	_, err = target.Write(targetVersion)
	if err != nil {
		return false, err
	}
	authType, err := readCon(target, 1)
	if err != nil {
		return false, err
	}
	if byte2int(authType) == 0 {
		return false, fmt.Errorf("negotiation failed: %d", authType)
	}
	f, err := readCon(target, byte2int(authType))
	if err != nil {
		return false, err
	}
	permittedAuthType := make([]int, 0)
	for _, b := range f {
		permittedAuthType = append(permittedAuthType, int(b))
	}
	err = target.Close()
	if err != nil {
		return false, err
	}
	return permittedAuthType[0] == VenEncrypt, nil
}

func readCon(c net.Conn, num int) ([]byte, error) {
	buf := make([]byte, num)
	length, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	if length != num {
		fmt.Printf("Incorrect read from socket, wanted %v bytes but got %v. Socket returned\n", num, length)
	}
	return buf, nil
}

func str2int(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func byte2int(data []byte) int {
	var ret = 0
	for i := 0; uint(i) < uint(len(data)); i++ {
		ret |= int(data[i]) << (i * 8)
	}
	return ret
}

func securityHandshake(target net.Conn) (net.Conn, error) {
	maj, _ := readCon(target, 1)
	min, _ := readCon(target, 1)
	majVer := byte2int(maj)
	minVer := byte2int(min)
	global.Log.Infof("Server sent VeNCrypt version %v.%v", majVer, minVer)
	if majVer != 0 || minVer != 2 {
		return nil, fmt.Errorf("only VeNCrypt version 0.2 is supported by this proxy, "+
			"but the server wanted to use version :%v.%v", majVer, minVer)
	}
	data := [2]byte{AuthStatusFail, AuthStatusPass}
	err := send(target, data)
	if err != nil {
		return nil, err
	}
	var isAccepted uint8
	err = receive(target, &isAccepted)
	if err != nil {
		return nil, err
	}
	if isAccepted > 0 {
		return nil, errors.New("Server could not use VeNCrypt version 0.2 ")
	}
	subTypesCnt, _ := readCon(target, 1)
	subAuthTypes := make([]int32, byte2int(subTypesCnt))
	err = receiveN(target, &subAuthTypes, byte2int(subTypesCnt))
	if err != nil {
		return nil, err
	}
	hasX509 := false
	for _, t := range subAuthTypes {
		if t == int32(X509None) {
			hasX509 = true
			break
		}
	}
	if !hasX509 {
		return nil, errors.New("Server does not support the x509None VeNCrypt ")
	}
	err = send(target, uint32(X509None))
	if err != nil {
		return nil, err
	}
	authAccepted, _ := readCon(target, 1)
	if byte2int(authAccepted) == 0 {
		return nil, errors.New("Server didn't accept the requested auth sub-type ")
	}
	// 这里使用不双向加密的vnc，不需要证书
	config := &tls.Config{
		InsecureSkipVerify: true, // nolint:gosec
	}
	conn := tls.Client(target, config)
	return conn, nil
}

func receiveN(c net.Conn, data interface{}, n int) error {
	if n == 0 {
		return nil
	}

	switch t := data.(type) {
	case *[]uint8:
		var v uint8
		for i := 0; i < n; i++ {
			if err := binary.Read(c, binary.BigEndian, &v); err != nil {
				return err
			}
			*t = append(*t, v)
		}
	case *[]int32:
		var v int32
		for i := 0; i < n; i++ {
			if err := binary.Read(c, binary.BigEndian, &v); err != nil {
				return err
			}
			*t = append(*t, v)
		}
	case *bytes.Buffer:
		var v byte
		for i := 0; i < n; i++ {
			if err := binary.Read(c, binary.BigEndian, &v); err != nil {
				return err
			}
			t.WriteByte(v)
		}
	default:
		return fmt.Errorf("unrecognized data type %v", reflect.TypeOf(data))
	}
	return nil
}

func receive(c net.Conn, data interface{}) error {
	if err := binary.Read(c, binary.BigEndian, data); err != nil {
		return err
	}
	return nil
}

// send a packet to the network.
func send(c net.Conn, data interface{}) error {
	if err := binary.Write(c, binary.BigEndian, data); err != nil {
		return err
	}
	return nil
}
