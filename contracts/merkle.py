import pyethereum
t = pyethereum.tester  # For t.sha3

# NOTE: t.sha3 is from this branch https://github.com/ethereum/pyethereum/pull/176

def eth_num(x):  # TODO wrong way around? (should be fixing t.sha3?)
    return x + 2**256 if x < 0 else x

def deth_num(x):
    assert x < 2**256
    return x - 2**256 if x > 2**255 else x

def eth_lt(x,y):
    return deth_num(x) <= deth_num(y)

def h2(x, y):
    return t.sha3([x,y] if eth_lt(x, y) else [y,x])

# Root from a path.
def path_w_root(sides):
    cur = sides[0]
    for el in sides[1:]:
        cur = h2(cur,el)
    return cur

# Element for moving root, with option of indicating interest(or lack)
class Element:
    def __init__(self, H, depth, up, left, right, data=None):
        self.H = H
        self.depth = depth
        self.up = up
        self.right = right
        self.left = left
        self.data = data
        if self.right is not None:
            assert self.left is not None
            assert eth_lt(self.left.H, self.right.H)

    def selftest(self):
        if self.right is not None:
            assert self.left is not None
            
            assert self.H == t.sha3([self.left.H, self.right.H])
            assert eth_lt(self.left.H, self.right.H)
        else:
            assert self.left is None
            
        if self.up is not True:
            assert (self is self.up.left) or (self is self.up.right)
            self.up.selftest()
        
    def el_path(self):
        return [self] + self.up._el_path(self)

    def _el_path(self, avoid):
        assert self.up  # Will only be a path to track if there is interest.
        assert (self.left is avoid) or (self.right is avoid)
        ret = [self.left if (self.right is avoid) else self.right]
        if self.up is True: # Top.
            return ret
        else:
            return ret + self.up._el_path(self)

    def path(self):
        return map(lambda(el):el.H, self.el_path())


def Element_from_two(x, y):
    l,r = ((x, y) if eth_lt(x.H, y.H) else (y, x))
    assert eth_lt(l.H, r.H), (eth_lt(x.H, y.H), eth_lt(y.H, x.H), x.H,y.H)
    return Element(t.sha3([l.H, r.H]), max(l.depth, r.depth) + 1, l.up or r.up, l, r)

# Moving root.
class MovingRoot:
    def __init__(self, _list=None):
        self.list = ([] if _list is None else _list)  # List of top-depth pairs.

    def copy(self):
        return MovingRoot(self.list)

    def incorporate(self, chunk, interest=False):
        el = Element(t.sha3(chunk), 0, interest, None, None, chunk)
        self.list.append(el)
        self.simplify()
        return el

    # Forcing makes "lobsized" trees with uneven depth.(for if element count != 2**n)
    def simplify(self, force=False):
        while len(self.list) >= 2:
            last = self.list[-1]
            prelast =  self.list[-2]
            if last.depth != prelast.depth and not force:  # Not equal depth, done..
                return
            
            # Combine equal depth tree tops.
            self.list = self.list[:-2]
            el = Element_from_two(prelast, last)
            self.list.append(el)
            if last.up:  # Refer on if interested.
                assert last.up is True
                last.up = el
            if prelast.up:
                assert prelast.up is True
                prelast.up = el

    # Force-combine lists, returning the current result.
    def finalize(self, dont_change=False):
        if dont_change:
            return self.copy().finalize()
        
        assert len(self.list) > 0
        self.simplify(True)
        assert len(self.list) == 1
        return self.list[0]
