import pyethereum
t = pyethereum.tester

import merkle
eth_num = merkle.eth_num
path_root = merkle.path_root

from random import randrange

s = t.state()
c = s.contract('merkle_root.se', t.k0)
c2 = s.contract('test/merkle_root_equiv.se', t.k0)

#c2 = s.contract("return(57896044618658097711785492504343953926634992332820282019728792003956564819968+1)", t.k0)
#print(s.send(t.k0, c2, 0))

def randval():
    return randrange(2**64)

# Doesnt test merkle trees, but instead just faked paths.
def test_random_path(n):
    path = map(lambda(x): randval(), range(n))
    root = path_root(path)
    r = s.send(t.k0,  c, 0, path)
    r2 = s.send(t.k0, c2, 0, path)
    assert len(r) == 1 and len(r2) == 1

    assert eth_num(r[0]) == root and eth_num(r2[0]) == root


for i in range(10):
    test_random_path(10)
    
