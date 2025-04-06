#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 8080
#define MAX_CACHE_SIZE 100
#define MAX_KEY_SIZE 50
#define MAX_VALUE_SIZE 256

typedef struct {
    char key[MAX_KEY_SIZE];
    char value[MAX_VALUE_SIZE];
} CacheEntry;

CacheEntry cache[MAX_CACHE_SIZE];
int cache_size = 0;

void store_in_cache(const char *key, const char *value) {
    if (cache_size < MAX_CACHE_SIZE) {
        strncpy(cache[cache_size].key, key, MAX_KEY_SIZE);
        strncpy(cache[cache_size].value, value, MAX_VALUE_SIZE);
        cache_size++;
    } else {
        printf("Cache is full!\n");
    }
}

const char* retrieve_from_cache(const char *key) {
    for (int i = 0; i < cache_size; i++) {
        if (strcmp(cache[i].key, key) == 0) {
            return cache[i].value;
        }
    }
    return NULL; // Not found
}

int main() {
    int server_fd, new_socket;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);
    char buffer[1024] = {0};

    // Create socket
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("Socket failed");
        exit(EXIT_FAILURE);
    }

    // Attach socket to the port
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt))) {
        perror("setsockopt");
        exit(EXIT_FAILURE);
    }

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    // Bind the socket
    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("Bind failed");
        exit(EXIT_FAILURE);
    }

    // Listen for incoming connections
    if (listen(server_fd, 3) < 0) {
        perror("Listen");
        exit(EXIT_FAILURE);
    }

    printf("Cache server listening on port %d\n", PORT);

    while (1) {
        // Accept a new connection
        if ((new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen)) < 0) {
            perror("Accept");
            exit(EXIT_FAILURE);
        }

        // Read the incoming request
        read(new_socket, buffer, 1024);
        char command[10], key[MAX_KEY_SIZE], value[MAX_VALUE_SIZE];

        // Parse the command
        sscanf(buffer, "%s %s %s", command, key, value);

        if (strcmp(command, "STORE") == 0) {
            store_in_cache(key, value);
            send(new_socket, "Stored\n", 8, 0);
        } else if (strcmp(command, "RETRIEVE") == 0) {
            const char *result = retrieve_from_cache(key);
            if (result) {
                send(new_socket, result, strlen(result), 0);
            } else {
                send(new_socket, "Not found\n", 10, 0);
            }
        }

        close(new_socket);
    }

    return 0;
}
