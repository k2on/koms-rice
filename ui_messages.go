package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

type MessagesComponent = *tview.List

func MakeMessages(client Client, state AppState) (MessagesComponent, UpdateStateFn) {
	messages := tview.NewList()
	UpdateMessagesStyle(messages, state)

	updateMessages := MakeMessagesUpdateFn(client, messages)

	updateMessages(state)

	return messages, updateMessages
}

func MakeMessagesUpdateFn(client Client, messages MessagesComponent) UpdateStateFn {
	return func(state AppState) {
		UpdateMessagesStyle(messages, state)
		messages.Clear()

		conversationMessages, exists := GetCacheMessages(state)
		if !exists { return }

		conversation := GetStateConversation(state)
		messagePos := GetStateMessagePos(state)
		JumpNumb := MakeMessageJumpNumbFn(len(conversationMessages), messagePos)

		for i, message := range conversationMessages {
			parsedMessage := ParseMessage(client, conversation, message)
			jumpNumb := "[grey]" + JumpNumb(i) + "[-] "
			parsedMessage = jumpNumb + parsedMessage
			messages.AddItem(parsedMessage, "", 0, nil)
		}

		messages.SetCurrentItem(messagePos)
	}
}

func MakeMessageJumpNumbFn(msgLen int, selectedPos int) func(int) string {
	PlaceValues := func(i int) int { return int(math.Log10(float64(i))) + 1 }
	Pad := func(len int) string { return strings.Repeat(" ", len) }

	placeValuesMax := PlaceValues(msgLen)

	return func(i int) string {
		distance := int(math.Abs(float64(i - selectedPos)))
		if distance == 0 { return Pad(placeValuesMax) }

		placeValues := PlaceValues(distance)
		padding := placeValuesMax - placeValues
		paddingStr := Pad(padding)
		return paddingStr + strconv.Itoa(distance)
	}
}

func AddContainerMessages(container *tview.Grid, messages *tview.List) {
	isFocused := true
	container.AddItem(
		messages,
		ROW_POS_MSGS,
		COLUMN_POS_MSGS,
		ROW_SPAN_MSGS,
		COLUMN_SPAN_MSGS,
		HEIGHT_MIN_MSGS,
		WIDTH_MIN_MSGS,
		isFocused,
	)
}
