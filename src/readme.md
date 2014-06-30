## merkletree/
The lib itself.

## common/
Just some shit i want to share between binary and testing sources..
(**TODO**.. better way of including that?)

## test_merkletree.go
Generates a merkle tree from chunks, randomly including/not including chunks as
interesting, and then testing the ones marked as interesting.(for those the tree
is complete enough to generate Merkle paths)

**TODO** test that if you mess with any bit of data, it goes wrong.
(Note: tests merely increase the surface area for detection of bugs. It might
miss a corner case. Review would be good.)

## path_chunk_n_root.go
Generates random data, outputting a *single* chunk, with the single merkle path,
and the merkle root. For use with commandline stuff or something.
