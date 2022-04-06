package bonusly

import (
	"testing"
)

func Test_newReason(t *testing.T) {
	type args struct {
		params *CreateBonusInput
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"one-receiver",
			args{params: &CreateBonusInput{
				GiverEmail: "test@example.com",
				Receivers:  []string{"bilbo.baggins"},
				Reason:     "Test Reason",
				Amount:     1337,
			}},
			"+1337 @bilbo.baggins Test Reason",
		},
		{
			"two-receiver",
			args{params: &CreateBonusInput{
				GiverEmail: "test@example.com",
				Receivers:  []string{"bilbo.baggins", "frodo.baggins"},
				Reason:     "Test Reason",
				Amount:     1337,
			}},
			"+1337 @bilbo.baggins @frodo.baggins Test Reason",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newReason(tt.args.params); got != tt.want {
				t.Errorf("newReason() = %v, want %v", got, tt.want)
			}
		})
	}
}
