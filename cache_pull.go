package cache

import (
	"io/ioutil"
	"path"
	"strings"
	"time"
)

func (c *cache) saveData(name string, data []byte) error {
	savePath := path.Join(c.DataDir, name+DATA_FILE_EXT)
	if err := ioutil.WriteFile(savePath, data, 0755); err != nil {
		c.log.Errorf("cache <%s>: cannot save the data file to <%s> -- %s", name, savePath, err.Error())
		return err
	}
	c.log.Infof("cache <%s>: saved to <%s>", name, savePath)
	return nil
}

func (c *cache) CachePullAll() []error {
	var errors []error

	for k, _ := range c.Items {
		if err := c.CachePull(k); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (c *cache) CachePull(name string) error {
	if ci, err := c.GetCacheItem(name); err != nil {
		c.log.Errorf("failed to pull <%s> -- %s", name, err.Error())
		return err
	} else {
		url := ci.Request.URL
		for k, v := range ci.Request.SubVar {
			url = strings.ReplaceAll(url, "{subvar:"+k+"}", v)
		}
		c.log.Debugf("CachePull(%s): new updated URL=<%s>", name, url)
		resp, err := c.conn.Req(ci.Request.Method, url, ci.Request.ID, ci.Request.Passwd, []byte(ci.Request.Body), ci.Request.ContentType)
		if err != nil {
			c.log.Errorf("cache <%s>: cannot get pull data from URL <%s> -- %s", name, url, err.Error())
			return err
		}
		if err:= c.saveData(name, resp); err != nil {
			return err
		}

		ci.Response.LastUpdated = c.now()
		ci.Response.Size = len(resp)

		return nil
	}
}

func(c *cache) now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
