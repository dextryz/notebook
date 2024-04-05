# NoteZero

View your nostr articles at [Ixian](https://ixian.me)

- Articles with no **headers** will not be shows.
- Articles with no **identifier** will not be shows.

## TODO

- Implement CLI to publish articles.
- Filter articles that have incorrect format and list them to user.

## Flow

I will have to create a Notebook event on nostr too. Maybe this is like a bookmark or curated list?
Dont need to store the content. Just the notebook ID. Each note in the notebook then should contain the notebook event uuid.

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

If NOSTR env var is set all nostr kind 30023 notes will be pulled into directory

List all articles with their identifier

```shell
> nz list --notebooks
nil
```

To start you have to initiate a notebook

```shell
> export NOSTR=~/.config/nostr/dextryz.json

> nz init --name slipbox --dir /tmp/slipbox
notebook 'slipbox' created at 2023-04-13 in dir '/tmp/slipbox'

> nz list --notebook slipbox
notebook dir: /tmp/slipbox
```

New should create a new file in the current notebook.
THen init it on nostr by creating an event with an empty content.
THis will set the title, tags, etc
THis confusing. I might want to just create a file without any commitment.

```shell
export NOTEBOOK=slipbox
> nz new
created file in notebook slipbox at:
/tmp/slipbox/202404041212.md
```

Update an article via their identifier

```shell
export NOTEBOOK=slipbox
> nz push --content 202402051756.md --title "Fake Knowledge" --tag nostr --tag bitcoin
```
```

## TODO

```shell
> nz pull --title "Fake Knowledge" > 2023.md
```
