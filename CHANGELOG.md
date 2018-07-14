# Bug / Issue

- (2018-07-13) Add allow headers and allow methods to `CORS` middleware

# Breaking Changes

- (2018-07-05) Replaced input params of func `Render` from `func() []byte` to `[]byte`
- (2018-07-09) `Raptor` func `Start` now return `error` instead
- (2018-07-09) Changed err handling structure, it return HTTPError instead, and added property `ErrorHandler` to `Raptor` struct
- (2018-07-10) Changed core structure of `Raptor`

# New

- (2018-07-10) Introduced new func `NewAPIError`
- (2018-07-10) Introduced new func `HTML`, `HTMLBlob`
- (2018-07-14) Introduced static gzip support
