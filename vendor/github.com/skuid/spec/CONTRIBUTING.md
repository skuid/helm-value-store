# Contributing

### Commit Message Conventions
Once you’ve implemented a bug fix or feature, please use the following commit
message format. In order to track and summarize the changes, we could use a
changelog automation tool called [changelog](https://github.com/skuid/changelog)
which scrapes information from commit messages. We follow the 'conventional'
commit message format.

As a summary, messages should be formatted like:

```
<type>(<scope>): <subject>
<empty line>
<body>
<empty line>
<footer>
```

####  Type

Type | Purpose
--------|------------
feat | A new feature. Please also link to the issue (in the body) if applicable. Causes a minor version bump.
fix | A bug fix. Please also link to the issue (in the body) if applicable.
docs | A documentation change.
style | A code change that does not affect the meaning of the code, (e.g. indentation).
refactor | A code change that neither fixes a bug or add a feature.
perf | A code change that improves performance.
chore | Changes to build process or auxiliary tools or libraries such as documentation generation.
config | Changes to configurations that have tangible effects on users, (e.g. renaming properties, changing defaults, etc).


#### Scope

The scope of the commit message indicates the area or feature the commit applies
to. For instance, if you were changing something in the `middlewares` package,
your commit message might look something like

```
feat(middlewares): Added a logging middleware
```

Or if you fixed the `redis` package:

```
fix(redis): Added missing colon
```
