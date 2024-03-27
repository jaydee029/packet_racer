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
Root privileges are required to run the binary, since certain methods within the project require root permissions.

To run the program globally add the binary to your globally path

## Functionality Available

### Optional Flags
- `-duration` : allows you to change the duration for which each method runs. By default this value is set to `1 Second`
- `-port` : allows the user to change the port at which the listening server spins off/ By default this value is set to `8080`
- `-size` : allows you to change the size of packets being sent across in the udp connection. By default it has its value set to `1470 bytes`

### How It Works
The program uses netcat to start the listening server at the given Port, the server runs for some given time, The udp methods are invoked synchronously on the same listener server, once the connnection is established the transfer of IP packets begin.

- `The UDP Method`: 


## About
## Motivation
