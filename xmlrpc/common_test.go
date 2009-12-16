package xmlrpc_test

import (. "xmlrpc"; "testing"; "strings")


func TestXMLReader(t *testing.T) {
    xml := "<int>4</int>";
    m, err := ParseMessage(strings.NewReader(xml));
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    i, ok := m.(IntValue);
    if !ok {
        t.Errorf("Unexpected item %T", m)
    }
    if int(i) != 4 {
        t.Errorf("Unexpected value: %d", int(i))
    }
}
