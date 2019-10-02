package main

type MockFileReader struct {
	ReadFileResult []byte
	ReadFileErr    error
}

func (r *MockFileReader) ReadFile(filename string) ([]byte, error) {
	if r.ReadFileErr != nil {
		return nil, r.ReadFileErr
	}
	return r.ReadFileResult, nil
}
