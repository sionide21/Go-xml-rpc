package xmlrpc_test

import (. "xmlrpc"; "testing"; "strings")

type comp func(MarshallUnmarshaller,MarshallUnmarshaller) bool;

func defaultCompare(a,b MarshallUnmarshaller) bool { return a == b }
func arrayCompare(a,b MarshallUnmarshaller) bool {
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
}

func runTest(xml string, value MarshallUnmarshaller, compare comp, t *testing.T) {
    m, err := ParseMessage(strings.NewReader(xml));
    if err != nil {
        t.Errorf("Could not parse (%s): %v", xml, err)
    }
    if !compare(m, value) {
        t.Errorf("Unexpected value: %v", m)
    }
}
