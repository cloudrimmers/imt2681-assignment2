package environment

import "testing"

func TestLoad(t *testing.T) {

	table := map[[2]string]bool{
		[2]string{"", "heroku"}:                 true,
		[2]string{"", "docker"}:                 true,
		[2]string{"", "local"}:                  true,
		[2]string{"", "test"}:                   false,
		[2]string{"", "tes"}:                    false,
		[2]string{"", "d"}:                      false,
		[2]string{"", "dskfjaølskfdiølkfsaødl"}: false,
	}

	for arg, want := range table {
		t.Run(arg[1], func(t *testing.T) {

			argslice := arg[0:]
			err := Load(argslice)
			if err == nil {
				if !want {
					t.Error("ARG validate, but we did not want it to!.")
				}
			} else {
				if want {
					t.Error("ARG did not validate, but we wanted it to", err.Error())
				}
			}
		})
	}
}
