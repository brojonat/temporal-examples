# temporal-examples

A collection of long lived business logic processes implemented with Temporal. Build the CLI with `go build -o cli cmd/*`, then simply `./cli --help`.

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
curl localhost:9090/metrics
# HELP temporal_samples_foo_counter_total temporal_samples_foo_counter_total counter
# TYPE temporal_samples_foo_counter_total counter
temporal_samples_foo_counter_total{activity_type="RunPromActivity",client_name="temporal_go",namespace="default",prom_test_label="foo",task_queue="temporal_examples",worker_type="none",workflow_type="RunPromWF"} 31
# HELP temporal_samples_foo_gauge temporal_samples_foo_gauge gauge
# TYPE temporal_samples_foo_gauge gauge
temporal_samples_foo_gauge{activity_type="RunPromActivity",client_name="temporal_go",namespace="default",prom_test_label="foo",task_queue="temporal_examples",worker_type="none",workflow_type="RunPromWF"} 1.727373674e+09
```

## Activity Heartbeats and Continue-As-New

Package `heart` provides and example implementation of a workflow with a very long running activity. When working with such Activities, you need to emit heartbeats to indicate to the Workflow that the Activity process isn't dead. Similarly, when you have very long running Workflows with lots of events, you may also want to use "Continue As New" to avoid history/memory overflow issues. This package demonstrates how to do both.

```bash
# in one terminal, start the HTTP server
./cli heart run-server
# in another terminal run the worker
./cli heart run-worker
# in another terminal, start the workflow
./cli heart start
# for this one, there's no fancy result to see, but you can open the
# temporal dashboard and see the activity running and eventually see the
# workflow continuing as a new. You can find the next workflow execution
# under the "relationships" tab.
```

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

## Trading Strategy: Market Maker

[TODO] Package `mm` provides an example implementation of a market maker. The workflow places limit bid (ask) orders of a fixed size below (above) the current market price.

## Trading Strategy: The Wheel

[TODO] Package `wheel` provides an example implementation of a derivatives trading strategy known as "the wheel", in which a trader continuously cycles equities for options in that equity. The strategy is delta neutral; it is designed to profit from the premium earned by selling an option that is backed by an underlying.
