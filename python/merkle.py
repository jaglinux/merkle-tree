from Crypto.Hash import keccak
import math

left = 0
right = 0

def get_hash(data):
    k = keccak.new(digest_bits=256)
    k.update(data)
    print(k.hexdigest())
    return k.hexdigest()

#root node index is 0, left child is 1 and right child is 2
def left_right_hash(left_hash, right_hash):
    return get_hash(left_hash + right_hash)

def get_left_child(index):
    return (2*index)+1

def get_right_child(index):
    return 2*(index+1)

class Merkle:
    def __init__(self):
        self.leaf_hashes = []
        self.root_hash = None
        self.all_hashes = []

    def get_txn_hashes(self, leaf_txn):
        for i in leaf_txn:
            self.leaf_hashes.append(i)

    def build_tree(self, leaf_hashes):
        leaf_len = len(self.leaf_hashes)
        if len(leaf_len) == 0:
            return
        if not leaf_len % 2:
            self.leaf_hashes.append(self.leaf_hashes[-1])
        leaf_len += 1
        height = math.log2(leaf_len)
        height = math.ceil(height)



