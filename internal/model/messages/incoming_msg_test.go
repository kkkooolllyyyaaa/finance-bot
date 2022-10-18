package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/kolya_cypandin/project-base/internal/mocks/messages"
)

const exampleID int64 = 123_456_789

func Test_OnStartCommand_ShouldAnswerWithMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSender := mocks.NewMockMessageSender(ctrl)

	model := New(mockSender)

	mockSender.EXPECT().SendMessage("Hello, World!", exampleID)

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: exampleID,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("Unknown command", exampleID)
	model := New(sender)

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: exampleID,
	})

	assert.NoError(t, err)
}
