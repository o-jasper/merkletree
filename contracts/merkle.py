import pyethereum
t = pyethereum.tester  # For t.sha3

# NOTE: t.sha3 is from this branch https://github.com/ethereum/pyethereum/pull/176

def eth_num(x):  # TODO wrong way around? (should be fixing t.sha3?)
    return x + 2**256 if x < 0 else x

def eth_lt(x,y):
    return x < y if x < 2**255 else x > y

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
    def __init__(self, H, depth, interest):
        self.H = H
        self.depth = depth
        self.interest = interest
        
    def el_path(self):
        assert self.interest  # Will only be a path to track if there is interest.
        if self.interest == True:  # End of the line.
            return [self]
        else:
            return [self] + self.interest.el_path()

    def path(self):
        return map(lambda(el):el.H, self.el_path())


# Moving root.
class MovingRoot:
    def __init__(self, _list=None):
        self.list = ([] if _list is None else _list)  # List of top-depth pairs.

    def copy(self):
        return MovingRoot(self.list)

    def incorporate(self, chunk, interest=False):
        return self.incorporate_H(t.sha3(chunk), interest)

    def incorporate_H(self, H, interest=False):
        el = Element(H, 1, interest)
        self.list.append(el)
        self.simplify()
        return el

    # Forcing makes "lobsized" trees with uneven depth.(for if element count != 2**n)
    def simplify(self, force = False):
        while len(self.list) >= 2:
            last = self.list[-1]
            prelast =  self.list[-2]
            if last.depth != prelast.depth and not force:  # Not equal depth, done..
                return
            # Combine equal depth tree tops.
            self.list = self.list[:-2]
            el = Element(h2(last.H, prelast.H), last.depth + 1,\
                         last.interest or prelast.interest)
            self.list.append(el)
            if last.interest:  # Refer on if interested.
                last.interest = el
            if prelast.interest:
                prelast.interest = el

    # Force-combine lists, returning the current result.
    def finalize(self, dont_change=False):
        if dont_change:
            return self.copy().finalize()
        
        assert len(self.list) > 0
        self.simplify(True)
        assert len(self.list) == 1
        return self.list[0].H
