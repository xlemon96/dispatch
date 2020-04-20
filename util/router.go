package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/navieboy/dispatch/constant"
	"github.com/navieboy/dispatch/model/api"
)

var (
	ErrTemp        = "{\"code\":%d, \"message\":\"%s\"}"
	ErrNoAction    = "{\"code\":11, \"message\":\"no such action\"}"
	ErrMarshalJson = "{\"code\":10, \"message\":\"marshal json err\"}"
	ErrReqParam    = "{\"code\":9, \"message\":\"request param invalid\"}"
	ErrDb          = "{\"code\":8, \"message\":\"db err\"}"
	ErrSystem      = "{\"code\":7, \"message\":\"system error\"}"
)

type CommonRequest struct {
	Action      string
	Content     []byte
	RequestId   string
	UserId      string
	UserPin     string
	ClientToken string
	Admin       string
	all         bool //设置为true代表查询所有数据，只有内部调用时允许设置
}

func (cr *CommonRequest) SetAll() {
	cr.all = true
}
func (cr *CommonRequest) IsAll() bool {
	return cr.all
}

type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

type Request struct {
	CommonRequest
	RawRequest      *http.Request
	BusinessRequest interface{} // 业务请求的request实例，可以为空，在做filter时需要注意判断nil
	ResponseWriter  http.ResponseWriter
}

type Response struct {
	Response *CommonResponse
}

// 包装器
type middleware struct {
	filter Filter
	next   *middleware
}

// 适配器。
func (m middleware) Adapter(request *Request, response *Response) {
	m.filter(request, response, m.next.Adapter)
}

// next类型
type NextFunc func(request *Request, response *Response)

// 调用链
type Filter func(request *Request, response *Response, next NextFunc)

//
type HandleFunc func(request *Request, response *Response) (err error)

func (h HandleFunc) Handle(request *Request, response *Response) (err error) {
	return h(request, response)
}

type IHandler interface {
	Handle(request *Request, response *Response) (err error)
}

type HandleContext struct {
	Action  string
	handler IHandler
	new     func(request *Request) (interface{}, error)
}

func (hc *HandleContext) prepare() bool {
	if hc.new == nil {
		return false
	}
	return true
}

func NewHandleContext(handler IHandler, new func(request *Request) (interface{}, error)) *HandleContext {
	return &HandleContext{handler: handler, new: new}
}

type Router struct {
	filters []Filter
	Handler map[string]*HandleContext
	// 拦截器
	middleware middleware
	Logger     *log.Logger
	hostname   string
}

func (t *Router) RegisterFilters(filters ...Filter) {
	if filters == nil {
		return
	}
	for _, f := range filters {
		t.Use(f)
	}
}

func (t *Router) Use(filter Filter) {
	t.filters = append(t.filters, filter)
	t.middleware = build(t.filters)
}

func NewRouter() *Router {
	hostname, _ := os.Hostname()
	return &Router{
		filters:  make([]Filter, 0),
		Handler:  make(map[string]*HandleContext),
		hostname: hostname,
	}
}

func build(handlers []Filter) middleware {
	var next middleware

	if len(handlers) == 0 {
		return createVoidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = createVoidMiddleware()
	}
	return middleware{handlers[0], &next}
}

func createVoidMiddleware() middleware {
	return middleware{
		filter: func(request *Request, response *Response, next NextFunc) {},
		next:   &middleware{},
	}
}

func (r *Router) RegisterHandler(action string, handler IHandler) *Router {
	r.Handler[action] = NewHandleContext(handler, nil)
	return r

}

func (r *Router) RegisterHandleFunc(action string, handleFunc HandleFunc) *Router {
	r.Handler[action] = NewHandleContext(handleFunc, nil)
	return r
}
func (r *Router) RegisterHandleContext(action string, handleContext *HandleContext) *Router {
	r.Handler[action] = handleContext
	return r
}

func (p *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Printf("[method: %v, url: %v, remote_addr:%v, panic: %v, x-forwarded-for: %v]，statck:%s", r.Method,
				r.URL.RequestURI(), r.RemoteAddr, err, r.Header.Get("X-Forwarded-For"), string(debug.Stack()))
			w.Write([]byte(ErrSystem))
		}
	}()

	//to caculate elapsed time
	start := time.Now()

	// add client_token
	var clientToken string
	clientTokens, _ := r.Header["Client-Token"]
	if len(clientTokens) > 0 && clientTokens[0] != "" {
		clientToken = clientTokens[0]
	}

	// add user id
	var userId string
	userIds, _ := r.Header["User-Id"]
	if len(userIds) > 0 && userIds[0] != "" {
		userId = userIds[0]
	}
	// add user pin
	var userPin string
	userPins, _ := r.Header["User-Pin"]
	if len(userPins) > 0 && userPins[0] != "" {
		userPin = userPins[0]
	}

	// admin user
	var admin string
	admins, _ := r.Header["Admin-Token"]
	if len(admins) > 0 && admins[0] != "" {
		admin = admins[0]
	}

	r.Header.Set("Content-Type", "application/json")

	ctx := r.Context()
	r = r.WithContext(ctx)

	//help to add default header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Requst-Trace", p.hostname)

	r.ParseForm()
	action := r.FormValue("Action")
	if action == "" {
		p.Logger.Printf("[method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"))
		w.Write([]byte(ErrReqParam))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		p.Logger.Printf("[method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"))
		w.Write([]byte(ErrReqParam))
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	request := &Request{
		RawRequest: r,
		CommonRequest: CommonRequest{
			Action:      action,
			Content:     body,
			UserId:      userId,
			UserPin:     userPin,
			Admin:       admin,
			ClientToken: clientToken,
		},
		ResponseWriter: w,
	}
	response := new(Response)
	p.middleware.filter(request, response, p.middleware.next.Adapter)

	if response.Response == nil {
		response.Response = new(CommonResponse)
	}
	jsonRsp, _ := json.Marshal(response.Response)
	w.Write(jsonRsp)
	elapsed := time.Since(start)

	p.Logger.Printf("[method: %v, url: %v, body: %v, remote_addr:%v, action: %v, x-forwarded-for: %v, elapsed: %v]", r.Method,
		r.URL.RequestURI(), strings.Replace(string(body), "\n", " ", -1), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"), elapsed)
}

func (r *Router) InvokerFilter(request *Request, response *Response, next NextFunc) {
	response.Response = &CommonResponse{}
	handler, ok := r.Handler[request.Action]
	if !ok {
		defaultHandler, ok := r.Handler["*"]
		if !ok {
			r.Logger.Printf("[method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", request.RawRequest.Method,
				request.RawRequest.URL.RequestURI(), request.RawRequest.RemoteAddr, request.Action, request.RawRequest.Header.Get("X-Forwarded-For"))
			response.Response.Code = -1
			response.Response.Data = ErrNoAction
			return
		}
		handler = defaultHandler
	}

	defer func() {
		if err := recover(); err != nil {
			r.Logger.Printf("[InvokerFilter] panic. err: %v, stack: %s ", err, string(debug.Stack()))
			response.Response.Code = -1
			response.Response.Message = "system.error"
			response.Response.Detail = "process panic"
		}
	}()

	err := handler.handler.Handle(request, response)
	if err != nil {
		response.Response.Code = -1
		response.Response.Message = err.Error()
		response.Response.Detail = err.Error()
		return
	}
	next(request, response)
}

func (r *Router) PrepareFilter(request *Request, response *Response, next NextFunc) {
	response.Response = &CommonResponse{}
	param := newParams(request.Action)
	r.Logger.Printf("prepair: action: %s, content: %s", request.Action, request.Content)
	if err := json.Unmarshal(request.Content, param); err != nil {
		r.Logger.Printf("unmarshal failed, content: %s, params: %+v, err: %s", request.Content, param, err.Error())
		return
	}
	request.BusinessRequest = param
	next(request, response)
}

func newParams(action string) interface{} {
	switch action {
	case constant.UpdateDAGInstance:
		return new(api.Test)
	default:
		return nil
	}
}