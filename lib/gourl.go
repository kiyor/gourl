/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : gourl.go

* Purpose :

* Creation Date : 01-02-2014

* Last Modified : Fri 21 Mar 2014 06:49:28 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package gourl

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Req struct {
	http.Request
	Url      string
	MyHeader []*MyHeader
}

type Resp struct {
	http.Response
}

type MyHeader struct {
	Key   string
	Value string
}

var (
	timeout   = 3 * time.Second
	transport = http.Transport{
		Dial:               dialTimeout,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client = http.Client{
		Transport: &transport,
	}
)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		// 		os.Exit(1)
	}
}

func (r *Req) getResp(method string) (*http.Response, error) {
	req, _ := http.NewRequest(method, r.Url, nil)

	for _, v := range r.MyHeader {
		req.Header.Add(v.Key, v.Value)
	}

	if r.Host != "" {
		req.Host = r.Host
	}

	resp, err := client.Do(req)
	return resp, err
}

func (r *Req) GetFull() (Resp, error) {
	resp, err := r.getResp("GET")
	if err != nil {
		return Resp{}, err
	}
	return Resp{*resp}, err
}

func (r *Resp) String() string {
	if r.StatusCode != 200 {
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		body = []byte("")
	}
	defer r.Body.Close()
	var s string
	if len(body) == 0 {
		s = ""
	} else {
		last := string(body)[len(string(body))-1 : len(string(body))]
		if last == "\n" || last == " " {
			s = string(body)[:len(string(body))-1]
		} else {
			s = string(body)
		}
	}
	return s
}

func (r *Resp) StringSlice() []string {
	if r.StatusCode != 200 {
		return []string{""}
	}
	s := r.String()
	slice := strings.Split(s, "\n")
	return slice[0 : len(slice)-1]
}

func (r *Req) GetString() string {
	resp, _ := r.GetFull()
	return resp.String()
}

func (r *Req) GetStringSlice() []string {
	s := r.GetString()
	slice := strings.Split(s, "\n")
	return slice[0 : len(slice)-1]
}

func (r *Req) GetHeader() (http.Header, error) {
	// 	resp, err := http.Head(r.Url)
	// 	checkErr(err)
	// 	return resp.Header
	resp, err := r.getResp("HEAD")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Header, err
}
