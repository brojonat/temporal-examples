# temporal-examples

A collection of long lived business logic processes implemented with Temporal. Build the CLI with `go build -o cli cmd/*`, then simply `./cli --help`.

## Auction

Package `auction` provides an example implementation of an auction clearing house.

```bash
# in one terminal, start the HTTP server
./cli auction run-server
# in another terminal run the worker
./cli auction run-worker
# in another terminal start an auction, issue some bids, and check the results
# after 20 min you should see a message in the server logs
# indicating the webhook was hit with the auction results.
./cli auction start --item foo --reserve-price 25 --duration 20m
./cli auction bid --item foo --email me@email.com --amount 50
./cli auction get-state --item foo
```

## Poll

Package `poll` provides an example implementation of a simple poll.

```bash
# in one terminal, start the HTTP server
./cli poll run-server
# in another terminal run the worker
./cli poll run-worker
# in another terminal start a poll and issue some votes;
# after 20 min you should see a message in the server logs
# indicating the webhook was hit with the poll results.
./cli poll start --poll foo --duration 20m -o "option 1" -o "option 2"
./cli poll vote --poll foo -o "option 1"
./cli poll get-state --poll foo
```

## Dead Man's Switch

Package `dms` provides an example implementation of a [Dead man's Switch](https://en.wikipedia.org/wiki/Dead_man%27s_switch).

```bash
# in one terminal, start the HTTP server
./cli dms run-server
# in another terminal run the worker
./cli dms run-worker
# in another terminal start a DMS and query the state;
# after 20 min you should see a message in the server logs
# indicating the webhook was hit with the DMS timeout.
./cli dms start --id foo --duration 20m --message 'oh no, switch timed out!'
./cli dms get-state --id foo
./cli dms deactivate --id foo
```

## Tontine

[TODO] Package `tontine` provides an example implementation of a [tontine](https://en.wikipedia.org/wiki/Tontine).

## Lottery

[TODO] Package `lotto` provides an example implementation of a lottery.

## Escrow

[TODO] Package `escrow` provides an example implementation of an escrow process.

## Prometheus

Package `prom` provides an example implementation of a workflow that emits Prometheus metrics. This is handy if you're interested in instrumenting your Temporal workers.

```bash
# in one terminal, start the HTTP server
./cli prom run-server
# in another terminal run the worker
./cli prom run-worker
# in another terminal, start the workflow
./cli prom start
# finally, you can hit the metrics endpoint and see your prometheus metrics
curl localhost:9090
```

## Activity Heartbeats and Continue-As-New

Package `heart` provides and example implementation of a workflow with a very long running activity. When working with such Activities, you need to emit heartbeats to indicate to the Workflow that the Activity process isn't dead. Under these circumstances, you may also want to use "Continue As New" to avoid history/memory overflow issues. This package demonstrates how to do both.

```bash
# in one terminal, start the HTTP server
./cli heart run-server
# in another terminal run the worker
./cli heart run-worker
# in another terminal, start the workflow
./cli heart start
# finally, for this one, there's not fancy result, but you can open the
# temporal dashboard and see the activity running and eventually see the
# workflow continuing as a new.
```
