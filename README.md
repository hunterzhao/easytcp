easytcp
==============

a good practice for golang & tcp

detail
==============
it is a custom tcp server. people can define their own 
- protocl & packet 
- callback for read new message & connection construct or close
- config for channal capacity


example
==============
cd example && go run echoServer.go
ps: echo/echoProto.go is customized by user

you can use telnet to test server ability

