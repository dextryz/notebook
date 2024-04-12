# Notebook

- A plain text note-taking assistant
- Leverage [Nostr](www.nostr.com) as an external database
- View your nostr published notes on [Ixian](https://ixian.me)

## Setup Nostr

*Notebook* interacts with nostr via long-form notes, as depicted in [NIP-23](https://github.com/nostr-protocol/nips/blob/master/23.md).

If you want to use [nostr](www.nostr.com) as an external database to store and propagate your notes, you have to create a config file and set the `NOSTR` env var:

Create your config file in `~/.config/nostr/dextryz.json` containing:

```
{
    "nsec": "nsec..."
    "relays": ["wss://relay.highlighter.com/", "wss://relay.damus.io/"],
}
```

and set the env var

```shell
> export NOSTR=~/.config/nostr/dextryz.json`
```

## Initialize a New Notebook

To start you have to initiate a notebook:

```shell
> export NOTEBOOk=/tmp/slipbox
> nz init
```

- Is nostr is setup, `notebook` will populate the `NOTEBOOK` path with all your kind `30023` notes.

## Create a New Note

The following command will create a new markdown file in your notebook directory.

```shell
> nz new
created file in notebook slipbox at:
/tmp/slipbox/202404041212.md
```

- The filename is the article `identifier` as specified in [NIP-23](https://github.com/nostr-protocol/nips/blob/master/23.md).

## Edit a Note

Open the note with your favourite editor (hopefully NeoVim) and update the content.

## Publish a Note to Nostr

If you have nostr setup, you can push your note with a title and set of tags. The note will be published to all relays specified in the `NOSTR` config file.

```shell
> nz push --content /tmp/slipbox/202404051040.md --title "Hello Friend" --tag nostr --tag bitcoin
```

## Search Notes

```
./nb search --title Intuition

./nb search --tag focus
```
