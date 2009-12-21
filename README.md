GO XML-RPC Library
==================

This is a library to bring [xml-rpc][1] functionality to the GO programming language.

[1]: <http://www.xmlrpc.com/> "XML-RPC Homepage"

Requirements
------------
The Makefile provided expects you to have setup the environment vairables `$GOROOT` and `$GOARCH` as instructed in the [Go setup guide][]. If you do not have these, you can still compile the project manually.

[Go setup guide]: http://golang.org/doc/install.html#tmp_17 "Go setup - Environment Variables"

Installation
------------

	Go into the 'xmlrpc' sub-directory.
	Type 'make install'

This will compile the library and install it in your `$GOROOT/pkg` directory. It is safe to run this to upgrade the library as well. To remove simply type `make nuke` in the same directory you typed install.

Usage
-----

Create a new `xmlrpc.RemoteMethod` with your endpoint and method and the `Call` it. For example to list all blogs a user has on a wordpress site:

	listBlogs := xmlrpc.RemoteMethod{
		Endpoint: "http://examlpe.com/wp/xmlrpc.php",
		Method: "wp.getUsersBlogs",
	}

	username := "testuser"
	password := "hunter2"
	
	result, error := listBlogs.Call(username, password)
	
`result` will be of type ParamValue and will require a type assertion to do much with. Read more about `ParamValue`s on the [wiki][].

[wiki]: http://wiki.github.com/sionide21/Go-xml-rpc/paramvalue "Lst of ParamValues"