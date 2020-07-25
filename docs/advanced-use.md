# Advanced Use

## Initialize A Graph

Now that Gateway is downloaded \(or built\), let's create our database. `init` is the subcommand to set up a database and the right indices.

You can set up a full [configuration file](configuration.md) if you'd prefer, but it will also work from the command line.

Examples for each backend can be found in `store.address` format from [config file](configuration.md).

Those two options \(db and dbpath\) are always going to be present. If you feel like not repeating yourself, setting up a configuration file for your backend might be something to do now. There's an example file, `gateway_example.yml` in the root directory.

You can repeat the `--db (-i)` and `--dbpath (-a)` flags from here forward instead of the config flag, but let's assume you created `gateway_overview.yml`

Note: when you specify parameters in the config file the config flags \(command line arguments\) are ignored.

## Load Data Into A Graph

After the database is initialized we load the data.

```bash
./gateway load -c gateway_overview.yml -i data/testdata.nq
```

And wait. It will load. If you'd like to watch it load, you can run

```bash
./gateway load -c gateway_overview.yml -i data/testdata.nq --alsologtostderr=true
```

And watch the log output go by.

If you plan to import a large dataset into Gateway and try multiple backends, it makes sense to first convert the dataset to Gateway-specific binary format by running:

```bash
./gateway conv -i dataset.nq.gz -o dataset.pq.gz
```

This will minimize parsing overhead on future imports and will compress dataset a bit better.

## Connect a REPL To Your Graph

Now it's loaded. We can use Gateway now to connect to the graph. As you might have guessed, that command is:

```bash
./gateway repl -c gateway_overview.yml
```

Where you'll be given a `gateway>` prompt. It's expecting Gizmo/JS, but that can also be configured with a flag.

New nodes and links can be added with the following command:

```bash
gateway> :a subject predicate object label .
```

Removing links works similarly:

```bash
gateway> :d subject predicate object .
```

This is great for testing, and ultimately also for scripting, but the real workhorse is the next step.

Go ahead and give it a try:

```text
// Simple math
gateway> 2 + 2

// JavaScript syntax
gateway> x = 2 * 8
gateway> x

// See all the entities in this small follow graph.
gateway> graph.Vertex().All()

// See only dani.
gateway> graph.Vertex("<dani>").All()

// See who dani follows.
gateway> graph.Vertex("<dani>").Out("<follows>").All()
```

## Serve Your Graph

Just as before:

```bash
./gateway http -c gateway_overview.yml
```

And you'll see a message not unlike

```bash
listening on :64210, web interface at http://localhost:64210
```

If you visit that address \(often, [http://localhost:64210](http://localhost:64210)\) you'll see the full web interface and also have a graph ready to serve queries via the [HTTP API](http.md)

### Access from other machines

When you want to reach the API or UI from another machine in the network you need to specify the host argument:

```bash
./gateway http --config=gateway.cfg.overview --host=0.0.0.0:64210
```

This makes it listen on all interfaces. You can also give it the specific the IP address you want Gateway to bind to.

**Warning**: for security reasons you might not want to do this on a public accessible machine.

