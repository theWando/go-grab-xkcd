package client

import (
	"fmt"
	"github.com/theWando/go-grab-xkcd/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func ExampleXKCDClient_BuildURL() {

	hc := NewXKCDClient()
	fmt.Println(hc.BuildURL(0))
	fmt.Println(hc.BuildURL(10))
	// Output:
	//https://xkcd.com/info.0.json
	//https://xkcd.com/10/info.0.json

}

func Test_xKCDClient_buildURL(t *testing.T) {
	type fields struct {
		client  *http.Client
		baseURL string
	}
	type args struct {
		n ComicNumber
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should be latest comic",
			fields: fields{
				baseURL: "http://mymocked.url.com",
			},
			args: args{0},
			want: "http://mymocked.url.com/info.0.json",
		},
		{
			name: "should be comic 9",
			fields: fields{
				baseURL: "http://mymocked.url.com",
			},
			args: args{9},
			want: "http://mymocked.url.com/9/info.0.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := &XKCDClient{
				client:  tt.fields.client,
				baseURL: tt.fields.baseURL,
			}
			if got := hc.BuildURL(tt.args.n); got != tt.want {
				t.Errorf("BuildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

const jsonResponse = `{
  "month": "9",
  "num": 1732,
  "link": "",
  "year": "2016",
  "news": "",
  "safe_title": "Earth Temperature Timeline",
  "transcript": "",
  "alt": "[Afters setting your car on fire] Listen, your car's temperature has changed before.",
  "img": "https://imgs.xkcd.com/comics/earth_temperature_timeline.png",
  "title": "Earth Temperature Timeline",
  "day": "12"
}`
const badJson = `{
  "month": "9,
  "day": "12"
}`

func mockServer(body string, statusCode int, headers map[string]string) *httptest.Server {

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		fmt.Fprintln(w, body)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func Test_xKCDClient_Fetch(t *testing.T) {

	defaultClient := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	type fields struct {
		client  *http.Client
		baseURL string
	}
	type args struct {
		n    ComicNumber
		save bool
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		fields  fields
		args    args
		want    model.Comic
		wantErr bool
	}{
		{
			name:   "Should return object",
			server: mockServer(jsonResponse, 200, map[string]string{"Content-Type": "application/json"}),
			fields: fields{
				client: defaultClient,
			},
			args: args{
				n:    9,
				save: false,
			},
			want: model.Comic{
				Title:       "Earth Temperature Timeline",
				Number:      1732,
				Date:        "12-9-2016",
				Description: "[Afters setting your car on fire] Listen, your car's temperature has changed before.",
				Image:       "https://imgs.xkcd.com/comics/earth_temperature_timeline.png",
			},
			wantErr: false,
		},
		{
			name:   "Should return an error when the json is invalid",
			server: mockServer(badJson, 200, map[string]string{"Content-Type": "application/json"}),
			fields: fields{
				client: defaultClient,
			},
			args: args{
				n:    0,
				save: false,
			},
			want:    model.Comic{},
			wantErr: true,
		},
		{
			name:   "Should return an error when the body is empty",
			server: mockServer("", 500, map[string]string{"Content-Type": "application/json"}),
			fields: fields{
				client: defaultClient,
			},
			args: args{
				n:    0,
				save: false,
			},
			want:    model.Comic{},
			wantErr: true,
		},
		{
			name: "Should return an error when server is unreachable",
			server: httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {

			})),
			fields: fields{
				client: defaultClient,
			},
			args: args{
				n:    0,
				save: false,
			},
			want:    model.Comic{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := &XKCDClient{
				client:  tt.fields.client,
				baseURL: tt.server.URL,
			}
			got, err := hc.Fetch(tt.args.n, tt.args.save)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
			tt.server.Close()
		})
	}
}
