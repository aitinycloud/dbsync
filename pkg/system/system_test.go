package system

import "testing"

func Test_GetIp(t *testing.T) {
	ip := GetIp()
	t.Log(" ip : ", ip)
	t.Log("ip : ", ip[len(ip)-1])
	t.Log("OK")
}

func Test_GetMac(t *testing.T) {
	mac := GetMac()
	t.Log(mac)
	if len(mac) > 0 {
		t.Log("mac : ", mac[len(mac)-1])
	}
}

func Test_GetCurrentDirectory(t *testing.T) {
	path := GetCurrentDirectory()
	t.Log(path)
}
