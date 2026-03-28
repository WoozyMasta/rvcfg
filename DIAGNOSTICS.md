# Diagnostics Registry

## Lexer

| Code   | Severity    | Description                 |
| :----: | :---------: | :-------------------------: |
| `1001` | **warning** | unexpected character        |
| `1002` | **error**   | unterminated string literal |
| `1003` | **error**   | unterminated block comment  |

## Parser

| Code   | Severity    | Description                                                                           |
| :----: | :---------: | :-----------------------------------------------------------------------------------: |
| `2001` | **error**   | unexpected token                                                                      |
| `2002` | **error**   | expected class name                                                                   |
| `2003` | **error**   | expected class body or semicolon                                                      |
| `2004` | **error**   | expected closing brace for class body                                                 |
| `2005` | **error**   | missing semicolon after class declaration                                             |
| `2006` | **error**   | expected name after delete                                                            |
| `2007` | **error**   | missing semicolon after delete declaration                                            |
| `2008` | **error**   | expected extern declaration name                                                      |
| `2009` | **error**   | missing semicolon after extern declaration                                            |
| `2010` | **error**   | expected right bracket in array assignment target                                     |
| `2011` | **error**   | expected assignment operator for array assignment                                     |
| `2012` | **error**   | missing semicolon after array assignment                                              |
| `2013` | **error**   | expected assignment operator                                                          |
| `2014` | **error**   | missing semicolon after assignment                                                    |
| `2015` | **error**   | expected value before end of file                                                     |
| `2016` | **error**   | expected value                                                                        |
| `2017` | **error**   | expected scalar value                                                                 |
| `2018` | **error**   | unterminated array literal                                                            |
| `2019` | **error**   | expected comma or right brace in array literal                                        |
| `2025` | **error**   | strict mode rejects class-like names starting with digit                              |
| `2020` | **warning** | autofix inserted missing class semicolon                                              |
| `2021` | **error**   | expected enum body                                                                    |
| `2022` | **error**   | expected enum item name                                                               |
| `2023` | **error**   | expected enum item delimiter                                                          |
| `2024` | **error**   | missing semicolon after enum declaration                                              |
| `2026` | **info**    | nested class without explicit inheritance in derived class may replace parent subtree |

## Processor

| Code   | Severity    | Description                                        |
| :----: | :---------: | :------------------------------------------------: |
| `3001` | **error**   | include target not found or unreadable             |
| `3002` | **error**   | unsupported config intrinsic                       |
| `3003` | **error**   | macro expansion failure                            |
| `3004` | **error**   | unterminated conditional block                     |
| `3005` | **error**   | invalid include syntax                             |
| `3006` | **error**   | unexpected elif directive                          |
| `3007` | **error**   | unexpected else directive                          |
| `3008` | **error**   | unexpected endif directive                         |
| `3009` | **error**   | error directive triggered                          |
| `3011` | **error**   | unsupported directive                              |
| `3012` | **error**   | missing macro name in define                       |
| `3013` | **error**   | invalid macro name in define                       |
| `3014` | **error**   | unterminated macro parameter list                  |
| `3015` | **warning** | macro redefinition                                 |
| `3016` | **error**   | unresolved macro-like invocation                   |
| `3017` | **error**   | unsupported __has_include in conditional directive |

