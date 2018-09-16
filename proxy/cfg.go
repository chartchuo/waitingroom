package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	fsnotify "gopkg.in/fsnotify.v1"
	yaml "gopkg.in/yaml.v2"
)

const CONFIG_FILE = "/etc/config/configmap-microservice-demo.yaml"
const BIND = "0.0.0.0:8080"

type Config struct {
	Message string `yaml:"message"`
}

func loadConfig(configFile string) *Config {
	conf := &Config{}
	configData, err := ioutil.ReadFile(configFile)
	check(err)

	err = yaml.Unmarshal(configData, conf)
	check(err)
	return conf
}

func cfg() {
	confManager := NewConfigManager(loadConfig(CONFIG_FILE))

	// Create a single GET Handler to print out our simple config message
	router := httprouter.New()
	router.GET("/", func(resp http.ResponseWriter, req *http.Request, p httprouter.Params) {
		conf := confManager.Get()
		fmt.Fprintf(resp, "%s", conf.Message)
	})

	// Watch the file for modification and update the config
	// manager with the new config when it's available
	watcher, err := WatchFile(CONFIG_FILE, time.Second, func() {
		fmt.Printf("Configfile Updated\n")
		conf := loadConfig(CONFIG_FILE)
		confManager.Set(conf)
	})
	check(err)

	// Clean up
	defer func() {
		watcher.Close()
	}()

	fmt.Printf("Listening on '%s'....\n", BIND)
	err = http.ListenAndServe(BIND, router)
	check(err)
}

func NewConfigManager(conf *Config) *ConfigManager {
	return &ConfigManager{conf, &sync.Mutex{}}
}

func (self *ConfigManager) Set(conf *Config) {
	self.mutex.Lock()
	self.conf = conf
	self.mutex.Unlock()
}

func (self *ConfigManager) Get() *Config {
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
	}()
	return self.conf
}

func (self *FileWatcher) run() {
	// Check for write events at this interval
	tick := time.Tick(self.interval)

	var lastWriteEvent *fsnotify.Event
	for {
		select {
		case event := <-self.fsNotify.Events:
			// If we see a Remove event, we know the config was updated
			if event.Op == fsnotify.Remove {
				// Since the symlink was removed, we must
				// re-register the file to be watched
				self.fsNotify.Remove(event.Name)
				self.fsNotify.Add(event.Name)
				lastWriteEvent = &event
			}
		case <-tick:
			// No events during this interval
			if lastWriteEvent == nil {
				continue
			}
			// Execute the callback
			self.callback()
			// Reset the last event
			lastWriteEvent = nil
		case <-self.done:
			goto Close
		}
	}
Close:
	close(self.done)
}
