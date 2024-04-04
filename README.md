# NoteZero

View your nostr articles at [Ixian](https://ixian.me)

- Articles with no **headers** will not be shows.
- Articles with no **identifier** will not be shows.

## TODO

- Implement CLI to publish articles.
- Filter articles that have incorrect format and list them to user.

## CLI Usage

Title, tag, and publish an article to relays listed in config.

```shell
> nz new 202402051756.md --title "Fake Knowledge" --tag nostr --tag bitcoin
```

List all articles with their identifier

```shell
> nz list
```

Update an article via their identifier

```shell
> nz update /tmp/zk/identifier.md
```

Pull all the articles into a directory

```shell
> nz pull /tmp/zk
```

```shell
> nz pull --title "Fake Knowledge" > 2023.md
```
