package biz

import "testing"

func Test_getFullName(t *testing.T) {
	type args struct {
		firstName string
		lastName  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: struct {
				firstName string
				lastName  string
			}{
				firstName: "一二三四五",
				lastName:  "上山打老虎",
			},
			want: "一二三四五 上山打老...",
		},
		{
			args: struct {
				firstName string
				lastName  string
			}{
				firstName: "一二三四五六七八九十",
			},
			want: "一二三四五六七八九十",
		},
		{
			args: struct {
				firstName string
				lastName  string
			}{
				firstName: "零一二三四五六七八九十",
			},
			want: "零一二三四五六七八九...",
		},
		{
			args: struct {
				firstName string
				lastName  string
			}{
				firstName: "0123四五六七八九十",
			},
			want: "0123四五六七八九...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFullName(tt.args.firstName, tt.args.lastName); got != tt.want {
				t.Errorf("getFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}
