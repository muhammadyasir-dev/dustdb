#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "DynamicArray.h"

// Create a new dynamic array with an initial capacity
DynamicArray* createDynamicArray(size_t initialCapacity) {
    DynamicArray *array = (DynamicArray*)malloc(sizeof(DynamicArray));
    array->data = (char*)malloc(initialCapacity * sizeof(char));
    array->size = 0;
    array->capacity = initialCapacity;
    return array;
}

// Append a character to the dynamic array
void appendCharacter(DynamicArray *array, char c) {
    // Resize the array if necessary
    if (array->size >= array->capacity) {
        array->capacity *= 2; // Double the capacity
        array->data = (char*)realloc(array->data, array->capacity * sizeof(char));
    }
    array->data[array->size++] = c; // Add the character and increment size
}

// Delete a character at a specific index
void deleteCharacter(DynamicArray *array, size_t index) {
    if (index < array->size) {
        for (size_t i = index; i < array->size - 1; i++) {
            array->data[i] = array->data[i + 1]; // Shift characters left
        }
        array->size--; // Decrease the size
    }
}

// Free the dynamic array
void freeDynamicArray(DynamicArray *array) {
    free(array->data);
    free(array);
}

// Print the contents of the dynamic array
void printDynamicArray(DynamicArray *array) {
    for (size_t i = 0; i < array->size; i++) {
        putchar(array->data[i]);
    }
    putchar('\n');
}
