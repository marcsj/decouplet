# decouplet

Encoding library for decoupling bytes using variable length keys in Go. 
This library takes bytes, looks at a key and takes deltas, 
and outputs a message as measurements of that key- 
effectively decoupling any meaning from a message without a key.

### Encoder Types

Type | Key | Delta
-----|-----|------
Image|image.Image|Pixel levels in RGBA, and CMYK
Byte |[]byte|Regular byte comparison with adds

### Uses

While this is not a typical encryption process, 
it does have similar value.
The typical use case would be with a small message 
or password, because while this process does 
effectively decouple its input to a high level, 
it also produces very large messages.

This can also be used with an already encrypted message,
or the output encrypted to further obfuscate a message.

### Installation

`go get -u github.com/marcsj/decouplet`

### Testing

Place images named `test.jpg` and `test.png` in images folder.
***
#### Credit

Idea based on *DVNC Whitepaper* by Joseph Lloyd under FDL1.3