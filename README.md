# tsteg-poc
## A steganography POC demonstrating encoding via transmission delays
This is a rough proof of concept showing a technique for encoding data into network traffic by strategically delaying transmission times. This technique requires no modification of bytes traveling over the wire and can encode a payload over arbitrary data. While the transmission rate may be comparatively slow, it is sufficient to issue remote shell commands.

### Building
* go install github.com/Uberroot/tsteg-poc/client
* go install github.com/Uberroot/tsteg-poc/server

### Usage
* Server: echo "some arbitrary evil data" | server &lt;harmless file&gt;
* Client: client &lt;server address&gt; 1&gt;evil 2&gt;harmless

### Notes
* Payload bits are manchester-encoded to prevent synchronization issues.
* The data rate currently encodes one bit per 5ms raw / 10ms manchester, capping transfers at 100bps. This rate has been tested on local networks and may need to be adjusted for tests over the internet.
* The "harmless" file must be at least popcount(evil) bytes in size
* This has only been tested on Linux systems.
