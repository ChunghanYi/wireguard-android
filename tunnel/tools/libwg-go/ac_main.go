/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2025 Slowboot(chunghan.yi@gmail.com). All Rights Reserved.
 */

package main

// #cgo LDFLAGS: -llog
// #include <android/log.h>
import "C"

import (
	"fmt"
	"time"
)

func acLogDebug(format string, args ...interface{}) {
	C.__android_log_write(C.ANDROID_LOG_DEBUG,
		cstring("WireGuard/GoBackend/AC/"),
		cstring(fmt.Sprintf(format, args...)))
}

func acLogError(format string, args ...interface{}) {
	C.__android_log_write(C.ANDROID_LOG_ERROR,
		cstring("WireGuard/GoBackend/AC/"),
		cstring(fmt.Sprintf(format, args...)))
}

//export acTurnOn
func acTurnOn(serverIp string, serverPort string, privatekey string, publickey string) *C.char {
	var serverAddr string
	if serverIp == "" || serverPort == "" {
		acLogDebug("[AC] serverIp or serverPort string is empty.")
		return nil 
	}
	serverAddr = serverIp + ":" + serverPort

	client := NewClient(serverAddr)

	var trycount int
	for {
		if client.connectServer(serverAddr) == nil {
			client.connected = true
			acLogDebug("[AC] Client connected to server successfully")
			break
		} else {
			acLogError("[AC] Failed to connect to server")

			trycount++
			if trycount >= 2 {
				client.connected = false
				break
			}
			time.Sleep(time.Second * 2)
			acLogDebug("[AC] Retrying to connect to server...")
		}
	}

	if !client.connected {
		acLogError("[AC] Connection to server is impossible.")
		client.Close()
		return nil 
	}

	var rmsg Message
	var client_vpnIP string
	var settings string
	trycount = 0
	for {
		if client.sendHelloMessage(publickey, &rmsg, &client_vpnIP) {
			settings = client.sendPingMessage(privatekey, publickey, &rmsg, client_vpnIP)
			if settings != "unknown" {
				client.Close()
				break
			}
		}

		trycount++
		if trycount >= 2 {
			acLogError("[AC] Auto Connection to server is impossible.")
			client.Close()
			break
		}
		time.Sleep(time.Second * 2)
	}

	if settings == "unknown" {
		return nil
	} else {
		return C.CString(settings)  
	}
}

//export acTurnOff
func acTurnOff(serverIp string, serverPort string, publickey string) int32 {
	var serverAddr string
	serverAddr = serverIp + ":" + serverPort

	client := NewClient(serverAddr)

	var trycount int
	for {
		if client.connectServer(serverAddr) == nil {
			client.connected = true
			acLogDebug("[AC] Client connected to server successfully")
			break
		} else {
			acLogError("[AC] Failed to connect to server")

			trycount++
			if trycount >= 2 {
				client.connected = false
				break
			}
			time.Sleep(time.Second * 2)
			acLogDebug("[AC] Retrying to connect to server...")
		}
	}

	if !client.connected {
		acLogError("[AC] Connection to server is impossible.")
		client.Close()
		return -1 
	}

	var rmsg Message
	var ret int32 
	trycount = 0
	for {
		if client.sendByeMessage(publickey, &rmsg) {
			client.Close()
			ret = 0
			break
		}

		trycount++
		if trycount >= 2 {
			acLogError("[AC] Auto Connection to server is impossible.")
			client.Close()
			ret = -1
			break
		}
		time.Sleep(time.Second * 2)
	}

	return ret
}
