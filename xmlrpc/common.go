package xmlrpc

import ("time"; "fmt"; "encoding/base64";  "strconv"; "strings"; "xml"; "os")

// The interface that all values that want to be transmitted must conform to.
type MarshallUnmarshaller interface {
    // Get trasmittable value of this item
    Marshall() string;
    // Populate this item based on the provided value
    Unmarshall(string) os.Error;
}

// Simple types
type IntValue int;
type BooleanValue bool;
type StringValue string;
type DoubleValue float;
type DateTimeValue time.Time;
type Base64Value []byte;
type StructValue map[string] MarshallUnmarshaller;
type ArrayValue []MarshallUnmarshaller;

// An error returned from the rpc server
type Fault struct {
    FaultCode int;
    FaultString string;
}

// Make fault fit os.Error
func (f Fault) String() string {
    return f.FaultString
}

type error string;

func (e error) String() string {
    return string(e)
}

// Marshal and Unmarshall functions
func (i IntValue) Marshall() string {
    return fmt.Sprintf("<int>%v</int>", i)
}

func (i IntValue) Unmarshall(s string) (err os.Error) {
    var tempInt int;
    tempInt, err = strconv.Atoi(s);
    i = IntValue(tempInt);
    return
}

func (b BooleanValue) Marshall() string {
    return fmt.Sprintf("<boolean>%v</boolean>", b)
}

func (b BooleanValue) Unmarshall(s string) os.Error {
    switch s {
        case "0":
            b = BooleanValue(false)
        case "1":
            b = BooleanValue(true)
        default:
            return error(fmt.Sprintf("Unrecognized boolean: %s", s))
    }
    return nil;
}

func (s StringValue) Marshall() string {
    return fmt.Sprintf("<string>%v</string>", s)
}

func (s StringValue) Unmarshall(val string) os.Error {
    s = StringValue(val);
    return nil;
}

func (d DoubleValue) Marshall() string {
    return fmt.Sprintf("<double>%v</double>", d)
}

func (d DoubleValue) Unmarshall(val string) (err os.Error) {
    var tempDouble float;
    tempDouble, err = strconv.Atof(val);
    d = DoubleValue(tempDouble);
    return
}

func (d DateTimeValue) Marshall() string {
    // TODO try to get ISO8601 in stdlib
    return fmt.Sprintf("<dateTime.iso8601>%s</dateTime.iso8601>", "NOT IMPLEMENTED")
}

func (d DateTimeValue) Unmarshall(val string) (err os.Error) {
    return error("date Not Implemented")
}

func (b Base64Value) Marshall() string {
    encLen := base64.StdEncoding.EncodedLen(len(b));
    enc := make([]byte, encLen);
    base64.StdEncoding.Encode(enc, b);
    return fmt.Sprintf("<base64>%s</base64>", string(enc));
}

func (b Base64Value) Unmarshall(s string) (err os.Error) {
    decLen := base64.StdEncoding.DecodedLen(len(s));
    b = Base64Value(make([]byte, decLen));
    _,err = base64.StdEncoding.Decode(b, strings.Bytes(s));
    return
}

func (s StructValue) Marshall() (ret string) {
    ret = "<struct>";
    for key, value := range s {
        ret += fmt.Sprintf("<member><name>%s</name><value>%s</value></member>", key, value.Marshall())
    }
    ret += "</struct>";
    return
}
func (s ArrayValue) Marshall() (ret string) {
    ret = "<array><data>";
    for _, value := range s {
        ret += fmt.Sprintf("<value>%s</value>", value.Marshall())
    }
    ret += "</data></array>";
    return
}
func (f Fault) Marshall() string {
    faultStruct := StructValue{"faultCode": IntValue(f.FaultCode), "faultString": StringValue(f.FaultString)};
    return fmt.Sprintf("<fault>%s</fault>", faultStruct.Marshall())
}
