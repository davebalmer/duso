# Bundling Applications into Custom Duso Binaries

Create a standalone Duso binary with your application scripts embedded. Your users get a single executable with everything built-in—no separate script files, no Duso installation required.

## Why Bundle?

- **Single executable** - Everything in one binary
- **Deploy anywhere** - Works across platforms with cross-compilation
- **Automatic startup** - Runs your script on launch
- **Easy distribution** - Share one file with your team

## Using bundle-duso

The `bundle-duso` script packages your application into a standalone Duso binary.

### Basic Usage

```bash
bundle-duso --bundle-dir <app-dir> --run-script <script> <output-binary>
```

**Example:**

```bash
bundle-duso --bundle-dir ./myapp --run-script main.du myapp
```

This creates `myapp`, a standalone binary that runs `main.du` automatically on launch.

### Your Application Structure

```
myapp/
  main.du            ← Your startup script (any name)
  lib/
    helpers.du
```

### Advanced Options

**Cross-compile:**

```bash
bundle-duso --bundle-dir ./myapp --run-script main.du myapp --target linux/amd64
```

**Custom library path:**

```bash
bundle-duso --bundle-dir ./myapp --run-script main.du myapp --add-lib /EMBED/myapp/lib
```

**Custom Duso repo:**

```bash
DUSO_REPO=/path/to/duso bundle-duso --bundle-dir ./myapp --run-script main.du myapp
```

## How It Works

1. Copies your app directory into the Duso build directory
2. Runs `go generate` to embed files
3. Builds a new Duso binary with your app embedded
4. Sets the startup script via build flags
