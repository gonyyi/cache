package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

const DATA_FILE_EXT = ".download"

var DEFAULT_LOG_OUTPUT = os.Stderr

type cache struct {
	log        cacheLogger
	conn       httpClient
	mu         sync.Mutex // global mutex
	configFile string
	DataDir    string                `json:"data_dir"`
	Items      map[string]*cacheItem `json:"items"`
}

func (c *cache) init() *cache {
	c.log = NewCacheLog(os.Stderr)
	c.conn = NewAReq()
	c.DataDir = "./tmp"
	c.Items = make(map[string]*cacheItem)
	return c
}

func (c *cache) IsCacheExists(name string) bool {
	if _, ok := c.Items[name]; ok {
		return true
	}
	return false
}
func (c *cache) GetCacheData(name string) ([]byte, error) {
	if _, ok := c.Items[name]; ok {
		dataPath := path.Join(c.DataDir, name+DATA_FILE_EXT)
		b, err := ioutil.ReadFile(dataPath)
		if err != nil {
			c.log.Errorf("cache item <%s>: cannot retreat data file <%s>", name, dataPath)
			return nil, err
		}
		c.log.Debugf("GetCacheData(%s): retreated data file <%s>", name, dataPath)
		return b, nil
	}
	c.log.Errorf("GetCacheData(%s): cache item does not exist", name)
	return nil, fmt.Errorf("cache item <%s> does not exist", name)
}
func (c *cache) GetCacheItem(name string) (*cacheItem, error) {
	if ci, ok := c.Items[name]; ok {
		c.log.Debugf("GetCacheItem(%s) -- found", name)
		return ci, nil
	}
	c.log.Errorf("cache item <%s> does not exist", name)
	return nil, fmt.Errorf("cache item <%s> does not exist", name)
}
func (c *cache) AddCache(name, method, url, id, passwd string) error {
	// check if already exist
	if c.IsCacheExists(name) {
		c.log.Errorf("cache item <%s> already exists", name)
		return fmt.Errorf("cache item <%s> already exists", name)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items[name] = newCacheItem(method, url, id, passwd)
	c.log.Infof("added cache item <%s> with method: <%s>, URL: <%s>, ID: <%s>, passwd: <%s>",
		name, method, url, id, strings.Repeat("*", len(passwd)))
	return nil
}

func (c *cache) RemoveCache(name string, removeData bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.log.Infof("removing cache <%s>", name)

	if _, ok := c.Items[name]; ok {
		if err := os.Remove(path.Join(c.DataDir, name+DATA_FILE_EXT)); err != nil {
			c.log.Errorf("failed to remove cache <%s>'s data file <%s> -- %s", name, name+DATA_FILE_EXT, err.Error())
			return err
		}
		delete(c.Items, name)
		c.log.Infof("removed cache <%s>", name)
		return nil
	}
	c.log.Errorf("failed to remove cache <%s> -- cache not exist", name)
	return fmt.Errorf("cache item <%s> not exist", name)
}

func (c *cache) Save() error {
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

func NewConfig(filename string, addExampleCache bool) error {
	c := new()
	c.configFile = filename
	if addExampleCache {
		c.AddCache("test", "GET", "https://gonyyi.com/copyright.txt", "", "")
	}
	if err := c.Save(); err != nil {
		return err
	}
	return nil
}

func new() *cache {
	r := cache{}
	return r.init()
}

func New(configPath string) (*cache, error) {
	r := new()
	r.configFile = configPath

	r.log.Infof("initiating cache <%s>", configPath)

	if b, err := ioutil.ReadFile(configPath); err != nil {
		r.log.Errorf("cannot open the config file <%s>: %s", configPath, err.Error())
		return nil, err
	} else if err = json.Unmarshal(b, &r); err != nil {
		r.log.Errorf("cannot unmarshal the config file <%s>: %s", configPath, err.Error())
		return nil, err
	}
	r.log.Infof("loaded the config file <%s>", configPath)

	// Create a data directory (DataDir) if not exist
	{
		if r.DataDir == "" {
			r.log.Warnf("config does not hav data_dir, use default <./tmp> instead")
			r.DataDir = "./tmp"
		}

		if _, err := os.Stat(r.DataDir); err != nil {
			if os.IsNotExist(err) {
				r.log.Warnf("data_dir <%s> not exist, creating the directory", r.DataDir)
				if os.MkdirAll(r.DataDir, 0755) != nil {
					r.log.Fatalf("failed to create data_dir <%s>", r.DataDir)
				}
			} else {
				r.log.Fatalf("cannot open data_dir <%s>", r.DataDir)
			}
		}
	}

	r.log.Infof("total cache items: %d", len(r.Items))
	count := 0
	for name, _ := range r.Items {
		count += 1
		r.log.Infof("cache item [%d]: <%s>", count, name)
	}
	return r, nil
}
