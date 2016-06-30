package decorator

import (
	"net/http"
	"testing"

	"github.com/lingualeo/kapusta"
	. "gopkg.in/check.v1"
	"golang.org/x/net/context"
)

type TestSuite struct {
	dummyClient *dummy
}

type dummy struct{}

func (d *dummy) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}

type callOrder []string

var (
	_ = Suite(&TestSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.dummyClient = &dummy{}
}

func (s *TestSuite) send(ctx context.Context, r *http.Request, decorator kapusta.DecoratorFunc) (*http.Response, error) {
	return kapusta.Decorate(s.dummyClient, decorator).Do(ctx, r)
}

func createDecorator(name string, order *callOrder) kapusta.DecoratorFunc {
	return func(c kapusta.Client) kapusta.Client {
		return kapusta.ClientFunc(func(ctx context.Context, r *http.Request) (*http.Response, error) {
			*order = append(*order, name)
			return c.Do(ctx, r)
		})
	}
}

func (s *TestSuite) TestDecorate(c *C) {
	var order callOrder
	client := kapusta.Decorate(s.dummyClient, createDecorator("A", &order), createDecorator("B", &order), createDecorator("C", &order))

	client.Do(context.Background(), &http.Request{})

	c.Assert(callOrder{"C", "B", "A"}, DeepEquals, order)
}

func (s *TestSuite) TestHeaderDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(context.Background(), r, HeaderDecorator("X-Foo", "123"))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
}

func (s *TestSuite) TestHeadersDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(context.Background(), r, HeadersDecorator(map[string]string{"X-Foo": "123", "X-Bar": "456"}))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
	c.Assert(res.Request.Header.Get("X-Bar"), Equals, "456")
}

func (s *TestSuite) TestPanicDecorator(c *C) {
	panicTriggerDecorator := func(c kapusta.Client) kapusta.Client {
		return kapusta.ClientFunc(func(ctx context.Context, r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	client := kapusta.Decorate(s.dummyClient, panicTriggerDecorator, RecoverDecorator())
	res, err := client.Do(context.Background(), r)

	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "*oops")
}
