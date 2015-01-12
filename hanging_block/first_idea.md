[Hanging_blocks](https://o-jasper.github.io/blog/ethereum,/hanging/blocks,/blockchain,/scalability/2014/06/03/hanging_blocks.html), 
but it is nicer if regular ethereum contracts can apply.

## *First* idea:

* Turn `sstore`s/`sload`s into calls-to-HB-contract to other contracts, submit to ethereum
* HB program turns the calls-to-HB back into sload/sstore

  The HB program seeing the contract accepted on the HB puts said contract in HB.

Auditing works by:

1. Check root to block. (otherwise you're at the 'data availability' problem)
2. Blocks consist of `Tx, Patricia root` pairs the root is repeated after
   each transaction. Run the code, and check all the roots.

If the roots do not fit, need to construct a proof it is wrong.

1. Construct the proof that `Root before, Tx, Root after` is in the block.
2. Re-run the transaction, now keeping track of what was set/accessed. Create a Patricia
   proof of the before accessed state, and then for one of the wrong after-state storage
   slots, a proof what the *claimed* value was.

Then you have something to submit to the HB contract. To check it, firstly the
HB contract checks the Merkle and Patricia proofs. Then it `sstore`s(yep.. costly? :/)
uses the Tx to whatever initial state was accessed.

Then the HB contract calls the HB-version of the contract that received the message.
In the HB version, `sstore`/`sload` is converted to calls to the HB contract. The HB
contract emulates `sstore`/`sload`. When all is said and done, it compares the
Patricia-proven end state value with the value at the end of running the contract.

If indeed wrong, the HB invalidates the block and subsequent blocks.(auch)
If correct, nothing happens. (Of course this can be simulated)

### Problems

1. The accessed state might be large, making proving it wrong too expensive.
2. `sstore`-ing a lot for state is costly.
3. `call`ing emulating `sstore` is costly, and it needs to emulate storage over
  multiple contracts.(more costly yet)
4. Still need to figure out how to have multiple instances of a contract on the
  HB? How does calling other contracts work?

## Other approaches?
Maybe a solution is to really compile checking code right into the
Ethereum-HB-end. HB program needs to get at the contract -as-it-runs there
aswel though.

(4) can be solved:

1. Simply not have them.. Each on-HB contract corresponds to one on-Ethereum.
    
   Monoliths that store the process of multiple.(not very handy sometimes..)
2. Mess with code more complicatedly, also altering `address`, messing with
   `calldataload` to stuff in the address somewhere, messing with `call` to
   correspond aswel.

Perhaps (3) can be solved, by converting the `sstore`/`sload`s locally
in-contract to specify the actual contract, and adding some extra
functionality allowing the HB contract to look at the relevant value
afterward.
