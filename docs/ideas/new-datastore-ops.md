# Datastore Path Support

## Design: Extend Existing Operations with Optional Paths

Instead of adding new builtins, extend existing operations to accept an array `[key, path]` as the first parameter for nested access.

**Pattern:**
```duso
// Access top-level
ds.get("counter")
ds.increment("counter", 1)

// Access nested field with dot-notation path
ds.get(["config", "db.host"])
ds.increment(["stats", "counters.requests"], 1)
ds.update(["config", "db"], {port=5432})
ds.wait(["config", "db.host"], "localhost", 5)
```

Benefits:
- No new builtins, consistent pattern across all operations
- Clear intent: array selector means "navigate this path within the key"
- Extensible for future operations
- Simple implementation: split path by dots, traverse nested objects

---

## Supported Operations with Paths

### `get([key, path])`
Get a nested field without copying entire object.

```duso
ds.set("config", {db={host="localhost", port=5432}, cache={ttl=3600}})

host = ds.get(["config", "db.host"])  // "localhost"
ttl = ds.get(["config", "cache.ttl"])  // 3600
missing = ds.get(["config", "db.ssl"])  // nil
```

### `increment([key, path] [, delta])`
Atomically increment a nested numeric field.

```duso
ds.set("stats", {counters={requests=0, errors=0}})

count = ds.increment(["stats", "counters.requests"])  // 1
count = ds.increment(["stats", "counters.requests"], 10)  // 11
```

### `update([key, path], updates)`
Update nested object directly.

```duso
ds.set("config", {db={host="localhost", port=5432}})

ds.update(["config", "db"], {port=3306, ssl=true})
// Result: {db={host="localhost", port=3306, ssl=true}}
```

### `wait([key, path] [, value] [, timeout])`
Wait on changes to a nested field.

```duso
ds.wait(["config", "db.host"], "new-host", 5)  // Block until db.host becomes "new-host"
```

---

## Implementation Notes

- Path traversal creates intermediate objects if missing (like `update()`)
- Non-existent paths return nil (for `get()`) or error (for operations expecting existing field)
- Dot notation supports arbitrary depth: `"a.b.c.d.e"`
- Each operation's atomicity applies to the entire key, not the path
- Broadcasts still target the key, watchers see key-level changes
