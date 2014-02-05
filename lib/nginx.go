/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : nginx.go

* Purpose :

* Creation Date : 01-10-2014

* Last Modified : Wed 29 Jan 2014 07:25:05 PM CST

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package gourl

import (
	// 	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	re1 = regexp.MustCompile(`(\d+)`)
	re2 = regexp.MustCompile(`Reading:\s(\d+)\sWriting:\s(\d+)\sWaiting:\s(\d+)\s`)
	re3 = regexp.MustCompile(`(\d+) (\d+) (\d+)`)
	uri = "/server-status"
)

type NginxStatus struct {
	Host            string
	Active          int
	Reading         int
	Writing         int
	Waiting         int
	LastUpdate      int64
	TTL             int64
	SinceLastUpdate int64
	StatusUri       string
	Server          NginxServer
}

type NginxServer struct {
	Accepts  int
	Handled  int
	Requests int
}

func (s *NginxServer) Init() {
	s.Accepts = -1
	s.Handled = -1
	s.Requests = -1
}

func (stat *NginxStatus) Init() {
	stat.Active = -1
	stat.Reading = -1
	stat.Writing = -1
	stat.Waiting = -1
	stat.LastUpdate = -1
	stat.TTL = 60
	stat.SinceLastUpdate = -1
	if stat.StatusUri == "" {
		stat.StatusUri = uri
	}
	var s NginxServer
	s.Init()
	stat.Server = s
}

func (stat *NginxStatus) UpdateSinceLastUpdate() {
	stat.SinceLastUpdate = time.Now().Unix() - stat.LastUpdate
}

func (stat *NginxStatus) Update() {
	url := "http://" + stat.Host + stat.StatusUri
	var r Req
	r.Url = url
	fullbody := r.GetString()
	if fullbody == "" {
		return
	}
	body := strings.Split(fullbody, "\n")
	activeLine := body[0]

	serverLine := body[2]

	mistLine := body[len(body)-1]
	active := re1.FindStringSubmatch(activeLine)[0]
	stat.Active, _ = strconv.Atoi(active)
	mist := re2.FindStringSubmatch(mistLine)

	server := re3.FindStringSubmatch(serverLine)

	stat.Reading, _ = strconv.Atoi(mist[1])
	stat.Writing, _ = strconv.Atoi(mist[2])
	stat.Waiting, _ = strconv.Atoi(mist[3])
	stat.LastUpdate = time.Now().Unix()

	var sv NginxServer

	sv.Accepts, _ = strconv.Atoi(server[1])
	sv.Handled, _ = strconv.Atoi(server[2])
	sv.Requests, _ = strconv.Atoi(server[3])

	stat.Server = sv
}
