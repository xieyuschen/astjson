# astjson

## motivation
Haskell json library [aeson](https://github.com/haskell/aeson) parsed AST 
first in its early version as a default way. However, it skips to convert AST 
first to speed up parsing as a new default way.

The way which parses AST first is slower than the opposite way, however, it enables 
more flexibilities.

## sources
- [json tokens](https://www.json.org/json-en.html)