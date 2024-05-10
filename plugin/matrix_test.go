package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thegeeklab/wp-matrix/plugin/mocks"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func TestMatrixMessageSend(t *testing.T) {
	tests := []struct {
		name       string
		messageOpt MatrixMessageOpt
		want       event.MessageEventContent
		wantErr    bool
	}{
		{
			name: "plain text message",
			messageOpt: MatrixMessageOpt{
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
			messageOpt: MatrixMessageOpt{
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
			messageOpt: MatrixMessageOpt{
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
			messageOpt: MatrixMessageOpt{
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
			messageOpt: MatrixMessageOpt{
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
			mockClient := mocks.NewMockIMatrixClient(t)
			m := &MatrixMessage{
				Opt:    tt.messageOpt,
				client: mockClient,
			}

			mockClient.
				On("SendMessageEvent", mock.Anything, tt.messageOpt.RoomID, event.EventMessage, tt.want).
				Return(&mautrix.RespSendEvent{}, nil)

			err := m.Send(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
