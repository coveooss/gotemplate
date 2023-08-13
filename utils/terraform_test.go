package utils

import "testing"

func TestIsTerraformFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: ".tf file",
			path: "somefile.tf",
			want: true,
		},
		{
			name: ".tf.json file",
			path: "somefile.tf.json",
			want: true,
		},
		{
			name: ".tfvars file",
			path: "somefile.tfvars",
			want: true,
		},
		{
			name: ".json file",
			path: "somefile.json",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTerraformFile(tt.path); got != tt.want {
				t.Errorf("IsTerraformFile(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
