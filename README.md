funcmock is function mocking library which replaces the lexical placeholder of the function with a mock. It is primarily intended to be used for testing purposes.

<!-- Documentation -->
<!-- ------------- -->

<!-- After installing, you can use `go doc` to get documentation: -->

<!--     go doc github.com/golang/mock/gomock -->

<!-- Alternatively, there is an online reference for the package hosted on GoPkgDoc -->
<!-- [here][gomock-ref]. -->


 Using funcmock
---------------

Example:

    package main
    
    import "funcmock"
    
    func funcToBeMocked() {
    	// I am computation intensive, hence take long
    	// Don't use me unnecessarily
    }
    
    func iUseFuncToBeMockedTenTimes(){
    
    	for _ : range 10{
    		funcToBeMocked()
    	}
    }
    
    func main() {
    
    	// The mock engine specific to a function
    	funcMockEngine = funcmock.Mock(&funcToBeMocked)
    
    	// call iUseFuncToBeMockedTenTimes which calls funcToBeMocked
    	iUseFuncToBeMockedTenTimes() 	// actual funcToBeMocked is not called
    
    	fmt.Println("Mock function called: ", funcMockEngine.CallCount())
    	// prints "Mock function called: 10"
    	
    }
    

<!-- The `mockgen` command is used to generate source code for a mock -->
<!-- class given a Go source file containing interfaces to be mocked. -->
<!-- It supports the following flags: -->

<!--  *  `-source`: The file containing interfaces to be mocked. You must -->
<!--     supply this flag. -->

<!--  *  `-destination`: A file to which to write the resulting source code. If you -->
<!--     don't set this, the code is printed to standard output. -->

<!--  *  `-package`: The package to use for the resulting mock class -->
<!--     source code. If you don't set this, the package name is `mock_` concatenated -->
<!--     with the package of the input file. -->

<!--  *  `-imports`: A list of explicit imports that should be used in the resulting -->
<!--     source code, specified as a comma-separated list of elements of the form -->
<!--     `foo=bar/baz`, where `bar/baz` is the package being imported and `foo` is -->
<!--     the identifier to use for the package in the generated source code. -->

<!--  *  `-aux_files`: A list of additional files that should be consulted to -->
<!--     resolve e.g. embedded interfaces defined in a different file. This is -->
<!--     specified as a comma-separated list of elements of the form -->
<!--     `foo=bar/baz.go`, where `bar/baz.go` is the source file and `foo` is the -->
<!--     package name of that file used by the -source file. -->

<!-- For an example of the use of `mockgen`, see the `sample/` directory. In simple -->
<!-- cases, you will need only the `-source` flag. -->


<!-- TODO: Brief overview of how to create mock objects and set up expectations, and -->
<!-- an example. -->

<!-- [golang]: http://golang.org/ -->
<!-- [golang-install]: http://golang.org/doc/install.html#releases -->
<!-- [gomock-ref]: http://godoc.org/github.com/golang/mock/gomock -->
