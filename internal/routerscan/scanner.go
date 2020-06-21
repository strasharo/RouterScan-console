package routerscan

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/dilap54/RouterScan-console/pkg/librouter"
)

type RouterScan struct {
	threadChannels  map[uint]chan tableData
	threadChanMutex *sync.Mutex
}

func (r *RouterScan) findFreeThread() uint {
	var i uint
	for i = 0; true; i++ {
		if _, ok := r.threadChannels[i]; !ok {
			return i
		}
	}
	panic("cannot find free thread")
}

func (r *RouterScan) getThreadChannel() (uint, chan tableData) {
	r.threadChanMutex.Lock()
	defer r.threadChanMutex.Unlock()
	threadID := r.findFreeThread()
	ch := make(chan tableData)
	r.threadChannels[threadID] = ch
	return threadID, ch
}

func (r *RouterScan) clearThreadChannel(threadID uint) {
	r.threadChanMutex.Lock()
	defer r.threadChanMutex.Unlock()
	close(r.threadChannels[threadID])
	delete(r.threadChannels, threadID)
}

func (r *RouterScan) handleSetTabledata(row uint, name string, value string) {
	r.threadChannels[row] <- tableData{name: name, value: value}
}
func (r *RouterScan) handleWriteLog(str string, verbosity int) {
	log.Printf("%d %s\n", verbosity, str)
}

func New(moduleOptions ModulesOptions, stOptions StOptions) (*RouterScan, error) {
	routerScan := RouterScan{
		threadChannels:  make(map[uint]chan tableData),
		threadChanMutex: &sync.Mutex{},
	}
	if err := librouter.Initialize(); err != nil {
		return nil, err
	}
	if err := librouter.SetSetTableDataCallback(routerScan.handleSetTabledata); err != nil {
		return nil, err
	}
	if err := librouter.SetWriteLogCallback(routerScan.handleWriteLog); err != nil {
		return nil, err
	}
	count, err := librouter.GetModuleCount()
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		info, err := librouter.GetModuleInfo(i)
		if err != nil {
			return nil, err
		}
		if strings.Contains(info.Name, "RouterScanRouter") && info.Enabled != moduleOptions.Scanrouter {
			if err := librouter.SwitchModule(i, moduleOptions.Scanrouter); err != nil {
				return nil, fmt.Errorf("cannot switch module RouterScanRouter to %t", moduleOptions.Scanrouter)
			}
		}
		if strings.Contains(info.Name, "ProxyCheckDetect") && info.Enabled != moduleOptions.Proxycheck {
			if err := librouter.SwitchModule(i, moduleOptions.Proxycheck); err != nil {
				return nil, fmt.Errorf("cannot switch module ProxyCheckDetect to %t", moduleOptions.Proxycheck)
			}
		}
		if strings.Contains(info.Name, "HNAPUse") && info.Enabled != moduleOptions.Hnap {
			if err := librouter.SwitchModule(i, moduleOptions.Hnap); err != nil {
				return nil, fmt.Errorf("cannot switch module HNAPUse to %t", moduleOptions.Hnap)
			}
		}
		if strings.Contains(info.Name, "SQLiteSQLite") && info.Enabled != moduleOptions.Sqlite {
			if err := librouter.SwitchModule(i, moduleOptions.Sqlite); err != nil {
				return nil, fmt.Errorf("cannot switch module SQLiteSQLite to %t", moduleOptions.Sqlite)
			}
		}
		if strings.Contains(info.Name, "HudsonHudson") && info.Enabled != moduleOptions.Hudson {
			if err := librouter.SwitchModule(i, moduleOptions.Hudson); err != nil {
				return nil, fmt.Errorf("cannot switch module HudsonHudson to %t", moduleOptions.Hudson)
			}
		}
		if strings.Contains(info.Name, "PMAphpMyAdmin") && info.Enabled != moduleOptions.Phpmyadmin {
			if err := librouter.SwitchModule(i, moduleOptions.Phpmyadmin); err != nil {
				return nil, fmt.Errorf("cannot switch module PMAphpMyAdmin to %t", moduleOptions.Phpmyadmin)
			}
		}
	}

	if err := librouter.SetParamBool(librouter.StEnableDebug, stOptions.EnableDebug); err != nil {
		return nil, err
	}
	if err := librouter.SetParamInt(librouter.StDebugVerbosity, stOptions.DebugVerbosity); err != nil {
		return nil, err
	}
	if err := librouter.SetParamString(librouter.StUserAgent, stOptions.UserAgent); err != nil {
		return nil, err
	}
	if err := librouter.SetParamString(librouter.StPairsBasic, stOptions.PairsBasic); err != nil {
		return nil, err
	}
	if err := librouter.SetParamString(librouter.StPairsDigest, stOptions.PairsDigest); err != nil {
		return nil, err
	}
	if err := librouter.SetParamString(librouter.StPairsForm, stOptions.PairsForm); err != nil {
		return nil, err
	}

	return &routerScan, nil
}

type tableData struct {
	name  string
	value string
}

type Result struct {
	IP       string
	Port     uint32
	Status   string
	Auth     string
	Type     string
	RadioOff string
	Hidden   string
	BSSID    string
	SSID     string
	Sec      string
	Key      string
	WPS      string
	LANIP    string
	LANMask  string
	WANIP    string
	WANMask  string
	WANGate  string
	DNS      string
	End      bool
}

func (r *RouterScan) Scan(target *Target) (*Result, error) {
	threadID, ch := r.getThreadChannel()
	defer r.clearThreadChannel(threadID)
	router, err := librouter.PrepareRouter(threadID, target.Ip, target.Port)
	if err != nil {
		return nil, err
	}
	result := Result{}
	go func() {
		for tableData := range ch {
			switch tableData.name {
			case "Status":
				result.Status = tableData.value
			case "Auth":
				result.Auth = tableData.value
			case "Type":
				result.Type = tableData.value
			case "RadioOff":
				result.RadioOff = tableData.value
			case "Hidden":
				result.Hidden = tableData.value
			case "BSSID":
				result.BSSID = tableData.value
			case "SSID":
				result.SSID = tableData.value
			case "Sec":
				result.Sec = tableData.value
			case "Key":
				result.Key = tableData.value
			case "WPS":
				result.WPS = tableData.value
			case "LANIP":
				result.LANIP = tableData.value
			case "LANMask":
				result.LANMask = tableData.value
			case "WANIP":
				result.WANIP = tableData.value
			case "WANMask":
				result.WANMask = tableData.value
			case "WANGate":
				result.WANGate = tableData.value
			case "DNS":
				result.DNS = tableData.value
			default:
				panic(fmt.Sprintf("unknown tableData.name %s", tableData.name))
			}
		}

	}()
	if err := router.Scan(); err != nil {
		return nil, err
	}
	if err := router.Free(); err != nil {
		return nil, err
	}
	return &result, nil
}
