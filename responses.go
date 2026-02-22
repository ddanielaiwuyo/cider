package main

// This is hardcoded at the moment
// but for now it should match the
// port the server is running on
const serverId userId = 4000

func CreateResponse(msg string) Response {
	return Response{
		From: serverId,
		Msg:  msg,
	}
}

