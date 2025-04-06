#ifndef SOCKET_SERVER_H
#define SOCKET_SERVER_H

#include <iostream>
#include <cstring>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>

#define PORT 8989
#define BUFFER_SIZE 1024

class SocketServer {
public:
    SocketServer();
    ~SocketServer();
    void start();

private:
    int server_fd;
    struct sockaddr_in address;
    int opt;
    int addrlen;
    char buffer[BUFFER_SIZE];

    void setupSocket();
    void acceptConnections();
    void handleClient(int client_socket);
};

#endif // SOCKET_SERVER_H
