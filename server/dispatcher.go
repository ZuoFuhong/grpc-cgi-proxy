package server

import (
	"encoding/json"
	"fmt"
	"github.com/ZuoFuhong/grpc-cgi-proxy/errcode"
	"github.com/ZuoFuhong/grpc-cgi-proxy/pkg/config"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var digitRegex = regexp.MustCompile("^\\d{1,10}$")

type Dispatcher struct {
	hmap map[string]*Handler
}

func NewDispatcher() *Dispatcher {
	cfg := config.GlobalConfig()
	hmap := make(map[string]*Handler)
	for _, service := range cfg.Services {
		for _, method := range service.Methods {
			key := fmt.Sprintf("%s-%s", method.CgiPath, method.Method)
			hmap[key] = NewHandler(service.ServiceName, service.Namespace, method.Cmd, method.Timeout)
		}
	}
	return &Dispatcher{
		hmap: hmap,
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprintf("%s-%s", r.URL.Path, r.Method)
	handler := d.hmap[key]
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 提取请求参数
	reqParams, err := getRequestParams(r)
	if err != nil {
		Error(w, errcode.BadRequestParam, err.Error())
		return
	}
	// 透传代理
	rsp, err := handler.handle(r.Context(), reqParams)
	if err != nil {
		Error(w, errcode.InternalServerError, err.Error())
		return
	}
	Ok(w, rsp)
}

func getRequestParams(r *http.Request) (map[string]interface{}, error) {
	reqParams := make(map[string]interface{})
	switch r.Method {
	case "GET", "DELETE":
		for k, v := range r.URL.Query() {
			if len(v) == 0 || v[0] == "" {
				continue
			}
			reqParams[k] = v[0]
			// 判断数字类型
			if digitRegex.MatchString(v[0]) {
				dv, _ := strconv.Atoi(v[0])
				reqParams[k] = dv
			}
		}
	case "POST", "PUT":
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(bodyBytes, &reqParams); err != nil {
				return nil, err
			}
		}
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				return nil, err
			}
			for k, v := range r.Form {
				if len(v) == 0 || v[0] == "" {
					continue
				}
				reqParams[k] = v[0]
				// 判断数字类型
				if digitRegex.MatchString(v[0]) {
					dv, _ := strconv.Atoi(v[0])
					reqParams[k] = dv
				}
			}
		}
	}
	return reqParams, nil
}
