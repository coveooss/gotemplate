package template

import (
	"reflect"
	"testing"
)

func Test_average(t *testing.T) {

	tests := []struct {
		name    string
		arg1    interface{}
		args    a
		want    interface{}
		wantErr bool
	}{
		{"Nil", nil, nil, nil, true},
		{"First nil", nil, l{1, 2}, nil, true},
		{"Single", 1, nil, int64(1), false},
		{"Two values", 1, l{2}, 1.5, false},
		{"With nil", 1, l{2, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := average(tt.arg1, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("average() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("average() = %[1]v (%[1]T), want %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	tests := []struct {
		name   string
		values l
		want   interface{}
	}{
		{"Nil", nil, nil},
		{"Zero", l{0}, int64(0)},
		{"Single float", l{1.1}, 1.1},
		{"Array of floats", l{1.1, 2.2, 3.3, 4}, 1.1},
		{"Mixed array", l{1.1, 2.2, 3.3, "hello"}, "1.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("min() = %[1]v (%[1]T), want %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	tests := []struct {
		name   string
		values l
		want   interface{}
	}{
		{"Nil", nil, nil},
		{"Zero", l{0}, int64(0)},
		{"Single float", l{1.1}, 1.1},
		{"Array of floats", l{1.1, 2.2, 3.3, 4}, int64(4)},
		{"Mixed array", l{1.1, 2.2, 3.3, "hello"}, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("max() = %[1]v (%[1]T), want %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}
