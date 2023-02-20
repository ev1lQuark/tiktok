package jwt

import (
	"testing"
)

const key = "12345678"
const userId = 1

func TestJWT(t *testing.T) {
	type args struct {
		secretKey string
		seconds   int64
		userId    int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				secretKey: key,
				seconds:   3600,
				userId:    userId,
			},
			want:    userId,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := Create(tt.args.secretKey, tt.args.seconds, tt.args.userId)
			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}
			if !Verify(tt.args.secretKey, tokenString) {
				t.Errorf("Verify() error")
				return
			}
			got, err := GetUserId(tt.args.secretKey, tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}
