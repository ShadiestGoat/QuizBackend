# Finale file
Finales are generated using markdown files. You'll need to specify a directory in which the finales are stored. All must end with `.md`.

Each file can consist of 2 `h1` headings `# Heading`, though they can repeat interchangeably & the data will be merged in the end. The headings are case-insensitive. The 2 headings are `faq` & `essay`. They represent 2 different parts of the finale.

## The FAQ section

The `faq` section consists of h2 elements (`## Heading 2`), which are interpreted as questions, and text elements, which are answers to the question in the `h2`. For example:

```md
# FAQ

## Here is a question

Here is an answer


## Here is another question

Here is another answer
And another line

And a 3d line!
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
```

## The essay portion

The essay is just a bunch of lines that can be shown the user.

```md
# Essay

This is an essay.
2nd line.

3d line. Thats a crazy one!
```

This will be shown to the user as the following:

```
This is an essay.
2nd line.
3d line. Thats a crazy one!
```

## The default finale

A file named `default_finale.md` is special - its the default finale. You **should** have this, otherwise a user that completes the quiz without a special key that you give to them will get a `401` response.

## Finale file interpretation

Spaces around between `h1`, `h2` and text elements are trimmed. Additionally, `h1` are merged together. So, in the following example:

```md
# Essay

Line 1

## FAQ

### Fake Question

Fake answer ig

# faq

## Real question

Real answer

# Essay

Another essay line

# fAq

## Question

# Essay

More lines here
Yay even more content

# Faq

Content here
```

This would me interpreted as the following yaml:

```yaml
essay: |-
  Line 1
  ## FAQ
  ### Fake Question
  Fake answer ig
  Another essay line
  More lines here
  Yay even more content
faq:
  - question: Real question
    answer: Real answer
```

## Name of the file

The name of the finale file - it signifies the key that a user needs to access this finale. The key is without the `.md` part.

