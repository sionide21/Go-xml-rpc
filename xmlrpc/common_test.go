package xmlrpc_test

import (. "xmlrpc"; "testing"; "strings"; "bytes")

type comp func(ParamValue,ParamValue) bool;

func defaultCompare(a,b ParamValue) bool { return a == b }
func arrayCompare(a,b ParamValue) bool {
    ap, ok := a.(Base64Value);
    if !ok {return false}
    bp, ok := b.(Base64Value);
    if !ok {return false}

    return string(ap) == string(bp)
}

func TestSimpleXMLReader(t *testing.T) {
    runTest("<int>4</int>", IntValue(4), defaultCompare, t);
    runTest("<string>Ben</string>", StringValue("Ben"), defaultCompare, t);
    runTest("<boolean>1</boolean>", BooleanValue(true), defaultCompare, t);
    runTest("<double>3.14</double>", DoubleValue(3.14), defaultCompare, t);
    runTest("<base64>SGVsbG8gV29ybGQ=</base64>", Base64Value(strings.Bytes("Hello World")), arrayCompare, t)
    //runTest("<?xml version=\"1.0\"?>\n<array> \n\t<data>\n\t\t<value><base64>SGVsbG8gV29ybGQ=</base64></value>\n\t</data>\n</array>", Base64Value(strings.Bytes("Hello World")), arrayCompare, t)
    //runTest("<?xml version=\"1.0\"?>\n<struct>\n\t<member>\t<name>Something</name><value><array> \n\t<data>\n\t\t<value><base64>SGVsbG8gV29ybGQ=</base64></value>\n\t</data>\n</array></value></member></struct>", Base64Value(strings.Bytes("Hello World")), arrayCompare, t)
}

func TestResponse(t *testing.T) {
    r := `<?xml version="1.0"?>
            <methodResponse>
               <params>
                  <param>
                     <value><string>South Dakota</string></value>
                     </param>
                  </params>
               </methodResponse>`
    resp, _ := ReadResponse(bytes.NewBufferString(r))
    defaultCompare(resp.Value, StringValue("South Dakota"))
}

func runTest(xml string, value ParamValue, compare comp, t *testing.T) {
    m, err := ParseMessage(strings.NewReader(xml));
    if err != nil {
        t.Errorf("Could not parse (%s): %v", xml, err)
    }
    if !compare(m, value) {
        t.Errorf("Unexpected value: %v", m)
    }
}
