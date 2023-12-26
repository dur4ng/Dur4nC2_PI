package implantpb

import (
	"google.golang.org/protobuf/proto"
)

const (
	MsgBeaconRegister = uint32(1 + iota)
	MsgBeaconTasks
	MsgRegister
	MsgBeaconID

	MsgWhoamiReq
	MsgWhoamiResp

	MsgRegisterExtensionReq
	MsgRegisterExtensionResp
	MsgListExtensionReq
	MsgListExtensionResp
	MsgCallExtensionReq
	MsgCallExtensionResp

	MsgExecuteAssemblyReq
	MsgExecuteAssemblyResp
	MsgExecuteShellcodeReq
	MsgExecuteShellcodeResp

	MsgUploadReq
	MsgUploadResp
	MsgDownloadReq
	MsgDownloadResp
)

// MsgNumber - Get a message number of type
func MsgNumber(request proto.Message) uint32 {
	switch request.(type) {

	case *Register:
		return MsgRegister
	case *BeaconTasks:
		return MsgBeaconTasks
	case *BeaconRegister:
		return MsgBeaconRegister
	}

	return uint32(0)
}
