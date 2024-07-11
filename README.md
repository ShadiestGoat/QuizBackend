# Quiz backend

This is a backend for a quiz.

## Features

- Configurable using a sections file
- Different finales for different keys
- Different section types (slide, question)
- Input validation
- Database-less deployment
- Easy markdown finale syntax
- Custom user flow based on different answers

## Configuration

### Sections file

This file is used to generate each *section*. A section is a screen shown to the user. A sections file has a schema thats stored in the `schema.json` file. An example file is in the `sections.example.yaml`

### Finale pages

The finale docs got a bit too large to keep in the README, [the docs for this are stored in the docs/finale.md file](docs/finale.md)

### Flags

|      flag(s)       | allowed multiple times? | default         | description                                                                                                                                                                        |
| :----------------: | :---------------------: | :-------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  `-f`, `--finale`  |           Yes           | "finale"        | The directory (or directories if specified multiple times) that finale files are stored in. If some files in the directories share the same name, then the latest one overrides it |
|   `-e`, `--env`    |           Yes           | NONE            | A list of `.env` files                                                                                                                                                             |
| `-s`, `--sections` |           No            | "sections.yaml" | Path to the sections file                                                                                                                                                          |

### ENV vars

Env vars can be loaded from `.env` file(s) (see `-e` flag)

|    var name     |      default       | description                                                                                                  |
| :-------------: | :----------------: | :----------------------------------------------------------------------------------------------------------- |
|      PORT       |        3000        | The port that the server runs on                                                                             |
| AUTH_SECRET_KEY | NONE, But required | The encryption key. Used for AES, so should be 16, 24 or 32 bytes. This is a secret, so store it accordingly |
