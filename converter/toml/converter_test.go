package toml

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/options"
)

var tomlStrings = []string{
	`title = "TOML Example"
# integers
int1 = +99
int2 = 42
int3 = 0
int4 = -17

# hexadecimal with prefix 0x
hex1 = 0xDEADBEEF
hex2 = 0xdeadbeef
hex3 = 0xdead_beef

# octal with prefix 0o
oct1 = 0o01234567
oct2 = 0o755

# binary with prefix 0b
bin1 = 0b11010110

# fractional
float1 = +1.0
float2 = 3.1415
float3 = -0.01

# exponent
float5 = 1e06
float6 = -2E-2

# both
float7 = 6.626e-34

# separators
float8 = 224_617.445_991_228

# infinity
infinite1 = inf # positive infinity
infinite2 = +inf # positive infinity
infinite3 = -inf # negative infinity

# not a number
not1 = nan
not2 = +nan
not3 = -nan 
[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00

[database]
enabled = true
ports = [ 8000, 8001, 8002 ]
data = [ ["delta", "phi"], [3.14] ]
temp_targets = { cpu = 79.5, case = 72.0 }
temp_targ = 12345
"" = 22
ggg = ""
ee = ""
123123123123123 = "-11"
a = [1,2,3,-1, "-1", -123, "-123123123"]

[servers]

[servers.alpha]
ip = "10.0.0.1"
role = "frontend"

[servers.beta]
ip = "10.0.0.2"
role = "backend"`,
}

var nixStrings = []string{
	`{
  id = "c7d8e9f0";
  users = [
    {
      name = "Alice";
      age = 28;
      pets = [
        {
          type = "cat";
          name = "Luna";
          toys = [
            "mouse"
            "ball"
          ];
        }
        {
          type = "dog";
          name = "Max";
        }
      ];
    }
    {
      name = "Bob";
      age = 34;
      age2 = -34;
      age3 = -3.45;
      ag2 = -0;
      "" = 123;
      helloooo = "";
      age12 = -0.45;
      ag2123 = 0.001;
      age4 = -0.0045;
      age5 = -0.00000000000000000000000045;
      pets = null;
    }
  ];
  settings = {
    theme = {
      dark = {
        primary = "#1a1a1a";
        accent = "#4287f5";
      };
      light = {
        primary = "#ffffff";
        accent = "#2196f3";
      };
    };
    notifications = true;
  };
  meta = {
    created = "2024-01-01";
    modified = {
      by = "system";
      timestamp = "2024-02-15T14:30:00Z";
    };
  };
}`,
}

// Not comparing anything because its using Go maps (unordered)
func TestTOMLToNix(t *testing.T) {
	options := converter.ConverterOptions{
		SortIterators: options.SortIterators{
			SortList:    true,
			SortHashmap: true,
		},
	}

	for _, tomlString := range tomlStrings {
		_, err := ToNix(tomlString, &options)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Not comparing anything because its using Go maps (unordered)
func TestTOMLFromNix(t *testing.T) {
	for _, nixString := range nixStrings {
		_, err := FromNix(nixString, converter.NewDefaultConverterOptions())
		if err != nil {
			t.Fatal(err)
		}
	}
}
