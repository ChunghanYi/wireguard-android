/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2025 Slowboot(chunghan.yi@gmail.com). All Rights Reserved.
 */

package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	remoteAddr *Addr    // remote address
	localAddr  *Addr    // local address
	conn       net.Conn // connect server obj, receive chan, send chan
	connected  bool     // is connect flag
}

func (c *Client) connectServer(address string) error {
	c.remoteAddr = NewAddr(address)

	//Set the timeout duration
	timeout := 3 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		acLogDebug("[AC] Dial failed:", err)
		return err
	}

	c.conn = conn
	c.connected = true
	c.localAddr = NewAddr(conn.LocalAddr().String())

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.connected = false
	}
}

func (c *Client) sendMsg(smsg *Message) bool {
	var s string
	s = smsg.Msg_type +
		smsg.Mac_addr +
		smsg.VpnIP +
		smsg.VpnNetmask +
		smsg.Public_key +
		smsg.EpIP +
		smsg.EpPort +
		smsg.Allowed_ips

	_, err := c.conn.Write([]byte(s))
	if err != nil {
		acLogError("c.conn.Write() failed.", err)
		return false
	}
	return true
}

func (c *Client) recvMsg(rmsg *Message) bool {
	data := make([]byte, 1024)
	_, err := c.conn.Read(data)
	if err != nil {
		acLogError("[AC] Error reading response:", err)
		return false
	}

	s := string(data)
	t := strings.Split(s, "\n")

	rmsg.Msg_type = t[0]
	rmsg.Mac_addr = t[1]
	rmsg.VpnIP = t[2]
	rmsg.VpnNetmask = t[3]
	rmsg.Public_key = t[4]
	rmsg.EpIP = t[5]
	rmsg.EpPort = t[6]
	rmsg.Allowed_ips = t[7]

	u := strings.Split(rmsg.Msg_type, ":=") //cmd:=HELLO
	switch u[1] {
	case "HELLO":
		acLogDebug("[AC] <<< cmd:=HELLO message received.")
	case "PONG":
		acLogDebug("[AC] <<< cmd:=PONG message received.")
	case "BYE":
		acLogDebug("[AC] <<< cmd:=BYE message received.")
	case "OK":
		acLogDebug("[AC] <<< cmd:=OK message received.")
	case "NOK":
		acLogDebug("[AC] <<< cmd:=NOK message received.")
	default:
		acLogDebug("[AC] <<< UNKNOWN message received.")
	}
	return true
}

func (c *Client) sendHelloMessage(publickey string, rmsg *Message, client_vpnIP *string) bool {
	var smsg Message

	smsg.Msg_type = "cmd:=HELLO\n"

	macaddr := make([]byte, 6)
	if GetMacAddress(macaddr) {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=%02X-%02X-%02X-%02X-%02X-%02X\n",
			macaddr[0], macaddr[1], macaddr[2], macaddr[3], macaddr[4], macaddr[0])
	} else {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=11-11-11-22-22-22\n")  //TBD
	}
	smsg.VpnIP = "vpnip:=0.0.0.0\n"
	smsg.VpnNetmask = "vpnnetmask:=0.0.0.0\n"
	smsg.Public_key = "publickey:=" + publickey + "\n"

	ipbytes := make([]byte, 4)
	if !GetLocalIpAddress(&smsg.EpIP, ipbytes) {
		acLogDebug("[AC] Failed to get local ip address.")
		smsg.EpIP = "0.0.0.0"
	}
	smsg.EpPort = "epport:=51820\n"

	smsg.Allowed_ips = fmt.Sprintf("allowedips:=10.1.0.0/16,%d.%d.0.0/16\n", ipbytes[0], ipbytes[1])

	if c.sendMsg(&smsg) {
		acLogDebug("[AC] >>> HELLO message sent.")
		if c.recvMsg(rmsg) {
			s := rmsg.VpnIP + "/32" //ex) 10.1.1.100/32
			t := []byte(s)
			*client_vpnIP = string(t[7:]) //7 => vpnip:=
			return true
		}
	}
	return false
}

func (c *Client) sendPingMessage(privatekey string, publickey string, rmsg *Message, client_vpnIP string) string {
	var smsg Message

	smsg.Msg_type = "cmd:=PING\n"

	macaddr := make([]byte, 6)
	if GetMacAddress(macaddr) {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=%02X-%02X-%02X-%02X-%02X-%02X\n",
			macaddr[0], macaddr[1], macaddr[2], macaddr[3], macaddr[4], macaddr[0])
	} else {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=11-11-11-22-22-22\n")  //TBD
	}

	smsg.VpnIP = rmsg.VpnIP + "\n"
	smsg.VpnNetmask = rmsg.VpnNetmask + "\n"
	smsg.Public_key = "publickey:=" + publickey + "\n"

	ipbytes := make([]byte, 4)
	if !GetLocalIpAddress(&smsg.EpIP, ipbytes) {
		acLogDebug("[AC] Failed to get local ip address.")
		smsg.EpIP = "0.0.0.0"
	}
	smsg.EpPort = "epport:=51820\n"

	smsg.Allowed_ips = fmt.Sprintf("allowedips:=10.1.0.0/16,%d.%d.0.0/16\n", ipbytes[0], ipbytes[1])

	if c.sendMsg(&smsg) {
		acLogDebug("[AC] >>> PING message sent.")
		if c.recvMsg(rmsg) {
			/*
				wg0.conf
				----------
				[Interface]
				PrivateKey =
				ListenPort =
				Address =

				[Peer]
				PublicKey =
				AllowedIPs =
				Endpoint =
			*/

			var unparsedConfig string

			data := "[Interface]\nPrivateKey = "
			unparsedConfig = data

			unparsedConfig += privatekey

			data = "\nListenPort = 51820\nAddress = " + client_vpnIP + "\n\n[Peer]\nPublicKey = "
			unparsedConfig += data

			t := []byte(rmsg.Public_key)
			publickeyData := string(t[11:]) //11 => publickey:=
			unparsedConfig += string(publickeyData)

			data = "\nAllowedIPs = "
			unparsedConfig += data

			t = []byte(rmsg.Allowed_ips)
			allowedData := string(t[12:]) //12 => allowedips:=
			unparsedConfig += string(allowedData)

			data = "\nEndpoint = "
			unparsedConfig += data

			e1 := []byte(rmsg.EpIP)
			e2 := []byte(rmsg.EpPort)
			endpoint := string(e1[6:]) + ":" + string(e2[8:]) + "\n" //6 => epip:=, 8 => epport:=
			unparsedConfig += endpoint

			acLogDebug("[AC] unparsedConfig ------> %s", unparsedConfig)
			return unparsedConfig
		}
	}
	return "unknown"
}

func (c *Client) sendByeMessage(publickey string, rmsg *Message) bool {
	var smsg Message

	smsg.Msg_type = "cmd:=BYE\n"

	macaddr := make([]byte, 6)
	if GetMacAddress(macaddr) {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=%02X-%02X-%02X-%02X-%02X-%02X\n",
			macaddr[0], macaddr[1], macaddr[2], macaddr[3], macaddr[4], macaddr[0])
	} else {
		smsg.Mac_addr = fmt.Sprintf("macaddr:=11-11-11-22-22-22\n")  //TBD
	}

	smsg.VpnIP = "0.0.0.0" + "\n"
	smsg.VpnNetmask = "0.0.0.0" + "\n"
	smsg.Public_key = "publickey:=" + publickey + "\n"

	ipbytes := make([]byte, 4)
	if !GetLocalIpAddress(&smsg.EpIP, ipbytes) {
		acLogDebug("[AC] Failed to get local ip address.")
		smsg.EpIP = "0.0.0.0"
	}
	smsg.EpPort = "epport:=51820\n"

	smsg.Allowed_ips = fmt.Sprintf("allowedips:=10.1.0.0/16,%d.%d.0.0/16\n", ipbytes[0], ipbytes[1])

	if c.sendMsg(&smsg) {
		acLogDebug("[AC] >>> BYE message sent.")
		if c.recvMsg(rmsg) {
			return true
		}
	}
	return false
}

func NewClient(address string) *Client {
	client := new(Client)
	return client
}
