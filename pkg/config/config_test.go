package config

import (
	"reflect"
	"testing"
)

func Test_config_GetCnfDefault(t *testing.T) {
	type args struct {
		key    string
		defCnf interface{}
	}
	tests := []struct {
		name string

		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{
			name: "",

			args: args{
				key:    "ok",
				defCnf: nil,
			},
			want: nil,
		},
		{
			name: "",
			args: args{
				key:    "test.name",
				defCnf: nil,
			},
			want: "tyltr",
		},
		{
			name: "",
			args: args{
				key:    "test.age",
				defCnf: nil,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			c.SetCnfFile("yaml", "test", ".")
			c.Load()
			got := c.GetCnfDefault(tt.args.key, tt.args.defCnf)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.GetCnfDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
