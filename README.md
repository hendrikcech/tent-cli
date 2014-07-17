tent-cli
======

Easily work with the Tent protocol from the command line. Discover URLs, get credentials and manage posts.

## Status
- [X] discovery
- [X] authentication
- [X] query support
- [X] get single posts
- [X] delete posts
- [X] create new posts
- [X] update existing posts
- [X] profile support
- [X] simple post schema support

## Usage
```
Usage:
  tent [command]

Available Commands:
  discover <url>                        Get the meta post that is associated with an url.
  auth [<entity>|<profile_name>]        Authorize a new app.
  create [<type> <content> | <json>]    Create a new post.
  update [<entity>] <post_id> <content> Update an existing post.
  get [<entity>] <post_id> [<version>]  Get a single post.
  query                                 Query the posts feed.
  delete <post_id> [<version_id>]       Delete a post.
  profiles                              Manage entity profiles.
  schemas                               Manage post schemas.

Use "tent help [command]" for more information about that command.
```

Profiles are used to save entity and credential configurations. Create a new profile with `tent profiles add entity https://entity.cupcake.is`. Run `tent auth entity` to add credentials to the profile.
Schemas are currently very simple. They basically just map a string to a post type uri to save you from typing these long uris over and over.
Some schemas are defined by default. ´app´ maps to ´https://tent.io/types/app/v0´ for example.

## Installation
Visit the [releases page](https://github.com/hendrikcech/tent-cli/releases/latest) and download a build for your system. If you're on OS X or Linux, unzip the file and copy `tent-cli` to `/usr/local/bin`.

Or, compile it yourself:
```
go get github.com/hendrikcech/tent-cli && go install github.com/hendrikcech/tent-cli
```


## License
The MIT License (MIT)

Copyright (c) 2014 Hendrik Cech

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
