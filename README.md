# Packet racer
A CLI project demonstrating packet transfer speed of various UDP (User Datagram Protocol) techniques.
This includes `net.Dial`, `AF_INET`, `AF_INET raw socket`, `AF_PACKET` methods. It then compares the performance of each method in terms of number of packets transferred and Mega Bytes transferred and arranges them in descending order.

## Build It Locally

Pre Requisites
- Currently works only on Ubuntu/Linux
- Needs GO already Installed

Clone the repository using the following command
```
git clone https://github.com/jaydee029/packet_racer
```

Move into the root folder, and build the project
```
cd packet_racer
go build -o packet
```

Now You can run the binary within the folder with or without flags li
```
sudo ./packet
```
Root privileges are required to run the binary, since certain methods such as raw sockets require root permissions.

To run the program globally add the binary to your globally path

## Functionality Available

### Optional Flags
- `-duration` : allows you to change the duration for which each method runs. By default this value is set to `1 Second`
- `-port` : allows the user to change the port at which the listening server spins off/ By default this value is set to `8080`
- `-size` : allows you to change the size of packets being sent across in the udp connection. By default it has its value set to `1470 bytes`

### How It Works
The program uses `netcat` to start the listening server at the given Port, the server runs for some given time, The udp methods are invoked synchronously on the same listener server, once the connnection is established the transfer of IP packets begin.

- `The net.Dial Method`: It is baseline implementation provided by the `net` library in golang , whoch allows connecting to a udp connection, and transfer ip packets by specifying the portnumber and the ip address.
- `AF_INET Method`: using the syscall method which is the underlying method used in the baseline implementation, we manually specify the use of AF_INET packet, specify the IP socket format and destination port and ip address. This method is generally faster than the baseline implementation.
- `AF_INET Raw Socket Method`: uses the syscall method similar to the AF_INET implementation, but instead using a predefined IP packet, an ip packet is constructed from scratch, source and destination ports as well as ip address, message size in bytes are specified while constructing the packet.
- `AF_Packet Method`: Uses the syscall method along with the AF_Packet method using a raw socket similar to the above method, it may require a source and destination mac address as well as an ethernet header along with the already specified options. This method is usually the fastest among the above methods.

