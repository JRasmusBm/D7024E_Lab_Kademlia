package main

import (
	"errors"
	"io"
	"testing"
	"time"
)

func TestSetupConnectionSuccess(t *testing.T) {
	var client Client = &RealClient{}
	var reader Reader = &MockReader{ReadStringResult: "1.2.3.4"}
	var network Network = &MockNetwork{}
	expected, _ := network.Dial("tcp", "1.2.3.4")
	actual, actualErr := client.SetupConnection(&reader, &network)
	if actualErr != nil {
		t.Error(actualErr)
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestMakeConnectionReader(t *testing.T) {
	var client Client = &RealClient{}
	var network Network = &MockNetwork{}
	conn, _ := network.Dial("tcp", "1.2.3.4")
	actual := client.MakeConnectionReader(&conn)
	if actual == nil {
		t.Errorf("Expected %#v not to be nil", actual)
	}
}

func TestListenToServerDoesNotCrashSuccess(t *testing.T) {
	var client Client = &RealClient{}
	var reader Reader = &MockReader{ReadStringResult: "Hello"}
	go client.ListenToServer(&reader)
	time.Sleep(300 * time.Millisecond)
}

func TestListenToServerDoesNotCrashErr(t *testing.T) {
	var client Client = &RealClient{}
	var reader Reader = &MockReader{ReadStringErr: errors.New("Random Error")}
	go client.ListenToServer(&reader)
	time.Sleep(5 * time.Millisecond)
}

func TestGetMessageSuccess(t *testing.T) {
	var client Client = &RealClient{}
	expected := "Hello!"
	var reader Reader = &MockReader{ReadStringResult: "Hello!"}
	actual, actualErr := client.GetMessageFromUser(&reader)
	if actualErr != nil {
		t.Error(actualErr)
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestHandlePingSuccess(t *testing.T) {
	var client Client = &RealClient{}
	expected := "ping 123;"
	var reader FileReader = &MockFileReader{}
	actual, actualErr := client.HandleMessage("ping 123", &reader)
	if actualErr != nil {
		t.Error(actualErr)
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestHandlePutSuccess(t *testing.T) {
	var client Client = &RealClient{}
	expected := "put 123;"
	var reader FileReader = &MockFileReader{ReadFileResult: []byte("123\n")}
	actual, actualErr := client.HandleMessage("put file.txt", &reader)
	if actualErr != nil {
		t.Error(actualErr)
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestHandlePutFail(t *testing.T) {
	var client Client = &RealClient{}
	var reader FileReader = &MockFileReader{
		ReadFileErr: errors.New("My perfect Error"),
	}
	_, actualErr := client.HandleMessage("put file.txt", &reader)
	if actualErr == nil {
		t.Errorf("Expected HandleMessage to throw an error")
	}
}

func TestSendMessageDoesNotCrash(t *testing.T) {
	var client Client = &RealClient{}
	var writer io.Writer = &MockWriter{}
	go client.SendMessage(&writer, "hello;")
	time.Sleep(300 * time.Millisecond)
}

func TestDial(t *testing.T) {
	var network Network = &RealNetwork{}
	go network.Dial("tcp", "0.0.0.0")
	time.Sleep(300 * time.Millisecond)
}

func TestReadFile(t *testing.T) {
	var fileReader FileReader = &RealFileReader{}
	_, err := fileReader.ReadFile("fjdslfhdsalguhewoireryiuqrhewkjrew.fads")
	if err == nil {
		t.Errorf("Expected to return error when file does not exist")
	}
}

func TestConnectionValidError(t *testing.T) {
	var reader Reader = &MockReader{ReadStringErr: errors.New("Could not read")}
	var writer io.Writer = &MockWriter{}
	var client Client = &RealClient{}
	if client.ConnectionValid(&writer, &reader) {
		t.Errorf("Should not be valid at read string error")
	}
}

func TestConnectionValidWrongResponse(t *testing.T) {
	var reader Reader = &MockReader{ReadStringResult: "boooyaah;"}
	var writer io.Writer = &MockWriter{}
	var client Client = &RealClient{}
	if client.ConnectionValid(&writer, &reader) {
		t.Errorf("Should not be valid when string is not ok;")
	}
}

func TestConnectionValidSuccess(t *testing.T) {
	var reader Reader = &MockReader{ReadStringResult: "ok;"}
	var writer io.Writer = &MockWriter{}
	var client Client = &RealClient{}
	if !client.ConnectionValid(&writer, &reader) {
		t.Errorf("Should be valid when string is ok;")
	}
}
