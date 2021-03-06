// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: login.proto

package user

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for LoginService service

type LoginService interface {
	// login
	Login(ctx context.Context, in *LoginModel, opts ...client.CallOption) (*LoginModel, error)
}

type loginService struct {
	c    client.Client
	name string
}

func NewLoginService(name string, c client.Client) LoginService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "user"
	}
	return &loginService{
		c:    c,
		name: name,
	}
}

func (c *loginService) Login(ctx context.Context, in *LoginModel, opts ...client.CallOption) (*LoginModel, error) {
	req := c.c.NewRequest(c.name, "LoginService.Login", in)
	out := new(LoginModel)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LoginService service

type LoginServiceHandler interface {
	// login
	Login(context.Context, *LoginModel, *LoginModel) error
}

func RegisterLoginServiceHandler(s server.Server, hdlr LoginServiceHandler, opts ...server.HandlerOption) error {
	type loginService interface {
		Login(ctx context.Context, in *LoginModel, out *LoginModel) error
	}
	type LoginService struct {
		loginService
	}
	h := &loginServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&LoginService{h}, opts...))
}

type loginServiceHandler struct {
	LoginServiceHandler
}

func (h *loginServiceHandler) Login(ctx context.Context, in *LoginModel, out *LoginModel) error {
	return h.LoginServiceHandler.Login(ctx, in, out)
}
