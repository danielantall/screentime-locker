# Screentime-locker

Set limits you can't circumvent.
Sets a macOS Screen Time passcode **you can't remember**.

Our phones can be addictive. Setting up screen time is great to set limits on our phone usage, but constantly pressing "Remind me in 15 minutes" is unproductive. Screetime-locker makes it a extremely long process to reset your passcode, making it easier for you to lock-in.

![demo](./demo.gif)

Generates a random 4-digit code, then walks you through typing it into System Settings one digit at a time, with fake digits and deletions mixed in so you never see the full passcode. Two rounds (enter + confirm), then it's gone from memory. Nothing is stored or logged.

## Run

```bash
brew install go
go run .
```

Open **System Settings → Screen Time → Lock Screen Time Settings** before starting.

## Reset

Change with iCloud Screen Time reset.


