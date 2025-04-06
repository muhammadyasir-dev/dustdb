#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct LMSNode {
    int key;
    struct LMSNode *left;
    struct LMSNode *right;
} LMSNode;

typedef struct LMS {
    LMSNode *root;
} LMS;

// Function to create a new LMS node
LMSNode* createNode(int key) {
    LMSNode *newNode = (LMSNode*)malloc(sizeof(LMSNode));
    newNode->key = key;
    newNode->left = NULL;
    newNode->right = NULL;
    return newNode;
}

// Function to insert a key into the LMS tree
void insert(LMS *tree, int key) {
    LMSNode *newNode = createNode(key);
    if (tree->root == NULL) {
        tree->root = newNode;
        return;
    }

    LMSNode *current = tree->root;
    LMSNode *parent = NULL;

    while (current != NULL) {
        parent = current;
        if (key < current->key) {
            current = current->left;
        } else {
            current = current->right;
        }
    }

    if (key < parent->key) {
        parent->left = newNode;
    } else {
        parent->right = newNode;
    }
}

// Function to search for a key in the LMS tree
LMSNode* search(LMS *tree, int key) {
    LMSNode *current = tree->root;
    while (current != NULL) {
        if (key == current->key) {
            return current;
        } else if (key < current->key) {
            current = current->left;
        } else {
            current = current->right;
        }
    }
    return NULL; // Key not found
}

// Function to find the minimum value node in the LMS tree
LMSNode* findMin(LMSNode *node) {
    while (node->left != NULL) {
        node = node->left;
    }
    return node;
}

// Function to delete a key from the LMS tree
LMSNode* deleteNode(LMSNode *root, int key) {
    if (root == NULL) {
        return root;
    }

    if (key < root->key) {
        root->left = deleteNode(root->left, key);
    } else if (key > root->key) {
        root->right = deleteNode(root->right, key);
    } else {
        // Node with only one child or no child
        if (root->left == NULL) {
            LMSNode *temp = root->right;
            free(root);
            return temp;
        } else if (root->right == NULL) {
            LMSNode *temp = root->left;
            free(root);
            return temp;
        }

        // Node with two children: Get the inorder successor (smallest in the right subtree)
        LMSNode *temp = findMin(root->right);
        root->key = temp->key; // Copy the inorder successor's content to this node
        root->right = deleteNode(root->right, temp->key); // Delete the inorder successor
    }
    return root;
}

// Function to print the LMS tree in-order
void inOrder(LMSNode *node) {
    if (node != NULL) {
        inOrder(node->left);
        printf("%d ", node->key);
        inOrder(node->right);
    }
}

// Main function to demonstrate the LMS tree
int main() {
    LMS tree;
    tree.root = NULL;

    // Insert keys into the LMS tree
    insert(&tree, 50);
    insert(&tree, 30);
    insert(&tree, 20);
    insert(&tree, 40);
    insert(&tree, 70);
    insert(&tree, 60);
    insert(&tree, 80);

    // Print the in-order traversal of the LMS tree
    printf("In-order traversal of the LMS tree: ");
    inOrder(tree.root);
    printf("\n");

    // Search for a key
    int keyToSearch = 40;
    LMSNode *foundNode = search(&tree, keyToSearch);
    if (foundNode) {
        printf("Key %d found in the LMS tree.\n", keyToSearch);
    } else {
        printf("Key %d not found in the LMS tree.\n", keyToSearch);
    }

    // Delete a key
    int keyToDelete = 20;
    tree.root = deleteNode(tree.root, keyToDelete);
    printf("In-order traversal after deleting %d: ", keyToDelete);
    inOrder(tree.root);
    printf("\n");

    return 0;
}
