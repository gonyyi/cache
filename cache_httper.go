// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package cache

import (
	"bytes"
	"fmt"
	"github.com/orangenumber/areq"
	"sync"
)

type httpClient interface {
	Req(METHOD, URL, ID, PWD string, BODY []byte, BODY_EXT string) (resp []byte, err error)
}

// =====================================================================================================================
// httpClient -- AReq adopter
// =====================================================================================================================
type cacheAReq struct {
	req *areq.AReq
	mu  sync.Mutex
}

func NewAReq() *cacheAReq {
	return &cacheAReq{
		req: areq.New(),
	}
}

func (c *cacheAReq) Req(Method, URL, ID, PWD string, BODY []byte, BODY_EXT string) (resp []byte, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	statusCode, err := c.req.Request(Method, URL,
		areq.Plugin.BasicAuth(ID, PWD),
		areq.Plugin.SetBody(bytes.NewReader(BODY), BODY_EXT),
		areq.Plugin.AcceptEncoding("gzip"))
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		err = fmt.Errorf(ERRF_BAD_STATUS_CODE_X, statusCode)
	}
	return c.req.Buf.Bytes(), err
}
