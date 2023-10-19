# astjson
[![codecov](https://codecov.io/gh/xieyuschen/astjson/graph/badge.svg?token=FFXVZQYUWF)](https://codecov.io/gh/xieyuschen/astjson)

## motivation
Haskell json library [aeson](https://github.com/haskell/aeson) parsed AST 
first in its early version as a default way. However, it skips to convert AST 
first to speed up parsing as a new default way.

Parsing AST first is slower than directly parsing json, however, it provides more flexibilities. 
Hence, inspired by aeson, I wrote a ast json to help parsing json to AST.

## sources
- [json tokens](https://www.json.org/json-en.html)