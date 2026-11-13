package xml

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/options"
)

// See https://gist.github.com/djangofan/3186223?permalink_comment_id=6148612#gistcomment-6148612
var xmlStrings = []string{
	`<statement name="INITIALIZATION">
  <block type="variables_set" id="setstake">
    <field name="VAR" id="stake">Stake</field>
    <value name="VALUE">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <next>
      <block type="variables_set" id="setinitialstake">
        <field name="VAR" id="initialstake">Initial Stake</field>
        <value name="VALUE">
          <block type="math_number">
            <field name="NUM">1</field>
          </block>
        </value>
        <next>
          <block type="variables_set" id="setmartingale">
            <field name="VAR" id="martingale">Martingale</field>
            <value name="VALUE">
              <block type="math_number">
                <field name="NUM">2</field>
              </block>
            </value>
            <next>
              <block type="variables_set" id="setdirection">
                <field name="VAR" id="direction">Direction</field>
                <value name="VALUE">
                  <block type="math_number">
                    <field name="NUM">1</field>
                  </block>
                </value>
                <next>
                  <block type="variables_set" id="setprofit">
                    <field name="VAR" id="takeprofit">Take Profit</field>
                    <value name="VALUE">
                      <block type="math_number">
                        <field name="NUM">10</field>
                      </block>
                    </value>
                    <next>
                      <block type="variables_set" id="setmaxloss">
                        <field name="VAR" id="maxloss">Max Loss</field>
                        <value name="VALUE">
                          <block type="math_number">
                            <field name="NUM">-20</field>
                          </block>
                        </value>
                        <next>
                          <block type="text_print" id="startmsg">
                            <value name="TEXT">
                              <shadow type="text">
                                <field name="TEXT">5 Minute Rise/Fall Switch Martingale Bot Started</field>
                              </shadow>
                            </value>
                          </block>
                        </next>
                      </block>
                    </next>
                  </block>
                </next>
              </block>
            </next>
          </block>
        </next>
      </block>
    </next>
  </block>
</statement>

<statement name="SUBMARKET">
  <block type="trade_definition_tradeoptions" id="tradeoptions">
    <mutation xmlns="http://www.w3.org/1999/xhtml" has_first_barrier="false" has_second_barrier="false" has_prediction="false"></mutation>
    <field name="DURATIONTYPE_LIST">m</field>
    <value name="DURATION">
      <shadow type="math_number_positive">
        <field name="NUM">5</field>
      </shadow>
    </value>
    <value name="AMOUNT">
      <shadow type="math_number_positive">
        <field name="NUM">1</field>
      </shadow>
      <block type="variables_get">
        <field name="VAR" id="stake">Stake</field>
      </block>
    </value>
  </block>
</statement>`,
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
            (-1)
            (-2)
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

func TestXMLToNix(t *testing.T) {
	t.Parallel()
	options := converter.ConverterOptions{
		SortIterators: options.SortIterators{
			SortList:    true,
			SortHashmap: true,
		},
	}

	for _, xmlString := range xmlStrings {
		_, err := ToNix(xmlString, &options)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestXMLFromNix(t *testing.T) {
	t.Parallel()
	for _, nixString := range nixStrings {
		_, err := FromNix(nixString, converter.NewDefaultConverterOptions())
		if err != nil {
			t.Fatal(err)
		}
	}
}
