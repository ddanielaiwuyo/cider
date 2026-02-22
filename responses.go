package main

import (
	"zod/protocol"
)

// This is hardcoded at the moment
// but for now it should match the
// port the server is running on
const serverId protocol.UserId = 4000

func CreateResponse(msg string, code protocol.MessageType) protocol.Response {
	return protocol.Response{
		From: serverId,
		Code: code,
		Msg:  msg,
	}
}


func ServerResponseMsg(msg string, code protocol.MessageType) protocol.Response {
	return protocol.Response{
		From: serverId,
		Code: code,
		Msg:  msg,
	}
}
