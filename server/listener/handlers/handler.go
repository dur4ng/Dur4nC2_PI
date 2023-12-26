package handlers

import "Dur4nC2/misc/protobuf/implantpb"

type ServerHandler func([]byte) *implantpb.Envelope

// GetHandlers - Returns a map of server-side msg handlers
func GetHandlers() map[uint32]ServerHandler {
	return map[uint32]ServerHandler{
		implantpb.MsgRegister:    beaconRegisterHandler,
		implantpb.MsgBeaconTasks: beaconTasksHandler,
	}
}
