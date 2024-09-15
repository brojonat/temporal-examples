# temporal-examples

A collection of long lived business logic processes implemented with Temporal.

## Auction

Package `auction` provides an example implementation of an auction clearing house.

```bash
# in one terminal, start the HTTP server
./cli auction run-server
# in another terminal run the worker
./cli auction run-worker
# in another terminal start an auction and issue some bids;
# after 20 min you should see a message in the server logs
# indicating the webhook was hit with the auction results.
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

## Escrow

Package `escrow` provides an example implementation of an escrow process.
