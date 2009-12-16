package xmlrpc

import ("time"; "fmt"; "encoding/base64";  "strconv"; "strings"; "xml"; "os"; "io")

// The interface that all values that want to be transmitted must conform to.
type MarshallUnmarshaller interface {
    // Get trasmittable value of this item
    Marshall() string;
    // Populate this item based on the provided value
    Unmarshall(*xml.Parser) (MarshallUnmarshaller,os.Error);
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
    return fmt.Sprintf("%s (%d)", f.FaultString, f.FaultCode)
}

type error string;

func (e error) String() string {
    return string(e)
}

// Marshal and Unmarshall functions
func (i IntValue) Marshall() string {
    return fmt.Sprintf("<int>%v</int>", i)
}

func (i IntValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller,os.Error) {
    s, er := readBody(p);
    if er != nil { return nil, er }
    tempInt, err := strconv.Atoi(s);
    i = IntValue(tempInt);
    return i,err
}

func (b BooleanValue) Marshall() string {
    return fmt.Sprintf("<boolean>%v</boolean>", b)
}

func (b BooleanValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller,os.Error) {
    s, er := readBody(p);
    if er != nil { return nil, er }
    switch s {
        case "0":
            b = BooleanValue(false)
        case "1":
            b = BooleanValue(true)
        default:
            return nil,error(fmt.Sprintf("Unrecognized boolean: %s", s))
    }
    return b,nil;
}

func (s StringValue) Marshall() string {
    return fmt.Sprintf("<string>%v</string>", s)
}

func (s StringValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller,os.Error) {
    val, er := readBody(p);
    if er != nil { return nil, er }
    s = StringValue(val);
    return s,nil;
}

func (d DoubleValue) Marshall() string {
    return fmt.Sprintf("<double>%v</double>", d)
}

func (d DoubleValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    val, er := readBody(p);
    if er != nil { return nil, er }
    tempDouble, err := strconv.Atof(val);
    d = DoubleValue(tempDouble);
    return d,err
}

func (d DateTimeValue) Marshall() string {
    // TODO try to get ISO8601 in stdlib
    return fmt.Sprintf("<dateTime.iso8601>%s</dateTime.iso8601>", "NOT IMPLEMENTED")
}

func (d DateTimeValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    return d,error("date Not Implemented")
}

func (b Base64Value) Marshall() string {
    encLen := base64.StdEncoding.EncodedLen(len(b));
    enc := make([]byte, encLen);
    base64.StdEncoding.Encode(enc, b);
    return fmt.Sprintf("<base64>%s</base64>", string(enc));
}

func (b Base64Value) Unmarshall(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    s, er := readBody(p);
    if er != nil { return nil, er }
    decLen := base64.StdEncoding.DecodedLen(len(s));
    b = Base64Value(make([]byte, decLen));
    rLen,err := base64.StdEncoding.Decode(b, strings.Bytes(s));
    b = b[0:rLen];
    return b,err
}

func (s StructValue) Marshall() (ret string) {
    ret = "<struct>";
    for key, value := range s {
        ret += fmt.Sprintf("<member><name>%s</name><value>%s</value></member>", key, value.Marshall())
    }
    ret += "</struct>";
    return
}

func (s StructValue) Unmarshall(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    for {
        var (
            name string;
            value MarshallUnmarshaller
        )
        // Skip <member>
        t, err := p.Token();
        if err != nil {return nil,err}
        end,ok := t.(xml.EndElement);
        if ok {
            if end.Name.Local == "struct" {
                return s, nil
            }
        }

        for x:=0;x<2;x++ {
            t, err = p.Token();
            if err != nil {return nil,err}
            field, ok := t.(xml.StartElement);
            if !ok {return nil, error("Unexpected Token")}
            switch field.Name.Local {
                case "name":
                    name,err = readBody(p);
                case "value":
                    value,err = parseMessage(p);
                    // Skip </value>
                    _,_ = p.Token();
            }
            if err != nil {return nil, err}
        }
        s[name] = value;
        // Skip </member>
        t, err = p.Token();
        if err != nil {return nil,err}
    }
    return s, nil
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

func (f Fault) Unmarshall(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    t, err := p.Token();
    if err != nil {return nil, err}
    start,ok := t.(xml.StartElement);
    if !ok || start.Name.Local != "value" {return nil, error("Unexpected symbol")}
    m, err := parseMessage(p);
    if err != nil {return nil, err}
    s,ok1 := m.(StructValue);
    msg, ok2 := s["faultString"].(StringValue);
    code, ok3 := s["faultCode"].(IntValue);
    if !(ok1 && ok2 && ok3) {return nil, error("Invalid fault response")}
    f.FaultString = string(msg);
    f.FaultCode = int(code);
    return f, nil
}

func ParseMessage(r io.Reader) (MarshallUnmarshaller, os.Error) {
    p := xml.NewParser(r);
    return parseMessage(p)
}

func parseMessage(p *xml.Parser) (MarshallUnmarshaller, os.Error) {
    t, err := p.Token();
    if err != nil {return nil, err}
    start,ok := t.(xml.StartElement);
    if !ok {return nil, error("Unexpected symbol")}
    switch start.Name.Local {
        case "int":
            fallthrough
        case "i4":
            return IntValue(0).Unmarshall(p)
        case "boolean":
            return BooleanValue(false).Unmarshall(p)
        case "string":
            return StringValue("").Unmarshall(p)
        case "double":
            return DoubleValue(0).Unmarshall(p)
        case "dateTime.iso8601":
            return DateTimeValue{}.Unmarshall(p)
        case "base64":
            return Base64Value(make([]byte,0)).Unmarshall(p)
        case "struct":
            return StructValue{}.Unmarshall(p)
        case "fault":
            return Fault{}.Unmarshall(p)
        default: return nil, error(fmt.Sprintf("Unknown type: %s", start.Name.Local))
    }
    return nil, nil
}

func readBody(p *xml.Parser) (string, os.Error) {
    ret := "";
    for {
        t,err := p.Token();
        if (err != nil) { return "", err }
        switch v := t.(type) {
            case xml.CharData:
                ret += string(v)
            case xml.EndElement:
                return ret, nil
            case xml.ProcInst:
            case xml.Comment:
            case xml.Directive:
            default:
                return "", error("Unexpected token")
        }
    }
    return "", error("UNREACHABLE CODE")
}
