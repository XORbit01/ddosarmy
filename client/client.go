package client

import "net/http"

type Client struct {
	http.Client
	DispatcherServer string // Dispatcher server http end point
}
