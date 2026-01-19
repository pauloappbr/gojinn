# 🛠️ Installation

**Gojinn** is a plugin for the Caddy web server. To use it, you need a version of the Caddy binary that includes the `http.handlers.gojinn` module compiled.

The easiest and recommended way to do this is using `xcaddy`.

## ✅ Prerequisites

- **Go (Golang):** Version 1.25 or higher installed
- **Terminal Access**

---

## 🚀 Method 1: Using xcaddy (Recommended)

`xcaddy` is the official tool for building custom versions of Caddy.

### 1. Install xcaddy

```bash
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

### 2. Compile Caddy with Gojinn

Run the command below to generate the `caddy` binary in the current directory:

```bash
xcaddy build \
    --with github.com/pauloappbr/gojinn
```

### 3. Verify the installation

Confirm that the module was installed correctly:

```bash
./caddy list-modules | grep gojinn
```

**Expected output:**
```
http.handlers.gojinn
```

## 💻 Method 2: Local Development

If you are contributing to the Gojinn source code or want to test a local unpublished version:

### 1. Clone the repository

```bash
git clone https://github.com/pauloappbr/gojinn.git
```

### 2. Compile with the local version

Use `replace` to point to your local folder:

```bash
xcaddy build \
    --with github.com/pauloappbr/gojinn=../path/to/your/local/repo
```

---

## 📚 Next Steps

Now that you have the Gojinn binary installed, the next step is to create your first function.

👉 [Quickstart (5 min)](../quickstart.md)