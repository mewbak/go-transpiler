package python

// converter represents a single go-type and provides the
// necessary functionality to convert variables of this
// type between go and c code
type converter interface {

    // GoType should return the name of the go type
    // that this converter is representing
    //GoType() string

    // GoTransitionType returns the argument type that should be used
    // on go functions exported to c (usually C.* types)
    GoTransitionType() string

    // CTransitionType returns the return value type that c should use to
    // define extern functions exported from go
    CTransitionType() string

    // ConvertGoFromC should return valid go code that
    // takes a c representation of this type and converts it
    // back to go code. The result of this function should be assignable
    ConvertGoFromC(varName string) string

    // ConvertGoToC should return valid go code that
    // takes a go representation of this type and converts it
    // to c code. The result of this function should be assignable
    ConvertGoToC(varName string) string

    // ConvertCFromGo should return valid C code that
    // takes a go representation of this type and converts it
    // back to c code. The result of this function should be assignable
    //
    // Because Go is considered the owner of all data, and manages
    // conversions through caches, this should
    // usually be a straight assignment or simple lookup, not
    // creating new data
    ConvertCFromGo(varName string) string

    // ConvertCToGo should return valid c code that
    // takes a c representation of this type and converts it
    // to the type expected by go code. The result of this function
    // should be assignable
    ConvertCToGo(varName string) string

    // ConvertPyFromC should return valid C code that
    // takes a PyObject* representation of this type and converts it
    // to a c-type. The result of this function should be assignable
    ConvertPyFromC(varName string) string

    // ConvertPyToC should return valid c code that
    // takes a c representation of this type and converts it
    // to a valid PyObject*. The result of this function
    // should be assignable
    ConvertPyToC(varName string) string

    // ValidatePyValue produces c code that ensures a PyObject*
    // is of the correct type for assignment. The result of this
    // function should function in the brackets of an if statement
    ValidatePyValue(varName string) string

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

    // CDeclarations should return decalrations for anything defined in
    // CDefinitions function. This will be included in the header of all
    // generated c files
    CDeclarations() string

    // CDefintions returns any global definitions that should be included
    // in the conversions.c file. These are utility functions or values
    // that are needed for this converter (this function will only called
    // once for each converter implementation)
    CDefinitions() string

    // GoDefintions returns any global definitions that should be included
    // in the conversions.go file. These are utility functions or values
    // that are needed for this converter (this function will only called
    // once for each converter implementation)
    GoDefinitions() string
}
