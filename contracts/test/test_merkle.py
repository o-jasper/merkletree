import pyethereum
t = pyethereum.tester

from merkle import deth_num, eth_num, path_w_root, MovingRoot, h2

from random import randrange, random

s = t.state()
c = s.contract('merkle_root.se', t.k0)
c2 = s.contract('test/merkle_root_equiv.se', t.k0)

#c2 = s.contract("return(57896044618658097711785492504343953926634992332820282019728792003956564819968+1)", t.k0)
#print(s.send(t.k0, c2, 0))

def randval():
    return randrange(2**64)

# Doesnt test merkle trees, but instead just faked paths.
def test_random_path(n):
    test_path(map(lambda(x): randval(), range(n)))

def test_path(path, expect_root=None):
    root = path_w_root(path)
    if expect_root is not None:
        assert root == expect_root
    r = s.send(t.k0,  c if random()<0.5 else c2, 0, path)
    assert len(r) == 1
    assert eth_num(r[0]) == root

def test_random_interest(n, interest_p=0.5, simulate_p=None, m=6):
    if simulate_p is None:
        simulate_p = 10.0/n
    mr = MovingRoot()
    interesting = []
    for _ in range(n):
        interest = (random() < interest_p)
        got = mr.incorporate(map(lambda _: randval(), range(randrange(m))), interest)
        if interest:
            interesting.append(got)
    root = mr.finalize()
    k = 0
    for el in interesting:
        k += 1
        el.selftest()
        path = el.path()
        assert path_w_root(path) == root.H
        s.mine()
        if random() < simulate_p:
            print("simulating %d" % k)
            test_path(path, root.H)

for _ in range(10):
    s.mine()    
    x, y = randval(), randval()
    assert h2(x,y)  == h2(y,x)
    assert deth_num(eth_num(x)) == x
    assert deth_num(eth_num(-x)) == -x

for _ in range(10):
    s.mine()    
    test_random_path(10)
    

for k in range(1):
    print(k)    
    test_random_interest(600)
