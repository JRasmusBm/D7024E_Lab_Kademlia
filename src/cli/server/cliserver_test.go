package cli

import (
	api_p "api"
	"errors"
	"testing"
	"time"
	networkutils "utils/network"
)

func TestServerInitDoesNotCrash(t *testing.T) {
	cliChannel := make(chan string)
	api := api_p.API{}
	var networkUtils networkutils.NetworkUtils = &networkutils.RealNetworkUtils{}
	go CliServerInit(api, &networkUtils, cliChannel)
	time.Sleep(100 * time.Millisecond)
}

func TestServerDoesNotCrash(t *testing.T) {
	cliChannel := make(chan string)
	var server Server = &MockServer{
		MessageParserResult: []string{""},
	}
	var network Network = &MockNetwork{}
	go CliServer(cliChannel, &network, &server)
	time.Sleep(100 * time.Millisecond)
}

func TestServerDoesNotCrashListenerError(t *testing.T) {
	cliChannel := make(chan string)
	var server Server = &MockServer{
		SetupListenerResult: nil,
		SetupListenerErr:    errors.New("Random Error"),
		MessageParserResult: []string{""},
	}
	var network Network = &MockNetwork{}
	go CliServer(cliChannel, &network, &server)
	time.Sleep(100 * time.Millisecond)
}

func TestServerDoesNotCrashListenForConnectionError(t *testing.T) {
	cliChannel := make(chan string)
	var server Server = &MockServer{
		ListenForConnectionResult: nil,
		ListenForConnectionErr:    errors.New("Random Error"),
		MessageParserResult:       []string{""},
	}
	var network Network = &MockNetwork{}
	go CliServer(cliChannel, &network, &server)
	time.Sleep(100 * time.Millisecond)
}

func TestServerDoesNotCrashCloseStatement(t *testing.T) {
	cliChannel := make(chan string)
	var server Server = &MockServer{
		MessageParserResult: []string{"close"},
	}
	var network Network = &MockNetwork{}
	go CliServer(cliChannel, &network, &server)
	time.Sleep(100 * time.Millisecond)
}
