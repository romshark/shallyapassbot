# shallyapassbot
Telegram anti-spam bot

## Usage
Build the bot executable:
```bash
cd cmd/shallyapassbot; go build
```

Run the bot executable:

```bash
TOKEN=<accesstoken> ./shallyapassbot [-admins <username>,...] [-logfmt (json|console)] [-debug] [-confirm-timeout <duration>] [-ban-period <duration>] [-text-fmt-welcome <template_string>] [-text-confirm <template_string>]
```

Example:

```bash
go build && TOKEN=7566345343:jKkL21xmsY_Qs_gHx60wqncQyL225sl2Y50 ./shallyapassbot -admins adminA,adminB -confirm-timeout 30s -ban-period 1m -debug -text-fmt-welcome 'Welcome, [{name}](tg://user?id={id})\\! Please confirm you are human by clicking the button below within {confirm-timeout}' -text-confirm 'Confirm'
```
