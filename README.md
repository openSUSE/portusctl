# portusctl

Under construction. To do:

- Integration tests
- Update method
- Validate method
- Refactoring
- Polish exec command

## Status of the API

| Endpoint                   | Portus | portusctl | observations                                  |
|----------------------------|--------|-----------|-----------------------------------------------|
| user #index                | Yes    | Yes       |                                               |
| user #create               | Yes    | Yes       |                                               |
| user #destroy              | Yes    | Yes        |                                               |
| user #show                 | Yes    | Yes       |                                               |
| user #update               | Yes    | No        |                                               |
| application token #index   | Yes    | Yes       |                                               |
| application token #create  | Yes    | Yes       |                                               |
| application token #destroy | Yes    | No        |                                               |
| registries #index          | Yes    | Yes       |                                               |
| registries #create         | No     | No        |                                               |
| registries #destroy        | No     | No        | unsupported                                   |
| registries #show           | No     | Yes       |                                               |
| registries #update         | No     | No        |                                               |
| teams #index               | Yes    | Yes       |                                               |
| teams #create              | Yes    | Yes       |                                               |
| teams #destroy             | No     | No        | unsupported                                   |
| teams #show                | Yes    | Yes       |                                               |
| teams #update              | No     | No        |                                               |
| teams namespaces           | Yes    | Yes       |                                               |
| teams members              | Yes    | Yes       |                                               |
| namespaces #index          | Yes    | Yes       |                                               |
| namespaces #create         | Yes    | Yes       |                                               |
| namespaces #destroy        | No     | No        | unsupported                                   |
| namespaces #show           | Yes    | Yes       |                                               |
| namespaces #update         | No     | No        |                                               |
| namespaces repositories    | Yes    | Yes       |                                               |
| repositories #index        | Yes    | Yes       |                                               |
| repositories #create       | No     | No        |                                               |
| repositories #destroy      | No     | No        | unsupported                                   |
| repositories #show         | Yes    | Yes       |                                               |
| repositories #update       | No     | No        |                                               |
| repositories tags          | Yes    | Yes       |                                               |
| repositories tags groupped | Yes    | Yes       | maybe as an option to the tags in portusctl ? |
| repositories show tag      | Yes    | Yes       |                                               |
| tags #index                | Yes    | Yes       |                                               |
| tags #create               | No     | No        | unsupported                                   |
| tags #destroy              | No     | No        | unsupported                                   |
| tags #show                 | Yes    | Yes       |                                               |
| tags #update               | No     | No        | unsupported                                   |

Missing resources on both sides:

- Webhooks
- Search and explore
- Ping and health
- Activities

Final considerations:

- What about typeahead kind of queries ? Should they be put in some special
  explore command ?

## License

This project is based on work we did for the
[Portus](https://github.com/SUSE/Portus) project. However, all the code has been
re-written from scratch, so the entire project is subject to the GPLv3 license:

```
Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
