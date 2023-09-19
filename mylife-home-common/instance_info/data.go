package instance_info

import "mylife-home-common/tools"

type InstanceInfo struct {
	typ            string
	hardware       map[string]string
	versions       map[string]string
	systemUptime   int64
	instanceUptime int64
	hostname       string
	capabilities   []string
	wifi           *WifiInfo
}

// 'ui' | 'studio' | 'core' | 'driver? (for arduino/esp/...)'
func (info *InstanceInfo) Type() string {
	return info.typ
}

// main: Raspberry ... | nodemcu | x64
// others are details like ram, cpu, ...
func (info *InstanceInfo) Hardware() tools.ReadonlyMap[string, string] {
	return tools.NewReadonlyMap(info.hardware)
}

// --- rpi
// os: linux-xxx
// node: 24.5
// mylife-home-core: 1.0.0
// mylife-home-common: 1.0.0
// --- esp/arduino
// mylife: 1.21.4
func (info *InstanceInfo) Versions() tools.ReadonlyMap[string, string] {
	return tools.NewReadonlyMap(info.versions)
}

func (info *InstanceInfo) SystemUptime() int64 {
	return info.systemUptime
}

func (info *InstanceInfo) InstanceUptime() int64 {
	return info.instanceUptime
}

func (info *InstanceInfo) Hostname() string {
	return info.hostname
}

func (info *InstanceInfo) Capabilities() tools.ReadonlySlice[string] {
	return tools.NewReadonlySlice(info.capabilities)
}

func (info *InstanceInfo) Wifi() *WifiInfo {
	return info.wifi
}

type WifiInfo struct {
	rssi int
}

func (info *WifiInfo) RSSI() int {
	return info.rssi
}

type instanceInfoData struct {
	Type           string
	Hardware       map[string]string
	Versions       map[string]string
	SystemUptime   int64
	InstanceUptime int64
	Hostname       string
	Capabilities   []string
	Wifi           *wifiInfoData
}

type wifiInfoData struct {
	RSSI int
}

func newInstanceInfo(data *instanceInfoData) *InstanceInfo {
	var wifi *WifiInfo
	wifiData := data.Wifi
	if wifiData != nil {
		wifi = &WifiInfo{
			rssi: wifiData.RSSI,
		}
	}

	return &InstanceInfo{
		typ:            data.Type,
		hardware:       data.Hardware,
		versions:       data.Versions,
		systemUptime:   data.SystemUptime,
		instanceUptime: data.InstanceUptime,
		hostname:       data.Hostname,
		capabilities:   data.Capabilities,
		wifi:           wifi,
	}
}

func extractData(info *InstanceInfo) *instanceInfoData {
	var wifiData *wifiInfoData
	wifiInfo := info.Wifi()
	if wifiInfo != nil {
		wifiData = &wifiInfoData{
			RSSI: wifiInfo.RSSI(),
		}
	}

	return &instanceInfoData{
		Type:           info.Type(),
		Hardware:       info.Hardware().Clone(),
		Versions:       info.Versions().Clone(),
		SystemUptime:   info.SystemUptime(),
		InstanceUptime: info.InstanceUptime(),
		Hostname:       info.Hostname(),
		Capabilities:   info.Capabilities().Clone(),
		Wifi:           wifiData,
	}
}
