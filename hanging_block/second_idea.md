Presumed that a before- and after-state is proved in by Patricia/Merkle
proofs, and the transactions existence in the HB is proven too.

A challenge disqualifying a block is successful if:
    
* Indeed all before- state used was proven.
* After-state is indeed different from the aledged after state. Or some arbitrary
  rule is violated, like gas-use-maximum being exceeded.

## Second approach (unfinished)
Found that gas is an issue, it is undesirable to need vetting before a contract
is accepted. Re-implementing gas in EVM, i.e. by adding a lot of adding operations,
is costly and messy. Afaik the only way to get it right is to run the same code
in auditing and checking a proof on the Ethereum HB side.

Then you simply measure gas, and any differences between HB and checking-HB are
accounted for.

It is necessary for the HB-eth to dictate some aspects of what the checking
contracts look like. Preferably the HB contract manipulates the code to find the
HB version itself. Have not played with `create` enough to be entirely sure it
can be done, but without this, you basically have to accept every contract into
the system. We want: 

1. Prepended part that allows HB to prepare it by filling in storage that was 
   Patricia-proven. The gas cost of this section must be known.
   (or duplicated on-HB)

2. Each HB-contract has multiple instances. Storage index is changed so the
   different instances are separate. Additionally.

   In storage, we need to know if a slot was properly proven. Initially i thought
   using the last 64 bits for indicating the HB-block number, but i really want 
   everything unadulterated Ethereum contracts. So instead, we'll just double 
   the blocks, the second one contains the HB block number. First argument
   of `call`s will contain this number.
   
   Unfortunately, we need it tested on the HB-contract and it doesnt need to be
   tested on the HB-program running. Maybe it uses the same gas, then not we need
   to account for it in gas price checking. Maybe not though, might be better to
   account for it.

3. Calls stay in HB-contracts. Call address is deconstructed to have that HB contract
   address and an *instance* address.
   
   They're modified to send back to a HB contract. The initial arguments need to be
   modified;
   
   * The address of the contract, aswel an an 'instance' address.
   * HB block number to compare.
   * Keep track of number of storage ops.

4. Returning values is done plainly, but then the `sload`-of-not-proven
   is handled by spinning off all gas.
   
   Or it prepends a didnt-do`sload`-not-proven to `return`s and returns
   `"unproven sload"` if it comes across on.

5. `address` needs to be modified to return the combination of the instance and
   Ethereum address against.

### Problems

* Are the modifications are possible?

* Are the `call` modifications doable.. How to get a large enough memory-spot? Just
  ditch the existing one, copy it to end, auch...
