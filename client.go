package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sgrimee/kizcool"
)

type Client struct {
	kiz           *kizcool.Kiz
	connected     bool
	workerContext context.Context
	workerClose   func()
	workedExited  chan struct{}
}

const tahomaUrlBase = "https://ha101-1.overkiz.com/enduser-mobile-web"
const refreshInterval = time.Minute
const pollInterval = time.Second * 2
const devicesRefreshInterval = time.Minute * 5

func MakeClient(user string, pass string) (*Client, error) {
	kiz, err := kizcool.New(user, pass, tahomaUrlBase, "")
	if err != nil {
		return nil, err
	}

	ctx, close := context.WithCancel(context.Background())

	client := &Client{
		kiz:           kiz,
		connected:     false,
		workerContext: ctx,
		workerClose:   close,
		workedExited:  make(chan struct{}),
	}

	go client.worker()

	return client, nil
}

func (client *Client) Terminate() {
	client.workerClose()
	<-client.workedExited
	client.kiz = nil
}

func (client *Client) setConnected(value bool) {
	if client.connected == value {
		return
	}

	client.connected = value
	fmt.Printf("CONNECTED = %t\n", client.connected)
}

func (client *Client) worker() {
	defer close(client.workedExited)

	refreshTimer := time.NewTicker(refreshInterval)
	pollTimer := time.NewTicker(pollInterval)
	devicesRefreshTimer := time.NewTicker(devicesRefreshInterval)

	client.devicesRefresh()
	client.refresh()

	for {
		select {
		case <-devicesRefreshTimer.C:
			client.devicesRefresh()

		case <-refreshTimer.C:
			client.refresh()

		case <-pollTimer.C:
			client.poll()

		case <-client.workerContext.Done():
			return
		}
	}
}

func (client *Client) afterReq(err error) {
	// consider after an error we are disconnected and after a success we are connected
	client.setConnected(err == nil)
}

func (client *Client) refresh() {
	err := client.kiz.RefreshStates()
	client.afterReq(err)

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
	}
}

func (client *Client) poll() {
	events, err := client.kiz.PollEvents()
	client.afterReq(err)

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
		return
	}

	for _, event := range events {
		fmt.Printf("  %+v\n", event)
	}
}

func (client *Client) devicesRefresh() {
	devices, err := client.kiz.GetDevices()
	client.afterReq(err)

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
		return
	}

	fmt.Printf("DEVICES\n")
	for _, device := range devices {
		fmt.Printf("  %s\n", device.Label)
	}

}
