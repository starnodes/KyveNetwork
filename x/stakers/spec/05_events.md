<!--
order: 5
-->

# Events

The `x/stakers` module contains the following events:

## EventCreateStaker

EventBundleProposed indicates that a new staker was created.

```protobuf
message EventCreateStaker {
  // staker is the account address of the protocol node.
  string staker = 1;
  // amount for inital self-delegation
  uint64 amount = 2;
}
```

It gets thrown from the following actions:

- MsgCreateStaker

## EventUpdateMetadata

EventUpdateMetadata is an event emitted when a protocol node updates their
metadata.

```protobuf
message EventUpdateMetadata {
  // staker is the account address of the protocol node.
  string staker = 1;
  // moniker ...
  string moniker = 2;
  // website ...
  string website = 3;
  // logo ...
  string logo = 4;
}
```

It gets thrown from the following actions:

- MsgUpdateMetadata

## EventUpdateCommission

EventUpdateCommission indicates that a staker has changes its commission.

```protobuf
message EventUpdateCommission {
  // staker is the account address of the protocol node.
  string staker = 1;
  // commission ...
  string commission = 2;
}
```

It gets thrown from the following actions:

- EndBlock

## EventJoinPool

EventClaimUploaderRole indicates that a staker has joined a pool.

```protobuf
message EventJoinPool {
  // pool_id is the pool the staker joined
  uint64 pool_id = 1;
  // staker is the address of the staker
  string staker = 2;
  // valaddress is the address of the protocol node which 
  // votes in favor of the staker
  string valaddress = 3;
  // amount is the amount of funds transferred to the valaddress
  uint64 amount = 4;
}
```

It gets thrown from the following actions:

- MsgJoinPool

## EventLeavePool

EventLeavePool indicates that a staker has left a pool.
Either by leaving or by getting kicked out for the following reasons:

- misbehaviour (usually together with a slash)
- all pool slots are taken and a node with more stake joined.

```protobuf
message EventLeavePool {
  // pool_id ...
  uint64 pool_id = 1;
  // staker ...
  string staker = 2;
}
```

It gets thrown from the following actions:

- EndBlock
- bundles/MsgSubmitBundleProposal
- MsgJoinPool