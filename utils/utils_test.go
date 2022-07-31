package utils

import "testing"

func TestGetIndexOf(t *testing.T) {
	type args struct {
		element string
		data    []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success",
			args: args{
				element: "cat",
				data:    []string{"dog", "cat", "bird"},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIndexOf(tt.args.element, tt.args.data); got != tt.want {
				t.Errorf("GetIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
