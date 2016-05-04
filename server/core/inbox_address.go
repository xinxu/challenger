package core

import (
	"fmt"
)

type InboxAddressType int

const (
	InboxAddressTypeUnknown           = iota
	InboxAddressTypeAdminDevice       // 管理员iPad
	InboxAddressTypeSimulatorDevice   // 模拟器
	InboxAddressTypeArduinoTestDevice //测试Arduino设备
	InboxAddressTypePostgameDevice    // 出口处iPad
	InboxAddressTypeWearableDevice    // 穿戴设备
	InboxAddressTypeMainArduinoDevice // Arduino主墙设备
	InboxAddressTypeSubArduinoDevice  // Arduino小墙设备
)

func (t InboxAddressType) IsPlayerControllerType() bool {
	return t == InboxAddressTypeSimulatorDevice || t == InboxAddressTypeWearableDevice
}

func (t InboxAddressType) IsArduinoControllerType() bool {
	return t == InboxAddressTypeMainArduinoDevice || t == InboxAddressTypeSubArduinoDevice
}

type InboxAddress struct {
	Type InboxAddressType `json:"type"`
	ID   string           `json:"id"`
}

func (addr InboxAddress) String() string {
	return fmt.Sprintf("%v:%v", addr.Type, addr.ID)
}
