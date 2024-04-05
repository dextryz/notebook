# Notebook

- A personal notebook on your local machine.
- Leverage nostr as an external database
- View your nostr article notes at [Ixian](https://ixian.me)

## Setup

The current notebook is set by defining the following env vars:

- `NOTEBOOK`: The notebook name
- `NOTEBOOK_DIR`: The path to where the notebook should store markdown files.

## Nostr

If you want to use [nostr](www.nostr.com) as an external database to store and propagate your notes,
you have to create a config file and set the `NOSTR` env var:

Create your config file in `~/.config/nostr/dextryz.json` containing:
```
{
    "nsec": "nsec..."
    "relays": ["wss://relay.highlighter.com/", "wss://relay.damus.io/"],
}
```
and the env var
```shell
export NOSTR=~/.config/nostr/dextryz.json`
```

## Initialize a New Notebook

To start you have to initiate a notebook:

```shell
> nz init --name slipbox --dir /tmp/slipbox
notebook 'slipbox' created at 2023-04-13 in dir '/tmp/slipbox'
```

If you have no nostr account setup then the directory will be empty. However, if nostr is setup, it'll pull all your Kind 30023 notes (articles) into the directory with the filename being that of the article identifier.

The `nz init` command will automatically set your `NOTEBOOK` and `NOTEBOOK_DIR` env vars.

## Create a New Note

The following command will create a new markdown file in your notebook directory.

```shell
> nz new
created file in notebook slipbox at:
/tmp/slipbox/202404041212.md
```

## Publish a Note to Nostr

If you have nostr setup, you can push your note with a title and set of tags. The note will be published to all relays specified in the `NOSTR` config file.

```shell
> nz push --content /tmp/slipbox/202404051040.md --title "Hello Friend" --tag nostr --tag bitcoin
```
