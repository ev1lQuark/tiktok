package logic

import (
	"bytes"
	"testing"
)

func Test_readFrameFromVideo(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				inputUrl: "http://192.168.0.210:9000/videos/1676811189258920932-345408069611110108320221129_185915.mp4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrame, err := readFrameFromVideo(tt.args.inputUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFrameFromVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFrame == nil {
				t.Errorf("readFrameFromVideo() gotFrame = nil")
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(gotFrame)
			t.Log(buf)
			t.Log(buf.Bytes())
			if len(buf.Bytes()) == 0 {
				t.Errorf("readFrameFromVideo() gotFrame length 0")
			}
		})
	}
}
