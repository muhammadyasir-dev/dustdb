#ifndef RING_HASH_H
#define RING_HASH_H

#include <string>
#include <map>
#include <vector>
#include <functional>

class RingHash {
public:
    RingHash(int replicas);
    void addNode(const std::string& node);
    std::string getNode(const std::string& key);

private:
    int replicas; // Number of virtual nodes per real node
    std::map<size_t, std::string> ring; // Hash ring
    std::vector<std::string> nodes; // Real nodes

    size_t hash(const std::string& key);
};

#endif // RING_HASH_H
