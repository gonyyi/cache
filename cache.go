// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Cache struct {
	log         Logger
	conn        httpClient
	mu          sync.Mutex // global mutex
	configFile  string
	DataFileExt string                `json:"data_file_ext"`
	DataDir     string                `json:"data_dir"`
	Items       map[string]*cacheItem `json:"items"`
}

func (c *Cache) init() *Cache {
	c.log = newDummyLogger()
	c.conn = NewAReq()
	c.DataFileExt = ".download"
	c.DataDir = "./tmp"
	c.Items = make(map[string]*cacheItem)
	return c
}

func (c *Cache) SetLogger(l Logger) error {
	if l != nil {
		c.log = l
		c.log.Debugf("logger updated")
		return nil
	}
	return errors.New(ERR_NIL_LOG)
}

func (c *Cache) IsCacheExists(name string) bool {
	if _, ok := c.Items[name]; ok {
		return true
	}
	return false
}
func (c *Cache) GetCacheData(name string) ([]byte, error) {
	if _, ok := c.Items[name]; ok {
		dataPath := path.Join(c.DataDir, name+c.DataFileExt)
		b, err := ioutil.ReadFile(dataPath)
		if err != nil {
			c.log.Errorf("Cache item <%s>: cannot retreat data file <%s>", name, dataPath)
			return nil, err
		}
		c.log.Debugf("GetCacheData(%s): retreated data file <%s>", name, dataPath)
		return b, nil
	}
	c.log.Errorf("GetCacheData(%s): Cache item does not exist", name)
	return nil, fmt.Errorf("Cache item <%s> does not exist", name)
}
func (c *Cache) GetCacheItem(name string) (*cacheItem, error) {
	if ci, ok := c.Items[name]; ok {
		c.log.Debugf("GetCacheItem(%s) -- found", name)
		return ci, nil
	}
	c.log.Errorf("Cache item <%s> does not exist", name)
	return nil, fmt.Errorf("Cache item <%s> does not exist", name)
}
func (c *Cache) AddCache(name, method, url, id, passwd string) error {
	// check if already exist
	if c.IsCacheExists(name) {
		c.log.Errorf("Cache item <%s> already exists", name)
		return fmt.Errorf("Cache item <%s> already exists", name)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items[name] = newCacheItem(method, url, id, passwd)
	c.log.Infof("added Cache item <%s> with method: <%s>, URL: <%s>, ID: <%s>, passwd: <%s>",
		name, method, url, id, strings.Repeat("*", len(passwd)))
	return nil
}

func (c *Cache) RemoveCache(name string, removeData bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.log.Infof("removing Cache <%s>", name)

	if _, ok := c.Items[name]; ok {
		if err := os.Remove(path.Join(c.DataDir, name+c.DataFileExt)); err != nil {
			c.log.Errorf("failed to remove Cache <%s>'s data file <%s> -- %s", name, name+c.DataFileExt, err.Error())
			return err
		}
		delete(c.Items, name)
		c.log.Infof("removed Cache <%s>", name)
		return nil
	}
	c.log.Errorf("failed to remove Cache <%s> -- Cache not exist", name)
	return fmt.Errorf("Cache item <%s> not exist", name)
}

func (c *Cache) Save() error {
	b, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(c.configFile, b, 0755); err != nil {
		c.log.Errorf("failed to save the config file <%s>", c.configFile)
		return err
	}
	c.log.Infof("saved the config file <%s>", c.configFile)
	return nil
}

func CreateNewConfig(filename string, addExampleCache bool) error {
	c := New()
	c.configFile = filename
	if addExampleCache {
		c.AddCache("test", "GET", "https://gonyyi.com/copyright.txt", "", "")
	}
	if err := c.Save(); err != nil {
		return err
	}
	return nil
}

func New() *Cache {
	r := Cache{}
	return r.init()
}

func (c *Cache) Open(configPath string) error {

	c.configFile = configPath

	c.log.Infof("initiating Cache <%s>", configPath)

	if b, err := ioutil.ReadFile(configPath); err != nil {
		c.log.Errorf("cannot open the config file <%s>: %s", configPath, err.Error())
		return err
	} else if err = json.Unmarshal(b, &c); err != nil {
		c.log.Errorf("cannot unmarshal the config file <%s>: %s", configPath, err.Error())
		return err
	}
	c.log.Infof("loaded the config file <%s>", configPath)

	// Create a data directory (DataDir) if not exist
	{
		if c.DataDir == "" {
			c.log.Warnf("config does not hav data_dir, use default <./tmp> instead")
			c.DataDir = "./tmp"
		}

		if _, err := os.Stat(c.DataDir); err != nil {
			if os.IsNotExist(err) {
				c.log.Warnf("data_dir <%s> not exist, creating the directory", c.DataDir)
				if os.MkdirAll(c.DataDir, 0755) != nil {
					c.log.Fatalf("failed to create data_dir <%s>", c.DataDir)
				}
			} else {
				c.log.Fatalf("cannot open data_dir <%s>", c.DataDir)
			}
		}
	}

	c.log.Infof("total Cache items: %d", len(c.Items))
	count := 0
	for name, _ := range c.Items {
		count += 1
		c.log.Infof("Cache item [%d]: <%s>", count, name)
	}
	return nil
}
