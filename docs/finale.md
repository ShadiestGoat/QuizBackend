# Finale config

Finales are the final section that is shown to a user when they have completed your quiz. The catch is that different people can get different finales through per-user (/per-group) *keys*. A *key* is a secret that is sent to the user, which authenticates them to view the finale

Finale sections are generated using a *finales config folder*, specified using the `-f` flag ([see flag descriptions in the README](../README.md#flags)). All *finale files* must end with `.md`.

The *finales config folder* consists of other folders. The name of each *finale folder* can take 3 forms:

1. `$[name]-[key]`
2. `$[name]`
3. `[key]`

Where square brackets are not needed. The *key* is describe previously. The *name* is a quick alias used for your viewing pleasure only. For obvious reasons, `-` are not allowed in names. Additionally, names are case-insensitive (keys are not!). Any name overlaps will be detected and the app will fail to start.
In the 2nd form, `$[name]` a key is missing. In order to give it a key, you must pass an env variable `KEYS_[NAME]`. (eg. for a folder named `$foobar`, the env var must be `KEYS_FOOBAR`)

Finally, there is a special name - `$default`. This is the default finale - what is shown to a user without a key. This *finale folder* is optional - but if not provided users without a key would go through the whole quiz without reward.

Each *finale folder* must contain at least 1 file - `essay.md` or `faq.md`. Both can be specified. They have slight differences in structure.

## Essay File

The essay file is pretty simple - its just a markdown file. It has no special rules, all headings are supported - add at `h1` headings if you want to

## FAQ File

This file has a special setup. It consists of `h1` elements (`# Heading 1`), which are interpreted as questions, and text elements, which are answers to the question in the `h1`. The text elements don't have to only include inline content, ie. they can also include headings (starting with `h2`, since `h1` is interpreted as a question). Any heading that is in text portion has 1 level added to it - so `## Heading` is interpreted as `# Heading`.

```md
# Here is a question

Here is an answer


# Here is another question

Here is another answer
And another line

And a 3d line!

## And a subheading
```

This file would be interpreted as the following:

```yaml
- question: Here is a question
  answer: Here is an answer

- question: Here is another question
  answer: -|
    Here is another answer
    And another line
    And a 3d line!
    \# And a subheading
```

## Finale file interpretation

Some markdown interpreters require 2 newlines to create a newline in the result. For removing this ambiguity, the app force parses any consecutive newlines as 1 newline. Additionally, extra whitespace **at the end** of lines is removed.
