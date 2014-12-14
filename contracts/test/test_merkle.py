import pyethereum
t = pyethereum.tester

from merkle import eth_num, path_w_root, MovingRoot, h2

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
    path = map(lambda(x): randval(), range(n))
    root = path_w_root(path)
    r = s.send(t.k0,  c, 0, path)
    r2 = s.send(t.k0, c2, 0, path)
    assert len(r) == 1 and len(r2) == 1

    assert eth_num(r[0]) == root and eth_num(r2[0]) == root

def test_random_interest(n, interest_p=0.5, m=6):
    mr = MovingRoot()
    interesting = []
    for _ in range(n):
        interest = (random() < interest_p)
        got = mr.incorporate(map(lambda _: randval(), range(randrange(m))), interest)
        if interest:
            interesting.append(got)
    root = mr.finalize()
    for el in interesting:
        assert path_w_root(el.path()[-1:]) == root

for _ in range(10):
    test_random_interest(60)

for _ in range(10):
    x, y = randval(), randval()
    assert h2(x,y)  == h2(y,x)

for _ in range(10):
    test_random_path(10)
    
