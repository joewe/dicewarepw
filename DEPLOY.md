# Deploying DicewarePW to Uberspace

This guide describes how to deploy the DicewarePW application to Uberspace.

## Prerequisites

- An Uberspace account.
- SSH access to your Uberspace.

## 1. Build the Application

Since Uberspace runs on Linux, you need to cross-compile the application if you are on a different OS (like macOS or Windows).

Run the following command in the root of your project:

```bash
GOOS=linux GOARCH=amd64 go build -o dicewarepw ./cmd/web
```

This will create a binary named `dicewarepw`.

## 2. Upload to Uberspace

Upload the binary and the `ui` directory to your Uberspace. You can use `scp` or `rsync`.

Replace `youruser` and `yourhost` with your actual Uberspace username and hostname.

```bash
# Create a directory for the app
ssh youruser@yourhost "mkdir -p ~/dicewarepw"

# Upload binary
scp dicewarepw youruser@yourhost:~/dicewarepw/

# Make binary executable
ssh youruser@yourhost "chmod +x ~/dicewarepw/dicewarepw"
```

## 3. Configure Web Backend

Login to your Uberspace via SSH:

```bash
ssh youruser@yourhost
```

Register a new web backend. This will assign a port to your application.

```bash
uberspace web backend set / --http --port <YOUR_PORT>
```
*Note: Replace `<YOUR_PORT>` with a port number between 1024 and 65535. For example, 8080. If the port is taken, try another one.*

## 4. Setup Supervisord

To keep your application running, use `supervisord`.

Create a configuration file:

```bash
nano ~/etc/services.d/dicewarepw.ini
```

Add the following content (replace `<YOUR_PORT>` with the port you chose):

```ini
[program:dicewarepw]
directory=%(ENV_HOME)s/dicewarepw
command=%(ENV_HOME)s/dicewarepw/dicewarepw
autostart=true
autorestart=true
environment=PORT="<YOUR_PORT>"
stdout_logfile=%(ENV_HOME)s/logs/supervisord/dicewarepw.out.log
stderr_logfile=%(ENV_HOME)s/logs/supervisord/dicewarepw.err.log
```

Save and exit (Ctrl+O, Enter, Ctrl+X).

Update supervisord to start the service:

```bash
supervisorctl reread
supervisorctl update
supervisorctl status
```

You should see `dicewarepw RUNNING`.

## 5. Configure Subdomain

If you want to run the app under a subdomain (e.g., `pw.yourdomain.tld`):

1.  **Set up the subdomain**:
    ```bash
    uberspace web domain add pw.yourdomain.tld
    ```

2.  **Configure the backend for the subdomain**:
    ```bash
    uberspace web backend set pw.yourdomain.tld --http --port <YOUR_PORT>
    ```
    *(Make sure to remove the previous root backend if you only want it on the subdomain: `uberspace web backend del /`)*

3.  **DNS**: Ensure the subdomain points to your Uberspace server (A record).

## 6. Verify

Open your browser and visit your domain (or subdomain). The DicewarePW app should be running!
