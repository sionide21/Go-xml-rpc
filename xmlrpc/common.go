package xmlrpc

import ("time")


// Make fault fit os.Error
func (f Fault) String() string {
    return f.faultString
}

// The interface that all values that want to be transmitted must conform to.
type MarshallUnmarshaller interface {
    // Get trasmittable value of this item
    func Marshall() string;
    // Populate this item based on the provided value
    func Unmarshall(string);
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
    faultCode int;
    faultString string;
}
