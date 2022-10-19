package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/kolya_cypandin/project-base/internal/mocks/messages"
)

const exampleID int64 = 123_456_789

func Test_OnMessage_ShouldSendIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSender := mocks.NewMockMessageSender(ctrl)

	model := New(mockSender)

	mockSender.EXPECT().SendMessage("Hello, World!", exampleID)

	err := model.Send(Message{
		Text:   "Hello, World!",
		UserID: exampleID,
	})

	assert.NoError(t, err)
}
