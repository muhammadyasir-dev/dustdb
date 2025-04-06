#include "RingHashWrapper.h"
#include "RingHash.h"

void* create_ring_hash(int replicas) {
    return new RingHash(replicas);
}

void add_node(void* ring_hash, const char* node) {
    static_cast<RingHash*>(ring_hash)->addNode(node);
}

const char* get_node(void* ring_hash, const char* key) {
    return const_cast<char*>(static_cast<RingHash*>(ring_hash)->getNode(key).c_str());
}

void destroy_ring_hash(void* ring_hash) {
    delete static_cast<RingHash*>(ring_hash);
}
