// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"
)

// GetModulePath ets the path of the project containing go.mod upwards from the current directory.
func GetModulePath(dir string) (string, error) {
	modPath := filepath.Join(dir, "go.mod")
	_, err := os.Stat(modPath)
	if err != nil {
		dir = filepath.Dir(dir)
		if dir == "/" || dir == "." {
			return "", errors.New("unavailable dir")
		}
		return GetModulePath(dir)
	}
	return dir, nil
}

// UnsafeToBytes converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func UnsafeToBytes(s string) []byte {
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// UnsafeToString converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func UnsafeToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GetLocalIP returns the local IP
func GetLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// Check if the IP is a loopback address.
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("local ip not found")
}

// IsLocalIP checks whether the IP address is local
func IsLocalIP(ip string) (bool, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}

	for i := range addrs {
		intf, _, err := net.ParseCIDR(addrs[i].String())
		if err != nil {
			return false, err
		}
		if net.ParseIP(ip).Equal(intf) {
			return true, nil
		}
	}
	return false, nil
}
