package message_test

import (
	"testing"

	"github.com/mgironi/operation-fire-quasar/message"
)

func TestGetMessage(t *testing.T) {
	type args struct {
		messages [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantMsg string
	}{
		{name: "TestEmpty", args: args{messages: [][]string{nil, nil, nil}}, wantMsg: ""},
		{name: "TestLessMessages", args: args{messages: [][]string{{"", "este", "es", "un", "mensaje"}, {"este", "", "un", "mensaje"}}}, wantMsg: " este es un mensaje"},
		{name: "TestMoreMessages", args: args{messages: [][]string{{"este", "", "", ""}, {"", "es", "", ""}, {"", "", "un", ""}, {"", "", "", "mensaje"}}}, wantMsg: "este es un mensaje"},
		{name: "TestCutedMessages", args: args{messages: [][]string{{"este", "", "", ""}, {"es", "", ""}, {"un", ""}, {"mensaje"}}}, wantMsg: "este es un mensaje"},
		{name: "TestEndTrimMessages", args: args{messages: [][]string{{"este"}, {"es"}, {"un"}, {"mensaje"}}}, wantMsg: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMsg := message.GetMessage(tt.args.messages...); gotMsg != tt.wantMsg {
				t.Errorf("GetMessage() = '%v', want '%v'", gotMsg, tt.wantMsg)
			}
		})
	}
}

func TestConsolidateMessage(t *testing.T) {
	type args struct {
		messages [][]string
	}
	tests := []struct {
		name                string
		args                args
		wantCompleteMessage string
		wantErr             bool
	}{

		{name: "testNils", args: args{[][]string{nil, nil, nil}}, wantCompleteMessage: "", wantErr: false},
		{name: "testSameSizes", args: args{[][]string{{"", "este", "es", "un", "mensaje"}, {"", "este", "", "un", "mensaje"}, {"", "", "es", "", "mensaje"}}}, wantCompleteMessage: " este es un mensaje", wantErr: false},
		{name: "testDiferentSizes", args: args{[][]string{{"", "este", "es", "un", "mensaje"}, {"este", "", "un", "mensaje"}, {"", "", "es", "", "mensaje"}}}, wantCompleteMessage: " este es un mensaje", wantErr: false},
		{name: "testMismatchWords", args: args{[][]string{{"", "este", "es", "un", "mensaje"}, {"este", "", "un", "mensaje"}, {"", "otra", "es", "", "mensaje"}}}, wantCompleteMessage: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCompleteMessage, err := message.ConsolidateMessage(tt.args.messages)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsolidateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCompleteMessage != tt.wantCompleteMessage {
				t.Errorf("ConsolidateMessage() = %v, want %v", gotCompleteMessage, tt.wantCompleteMessage)
			}
		})
	}
}
