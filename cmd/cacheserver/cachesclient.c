#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 8080

int main() {
    int sock = 0;
    struct sockaddr_in serv_addr;
    char buffer[1024] = {0};

    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        printf("\n Socket creation error \n");
        return -1;
    }

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(PORT);

    // Convert IPv4 and IPv6 addresses from text to binary form
    if (inet_pton(AF_INET, "127.0.0.1", &serv_addr.sin_addr) <= 0) {
        printf("\nInvalid address/ Address not supported \n");
        return -1;
    }

    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
        printf("\nConnection Failed \n");
        return -1;
    }

    // Example of storing a value
    const char *store_command = "STORE key1 value1";
    send(sock, store_command, strlen(store_command), 0);
    read(sock, buffer, 1024);
    printf("%s\n", buffer);

    // Example of retrieving a value
    const char *retrieve_command = "RETRIEVE key1";
    send(sock, retrieve_command, strlen(retrieve_command), 0);
    read(sock, buffer, 1024);
    printf("%s\n", buffer);

    close(sock);
    return 0;
}
