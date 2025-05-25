# mfa 🔐

TOTP and MFA in Go.

## Installation 🚀

```bash
go install github.com/gipuv/mfa@latest
```

After installation, the `mfa` executable will be placed in your `$GOPATH/bin` or `$HOME/go/bin` directory.

Make sure to add this directory to your `PATH` environment variable to use it globally:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage Examples 💡

### Add or Update a Secret 🔑

```bash
mfa -op add -name github -secret JBSWY3DPEHPK3PXP
```

If the name already exists, the program will prompt you whether to replace the secret.

### Get Current TOTP Code 🎫

```bash
mfa -op get -name github
```

### Interactive Mode 🤝

```bash
mfa github
```

or simply

```bash
mfa
```

The program will then prompt you to enter the name and secret interactively.

## View Database Contents 📂

The tool stores secrets in a SQLite database (`.db` file).
To preview or manage the `.db` file, you can use the following tool:

🔍 **DB Browser for SQLite**
Website: [https://sqlitebrowser.org/dl/](https://sqlitebrowser.org/dl/)
Download and open your `.db` file for easy visual management.

## Notes 📝

* The secret must be a valid Base32 encoded string.
* The default TOTP code validity period is 30 seconds.
