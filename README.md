# DDOS ARMY
![Banner Image](img/ddosarmydemo.gif)



## What is DDOS ARMY?
ddosarmy is a tool that allows
you to launch a DDoS attack on IP victim address,
with other machines synchronously.

### ⚠️ DISCLAIMER ⚠️:
️ It is intended to be used for educational purposes only.
I am not responsible for any misuse or damage caused by this program.
## Features
* 🏕️ Support clustering attack concept 
* 💣 Support ICMP, SYN flood attack (more to come)
* ⛓️ Support concurrency and parallelism
* Cross-platform (Linux 💀, Windows 🪟, Mac OS X🍏, BSD)
* ⚙️ Easy Army management (stop, start, shutdown, change attack strategy, etc.)
* ✅ Friendly terminal user interface
* 📈 Profiling and benchmarking requests for each soldier (node)
## Installation

### binary (recommended)
```bash
go install github.com/XORbit01/ddosarmy@latest 
```
you can download the binary from 
[Releases](https://github.com/XORbit01/DDOS-ARMY/releases/tag/v0.0.5)


### From source
```bash
git clone https://github.com/XORbit01/DDOS-ARMY 
cd ddosarmy
go build -o ddosarmy 
```

## Usage
```bash
ddosarmy -h
```

## steps
### 1. create army camp
this must be hosted on a server with public ip address, so that the soldiers can connect to it.
```bash
ddosarmy camp -v <victim_ip>
```
### 2. leader access camp
leader is the node that will orchestrate the army and send the attack commands to the soldiers.
```bash
ddosarmy leader -c <camp_ip> -s password
```
the default password is `password`, this password is set when you created the camp
### 3. soldiers access camp
the soldiers are the nodes responsible for sending dummy packets flooding the victim ip address.
it will listen to the leader node and execute the attack commands.
```bashq
ddosarmy soldier -c <camp_ip>  
```


## How it works
as you can see in the image below, the leader node is responsible for sending the attack commands to the soldiers.
the dispatcher (camp) is responsible for sending the attack commands to the soldiers.
he is not responsible for sending a single packet to the victim ip address.

![How it works](img/howitworks.jpg)


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
you can also contact me on discord:`XORbit#5945`


## Powered By Malwarize
[![image](https://user-images.githubusercontent.com/130087473/232165094-73347c46-71dc-47c0-820a-1eb36657a8c0.png)](https://discord.gg/g9y7D3xCab)


 
