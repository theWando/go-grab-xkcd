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
			hc := &xKCDClient{
				client:  tt.fields.client,
				baseURL: tt.fields.baseURL,
			}
			if got := hc.buildURL(tt.args.n); got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
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

func mockServer(body string) *httptest.Server {

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func Test_xKCDClient_Fetch(t *testing.T) {

	server := mockServer(jsonResponse)
	defer server.Close()

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
		fields  fields
		args    args
		want    model.Comic
		wantErr bool
	}{
		{
			name: "Should return object",
			fields: fields{
				client: &http.Client{
					Timeout: 10 * time.Second,
				},
				baseURL: server.URL,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := &xKCDClient{
				client:  tt.fields.client,
				baseURL: tt.fields.baseURL,
			}
			got, err := hc.Fetch(tt.args.n, tt.args.save)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
