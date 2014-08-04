package linode

import (
	"os"
	"testing"

	. "github.com/motain/gocheck"
	"github.com/pearkes/linode/testutil"
)

type S struct {
	client *Client
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("foobar")
	s.client.URL = "http://localhost:4444"
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func makeClient(t *testing.T) *Client {
	client, err := NewClient("foobarkey")

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if client.Key != "foobarkey" {
		t.Fatalf("key not set on client: %s", client.Key)
	}

	return client
}

func Test_NewClient_env(t *testing.T) {
	os.Setenv("LINODE_KEY", "bar")
	client, err := NewClient("")

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if client.Key != "bar" {
		t.Fatalf("key not set on client: %s", client.Key)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := makeClient(t)

	params := map[string]string{
		"foo": "bar",
		"baz": "bar",
	}

	req, err := c.NewRequest("GET", []map[string]string{params})

	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	if req.URL.String() != "https://api.linode.com?api_action=batch&api_key=foobarkey&api_requestArray=%5B%7B%22baz%22%3A%22bar%22%2C%22foo%22%3A%22bar%22%7D%5D" {
		t.Fatalf("bad base url: %v", req.URL.String())
	}

	if req.Method != "GET" {
		t.Fatalf("bad method: %v", req.Method)
	}
}
