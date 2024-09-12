# temporal-examples

A collection of long lived business logic processes implemented with Temporal.

## Auction

Package `auction` provides an example implementation of an auction clearing house.

```bash
./cli auction run-server
./cli auction run-worker
./cli auction start-auction --item foo --reserve-price 25 --duration 20m
./cli action place-bid --item foo --email me@email.com --amount 50
```

## Dead Man's Switch

Package `dms` provides an example implementation of a [Dead man's Switch](https://en.wikipedia.org/wiki/Dead_man%27s_switch).

## Tontine

Package `tontine` provides an example implementation of a [tontine](https://en.wikipedia.org/wiki/Tontine).

## Poll

Package `poll` provides an example implementation of a simple poll.

## Lottery

Package `lotto` provides an example implementation of a lottery.
