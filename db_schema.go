package main

type UserTable struct {
	username string
	uuid     string
	email    string
}

/*
A user can have many friends, so we have
a one-to-many relationship here
*/

type ContactsTable struct {
	userId    string
	contactId string
	contacts  []string // not sure if we're just going to store userIds here
}

type ConversationTable struct {
	id    string
	peers []string
	topic string // topic of conversation?
}

type Message struct {
	id          string
	senderId    string
	recipientId string
	content     string
	timeStamp   string
}

type MessageTable struct {
	uuid        string
	senderId    string
	recipientId string
	message     string
	/*
	 we hit this table to find the the recipientId who the sender is looking for
	 and look for a converationId where both of them are involved. 
	 And then we can query this table again but based on convoId
	 */
	converationId string
}
