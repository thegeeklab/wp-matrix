package matrix

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thegeeklab/wp-matrix/matrix/mocks"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func TestMessageSend(t *testing.T) {
	tests := []struct {
		name       string
		messageOpt MessageOptions
		want       event.MessageEventContent
		wantErr    bool
	}{
		{
			name: "plain text message",
			messageOpt: MessageOptions{
				RoomID:  "test-room",
				Message: "hello world",
			},
			want: event.MessageEventContent{
				MsgType: "m.text",
				Body:    "hello world",
			},
		},
		{
			name: "markdown message",
			messageOpt: MessageOptions{
				RoomID:  "test-room",
				Message: "**hello world**",
			},
			want: event.MessageEventContent{
				MsgType:       "m.text",
				Body:          "**hello world**",
				Format:        "org.matrix.custom.html",
				FormattedBody: "<strong>hello world</strong>",
			},
		},
		{
			name: "html message",
			messageOpt: MessageOptions{
				RoomID:         "test-room",
				Message:        "hello<br>world",
				TemplateUnsafe: true,
			},
			want: event.MessageEventContent{
				MsgType:       "m.text",
				Body:          "hello\nworld",
				Format:        "org.matrix.custom.html",
				FormattedBody: "hello<br>world",
			},
		},
		{
			name: "safe html message",
			messageOpt: MessageOptions{
				RoomID:         "test-room",
				Message:        "hello world<script>alert('XSS')</script>",
				TemplateUnsafe: false,
			},
			want: event.MessageEventContent{
				MsgType:       "m.text",
				Body:          "hello world<script>alert('XSS')</script>",
				Format:        "org.matrix.custom.html",
				FormattedBody: "hello world&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
			},
		},
		{
			name: "unsafe html message",
			messageOpt: MessageOptions{
				RoomID:         "test-room",
				Message:        "hello world<script>alert('XSS')</script>",
				TemplateUnsafe: true,
			},
			want: event.MessageEventContent{
				MsgType:       "m.text",
				Body:          "hello world",
				Format:        "org.matrix.custom.html",
				FormattedBody: "hello world",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockClient := mocks.NewMockAPIClient(t)
			m := &Message{
				Opt:    tt.messageOpt,
				client: mockClient,
			}

			mockClient.
				On("SendMessageEvent", mock.Anything, tt.messageOpt.RoomID, event.EventMessage,
					mock.MatchedBy(func(content event.MessageEventContent) bool {
						tt.want.Mentions = &event.Mentions{}

						return assert.EqualValues(t, tt.want, content)
					})).
				Return(&mautrix.RespSendEvent{}, nil)

			err := m.Send(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
