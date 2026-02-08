# Diagnostics Registry

## Lexer

| Code     | Severity    | Summary                     |
| :------: | :---------: | :-------------------------: |
| `LEX001` | **warning** | unexpected character        |
| `LEX002` | **error**   | unterminated string literal |
| `LEX003` | **error**   | unterminated block comment  |

## Parser

| Code     | Severity    | Summary                                                  |
| :------: | :---------: | :------------------------------------------------------: |
| `PAR001` | **error**   | unexpected token                                         |
| `PAR002` | **error**   | expected class name                                      |
| `PAR003` | **error**   | expected class body or semicolon                         |
| `PAR004` | **error**   | expected closing brace for class body                    |
| `PAR005` | **error**   | missing semicolon after class declaration                |
| `PAR006` | **error**   | expected name after delete                               |
| `PAR007` | **error**   | missing semicolon after delete declaration               |
| `PAR008` | **error**   | expected extern declaration name                         |
| `PAR009` | **error**   | missing semicolon after extern declaration               |
| `PAR010` | **error**   | expected right bracket in array assignment target        |
| `PAR011` | **error**   | expected assignment operator for array assignment        |
| `PAR012` | **error**   | missing semicolon after array assignment                 |
| `PAR013` | **error**   | expected assignment operator                             |
| `PAR014` | **error**   | missing semicolon after assignment                       |
| `PAR015` | **error**   | expected value before end of file                        |
| `PAR016` | **error**   | expected value                                           |
| `PAR017` | **error**   | expected scalar value                                    |
| `PAR018` | **error**   | unterminated array literal                               |
| `PAR019` | **error**   | expected comma or right brace in array literal           |
| `PAR025` | **error**   | strict mode rejects class-like names starting with digit |
| `PAR020` | **warning** | autofix inserted missing class semicolon                 |
| `PAR021` | **error**   | expected enum body                                       |
| `PAR022` | **error**   | expected enum item name                                  |
| `PAR023` | **error**   | expected enum item delimiter                             |
| `PAR024` | **error**   | missing semicolon after enum declaration                 |

## Processor

| Code    | Severity    | Summary                                            |
| :-----: | :---------: | :------------------------------------------------: |
| `PP001` | **error**   | include target not found or unreadable             |
| `PP002` | **error**   | unsupported config intrinsic                       |
| `PP003` | **error**   | macro expansion failure                            |
| `PP004` | **error**   | unterminated conditional block                     |
| `PP005` | **error**   | invalid include syntax                             |
| `PP006` | **error**   | unexpected elif directive                          |
| `PP007` | **error**   | unexpected else directive                          |
| `PP008` | **error**   | unexpected endif directive                         |
| `PP009` | **error**   | error directive triggered                          |
| `PP011` | **error**   | unsupported directive                              |
| `PP012` | **error**   | missing macro name in define                       |
| `PP013` | **error**   | invalid macro name in define                       |
| `PP014` | **error**   | unterminated macro parameter list                  |
| `PP015` | **warning** | macro redefinition                                 |
| `PP016` | **error**   | unresolved macro-like invocation                   |
| `PP017` | **error**   | unsupported __has_include in conditional directive |
