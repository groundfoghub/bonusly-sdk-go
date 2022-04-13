package bonusly

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type errReadCloser struct {
	rerr error
	cerr error
}

func (r errReadCloser) Read(p []byte) (n int, err error) {
	return 0, r.rerr
}

func (r errReadCloser) Close() error {
	return r.cerr
}

func Test_readAndCloseBody(t *testing.T) {
	type args struct {
		r *http.Response
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"nil-response",
			args{r: nil},
			nil,
			true,
		},
		{
			"ok",
			args{r: &http.Response{Body: ioutil.NopCloser(strings.NewReader("Test"))}},
			[]byte("Test"),
			false,
		},
		{
			"error-read",
			args{r: &http.Response{Body: errReadCloser{rerr: errors.New("read error")}}},
			nil,
			true,
		},
		{
			"error-read-close",
			args{r: &http.Response{Body: errReadCloser{rerr: errors.New("read error"), cerr: errors.New("close error")}}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readAndCloseBody(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndCloseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAndCloseBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
