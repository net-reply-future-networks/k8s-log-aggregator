package mocks

import "os"

type OsMock struct {
	MockExecPs   *MockExecPs
	MockOpenFile MockOpenFile
}

type MockExecPs struct {
	OutBytes            []byte
	OutErr              error
	ExpectedInvocations int
	ActualInvocations   int
}

type MockOpenFile struct {
	OutFile *os.File
	OutErr  error
}

func (o *OsMock) ExecPs() ([]byte, error) {
	o.MockExecPs.ActualInvocations++
	return o.MockExecPs.OutBytes, o.MockExecPs.OutErr
}

func (o *OsMock) OpenFile(name string) (*os.File, error) {
	return o.MockOpenFile.OutFile, o.MockOpenFile.OutErr
}
