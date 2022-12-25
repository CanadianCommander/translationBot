package slackutil

import "github.com/slack-go/slack"

//==========================================================================
// Public
//==========================================================================

// MsgBlocks extracts block slice from slack message
// #### params
// message - message to extract blocks from
func MsgBlocks(message slack.Message) []slack.Block {
	return message.Blocks.BlockSet
}

// SlackMessageToMsgOption converts a Message in to a MsgOption. This is useful when sending messages via slack API.
// #### params
// message - the message to convert
func SlackMessageToMsgOption(message slack.Message) slack.MsgOption {
	return slack.MsgOptionBlocks(MsgBlocks(message)...)
}
