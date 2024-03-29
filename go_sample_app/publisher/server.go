/*
Copyright (c) 2021 Intel Corporation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	eiicfgmgr "ConfigMgr/eiiconfigmgr"
	eiimsgbus "EIIMessageBus/eiimsgbus"
	"fmt"

	"github.com/golang/glog"
)

func start_server() {

	configmgr, err := eiicfgmgr.ConfigManager()
	if err != nil {
		glog.Fatal("Config Manager initialization failed...")
	}
	defer configmgr.Destroy()

	numOfServers, _ := configmgr.GetNumServers()
	if numOfServers == -1 {
		glog.Errorf("No server instances found, exiting...")
		return
	}
	// serverctx, err := configmgr.GetServerByName("echo_service")
	serverctx, err := configmgr.GetServerByIndex(0)
	if err != nil {
		glog.Fatal("GetServerByIndex is failed")
	}
	defer serverctx.Destroy()

	interfaceVal, err := serverctx.GetInterfaceValue("Name")
	if err != nil {
		fmt.Printf("Error to GetInterfaceValue of 'Name': %v\n", err)
		return
	}

	serviceName, err := interfaceVal.GetString()
	if err != nil {
		fmt.Printf("Error to GetString value of 'Name'%v\n", err)
		return
	}

	config, err := serverctx.GetMsgbusConfig()
	if err != nil {
		glog.Fatal("Error occured with error:%v", err)
	}

	sClient, err := eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		fmt.Printf("-- Error initializing message bus context: %v\n", err)
		return
	}
	defer sClient.Close()

	fmt.Printf("-- Initializing service %s\n", serviceName)
	service, err := sClient.NewService(serviceName)
	if err != nil {
		fmt.Printf("-- Error initializing service: %v\n", err)
	}
	defer service.Close()

	fmt.Println("-- Running server...")

	for {
		req, err := service.ReceiveRequest(-1)
		if err != nil {
			fmt.Printf("-- Error receiving request: %v\n", err)
		}
		fmt.Printf("--Server received request: %v\n", req)
		service.Response(map[string]interface{}{"Response": "OK"})
	}
}
