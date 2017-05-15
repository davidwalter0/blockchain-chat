## blockchain-chat

The application allows you to connect to a local p2p network and exchange encrypted anonymous messages. Anonymity means:

Lack of information about the IP-addresses of the sender and the recipient in the message
Absence of public keys of the sender and the recipient in the message
At the same time, the correctness of the dialogue (absence of message substitution) can always be confirmed by using a technology similar to blockchain.

Installation and startup

Installation:

```
go get github.com/poslegm/blockchain-chat
make install
```

Running:

```
make run
```

After that, when switching to localhost: 8080, the browser should open a graphical user interface.
To get started, you need to add the known IP addresses of network members to the ips.txt file in the root directory of the application. If this file is empty, then the application will work in the standby mode of the request to connect to the network.
In addition, to send messages, you need to add a couple of GPG keys to the local database. You can do this via the graphical interface, by clicking on the link "Add a key pair"


## blockchain-chat

The app allows you to connect to your local p2p-network and to exchange encrypted anonymous messaging. By anonymity means:

lack of information about IP-addresses of the sender and recipient of the message
the absence of a public key of the sender and recipient of the message
In this case, the correctness of the dialogue (no substitution messages) can always be confirmed through the use of technology similar blockchain.

Install and run

Setting:

```
go get github.com/poslegm/blockchain-chat
make install
```

Running:

```
make run
```

After that, when you go to localhost: 8080 in a browser should open a graphical user interface. 
To get started you need to add in ips.txt file in the root directory of the application known IP-addresses of network members. If this file is empty, the application will run in the request standby for connection to the network. 
In addition, to send messages to add to the local database a GPG key. This can be done through a graphical interface by clicking on the link "Add a pair of keys"

