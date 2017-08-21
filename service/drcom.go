package service

import (
	"bytes"
	"errors"
	"fmt"
	"time"
	"math/rand"
	"encoding/hex"
	"encoding/binary"

    "../utils"
)


func (s *Service) Challenge(tryTimes int) (err error) {
	var (
		r   []byte
		buf = []byte{0x01, (byte)(0x02 + tryTimes),
			byte(rand.Int()), byte(rand.Int()), 0x6a,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00}
		conn = s.conn
	)

	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	log.Debug("(send) request %s", hex.EncodeToString(buf))

	r = make([]byte, 128)
	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	log.Debug("(recv) response %s", hex.EncodeToString(r))

	if r[0] == 0x02 {
		copy(s.salt, r[4:8])
		copy(s.clientIP, r[20:24])
		return
	}
	err = errors.New("challenge receive head is not correct")
	return
}

func (s *Service) Login() (err error) {
	var (
		r    []byte
		buf  []byte
		conn = s.conn
	)
	if buf, err = s.bufIn(); err != nil {
		return fmt.Errorf("service.bufIn() error(%v)", err)
	}
	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	log.Debug("(send) login %s", hex.EncodeToString(buf))

	r = make([]byte, 512)
	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	log.Debug("(recv) auth_info %s", hex.EncodeToString(r))
	if r[0] != 0x04 {
		if r[0] == 0x05 {
			if r[4] == 0x0B {
				err = errors.New("invalid mac address, please select the address registered in http://192.168.161.2/")
			} else {
				err = errors.New("invalid username or password")
			}
		} else {
			err = errors.New("fmtin failed: unknown error")
		}
		return
	}
	// 保存 tail1. 构造 keep38 要用 md5a(在mkptk中保存) 和 tail1
	// 注销也要用 tail1
	copy(s.tail1, r[23:39])

	// 记录成功信息
	//log.Critical("%s\t%s\tUsed Time:%d\tBalance:%.2f", s.conf.Username, s.conf.Password, binary.LittleEndian.Uint32(r[5:9]), math.Float32frombits(binary.LittleEndian.Uint32(r[13:17])))
	log.Critical("%s\t%s\tUsed Time:%d\tBalance:%.2f", s.conf.Username, s.conf.Password, binary.LittleEndian.Uint32(r[5:9]), float32(binary.LittleEndian.Uint32(r[13:17]))/100)
	return
}

func (s *Service) bufIn() (buf []byte, err error) {
	var (
		md5a, md5b, md5c, mac []byte
	)
	buf = make([]byte, 0, 334+(len(s.conf.Password)-1)/4*4)
	buf = append(buf, _codeIn, _type, _eof,
		byte(len(s.conf.Username)+20)) // [0:4]
	// md5a
	md5a = s.md5([]byte{_codeIn, _type}, s.salt, []byte(s.conf.Password))
	copy(s.md5a, md5a)
	buf = append(buf, md5a...) // [4:20]
	// username
	user := make([]byte, 36)
	copy(user, []byte(s.conf.Username))
	buf = append(buf, user...)                    // [20:56]
	buf = append(buf, _controlCheck, _adapterNum) //[56:58]
	// md5a xor mac
	if mac, err = s.mac(); err != nil {
		return buf, fmt.Errorf("service.mac() error(%v)", err)
	}
	for i := 0; i < 6; i++ {
		mac[i] = mac[i] ^ s.md5a[i]
	}
	buf = append(buf, mac...) // [58:64]
	// md5b
	md5b = s.md5([]byte{0x01}, []byte(s.conf.Password), []byte(s.salt), []byte{0x00, 0x00, 0x00, 0x00})
	buf = append(buf, md5b...)                      // [64:80]
	buf = append(buf, byte(0x01))                   // [80:81]
	buf = append(buf, s.clientIP...)                // [81:85]
	buf = append(buf, bytes.Repeat(_emptyIP, 3)...) // [85:97]
	// md5c
	tmp := make([]byte, len(buf))
	copy(tmp, buf)
	tmp = append(tmp, []byte{0x14, 0x00, 0x07, 0x0b}...)
	md5c = s.md5(tmp)
	buf = append(buf, md5c[:8]...)   // [97:105]
	buf = append(buf, _ipDog)        // [105:106]
	buf = append(buf, _delimiter...) // [106:110]
	hostname := make([]byte, 32)
	copy(hostname, []byte(s.conf.Hostname))
	buf = append(buf, hostname...)                       // [110:142]
	buf = append(buf, _primaryDNS...)                    // [142:146]
	buf = append(buf, _dhcpServer...)                    // [146:150]
	buf = append(buf, _emptyIP...)                       // secondary dns, [150:154]
	buf = append(buf, bytes.Repeat(_delimiter, 2)...)    // [154,162]
	buf = append(buf, []byte{0x94, 0x00, 0x00, 0x00}...) // [162,166]
	buf = append(buf, []byte{0x06, 0x00, 0x00, 0x00}...) // [166,170]
	buf = append(buf, []byte{0x02, 0x00, 0x00, 0x00}...) // [170,174]
	buf = append(buf, []byte{0xf0, 0x23, 0x00, 0x00}...) // [174,178]
	buf = append(buf, []byte{0x02, 0x00, 0x00, 0x00}...) // [178,182]
	buf = append(buf, []byte{
		0x44, 0x72, 0x43, 0x4f,
		0x4d, 0x00, 0xcf, 0x07}...) // [182,190]
	buf = append(buf, 0x6a)                              // [190,191]
	buf = append(buf, bytes.Repeat([]byte{0x00}, 55)...) // [191:246]
	buf = append(buf, []byte{
		0x33, 0x64, 0x63, 0x37,
		0x39, 0x66, 0x35, 0x32,
		0x31, 0x32, 0x65, 0x38,
		0x31, 0x37, 0x30, 0x61,
		0x63, 0x66, 0x61, 0x39,
		0x65, 0x63, 0x39, 0x35,
		0x66, 0x31, 0x64, 0x37,
		0x34, 0x39, 0x31, 0x36,
		0x35, 0x34, 0x32, 0x62,
		0x65, 0x37, 0x62, 0x31,
	}...) // [246:286]
	buf = append(buf, bytes.Repeat([]byte{0x00}, 24)...) // [286:310]
	buf = append(buf, _authVersion...)                   // [310:312]
	buf = append(buf, 0x00)                              // [312:313]
	pwdLen := len(s.conf.Password)
	if pwdLen > 16 {
		pwdLen = 16
	}
	buf = append(buf, byte(pwdLen)) // [313:314]
	ror := s.ror(s.md5a, []byte(s.conf.Password))
	buf = append(buf, ror[:pwdLen]...)       // [314:314+pwdLen]
	buf = append(buf, []byte{0x02, 0x0c}...) // [314+l:316+l]
	tmp = make([]byte, 0, len(buf))
	copy(tmp, buf)
	tmp = append(tmp, []byte{0x01, 0x26, 0x07, 0x11, 0x00, 0x00}...)
	tmp = append(tmp, mac[:4]...)
	sum := s.checkSum(tmp)
	buf = append(buf, sum[:4]...)            // [316+l,320+l]
	buf = append(buf, []byte{0x00, 0x00}...) // [320+l,322+l]
	buf = append(buf, mac...)                // [322+l,328+l]
	zeroCount := (4 - pwdLen%4) % 4
	buf = append(buf, bytes.Repeat([]byte{0x00}, zeroCount)...)
	for i := 0; i < 2; i++ {
		buf = append(buf, byte(rand.Int()))
	}
	return
}

func (s *Service) Alive() (err error) {
	var (
		r, buf []byte
		conn   = s.conn
	)
	buf = s.buf38()
	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	log.Debug("(send) keepalive_38 %s", hex.EncodeToString(buf))

	r = make([]byte, 128)
	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	log.Debug("(recv) response_38 %s", hex.EncodeToString(r))
//	s.keepAliveVer[0] = r[28]
//	s.keepAliveVer[1] = r[29]
	time.Sleep(2 * time.Second)
//	if s.extra() {
//		buf = s.buf40(true, true)
//		if _, err = conn.Write(buf); err != nil {
//			return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
//		}
//		r = make([]byte, 512)
//		if _, err = conn.Read(r); err != nil {
//			return fmt.Errorf("conn.Read() error(%v)", err)
//		}
//		s.Count++
//	}

	// 40_1
	buf = s.buf40(true, false)
	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	log.Debug("(send) keepalive_40_1 %s", hex.EncodeToString(buf))

	r = make([]byte, 40)
	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	log.Debug("(recv) keepalive_40_2 %s", hex.EncodeToString(r))
	s.Count++
	copy(s.tail2, r[16:20])
	// 40_2
	buf = s.buf40(false, false)
	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	log.Debug("(send) keepalive_40_3 %s", hex.EncodeToString(buf))

	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	log.Debug("(recv) keepalive_40_4 %s\n", hex.EncodeToString(r))
	s.Count++
	return
}

func (s *Service) buf38() (buf []byte) {
	buf = make([]byte, 0, 38)
	buf = append(buf, byte(0xff))                       // [0:1]
	buf = append(buf, s.md5a...)                        // [1:17]
	buf = append(buf, bytes.Repeat([]byte{0x00}, 3)...) // [17:20]
	buf = append(buf, s.tail1...)                       // [20:36]
	for i := 0; i < 2; i++ {                            // [36:38]
		buf = append(buf, byte(rand.Int()))
	}
	return
}

func (s *Service) buf40(first, extra bool) (buf []byte) {
	buf = make([]byte, 0, 40)
	buf = append(buf, []byte{0x07, byte(s.Count), 0x28, 0x00, 0x0b}...) // [0:5]
	// keep40_1   keep40_2
	// 发送  接收  发送  接收
	// 0x01 0x02 0x03 0xx04
	// [5:6]
	if first || extra { //keep40_1 keep40_extra 是 0x01
		buf = append(buf, byte(0x01))
	} else {
		buf = append(buf, byte(0x03))
	}
	// [6:8]
	if extra {
		buf = append(buf, []byte{0x0f, 0x27}...)
	} else {
		buf = append(buf, []byte{s.keepAliveVer[0], s.keepAliveVer[1]}...)
	}
	// [8:10]
	for i := 0; i < 2; i++ {
		buf = append(buf, byte(rand.Int()))
	}
	buf = append(buf, bytes.Repeat([]byte{0x00}, 6)...) //[10:16]
	buf = append(buf, s.tail2...)                       // [16:20]
	buf = append(buf, bytes.Repeat([]byte{0x00}, 4)...) //[20:24]
	if !first {
		tmp := make([]byte, len(buf))
		copy(tmp, buf)
		tmp = append(tmp, s.clientIP...)
		sum := s.crc(tmp)
		buf = append(buf, sum...)                           // [24:28]
		buf = append(buf, s.clientIP...)                    // [28:32]
		buf = append(buf, bytes.Repeat([]byte{0x00}, 8)...) //[32:40]
	}
	if len(buf) < 40 {
		buf = append(buf, bytes.Repeat([]byte{0x00}, 40-len(buf))...)
	}
	return
}

func (s *Service) Logout() (err error) {
	var (
		r, buf []byte
		conn   = s.conn
	)
	if buf, err = s.bufOut(); err != nil {
		return fmt.Errorf("service.bufOut() error(%v)", err)
	}
	if _, err = conn.Write(buf); err != nil {
		return fmt.Errorf("conn.Write(%v) error(%v)", buf, err)
	}
	r = make([]byte, 512)
	if _, err = conn.Read(r); err != nil {
		return fmt.Errorf("conn.Read() error(%v)", err)
	}
	if r[0] != 0x04 {
		err = errors.New("failed to fmtout: unknown error")
	}
	return
}

func (s *Service) bufOut() (buf []byte, err error) {
	var (
		md5, mac []byte
	)
	buf = make([]byte, 0, 80)
	buf = append(buf, _codeOut, _type, _eof, byte(len(s.conf.Username)+20))
	// md5
	md5 = s.md5([]byte{_codeOut, _type}, s.salt, []byte(s.conf.Password))
	buf = append(buf, md5...)
	tmp := make([]byte, 36)
	copy(tmp, []byte(s.conf.Username))
	buf = append(buf, tmp...)
	buf = append(buf, _controlCheck, _adapterNum)
	// md5 xor mac
	if mac, err = s.mac(); err != nil {
		buf = nil
		return buf, fmt.Errorf("service.mac() error(%v)", err)
	}
	for i := 0; i < 6; i++ {
		mac[i] = mac[i] ^ md5[i]
	}
	buf = append(buf, mac...)     // [58:64]
	buf = append(buf, s.tail1...) // [64:80]
	return
}
