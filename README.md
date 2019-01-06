## GRPC TCP Tunnel

This piece of software is able to make a TCP tunnel over GRPC, using GRPC streaming (HTTP/2 channels).

### Examples

Listen for GRPC tunnel clients on port 22223
```
go run . server :22223
```

Listen on port 18080 using ncat, tunneling through the client part of this software.
The client will connect to the GRPC server @ ':22223', making it connect to example.com port 80.
```
while true; do ncat -l 18080 -c 'go run . client :22223 example.com 80'; done
```