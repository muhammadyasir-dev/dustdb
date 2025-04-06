#include <stdio.h>
#include <stdlib.h>
#include "HashSet.h"

// Hash function
unsigned int hash(int key) {
    return key % HASHSET_SIZE;
}

// Create a new hash set
HashSet* createHashSet() {
    HashSet *set = (HashSet*)malloc(sizeof(HashSet));
    for (int i = 0; i < HASHSET_SIZE; i++) {
        set->buckets[i] = NULL;
    }
    return set;
}

// Insert a key into the hash set
void insert(HashSet *set, int key) {
    unsigned int index = hash(key);
    Node *newNode = (Node*)malloc(sizeof(Node));
    newNode->key = key;
    newNode->next = set->buckets[index];
    set->buckets[index] = newNode;
}

// Check if a key exists in the hash set
bool contains(HashSet *set, int key) {
    unsigned int index = hash(key);
    Node *current = set->buckets[index];
    while (current != NULL) {
        if (current->key == key) {
            return true;
        }
        current = current->next;
    }
    return false;
}

// Remove a key from the hash set
void removeKey(HashSet *set, int key) {
    unsigned int index = hash(key);
    Node *current = set->buckets[index];
    Node *prev = NULL;

    while (current != NULL) {
        if (current->key == key) {
            if (prev == NULL) {
                set->buckets[index] = current->next;
            } else {
                prev->next = current->next;
            }
            free(current);
            return;
        }
        prev = current;
        current = current->next;
    }
}

// Free the hash set
void freeHashSet(HashSet *set) {
    for (int i = 0; i < HASHSET_SIZE; i++) {
        Node *current = set->buckets[i];
        while (current != NULL) {
            Node *temp = current;
            current = current->next;
            free(temp);
        }
    }
    free(set);
}
