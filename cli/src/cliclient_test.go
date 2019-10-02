package main

import (
	"errors"
	"testing"
	"time"
)

func TestMainDoesNotCrash(t *testing.T) {
	go main()
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrash(t *testing.T) {
	var ioReader Reader = &MockReader{ReadStringResult: "Hello"}
	var client Client = &MockClient{
		ConnectionValidResult:      true,
		SetupConnectionResult:      nil,
		MakeConnectionReaderResult: &ioReader,
	}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrashSetupConnectionErr(t *testing.T) {
	var client Client = &MockClient{
		SetupConnectionErr:    errors.New("Random Error"),
		ConnectionValidResult: true,
	}
	var ioReader Reader = &MockReader{}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrashFailedConnect(t *testing.T) {
	var client Client = &MockClient{ConnectionValidResult: false}
	var ioReader Reader = &MockReader{}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrashGetMessageErr(t *testing.T) {
	var client Client = &MockClient{
		ConnectionValidResult: true,
		GetMessageFromUserErr: errors.New("Random Error"),
	}
	var ioReader Reader = &MockReader{}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrashHandleMessageErr(t *testing.T) {
	var client Client = &MockClient{
		ConnectionValidResult: true,
		HandleMessageErr:      errors.New("Random Error"),
	}
	var ioReader Reader = &MockReader{}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}

func TestCLIDoesNotCrashClose(t *testing.T) {
	var client Client = &MockClient{
		HandleMessageResult:   "close;",
		ConnectionValidResult: true,
	}
	var ioReader Reader = &MockReader{}
	var network Network = &MockNetwork{}
	var fileReader FileReader = &RealFileReader{}
	go cliClient(&client, &ioReader, &network, &fileReader)
	time.Sleep(300 * time.Millisecond)
}
