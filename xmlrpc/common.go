// This package provides low level xmlrpc functionality.
package xmlrpc

import ("time"; "fmt"; "encoding/base64";  "strconv"; "strings"; "xml"; "os"; "io")

// The interface that all values that want to be transmitted must conform to.
type ParamValue interface {
    // Get trasmittable value of this item
    ToXML() string;
    // Populate this item based on the provided value
    LoadXML(*xml.Parser) (ParamValue,os.Error);
}

// Allows for easier handling of xml tokens
type tokenStream struct { *xml.Parser }
func (t tokenStream) next(allowChar bool) (xml.Token, os.Error) {
    tok, err := t.Token()
    if err != nil { return nil, err }
    switch v := tok.(type) {
        case xml.StartElement: return v, nil
        case xml.EndElement:   return v, nil
        case xml.CharData:
            if !allowChar {
                return t.next(allowChar)
            }
            return v, nil
        case xml.Comment:      return t.next(allowChar)
        case xml.ProcInst:     return t.next(allowChar)
        case xml.Directive:    return t.next(allowChar)
    }
    return nil, error("Unknown Token Type")
}

// Simple types
type IntValue int;
type BooleanValue bool;
type StringValue string;
type DoubleValue float;
type DateTimeValue time.Time;
type Base64Value []byte;
type StructValue map[string] ParamValue;
type ArrayValue []ParamValue;

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

// ToXML and LoadXML functions
func (i IntValue) ToXML() string {
    return fmt.Sprintf("<int>%v</int>", i)
}

func (i IntValue) LoadXML(p *xml.Parser) (ParamValue,os.Error) {
    s, er := readBody(p);
    if er != nil { return nil, er }
    tempInt, err := strconv.Atoi(s);
    i = IntValue(tempInt);
    return i,err
}

func (b BooleanValue) ToXML() string {
    return fmt.Sprintf("<boolean>%v</boolean>", b)
}

func (b BooleanValue) LoadXML(p *xml.Parser) (ParamValue,os.Error) {
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

func (s StringValue) ToXML() string {
    return fmt.Sprintf("<string>%v</string>", s)
}

func (s StringValue) LoadXML(p *xml.Parser) (ParamValue,os.Error) {
    val, er := readBody(p);
    if er != nil { return nil, er }
    s = StringValue(val);
    return s,nil;
}

func (d DoubleValue) ToXML() string {
    return fmt.Sprintf("<double>%v</double>", d)
}

func (d DoubleValue) LoadXML(p *xml.Parser) (ParamValue, os.Error) {
    val, er := readBody(p);
    if er != nil { return nil, er }
    tempDouble, err := strconv.Atof(val);
    d = DoubleValue(tempDouble);
    return d,err
}

func (d DateTimeValue) ToXML() string {
    // TODO try to get ISO8601 in stdlib
    return fmt.Sprintf("<dateTime.iso8601>%s</dateTime.iso8601>", "NOT IMPLEMENTED")
}

func (d DateTimeValue) LoadXML(p *xml.Parser) (ParamValue, os.Error) {
    return d,error("date Not Implemented")
}

func (b Base64Value) ToXML() string {
    encLen := base64.StdEncoding.EncodedLen(len(b));
    enc := make([]byte, encLen);
    base64.StdEncoding.Encode(enc, b);
    return fmt.Sprintf("<base64>%s</base64>", string(enc));
}

func (b Base64Value) LoadXML(p *xml.Parser) (ParamValue, os.Error) {
    s, er := readBody(p);
    if er != nil { return nil, er }
    decLen := base64.StdEncoding.DecodedLen(len(s));
    b = Base64Value(make([]byte, decLen));
    rLen,err := base64.StdEncoding.Decode(b, strings.Bytes(s));
    b = b[0:rLen];
    return b,err
}

func (s StructValue) ToXML() (ret string) {
    ret = "<struct>";
    for key, value := range s {
        ret += fmt.Sprintf("<member><name>%s</name><value>%s</value></member>", key, value.ToXML())
    }
    ret += "</struct>";
    return
}

func (s StructValue) LoadXML(parser *xml.Parser) (ParamValue, os.Error) {
    p := tokenStream{parser}
    for {
        var (
            name string;
            value ParamValue
        )
        // Skip <member>
        t, err := p.next(false);
        if err != nil {return nil,err}
        end,ok := t.(xml.EndElement);
        if ok {
            if strings.ToLower(end.Name.Local) == "struct" {
                return s, nil
            }
        }

        for x:=0;x<2;x++ {
            t, err = p.next(false);
            if err != nil {return nil,err}
            field, ok := t.(xml.StartElement);
            if !ok {return nil, error(fmt.Sprintf("Unexpected Token: %v, (%+T)", t, t, x)) }
            switch field.Name.Local {
                case "name":
                    name,err = readBody(p.Parser);
                case "value":
                    value,err = parseMessage(p.Parser);
                    // Skip </value>
                    p.next(false);
            }
            if err != nil {return nil, err}
        }
        s[name] = value;
        // Skip </member>
        t, err = p.next(false);
        if err != nil {return nil,err}
    }
    return s, nil
}

func (s ArrayValue) ToXML() (ret string) {
    ret = "<array><data>";
    for _, value := range s {
        ret += fmt.Sprintf("<value>%s</value>", value.ToXML())
    }
    ret += "</data></array>";
    return
}

func (a ArrayValue) LoadXML(parser *xml.Parser) (ParamValue, os.Error) {
    if len(a) == 0 { a = make(ArrayValue, 2) }
    var x int
    p := tokenStream{parser}
    // Skip <data>
    t, err := p.next(false);
    if err != nil {return nil,err}
    for x=0;;x++{
        var value ParamValue;
        t, err = p.next(false);
        if err != nil {return nil,err}
        end,ok := t.(xml.EndElement);
        if ok {
            if end.Name.Local == "data" {
                // skip </array>
                p.next(false)
                return a[0:x], nil
            }
        }
        value,err = parseMessage(p.Parser);
        if err != nil {return nil, err}
        if cap(a) <= x {
            b := make(ArrayValue, 2*x);
            for i,c := range a {
                b[i] = c
            }
            a = b;
        }
        a[x] = value;
        // Skip </value>
        t, err = p.next(false);
        if err != nil {return nil,err}
    }
    return a[0:x], nil

}

func (f Fault) ToXML() string {
    faultStruct := StructValue{"faultCode": IntValue(f.FaultCode), "faultString": StringValue(f.FaultString)}
    return faultStruct.ToXML()
}

func (f Fault) LoadXML(parser *xml.Parser) (ParamValue, os.Error) {
    p := tokenStream{parser}
    t, err := p.next(false);
    if err != nil {return nil, err}
    start,ok := t.(xml.StartElement);
    if !ok || start.Name.Local != "value" {return nil, error("Unexpected symbol")}
    m, err := parseMessage(p.Parser);
    if err != nil {return nil, err}
    s,ok1 := m.(StructValue);
    msg, ok2 := s["faultString"].(StringValue);
    code, ok3 := s["faultCode"].(IntValue);
    if !(ok1 && ok2 && ok3) {return nil, error("Invalid fault response")}
    f.FaultString = string(msg);
    f.FaultCode = int(code);
    return f, nil
}

func ParseMessage(r io.Reader) (ParamValue, os.Error) {
    p := xml.NewParser(r);
    return parseMessage(p)
}

func parseMessage(p *xml.Parser) (ParamValue, os.Error) {
    t, err := tokenStream{p}.next(false);
    if err != nil {return nil, err}
    start,ok := t.(xml.StartElement);
    if !ok {return nil, error(fmt.Sprintf("Unexpected symbol: %v", t))}
    switch strings.ToLower(start.Name.Local) {
        case "int":
            fallthrough
        case "i4":
            return IntValue(0).LoadXML(p)
        case "boolean":
            return BooleanValue(false).LoadXML(p)
        case "string":
            return StringValue("").LoadXML(p)
        case "double":
            return DoubleValue(0).LoadXML(p)
        case "datetime.iso8601":
            return DateTimeValue{}.LoadXML(p)
        case "base64":
            return make(Base64Value,0).LoadXML(p)
        case "struct":
            return StructValue{}.LoadXML(p)
        case "fault":
            return Fault{}.LoadXML(p)
        case "array":
            return make(ArrayValue,0).LoadXML(p)
        default: return nil, error(fmt.Sprintf("Unknown type: %s", start.Name.Local))
    }
    return nil, nil
}

func readBody(p *xml.Parser) (string, os.Error) {
    ret := "";
    for {
        t,err := tokenStream{p}.next(true);
        if (err != nil) { return "", err }
        switch v := t.(type) {
            case xml.CharData:
                ret += string(v)
            case xml.EndElement:
                return ret, nil
            default:
                return "", error("Unexpected token")
        }
    }
    return "", error("UNREACHABLE CODE")
}
