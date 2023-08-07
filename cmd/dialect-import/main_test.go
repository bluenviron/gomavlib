package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDialect = `<?xml version="1.0"?>
<mavlink>
  <version>0</version>
  <dialect>0</dialect>
  <enums>
    <enum name="A_TYPE">
      <description>Detected Anomaly Types.</description>
      <entry value="0" name="A">
        <description>A.</description>
      </entry>
      <entry value="1" name="B">
        <description>B.</description>
      </entry>
      <entry value="2" name="C">
        <description>C.</description>
      </entry>
      <entry value="3" name="D">
        <description>D.</description>
      </entry>
      <entry value="4" name="E">
        <description>E</description>
      </entry>
    </enum>
  </enums>
  <messages>
    <!-- Detected anomaly info measured by onboard sensors and actuators -->
    <message id="43000" name="A_MESSAGE">
      <description>Detected anomaly info measured by onboard sensors and actuators. </description>
      <field type="uint64_t" name="timestamp" units="us">Timestamp (UNIX epoch time)</field>
      <field type="uint8_t" name="a_field" instance="true">whether anomaly has been detected or not</field>
      <field type="uint8_t" name="b_field" enum="A_TYPE">which anomaly has been detected.</field>
    </message>
  </messages>
</mavlink>
`

func TestRun(t *testing.T) {
	err := os.WriteFile("testdialect.xml", []byte(testDialect), 0o644)
	require.NoError(t, err)
	defer os.Remove("testdialect.xml")

	err = run([]string{"", "testdialect.xml"})
	require.NoError(t, err)

	_, err = os.ReadFile("testdialect/message_a_message.go")
	require.NoError(t, err)

	os.RemoveAll("testdialect")
}
