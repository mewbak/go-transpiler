package python

// converter represents a single go-type and provides the
// necessary functionality to convert variables of this
// type between go and c code
type converter interface {

    // GoType should return the name of the go type
    // that this converter is representing
    GoType() string

    // CMemberType should return the c type name
    // for this converters go type. This is the type of variable
    // that should be stored in a c struct to define the underlying
    // data of a python object
    CMemberType() string

    // GoIncomingArgType returns the argument type that should be used
    // on go functions exported to c (usually C.* types)
    GoIncomingArgType() string

    // COutgoingArgType returns the argument type that c should use to
    // define extern functions exported from go
    COutgoingArgType() string

    // ConvertFromCValue should return valid go code that
    // takes a c representation of this type and converts it
    // back to go code. The result of this function should be assignable
    ConvertFromCValue(varName string) string

    // ConvertToCValue should return valid go code that
    // takes a go representation of this type and converts it
    // to c code. The result of this function should be assignable
    ConvertToCValue(varName string) string

    // ConvertFromGoValue should return valid C code that
    // takes a go representation of this type and converts it
    // back to c code. The result of this function should be assignable
    //
    // Because Go is considered the owner of all data, and manages
    // conversions through caches, this should
    // usually be a straight assignment or simple lookup, not
    // creating new data
    ConvertFromGoValue(varName string) string

    // ConvertToGoValue should return valid c code that
    // takes a c representation of this type and converts it
    // to the type expected by go code. The result of this function
    // should be assignable
    ConvertToGoValue(varName string) string

    // PyMemberDefTypeEnum returns the python enum that properly
    // describes this type in a member definition (eg PY_OBJECT_EX)
    PyMemberDefTypeEnum() string

    // PyTupleTarget returns a string that defines target variables
    // to be used when parsing this varaible type out of a tuple set.
    // The given identifier string should be appended to variable names
    // in order to ensure uniqueness
    PyTupleTarget(ident int) string

    // PyParseTupleArgs returns a string that are arguments
    // to be passed to tuple argument parding functions in python api.
    // The given identifier string should be appended to variable names
    // in order to match what would be returned from PythonTupleTarget
    PyParseTupleArgs(ident int) string

    // PyTupleResult returns the name of the var generated for PyTupleTarget
    PyTupleResult(ident int) string

    // PyTupleFormat returns the set of format character(s) that
    // define this variable type as represented in python tuples
    // ex (int vars would return "i")
    PyTupleFormat() string
}
