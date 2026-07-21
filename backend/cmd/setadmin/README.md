# Set admin Firebase tool

Used to update the user's admin status (i.e. make user admin/remove admin status) on Firebase.

### Prerequisites

- A Firebase user (manually created on the Firebase console)
- `FIREBASE_CREDS` set in `.env` file.

### Usage

Make user admin:

```bash
go run cmd/setadmin/main.go -set-admin=true
```

Remove user admin status:

```bash
go run cmd/setadmin/main.go
```

