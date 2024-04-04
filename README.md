# NoteZero

View your nostr articles at [Ixian](https://ixian.me)

- Articles with no **headers** will not be shows.
- Articles with no **identifier** will not be shows.

## TODO

- Implement CLI to publish articles.
- Filter articles that have incorrect format and list them to user.

## Flow

Init creates a new notebook in the given directory. If the NOSTR env var is not 
set, then we assume the user is not using nostr. If it is set, pull the user notes from nostr.
 
A notebook is created when:
    - Explicity instance one using the CLI
    - When you pull all your nostr articles into a defined directory.

A note can only be created within a notebook instance.
An env var defines the current notebook instance path.
If this env is empty, no notebook is set.

This means that pulling nostr artciles shoudl set this env.

How will nostr keep track of notebooks? Hashtag or App specific tag like a TODO list?

We only interact with a notebook. A notebook has to interact with nostr and the underlying database.

## CLI Usage

Title, tag, and publish an article to relays listed in config.

Content can be a filename or a string containing the literally content.

The notebook is set if the NOTEBOOK is specified and NOTEBOOK_DIR

```shell
export NOSTR=~/.config/nostr/dextryz.json
> nz init --name slipbox --dir /tmp/slipbox
```

```shell
export NOTEBOOK=slipbox
> nz new --content 202402051756.md --title "Fake Knowledge" --tag nostr --tag bitcoin
```

List all articles with their identifier

```shell
> nz list --notebook slipbox
```

Update an article via their identifier

```shell
> nz update /tmp/zk/identifier.md
> nz push --content identifier.md
```

Pull all the articles into a directory

```shell
> nz pull /tmp/zk
```

```shell
> nz pull --title "Fake Knowledge" > 2023.md
```
