# decouplet

Transcoder library for decoupling bytes from source objects in Go. 
This library takes bytes, transcodes using a transcoder function, 
and outputs final bytes as the original message measurements of a source object.

Currently, this library only contains an image transcoder.
In the future, there are plans for more versions of the same concept, 
and improvements to efficiency and better usage of source objects.

### Installation

`go get -u github.com/marcsj/decouplet`

### Testing

Place images named `test.jpg` and `test.png` in images folder.