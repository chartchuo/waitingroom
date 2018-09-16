package main

import (
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/fsnotify.v1"
	yaml "gopkg.in/yaml.v2"
)

//todo refactor to library

/*
 Watches a file on a set interval, and preforms de-duplication of write
 events such that only 1 write event is reported even if multiple writes
 happened during the specified duration.
*/
type FileWatcher struct {
	fsNotify *fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	callback func()
}

/*
 Begin watching a file with a specific interval and action
*/
func WatchFile(path string, interval time.Duration, action func()) (*FileWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Add the file to be watched
	fsWatcher.Add(path)

	watcher := &FileWatcher{
		fsWatcher,
		interval,
		make(chan struct{}, 1),
		action,
	}
	// Launch a go thread to watch the file
	go watcher.run()

	return watcher, err
}

func (fw *FileWatcher) run() {
	// Check for write events at this interval
	tick := time.Tick(fw.interval)

	var lastWriteEvent *fsnotify.Event
	for {
		select {
		case event := <-fw.fsNotify.Events:
			// When a ConfigMap update occurs kubernetes AtomicWriter() creates a new directory;
			// writing the updated ConfigMap contents to the new directory. Once the write is
			// complete it removes the original file symlink and replaces it with a new symlink
			// pointing to the contents of the newly created directory. It does this to achieve
			// atomic ConfigMap updates. But it also means the file we were monitoring for write
			// events never got them and was instead deleted.

			// The correct way to handle this would be to monitor the symlink instead of the
			// actual file for events. However fsnotify.v1 does not allow us to pass in the
			// IN_DONT_FOLLOW flag to inotify which would allow us to monitor the
			// symlink for changes instead of the de-referenced file. This is not likely to
			// change as fsnotify is designed as cross platform and not all platforms support
			// symlinks.

			if event.Op == fsnotify.Remove {
				// Since the symlink was removed, we must
				// re-register the file to be watched
				fw.fsNotify.Remove(event.Name)
				fw.fsNotify.Add(event.Name)
				lastWriteEvent = &event
			}

			// If it was a write event
			if event.Op == fsnotify.Write {
				lastWriteEvent = &event
			}
		case <-tick:
			// No events during this interval
			if lastWriteEvent == nil {
				continue
			}
			// Execute the callback
			fw.callback()
			// Reset the last event
			lastWriteEvent = nil
		case <-fw.done:
			goto Close
		}
	}
Close:
	close(fw.done)
}

func (fw *FileWatcher) Close() {
	fw.done <- struct{}{}
	fw.fsNotify.Close()
}

/*
 Simple interface that allows us to switch out both implementations of the Manager
*/
type ConfigManager interface {
	Set(*Config)
	Get() *Config
	Close()
}

/*
 This struct manages the configuration instance by
 preforming locking around access to the Config struct.
*/
type MutexConfigManager struct {
	conf  *Config
	mutex *sync.Mutex
}

func NewMutexConfigManager(conf *Config) *MutexConfigManager {
	return &MutexConfigManager{conf, &sync.Mutex{}}
}

func (mcm *MutexConfigManager) Set(conf *Config) {
	mcm.mutex.Lock()
	mcm.conf = conf
	mcm.mutex.Unlock()
}

func (mcm *MutexConfigManager) Get() *Config {
	mcm.mutex.Lock()
	temp := mcm.conf
	mcm.mutex.Unlock()
	return temp
}

func (mcm *MutexConfigManager) Close() {
	//Do Nothing
}

/*
 This struct manages the configuration instance by feeding a
 pointer through a channel whenever the user calls Get()
*/
type ChannelConfigManager struct {
	conf *Config
	get  chan *Config
	set  chan *Config
	done chan bool
}

func NewChannelConfigManager(conf *Config) *ChannelConfigManager {
	parser := &ChannelConfigManager{conf, make(chan *Config), make(chan *Config), make(chan bool)}
	parser.Start()
	return parser
}

func (ccm *ChannelConfigManager) Start() {
	go func() {
		defer func() {
			close(ccm.get)
			close(ccm.set)
			close(ccm.done)
		}()
		for {
			select {
			case ccm.get <- ccm.conf:
			case value := <-ccm.set:
				ccm.conf = value
			case <-ccm.done:
				return
			}
		}
	}()
}

func (ccm *ChannelConfigManager) Close() {
	ccm.done <- true
}

func (ccm *ChannelConfigManager) Set(conf *Config) {
	ccm.set <- conf
}

func (ccm *ChannelConfigManager) Get() *Config {
	return <-ccm.get
}

//  Simple Yaml Config file loader
func loadConfig(configFile string) (*Config, error) {
	conf := &Config{}
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
