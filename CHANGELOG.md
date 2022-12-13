# 1.4.2 (2022/10/22)

### Features

- feat: add logger.WithFields (#2546)
- chore: refactor (#2545)
- feat(trace): support for disabling tracing of specified `spanName` (#2363)
- chore: adjust rpc comment format (#2501)
- feat: remove info log when disable log (#2525)
- chore(action): upgrade action (#2521)
- feat: support uuid.UUID in mapping (#2537)
- chore: add more tests (#2536)
- Fix typo (#2531)
- chore(deps): bump google.golang.org/grpc from 1.50.0 to 1.50.1 (#2527)
- Fix the wrong key about FindOne in mongo of goctl. (#2523)
- chore: add golangci-lint config file (#2519)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2514)
- Fix mongo insert tpl (#2512)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2511)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2510)
- chore: sqlx's metric name is different from redis (#2505)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.10.0 to 1.11.0 (#2504)
- chore: remove unnecessary code (#2499)
- token limit support context (#2335)
- feat(goctl): better generate the api code of typescript (#2483)
- chore(deps): bump google.golang.org/grpc from 1.49.0 to 1.50.0 (#2487)
- chore: remove init if possible (#2485)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.2 to 1.10.3 (#2484)
- chore: refactor to reduce duplicated code (#2477)
- chore: better shedding algorithm, make sure recover from shedding (#2476)
- feat(redis):add timeout method to extend blpop (#2472)
- chore: sort methods (#2470)
- feat: add logc package, support AddGlobalFields for both logc and logx. (#2463)
- feat: add string to map in httpx parse method (#2459)
- Readme Tweak (#2436)
- refactor: adjust http request slow log format (#2440)
- chore: gofumpt (#2439)
- feat: add color to debug (#2433)
- chore: replace fmt.Fprint (#2425)
- cleanup: deprecated field and func (#2416)
- refactor: redis error for prometheus metric label (#2412)
- feat: add log debug level (#2411)
- chore: add more tests (#2410)
- feat(goctl):Add ignore-columns flag (#2407)
- chore: add more tests (#2409)
- chore: update go-zero to v1.4.1
- feat: support caller skip in logx (#2401)

### Bug Fixes

- fix: redis's pipeline logs are not printed completely (#2538)
- fix(goctl): Fix issues (#2543)
- chore: fix lint errors (#2520)
- chore: fix naming problem (#2500)
- fix: replace Infof() with Errorf() in DurationInterceptor (#2495) (#2497)
- fix a few function names on comments (#2496)
- fix(mongo): fix file name generation errors (#2479)
- fix: etcd reconnecting problem (#2478)
- fix #2343 (#2349)
- fix: add more tests (#2473)
- fix(goctl): fix the unit test bug of goctl (#2458)
- fix #2435 (#2442)
- fix: fix log output (#2424)
- fix goctl help message (#2414)

# 1.4.1 (2022/09/17)

### Features

- chore: refactor the imports (#2406)
- feat: mysql and redis metric support (#2355)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2402)
- chore(deps): bump go.etcd.io/etcd/client/v3 from 3.5.4 to 3.5.5 (#2395)
- chore(deps): bump go.etcd.io/etcd/api/v3 from 3.5.4 to 3.5.5 (#2394)
- chore(deps): bump github.com/jhump/protoreflect from 1.12.0 to 1.13.0 (#2393)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2389)
- chore: refactor (#2388)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2385)
- feat: add grpc export (#2379)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 (#2383)
- feat: support targetPort option in goctl kube (#2378)
- Update readme-cn.md
- chore(deps): bump github.com/lib/pq from 1.10.6 to 1.10.7 (#2373)
- feat: support baggage propagation in httpc (#2375)
- chore(deps): bump go.uber.org/goleak from 1.1.12 to 1.2.0 (#2371)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.1 to 1.10.2 (#2370)
- chore: refactor (#2365)
- correct test case (#2340)
- Hidden java (#2333)
- make logx#getWriter concurrency-safe (#2233)
- generates nested types in doc (#2201)
- Add strict flag (#2248)
- improve: number range compare left and righ value (#2315)
- refactor: sequential range over safemap (#2316)
- safemap add Range method (#2314)
- chore: remove unused packages (#2312)
- Fix/del server interceptor duplicate copy md 20220827 (#2309)
- chore: refactor gateway (#2303)
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.3 to 2.0.5 (#2305)
- chore: refactor stat (#2299)
- Initialize CPU stat code only if used (#2020)
- chore(deps): bump google.golang.org/grpc from 1.48.0 to 1.49.0 (#2297)
- feat(redis): add ZaddFloat & ZaddFloatCtx (#2291)
- doc(readme): add star history (#2275)
- feat: rpc add health check function configuration optional (#2288)
- Update readme-cn.md
- Update issues.yml
- chore: Update readme (#2280)
- Update readme-cn.md
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.2 to 2.0.3 (#2267)
- chore: refactor logx (#2262)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.22.0 to 2.23.0 (#2260)
- test: add more tests (#2261)
- chore(deps): bump github.com/prometheus/client_golang (#2244)
- chore(deps): bump github.com/fullstorydev/grpcurl from 1.8.6 to 1.8.7 (#2245)
- chore: release action for goctl (#2239)

### Bug Fixes

- fix #2364 (#2377)
- fix: issue #2359 (#2368)
- fix:trace graceful stop,pre loss trace (#2358)
- fix:etcd get&watch not atomic (#2321)
- fix: thread-safe in getWriter of logx (#2319)
- fix: range validation on mapping (#2317)
- fix: handle the scenarios that content-length is invalid (#2313)
- fix: more accurate panic message on mapreduce (#2311)
- fix #2301,package conflict generated by ddl (#2307)
- fix: logx disable not working in some cases (#2306)
- fix:duplicate copy MD (#2304)
- fix resource manager dead lock (#2302)
- fix #2163 (#2283)
- fix(logx): display garbled characters in windows(DOS, Powershell) (#2232)
- fix #2240 (#2271)
- fix #2240 (#2263)
- fix: unsignedTypeMap type error (#2246)
- fix: test failure, due to go 1.19 compatibility (#2256)
- fix: time repr wrapper (#2255)

# 1.4.0 (2022/08/07)

### Features

- Update readme-cn.md
- Update readme.md
- feat: more meaningful error messages, close body on httpc requests (#2238)
- Update readme.md
- Update readme.md
- docs: update docs for gateway (#2236)
- chore: renaming configs (#2234)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.0 to 1.10.1 (#2225)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2222)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2223)
- chore: refactor redislock (#2210)
- feat(redislock): support set context (#2208)
- chore(deps): bump google.golang.org/protobuf from 1.28.0 to 1.28.1 (#2205)
- support mulitple protoset files (#2190)
- chore: let logx.SetWriter can be called anytime (#2186)
- chore: refactoring (#2182)
- chore: refactoring logx (#2181)
- feat: logx support logs rotation based on size limitation. (#1652) (#2167)
- Update goctl version (#2178)
- feat: Support for multiple rpc service generation and rpc grouping (#1972)
- Update readme-cn.md
- Update api template (#2172)
- chore: refactoring mapping name (#2168)
- feat: support customized header to metadata processor (#2162)
- feat: support google.api.http in gateway (#2161)
- feat: set content-type to application/json (#2160)
- feat: verify RpcPath on startup (#2159)
- feat: support form values in gateway (#2158)
- feat: export gateway.Server to let users add middlewares (#2157)
- Update readme-cn.md
- Update readme.md
- feat: restful -> grpc gateway (#2155)
- chore(deps): bump google.golang.org/grpc from 1.47.0 to 1.48.0 (#2147)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.9.1 to 1.10.0 (#2150)
- chore: remove unimplemented gateway (#2139)
- docs: update goctl readme (#2136)
- chore: refactor (#2130)
- chore: add more tests (#2129)
- feat:Add `Routes` method for server (#2125)
- feat: support logx.WithFields (#2128)
- feat: add Wrap and Wrapf in errorx (#2126)
- chore: coding style (#2120)

### Bug Fixes

- fix: #2216 (#2235)
- fix(logx): need to wait for the first caller to complete the execution. (#2213)
- fix: fix comment typo (#2220)
- fix: handling rpc error on gateway (#2212)
- fix: only setup logx once (#2188)
- fix: logx test foo (#2144)
- fix(httpc): fix typo errors (#2189)
- fix: remove invalid log fields in notLoggingContentMethods (#2187)
- fix:duplicate route check (#2154)
- fix: fix switch doesn't work bug (#2183)
- fix: fix #2102, #2108 (#2131)
- fix: goctl genhandler duplicate rest/httpx & goctl genhandler template support custom import httpx package (#2152)
- fix: generated sql query fields do not match template (#2004)
- fix goctl rpc protoc strings.EqualFold Service.Name GoPackage (#2046)

# 1.3.5 (2022/07/09)

### Features

- chore: remove blank lines (#2117)
- feat:`goctl model mongo ` add `easy` flag  for easy declare. (#2073)
- refactor:remove duplicate codes (#2101)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2115)
- feat: add method to jsonx (#2049)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2112)
- chore: update goctl version to 1.3.9 (#2111)
- chore: refactor (#2087)
- remove legacy code (#2086)
- chore: refactor (#2085)
- Update readme-cn.md
- remove legacy code (#2084)
- Update readme.md
- feat: CompareAndSwapInt32 may be better than AddInt32 (#2077)
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.1 to 2.0.2 (#2072)
- chore(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.4.2 (#2066)
- chore(deps): bump github.com/stretchr/testify from 1.7.2 to 1.8.0 (#2068)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.21.0 to 2.22.0 (#2067)
- chore(deps): bump github.com/ClickHouse/clickhouse-go/v2 (#2064)
- Create dependabot.yml
- [ci skip] Fix dead doc link (#2047)
- chore: remove lifecycle preStop because sh not exist in scratch (#2042)
- Update readme-cn.md
- Update readme-cn.md
- chore: Add command desc & color commands (#2013)
- Update readme.md
- Update readme.md
- feat: rest.WithChain to replace builtin middlewares (#2033)
- Update readme-cn.md
- chore: refactor to simplify disabling builtin middlewares (#2031)
- add user middleware chain function (#1913)
- Add fig (#2008)
- feat: support build Dockerfile from current dir (#2021)
- chore: upgrade action version (#2027)
- feat: add trace in httpc (#2011)
- Update readme-cn.md
- chore: coding style (#2012)
- Update readme.md
- Update readme.md
- feat: Replace mongo package with monc & mon (#2002)
- feat: convert grpc errors to http status codes (#1997)
- chore: rename methods (#1998)
- periodlimit new function TakeWithContext (#1983)
- typo: add type keyword (#1992)
- feat: add 'imagePullPolicy' parameter for 'goctl kube deploy' (#1996)
- chore: update dependencies (#1985)
- chore: update goctl version to 1.3.8 (#1981)
- Fix pg subcommand level error (#1979)
- chore: make methods consistent in signatures (#1971)
- test: make tests stable (#1968)
- chore: make print pretty (#1967)
- chore: better mongo logs (#1965)
- feat: print routes (#1964)
- chore: update dependencies (#1963)
- Update readme-cn.md
- Update readme.md
- Chore/goctl version (#1962)
- chore: update version (#1961)

### Bug Fixes

- fix #2109 (#2116)
- fix: type matching supports string to int (#2038)
- fix concurrent map writes (#2079)
- fix 当表有唯一键时，update()的形参和实参不一致 (#2010)
- fix: `\u003cnil\u003e` log output when http server shutdown. (#2055)
- fix:typo in readme.md (#2061)
- fix: quickstart wrong package when go.mod exists in parent dir (#2048)
- fix #1977 (#2034)
- fix: 修复 clientinterceptors/tracinginterceptor.go 显示接受消息字节为0 (#2003)
- fix goctl api clone template fail (#1990)
- fix: generate bad Dockerfile on given dir (#1980)
- fix: The validation of tag "options" is not working with int/uint type (#1969)
- test: fix fails (#1970)
- 🐞 fix: fixed typo (#1916)

# 1.3.4 (2022/06/03)

### Features

- Update readme-cn.md
- Update readme.md
- Update readme.md
- docs: add docs for logx (#1960)
- chore: refactoring mapping string to slice (#1959)
- chore: update roadmap (#1948)
- Delete duplicated crash recover logic. (#1950)
- Update readme-cn.md
- Update readme-cn.md
- chore: refine docker for better compatible with package main (#1944)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Update readme-cn.md
- Update readme.md
- Update readme-cn.md
- 优化代码 (#1939)
- core/mr:a little optimization for collector initialization in ForEach function (#1937)
- chore(action): simplified release configuration (#1935)
- chore: add release action to auto build binaries (#1884)
- feat: update docker alpine package mirror (#1924)
- Support built-in shorthand flags (#1925)
- feat: set default connection idle time for grpc servers (#1922)
- chore: update k8s.io/client-go for security reason, go is upgrade to 1.16 (#1912)
- Update readme-cn.md
- Update readme-cn.md
- Fix process blocking problem during check (#1911)
- feat: support WithStreamClientInterceptor for zrpc clients (#1907)
- Update FUNDING.yml
- chore: use get for quickstart, plain logs for easy understanding (#1905)
- use goproxy properly, remove files (#1903)
- feat: add toml config (#1899)
- chore: coding style for quickstart (#1902)
- refactor: refactor trace in redis & sql & mongo (#1865)
- feat: Add goctl quickstart (#1889)
- Fix code generation (#1897)
- Update readme-cn.md
- Update readme-cn.md
- chore: improve codecov (#1878)
- chore: update some logs (#1875)
- feat: logx with color (#1872)
- feat: Replace cli to cobra (#1855)
- Update readme.md
- Update readme-cn.md
- Update readme.md
- add conf documents (#1869)
- chore: refine tests (#1864)
- test: add codecov (#1863)
- test: add codecov (#1861)
- chore: use time.Now() instead of timex.Time() because go optimized it (#1860)
- feat: add fields with logx methods, support using third party logging libs. (#1847)
- test: add more tests (#1856)
- docs: update readme (#1849)

### Bug Fixes

- fix: panic on convert to string on fillSliceFromString() (#1951)
- fix: Useless delete cache logic in update (#1923)
- fix ts tpl (#1879)
- fix:tools/goctl/rpc/generator/template_test.go file has wrong parameters (#1882)
- chore: fix deprecated usages (#1871)
- fix time, duration, slice types on logx.Field (#1868)
- fix typo (#1857)

# 1.3.3 (2022/04/28)

### Features

- refactor: move json related header vars to internal (#1840)
- feat: Support model code generation for multi tables (#1836)
- Update readme-cn.md
- refactor: simplify the code (#1835)
- feat: support sub domain for cors (#1827)
- feat: upgrade grpc to 1.46, and remove the deprecated grpc.WithBalancerName (#1820)
- chore: optimize code (#1818)
- chore: remove gofumpt -s flag, default to be enabled (#1816)
- chore: refactor (#1814)
- feat: add trace in redis & mon & sql (#1799)
- chore: Embed unit test data (#1812)
- add go-grpc_opt and go_opt for grpc new command (#1769)
- feat: use mongodb official driver instead of mgo (#1782)
- update rpc generate sample proto file (#1709)
- feat(goctl): go work multi-module support (#1800)
- docs(goctl): goctl 1.3.4 migration note (#1780)
- chore: use grpc.WithTransportCredentials and insecure.NewCredentials() instead of grpc.WithInsecure (#1798)
- goctl api new should given a service_name explictly (#1688)
- show help when running goctl api without any flags (#1678)
- show help when running goctl docker without any flags (#1679)
- show help when running goctl rpc protoc without any flags (#1683)
- improve goctl rpc new (#1687)
- revert postgres package refactor (#1796)
- refactor: move postgres to pg package (#1781)
- feat: add httpc.Do & httpc.Service.Do (#1775)
- feat(goctl): supports api  multi-level importing (#1747)
- show help when running goctl rpc template without any flags (#1685)
- chore: avoid deadlock after stopping TimingWheel (#1768)
- Fix #1765 (#1767)
- chore: remove legacy code (#1766)
- chore: add doc (#1764)
- add more tests (#1763)
- feat: add goctl docker build scripts (#1760)
- Update readme-cn.md
- Update readme.md
- Update readme.md
- Update readme-cn.md
- feat: support ctx in kv methods (#1759)
- feat: use go:embed to embed templates (#1756)
- Support `goctl env install` (#1752)
- chore: update go-zero to v1.3.2 in goctl (#1750)

### Bug Fixes

- fix #1838 (#1839)
- fix: remove deprecated dependencies (#1837)
- fix #1806 (#1833)
- fix: rest: WriteJson get 200 when Marshal failed. (#1803)
- fix: Fix issue #1810 (#1811)
- fix: ignore timeout on websocket (#1802)
- fix: Hdel check result & Pfadd check result (#1801)
- chhore: fix usage typo (#1797)
- fix marshal ptr in httpc (#1789)
- fix(goctl): api/new/api.tpl (#1788)
- fix #1729 (#1783)
- fix bug: crash when generate model with goctl. (#1777)
- fix nil pointer if group not exists (#1773)
- fix: model unique keys generated differently in each re-generation (#1771)
- fix #1754 (#1757)

# 1.3.2 (2022/04/03)

### Features

- feat: simplify httpc (#1748)
- feat: return original value of setbit in redis (#1746)
- chore: update goctl version to 1.3.4 (#1742)
- feat: let model customizable (#1738)
- Fix zrpc code generation error with --remote (#1739)
- chore: refactor to use const instead of var (#1731)
- feat(goctl): supports model code 'DO NOT EDIT' (#1728)
- Fix unit test (#1730)
- refactor: guard timeout on API files (#1726)
- Added support for setting the parameter size accepted by the interface and custom timeout and maxbytes in API file (#1713)
- chore: refactor code (#1708)
- feat: remove reentrance in redislock, timeout bug (#1704)
- chore: refactor code (#1700)
- chore: refactor code (#1699)
- feat: add getset command in redis and kv (#1693)
- feat: add httpc.Parse (#1698)
- Add verbose flag (#1696)
- Update LICENSE
- refactor: simplify the code (#1670)
- Remove debug log (#1669)
- feat: support -base to specify base image for goctl docker (#1668)
- Remove unused code (#1667)
- feat: add Dockerfile for goctl (#1666)
- feat: Remove  command `goctl rpc proto`  (#1665)
- Mkdir if not exists (#1659)
- chore: remove unnecessary env (#1654)
- refactor: httpc package for easy to use (#1645)
- refactor: httpc package for easy to use (#1643)
- FindOneBy 漏 Context (#1642)
- feat: add httpc/Service for convinience (#1641)
- feat: add httpc/Get httpc/Post (#1640)
- feat: add rest/httpc to make http requests governacible (#1638)
- Update ROADMAP.md
- Update ROADMAP.md
- feat: support cpu stat on cgroups v2 (#1636)
- feat: support oracle :N dynamic parameters (#1552)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Support for referencing types in different API files using format (#1630)
- feat: support scratch as the base docker image (#1634)
- chore: reduce the docker image size (#1633)
- Fix #1585 #1547 (#1624)
- chore: update goctl version to 1.3.3, change docker build temp dir (#1621)
- Fix #1609 (#1617)
- Fix #1614 (#1616)
- chore: refactor code (#1613)
- chore: add unit tests (#1615)
- model中db标签增加'-'符号以支持数据库查询时忽略对应字段. (#1612)
- feat(goctl): api dart support flutter v2 (#1603)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- test: add more tests (#1604)
- Update readme-cn.md
- chore: update go-zero to v1.3.1 in goctl (#1599)

### Bug Fixes

- fix: model generation bug on with cache (#1743)
- fix(goctl): api format with reader input (#1722)
- fix: empty slice are set to nil (#1702)
- fix -cache=true insert no clean cache (#1672)
- chore: fix lint issue (#1694)
- fix: the new  RawFieldNames considers the tag with options. (#1663)
- fix(goctl): model method FindOneCtx should be FindOne (#1656)
- typo (#1655)
- fix: typo (#1646)
- fix: Update unix-like path regex (#1637)
- fix(goctl): kotlin code generation (#1632)
- fix(goctl): dart gen user defined struct array (#1620)
- fix: HitQuota should be returned instead of Allowed when limit is equal to 1. (#1581)
- fix: fix(gctl): apiparser_parser auto format (#1607)

# 1.3.1 (2022/03/01)

### Features

- chore: upgrade etcd (#1597)
- Update readme.md
- build: update goctl dependency ddl-parser to v1.0.3 (#1586)
- test: add testcase for FIFO Queue in collection module (#1589)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Update readme.md
- Fix bug int overflow while build goctl on arch 386 (#1582)
- chore: add goctl command help (#1578)
- Update readme.md
- Update readme.md
- feat: supports `importValue` for more path formats (#1569)
- update goctl to go 1.16 for io/fs usage (#1571)
- feat: support pg serial type for auto_increment (#1563)
- Feature: Add goctl env (#1557)
- feat: log 404 requests with traceid (#1554)
- feat: support ctx in sql model generation (#1551)
- feat: support ctx in sqlx/sqlc, listed in ROADMAP (#1535)
- docs: add go-zero users (#1546)
- ignore context.Canceled for redis breaker (#1545)
- chore: update help message (#1544)
- add the serviceAccount of deployment (#1543)
- chore:use struct pointer (#1538)
- docs: update roadmap (#1537)
- Update readme-cn.md
- chore: refactor cache (#1532)
- feat: support ctx in `Cache` (#1518)
- chore: goctl format issue (#1531)
- upgrade grpc version (#1530)
- chore: update goctl version to 1.3.2 (#1524)
- refactor: refactor yaml unmarshaler (#1517)
- chore: optimize yaml unmarshaler (#1513)
- chore: make error clearer (#1514)
- feat: update go-redis to v8, support ctx in redis methods (#1507)
- Update readme-cn.md
- feature: Add `goctl completion` (#1505)
- test: change fuzz tests (#1504)
- ci: add test for win (#1503)
- chore: update command comment (#1501)
- docs: update tal-tech to zeromico in docs (#1498)

### Bug Fixes

- Revert "🐞 fix(gen): pg gen of insert (#1591)" (#1598)
- 🐞 fix(gen): pg gen of insert (#1591)
- fix: goctl api dart support `form` tag (#1596)
- chore: fix data race (#1593)
- fix #1541 (#1542)
- fix issue of default migrate version (#1536)
- fix #1525 (#1527)
- fix: fix a typo (#1522)
- fixes typo (#1511)
- fix typo: goctl protoc usage (#1502)
- fix: goctl not compile on windows (#1500)

# 1.3.0 (2022/02/01)

### Features

- chore: update goctl version (#1497)
- feat: add runtime stats monitor (#1496)
- feat: handling panic in mapreduce, panic in calling goroutine, not inside goroutines (#1490)
- Update readme-cn.md
- chore: improve migrate confirmation (#1488)
- chore: update warning message (#1487)
- patch: goctl migrate (#1485)
- chore: update go version for goctl (#1484)
- refactor: rename from tal-tech to zeromicro for goctl (#1481)

### Bug Fixes

- fix: goroutine stuck on edge case (#1495)

# 1.3.0-beta (2022/01/25)

### Features

- Feature/trie ac automation (#1479)
- chore: optimize string search with Aho–Corasick algorithm (#1476)
- Polish the words in readme.md (#1475)
- docs: add go-zero users (#1473)
- chore: update unauthorized callback calling order (#1469)
- Fix/issue#1289 (#1460)
- patch: save missing templates to disk (#1463)
- Fix/issue#1447 (#1458)
- feat: implement console plain output for debug logs (#1456)
- chore: check interface satisfaction w/o allocating new variable (#1454)
- chore: remove jwt deprecated (#1452)
- feat: 支持redis的LTrim方法 (#1443)
- chore: upgrade dependencies (#1444)
- ci: add translator action (#1441)
- Feature rpc protoc (#1251)
- chore: refactor periodlimit (#1428)
- docs: add go-zero users (#1425)
- docs: add go-zero users (#1424)
- update docs (#1421)
- Fix pg model generation without tag (#1407)
- feat: Add migrate (#1419)
- docs: update install readme (#1417)

### Bug Fixes

- fix #1468 (#1478)
- chore: fix typo (#1437)
- remove unnecessary drain, fix data race (#1435)
- fix: mr goroutine leak on context deadline (#1433)
- fix: golint issue (#1423)

# 1.3.0-alpha (2022/01/05)

### Features

- chore: refactor rest/timeouthandler (#1415)
- feat: rename module from tal-tech to zeromicro (#1413)
- chore: update go-zero to v1.2.5 (#1410)
- refactor file|path (#1409)

# 1.2.5 (2022/01/02)

### Features

- docs: update roadmap (#1405)
- feat: support tls for etcd client (#1390)
- refactor: optimize fx (#1404)
- docs: update goctl installation command (#1403)
- feat: implement fx.NoneMatch, fx.First, fx.Last (#1402)
- chore: simplify mapreduce (#1401)
- chore: update goctl version (#1394)
- ci: remove 386 binaries (#1393)
- ci: remove windows 386 binary (#1392)
- test: add more tests (#1391)
- feat: Add --remote (#1387)
- feat: support array in default and options tags (#1386)
- docs: add go-zero users (#1381)
- docs: update slack invitation link (#1378)
- chore: update goctl version to 1.2.4 for release tools/goctl/v1.2.4 (#1372)
- Updated MySQL生成表结构体遇到关键字db部分保持原字段名定义 (#1369)
- ci: add release action to auto build binaries (#1371)
- docs: update goctl markdown (#1370)
- feat: support context in MapReduce (#1368)
- chore: add 1s for tolerance in redislock (#1367)
- chore: coding style and comments (#1361)
- chore: optimize `ParseJsonBody` (#1353)
- chose: cancel the assignment and judge later (#1359)
- chore: put error message in error.log for verbose mode (#1355)
- test: add more tests (#1352)
- Revert "排除客户端中断导致的503错误 (#1343)" (#1351)
- feat: treat client closed requests as code 499 (#1350)
- 排除客户端中断导致的503错误 (#1343)
- Update FUNDING.yml
- chore: add tests & refactor (#1346)
- Feature: support adding custom cache to mongoc and sqlc (#1313)
- chore: add comments (#1345)
- chore: update goctl version to 1.2.5 (#1337)
- Update template (#1335)
- Feat goctl bug (#1332)
- feat: tidy mod, update go-zero to latest (#1334)
- feat: tidy mod, update go-zero to latest (#1333)
- chore: refactor (#1331)
- Update types.go (#1314)
- feat: tidy mod, add go.mod for goctl (#1328)
- chore: format code (#1327)
- commit  missing method for redis (#1325)
- docs: add go-zero users (#1323)
- style: format code (#1322)
- chore: rename service context from ctx to svcCtx (#1299)
- docs: add go-zero users (#1294)
- Revert "feat: reduce dependencies of framework by add go.mod in goctl (#1290)" (#1291)
- feat: reduce dependencies of framework by add go.mod in goctl (#1290)
- chore: update cli version (#1287)

### Bug Fixes

- fix #1330 (#1382)
- chore: fix golint issues (#1396)
- fix readme-cn (#1388)
- fix #1070 (#1389)
- fix #1376 (#1380)
- fix:  command system info missing go version (#1377)
- fix redis try-lock bug (#1366)
- go-zero tools ,fix a func,api new can not choose style (#1356)
- fix: #1318 (#1321)
- fix #1309 (#1317)
- fix #1305 (#1307)
- fix: go issue 16206 (#1298)
- fix #1288 (#1292)

# 1.2.4 (2021/12/01)

### Features

- feat: support third party orm to interact with go-zero (#1286)
- Feature api root path (#1261)
- chore: cleanup zRPC retry code (#1280)
- feature(retry): Delete retry mechanism (#1279)
- feat: support %w in logx.Errorf (#1278)
- chore: only allow cors middleware to change headers (#1276)
- chore: avoid superfluous WriteHeader call errors (#1275)
- feat: add rest.WithCustomCors to let caller customize the response (#1274)
- Cli (#1272)
- docs: update readme to use goctl@cli (#1255)
- chore: update goctl version (#1250)
- Revert "Revert "feat: enable retry for zrpc (#1237)"" (#1246)

### Bug Fixes

- fixes #987 (#1283)
- feat: add etcd resolver scheme, fix discov minor issue (#1281)
- fixes #1257 (#1271)

# 1.2.3 (2021/11/15)

### Features

- Revert "feat: enable retry for zrpc (#1237)" (#1245)
- Duplicate temporary variable (#1244)
- Update template (#1243)
- feat: enable retry for zrpc (#1237)
- chore: remove conf.CheckedDuration (#1235)
- reset link goctl (#1232)
- feat: disable grpc retry, enable it in v1.2.4 (#1233)
- chore: refactor, better goctl message (#1228)
- chore: remove unused const (#1224)
- Update FUNDING.yml
- chore: redislock use stringx.randn replace randomStr func (#1220)
- feat: exit with non-zero code on errors (#1218)
- feat: support CORS, better implementation (#1217)
- Create FUNDING.yml
- chore: refine code (#1215)
- docs: add go-zero users (#1214)
- Fix issue 1205 (#1211)
- feat: support CORS by using rest.WithCors(...) (#1212)
- update dependencies. (#1210)
- test: add more tests (#1209)
- goctl docker command add -version (#1206)
- feat: support customizing timeout for specific route (#1203)
- feat: add NewSessionFromTx to interact with other orm (#1202)
- feat: simplify the grpc tls authentication (#1199)
- feat: use WithBlock() by default, NonBlock can be set in config or WithNonBlock() (#1198)
- chore: remove semicolon for routes of services in api files (#1195)
- chore: update goctl version to 1.2.3, prepare for release (#1193)
- feat: slow threshold customizable in zrpc (#1191)
- feat: slow threshold customizable in rest (#1189)
- feat: slow threshold customizable in sqlx (#1188)
- feat: slow threshold customizable in redis (#1187)
- feat: slow threshold customizable in mongo (#1186)
- feat: slow threshold customizable in redis (#1185)
- docs: update roadmap (#1184)
- feat: support multiple trace agents (#1183)
- feat: let different services start prometheus on demand (#1182)
- refactor: simplify tls config in rest (#1181)
- [update] add plugin config (#1180)
- test: add more tests (#1179)
- docs: update roadmap (#1178)
- docs: add go-zero users (#1176)
- feat: support auth account for etcd (#1174)
- feat: support ssl on zrpc, simplify the config (#1175)
- support RpcClient Vertify With Unilateralism and Mutual (#647)
- docs: add go-zero users (#1172)
- Feature add template version (#1152)
- test: add more tests (#1166)
- chore: reorg imports, format code, make MaxRetires default to 0 (#1165)
- Add grpc retry (#1160)
- test: add more tests (#1163)
- chore: reverse the order of stopping services (#1159)
- docs: update qr code (#1158)
- test: add more tests (#1154)
- test: add more tests (#1150)
- Mark deprecated syntax (#1148)
- test: add more tests (#1149)
- test: add more tests (#1147)
- docs: add go-zero users (#1141)
- test: add more tests (#1138)
- test: add more tests (#1137)
- docs: add go-zero users (#1135)
- test: add more tests (#1134)
- Fix issue #1127 (#1131)
- docs: add go-zero users (#1130)

### Bug Fixes

- fixes #1169 (#1229)
- fixes #1222 (#1223)
- feat: ignore rest.WithPrefix on empty prefix (#1208)
- Generate route with prefix (#1200)
- feat: add rest.WithPrefix to support route prefix (#1194)
- fix the package name of grpc client (#1170)
- fix(goctl): repeat creation protoc-gen-goctl symlink (#1162)

# 1.2.2 (2021/10/12)

### Features

- chore: refine rpc template in goctl (#1129)
- go-zero/core/hash/hash_test.go  增加测试 TestMd5Hex (#1128)
- Add `opts ...grpc.CallOption` in grpc client (#1122)
- Add request method in http log (#1120)
- update goctl version to 1.2.2 (#1125)
- add cncf landscape (#1123)
- test: add more tests (#1119)
- add more tests (#1115)
- add more tests (#1114)
- add more tests (#1113)
- test: add more tests (#1112)
- feat: opentelemetry integration, removed self designed tracing (#1111)
- docs: update roadmap (#1110)
- feat: reflection grpc service (#1107)
- test: add more tests (#1106)
- Fix the `resources` variable not reset after the resource manager is closed (#1105)
- chore: replace redis.NewRedis with redis.New (#1103)
- chore: mark redis.NewRedis as Deprecated, use redis.New instead. (#1100)
- update grpc package (#1099)
- Update Makefile (#1098)
- ci: add reviewdog (#1096)
- docs: update roadmap (#1094)
- docs: update roadmap (#1093)
- ci: accurate error reporting on lint check (#1089)
- update zero-doc links in readme (#1088)
- docs: change organization from tal-tech to zeromicro in readme (#1087)
- ci: add Lint check on commits (#1086)
- Revert "chore: run unit test with go 1.14 (#1084)" (#1085)
- chore: run unit test with go 1.14 (#1084)
- update goctl api (#1052)
- chore: when run goctl-rpc, the order of proto message aliases should be (#1078)
- coding style (#1083)
- we can use otel.ErrorHandlerFunc instead of custom struct when we update OpenTelemetry to 1.0.0 (#1081)
- update go.mod (#1079)
- Create a symbol link file named protoc-gen-goctl from goctl (#1076)
- update OpenTelemetry to 1.0.0 (#1075)
- update issue templates (#1074)
- Update issue templates
- downgrade golang-jwt to support go 1.14 (#1073)
- Add MustTempDir (#1069)
- upgrade grpc version & replace github.com/golang/protobuf/protoc-gen-go with google.golang.org/protobuf (#1065)
- add repo moving notice (#1062)
- add go-zero users (#1061)
- chore: make comment accurate (#1055)
- mention cncf landscape (#1054)
- add go-zero users (#1051)

### Bug Fixes

- fix: opentelemetry traceid not correct (#1108)
- fix bug: generating dart code error (#1090)
- fix AtomicError panic when Set nil (#1049) (#1050)
- fix jwt security issue by using golang-jwt package (#1066)
- fix #1058 (#1064)
- chore: fix comment issues (#1056)
- fix test error on ubuntu (#1048)
- fix typo parse.go error message (#1041)

# 1.2.1 (2021/09/14)

### Features

- update goctl version to 1.2.1 (#1042)
- remove goctl config command (#1035)

# 1.2.0 (2021/09/13)

### Features

- update k8s.io/client-go etc to use go 1.15 (#1031)
- reorg imports, format code (#1024)
- revert changes
- add go-zero users (#1022)
- rename sharedcalls to singleflight (#1017)
- refactor for better error reporting on sql error (#1016)
- expose sql.DB to let orm operate on it (#1015)
- update codecov settings (#1010)
- refactoring tracing interceptors. (#1009)
- use sdktrace instead of trace for opentelemetry to avoid conflicts (#1005)
- api文件中使用group时生成的handler和logic的包名应该为group的名字 (#545)
- add api template file (#1003)
- add opentelemetry test (#1002)
- reorg imports, format code (#1000)
- 开启otel后，tracelog自动获取otle的traceId和spanId (#946)
- optimize unit test (#999)
- rest log with context (#998)
- refactor to shorter config name (#997)
- feat: change logger to traceLogger for getting traceId when recovering (#374)
- 修复使用 postgres 数据库时，位置参数重复，导致参数与值不对应的问题。 (#960)
- move opentelemetry into trace package, and refactoring (#996)
- Feature goctl error wrap (#995)
- disable codecov github checks (#993)
- implement k8s service discovery (#988)
- format code, and reorg imports (#991)
- Fix filepath (#990)
- remote handler blank line when .HasRequest is false (#986)
- use codecov action v1 (#985)
- format coding style (#983)
- httpx.Error response without body (#982)
- format code (#979)
- configurable for load and stat statistics logs (#980)
- add go-zero users (#978)
- refactor (#977)
- update go version to 1.14 for github workflow (#976)
- Update Codecov `action` (#974)
- format coding style (#970)
- Fix context error in grpc (#962)
- format coding style (#969)
- Fix issues (#965)
- Add a test case for database code generation `tool` (#961)
- 修复stream拦截器tracer名问题 (#944)
- rest otel support (#943)
- add the opentelemetry tracing (#908)
- make sure setting code happen before callback in rest (#936)
- Fix issues (#931)
- update slack invite url (#933)
- add go-zero users (#928)
- Optimize model naming (#910)
- add unit test (#921)
- refactor (#920)
- Add traceId to the response headers (#919)
- add stringx.FirstN with ellipsis (#916)
- redis.go，type StringCmd = red.StringCmd (#790)
- add stringx.FirstN (#914)
- export pathvar for user-defined routers (#911)
- add Errorv/Infov/Slowv (#909)
- optimize grpc generation env check (#900)
- add workflow for closing inactive issues (#906)
- Update readme for better description. (#904)
- format coding style (#905)
- 带下划线的项目,配置文件名字错误。 (#733)
- refactor goctl (#902)
- refactor (#878)
- remove unnecessary chars. (#898)
- simplify type definition in readme (#896)
- refactor rest code (#895)
- add logx.DisableStat() to disable stat logs (#893)
- better text rendering (#892)
- format code (#888)
- format code (#884)
- add go-zero users (#883)

### Bug Fixes

- fix symlink issue on windows for goctl (#1034)
- fix golint issues (#1027)
- fix proc.Done not found in windows (#1026)
- fix #1014 (#1018)
- fix golint issues, update codecov settings. (#1011)
- fix #1006 (#1008)
- fix golint issues (#992)
- fix #971 (#972)
- fix #957 (#959)
- fix #556 (#938)
- fix lint errors (#937)
- fix #820 (#934)
- fix missing `updateMethodTemplateFile` (#924)
- fix #915 (#917)
- fix #889 (#912)
- refactor goctl, fix golint issues (#903)
- fix: 解决golint 部分警告 (#897)
- fix golint issues (#899)
- fix http header binding failure bug #885 (#887)
- optimize mongo generation without cache (fix #881) (#882)

# 1.1.10 (2021/08/04)

### Features

- update goctl version to 1.1.10 (#874)
- add goctl rpc template home flag (#871)

### Bug Fixes

- fix #792 (#873)
- fix context missing (#872)
- fix bug that proc.SetTimeToForceQuit not working in windows (#869)

# 1.1.9 (2021/08/01)

### Features

- add correct example for pg's url (#857)
- optimize typo (#855)
- upgrade grpc package (#845)
- simplify timeoutinterceptor (#840)
- Fixed http listener error. (#843)
- Feature model postgresql (#842)
- remove faq for old versions (#828)
- add go-zero users, update faq (#827)
- Fix the error stream method name (#826)
- go format with extra rules (#821)
- Fix issues: #725, #740 (#813)
- optimized (#819)
- Fix rpc generator bug (#799)
- Add --go_opt flag to adapt to the version after 1.4.0 of protoc-gen-go (#767)
- [WIP]Add parse headers info (#805)

### Bug Fixes

- fix issue #861 (#862)
- fix issue #831 (#850)
- fix issue #836 (#849)
- fix #796 (#844)
- Added database prefix of cache key. (#835)
- To generate grpc stream, fix issue #616 (#815)
- fix bug that empty query in transaction (#801)

# 1.1.8 (2021/06/23)

### Features

- refactor mapping (#782)
- Add Sinter,Sinterstore & Modify TestRedis_Set (#779)
- update readme images (#776)
- update image rendering in readme (#775)
- disable load & stat logs for goctl (#773)
- upgrade grpc & etcd dependencies (#771)
- Fix issue #747 (#765)
- remove useless annotation (#761)
- add roadmap (#764)
- add contributing guid (#762)
- add code of conduct (#760)
- refactor fx (#759)
- Add some stream features (#712)
- add go-zero users (#756)
- add go-zero users (#751)
- replace cache key with colon (#746)
- add go-zero users (#739)
- Fix a typo (#729)
- add go-zero users, update slack invite link (#728)
- add go-zero users (#726)
- add go-zero users. (#723)
- optimize nested conditional (#709)
- optimize nested conditional (#708)
- Add document & comment for spec (#703)
- print entire sql statements in logx if necessary (#704)
- chore(format): change by gofumpt tool (#697)
- update goctl version to 1.1.8 (#696)
- add go-zero users (#688)
- resolve #610 (#684)
- Optimize model nl (#686)
- update readme for documents links (#681)

### Bug Fixes

- fix: Fix problems with non support for multidimensional arrays and basic type pointer arrays (#778)
- fix bug that etcd stream cancelled without re-watch (#770)
- fix broken link (#763)
- fix #736 (#738)
- fix golint issues, and optimize code (#705)
- fix #683 (#690)
- fix invalid link (#689)
- fix #676 (#682)
- fix zh_cn document url (#678)
- fix issue: https://github.com/zeromicro/goctl-swagger/issues/6 (#680)
- fix some typo (#677)

# 1.1.7 (2021/05/08)

### Features

- update readme (#673)
- replace antlr module (#672)
- modify the order of PrometheusHandler (#670)
- update wechat qrcode (#665)
- disable prometheus if not configured (#663)
- add go-zero users (#643)
- update readme (#640)
- update readme (#638)
- optimize code (#637)
- spelling mistakes (#634)
- chore: update code format. (#628)
- update go-zero users (#623)
- add syncx.Guard func (#620)
- update readme (#617)
- add code coverage (#615)
- add FAQs in readme (#612)
- update go-zero users (#611)
- update go-zero users (#609)
- add go-zero users registry notes (#608)
- add go-zero users (#607)
- simplify redis tls implementation (#606)
- redis增加tls支持 (#595)
- refactor - remove ShrinkDeadline, it's the same as context.WithTimeout (#599)
- Replace contextx.ShrinkDeadline with context.WithTimeout (#598)
- Simplify contextx.ShrinkDeadline (#596)
- update regression test comment (#590)
- update regression test comment (#589)
- remove rt mode log (#587)

### Bug Fixes

- fix antlr mod (#669)
- fix some typo (#667)
- fix comment function names (#649)
- doc: fix spell mistake (#627)
- fix (#592)
- fix a simple typo (#588)
- fix typo (#586)
- fix typo (#585)
- fix golint issues (#584)

# 1.1.6 (2021/03/27)

### Features

- optimize code (#579)
- support postgresql (#583)
- avoid goroutine leak after timeout (#575)
- gofmt logs (#574)
- add timezone and timeformat (#572)
- zrpc timeout & unit tests (#573)
- make hijack more stable (#565)
- refactor, and add comments to describe graceful shutdown (#564)
- Feature mongo gen (#546)
- Hdel support for multiple key deletion (#542)
- add important notes in readme (#560)
- add http hijack methods (#555)
- update doc link (#552)
- rename (#543)
- 暴露redis EvalSha  以及ScriptLoad接口 (#538)

### Bug Fixes

- fix golint issues (#561)
- fix spelling (#551)
- fix golint issues (#540)
- fix collection breaker (#537)

# 1.1.5 (2021/03/02)

### Features

- patch 1.1.5 (#530)
- feature 1.1.5 (#411)
- golint core/discov (#525)
- 修正http转发头字段值错误 (#521)
- Code optimized (#523)
- Code optimized (#493)
- add redis bitmap command (#490)
- redis add bitcount (#483)
- prevent negative timeout settings (#482)
- zrpc client support block (#412)
- Code optimized (#474)
- add more tests for service (#463)
- add more tests for rest (#462)
- Update serviceconf.go (#460)
- move examples into zero-examples (#453)
- remove images, use zero-doc instead (#452)
- add api doc (#449)
- add discov tests (#448)
- remove etcd facade, added for testing purpose (#447)
- add more tests for stores (#446)
- add more tests for stores (#445)
- add more tests for mongoc (#443)
- add more tests for sqlx (#442)
- add more tests for zrpc (#441)
- add more tests for sqlx (#440)
- add more tests for proc (#439)
- ring struct add lock (#434)
- Update readme.md
- update readme for broken links (#432)
- Support redis command Rpop (#431)
- support hscan in redis (#428)
- use english readme as default, because of github ranking (#427)
- Modify the http content-length max range : 30MB --> 32MB (#424)
- modify the maximum content-length to 30MB (#413)
- optimize code (#417)
- support zunionstore in redis (#410)
- use env if necessary in loading config (#409)
- update goctl version to 1.1.3 (#402)

### Bug Fixes

- fix golint issues (#535)
- fix golint issues (#533)
- fix golint issues (#532)
- fix golint issues in zrpc (#531)
- fix golint issues in rest (#529)
- fix broken build (#528)
- fix golint issues in core/stores (#527)
- fix golint issues in core/syncx (#526)
- fix golint issues in core/threading (#524)
- fix golint issues in core/utils (#520)
- fix golint issues in core/timex (#517)
- fix golint issues in core/stringx (#516)
- fix golint issues in core/stat (#515)
- fix misspelling (#513)
- fix golint issues in core/service (#512)
- fix golint issues in core/search (#509)
- fix golint issues in core/rescue (#508)
- fix golint issues in core/queue (#507)
- fix golint issues in core/prometheus (#506)
- fix broken links in readme (#505)
- fix golint issues in core/prof (#503)
- fix golint issues in core/proc (#502)
- fix golint issues in core/netx (#501)
- fix golint issues in core/mr (#500)
- fix golint issues in core/metric (#499)
- fix golint issues in core/mathx (#498)
- fix golint issues in core/mapping (#497)
- fix golint issues in core/logx (#496)
- fix golint issues in core/load (#495)
- fix golint issues in core/limit (#494)
- fix golint issues in core/lang (#492)
- fix golint issues in core/jsonx (#491)
- fix golint issues in core/jsontype (#489)
- fix golint issues in core/iox (#488)
- fix golint issues in core/hash (#487)
- fix golint issues in core/fx (#486)
- fix golint issues in core/filex (#485)
- fix golint issues in core/executors (#484)
- fix golint issues in core/errorx (#480)
- fix golint issues in core/discov (#479)
- fix golint issues in core/contextx (#477)
- fix golint issues in core/conf (#476)
- fix golint issues in core/collection, refine cache interface (#475)
- fix golint issues in core/codec (#473)
- fix issue #469 (#471)
- fix gocyclo warnings (#468)
- fix golint issues in core/cmdline (#467)
- fix golint issues in core/breaker (#466)
- fix golint issues in core/bloom (#465)
- fix golint issues (#459)
- fix golint issues (#458)
- fix golint issues, else blocks (#457)
- fix golint issues, redis methods (#455)
- fix golint issues, package comments (#454)
- fix golint issues, exported doc (#451)
- fix golint issues (#450)
- fixes issue #425 (#438)
- fix readme.md error (#429)

# 1.1.4 (2021/01/16)

### Features

- optimized (#392)
- add more tests for codec (#391)
- update readme (#390)
- Update periodicalexecutor.go (#389)
- format code (#386)

### Bug Fixes

- fix type convert error (#395)

# 1.1.3 (2021/01/13)

# 1.1.3-pre (2021/01/13)

# 1.1.3-beta (2021/01/13)

### Features

- simplify cgroup controller separation (#384)
- make sure unlock safe even if listeners panic (#383)
- code optimized (#382)
- Java optimized (#376)
- feature: refactor api parse to g4 (#365)
- add more tests for conf (#371)
- update doc to use table to render plugins (#370)
- remove duplicated code in goctl (#369)
- update goctl version to 1.1.3 (#364)

### Bug Fixes

- fix cgroup bug (#380)
- fix server.start return nil points (#379)
- f-fix spell (#381)
-  fix inner type generate error (#377)
- fix return in for (#367)
- Feature model fix (#362)

# 1.1.2 (2021/01/05)

### Features

- make sure offset less than size even it's checked inside (#354)
- add godoc for RollingWindow (#351)
- simple rolling windows code (#346)
- optimized goctl format (#336)
- close issue of #337 (#347)
- align bucket boundary to interval in rolling window (#345)
- simplify rolling window code, and make tests run faster (#343)
- set guarded to false only on quitting background flush (#342)
- simplify periodical executor background routine (#339)
- add discord chat group in readme
- modify the goctl gensvc template (#323)
- Java (#327)
- simplify http.Flusher implementation (#326)
- simplify code with http.Flusher type conversion (#325)
- The ResponseWriters defined in rest.handler add Flush interface. (#318)
- add more tests for prof (#322)
- add more tests for zrpc (#321)
- add more tests (#320)
- add more tests (#319)
- add go report card back (#313)
- Update codeql-analysis.yml
- format code (#312)
- add config load support env var (#309)
- Update readme.md
- add wechat micro practice qrcode image (#289)
- Update readme.md
- Update readme-en.md
- Update readme.md
- Create codeql-analysis.yml

### Bug Fixes

- fix potential data race in PeriodicalExecutor (#344)
- fix rolling window bug (#340)
- optimize code that fixes issue #317 (#338)
- fix bug #317 (#335)
- fix issue #317 (#331)
- fix broken link.
- fix broken doc link
- fixes #286 (#315)
- feature model fix (#296)

# 1.1.1 (2020/12/12)

### Features

- optimize dockerfile generation (#284)
- refactor (#283)
- format dockerfile on non-chinese mode (#282)
- Update readme-en.md
- add EXPOSE in dockerfile generation (#281)
- optimize test case of TestRpcGenerate (#279)
- add category docker & kube (#276)
- optimize dockerfile (#272)
- fmt code (#270)

### Bug Fixes

- fix gocyclo warnings (#278)
- fix dockerfile generation bug (#277)
- fix issue #266 (#275)
- fix tracelogger_test TestTraceLog (#271)

# 1.1.0 (2020/12/09)

### Features

- require go 1.14 (#263)
- feature plugin custom flag (#251)
- optimized parse tag (#256)
- refactor & format code (#255)
- Feature bookstore update (#253)
- improve data type conversion (#236)
- goctl add plugin support (#243)
- support k8s deployment yaml generation (#247)
- optimize docker file generation, make docker build faster (#244)
- optimization (#241)
- some optimize by kevwan and benying (#240)
- simplify code (#234)

### Bug Fixes

- fix lint errors (#249)

# 1.0.29 (2020/11/28)

### Features

- simplify code, format makefile (#233)
- optimization (#221)
- modify the service name from proto (#230)
- Improve Makefile robustness (#224)
- set default handler value (#228)
- update version (#226)
- check go.mod before build docker image (#225)
- feature model interface (#222)
- feature: file namestyle (#223)
- format import
- 1.use local variable i; 2.make sure limiter larger than timer period (#218)

### Bug Fixes

- fix doc errors

# 1.0.28 (2020/11/19)

### Features

- optimize api new (#216)
- patch model&rpc (#207)
- update readme
- modify image url (#213)
- type should not define nested (#212)
- add error handle tests
- support error customization
- support type def without struct token (#210)
- add redis geospatial (#209)
- optimize parser (#206)
- update goctl readme
- update example
- update example

### Bug Fixes

- fix issue #205

# 1.0.27 (2020/11/13)

### Features

- refactor parser and remove deprecated code (#204)
- 1. group support multi level folder 2. remove force flag (#203)
- api support for comment double slash // (#201)
- format code
- change grpc interceptor to chain interceptor (#200)
- update etcd yaml to avoid no such nost resolve problem
- no default metric (#199)
- add dockerfile into template
- format service and add test (#197)
- rename postgres
- default metric host (#196)
- rewrite (#194)

# 1.0.26 (2020/11/08)

### Features

- add dockerfile generator
- simplify http server starter
- graceful shutdown refined
- update doc (#193)
- Close the process when shutdown is finished (#157)
- break generator when happen error (#192)
- update cli package
- add more test (#189)
- refine code style

### Bug Fixes

- fix issue #186

# 1.0.25 (2020/11/05)

### Features

- faster the tests

### Bug Fixes

- rpc generation fix (#184)
- fix duplicate alias (#183)

# 1.0.24 (2020/11/05)

### Features

- reactor rpc (#179)
- add https listen and serve
- move redistest into redis package
- refine tests
- update doc
- route support no request and response (#178)
- add more tests
- add more tests
- Update sharedcalls.go (#174)
- add gitee url
- update doc
- update bookstore example for generation prototype
- update doc
- remove wechat image
- update wechat qrcode
- support https in rest
- update wechat qrcode
- update readme
- add images back because of gitee not showing
- add images back because of gitee not showing

### Bug Fixes

- fix url 404 (#180)

# 1.0.23 (2020/10/28)

### Features

- goctl add stdin flag (#170)
- update doc using raw images
- update deployment version (#165)
- refactor middleware generator (#159)
- gen api svc add middleware implement temp code (#151)
- add logo in readme
- api handler generate incompletely while has no request  (#158)
- update api template (#156)
- add vote link
- docs: format markdown and add go mod in demo (#155)
- ignore blank between bracket and service tag (#154)
- model support globbing patterns (#153)

### Bug Fixes

- model template fix (#169)
- spell fix (#167)
- fix bug: generate incomplete model code  in case findOneByField (#160)

# 1.0.22 (2020/10/21)

### Features

- make tests faster
- update wechat info
- can only specify one origin in cors
- make tests faster
- gozero template (#147)
- support cors in rest server
- optimized generator formatted code (#148)
- export WithUnaryClientInterceptor
- let balancer to be customizable
- rename NewPatRouter to NewRouter
- use goctl template to generate all kinds of templates
- api add middleware support (#140)
- add more tests
- add redis Zrevrank (#137)

### Bug Fixes

- fix zrpc client interceptor calling problem
- fix name typo and format with newline (#143)

# 1.0.21 (2020/10/17)

### Features

- update breaker doc
- to correct breaker interface annotation (#136)
- update doc
- update doc
- add logx.Alert
- add fx.Split
- add anonymous annotation (#134)
- update goctl rpc template log print url (#133)
- print more message when parse error (#131)
- delete goctl rpc main tpl no use import (#130)
- update doc
- support api templates
- assert len > 0
- fail fast when rolling window size is zero
- simplify code generation
- stop rpc server when main function exit (#120)
- faster the tests
- Gozero sqlgen patch (#119)
- update readme
- add qq qrcode
- update readme
- make tests faster
- Goctl rpc patch (#117)
- make tests race-free

### Bug Fixes

- fix golint issues
- fix golint issues
- fix: fx/fn.Head func will forever block when n is less than 1 (#128)
- fix syncx/barrier test case (#123)
- fix: template cache key (#121)

# 1.0.20 (2020/10/10)

### Features

- avoid bigint converted into float64 when unmarshaling
- add more tests
- add more tests
- parser ad test (#116)
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add fx.Count
- add more tests

### Bug Fixes

- fix data race in tests
- fix data race in tests

# 1.0.19 (2020/10/04)

### Features

- add more tests
- breaker: remover useless code (#114)
- update wechat qrcode
- update codecov settings
- add more tests
- add more tests
- add more tests
- add more tests
- remove markdown linter
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- GetBreaker need double-check (#112)
- add more tests
- perfect the bookstore and shorturl doc (#109)
- better lock practice in sharedcalls

### Bug Fixes

- fix data race
- fix data race
- fix int64 primary key problem
- update: fix wrong word (#110)

# 1.0.18 (2020/09/29)

### Features

- update shorturl doc
- export cache package, add client interceptor customization
- add zrpc client interceptor
- doc: update sharedcalls.md layout (#107)
- refactor and rename folder to group (#106)
- add api doc link
- use default mongo db (#103)
- unmarshal should be struct
- Revert "goreportcard not working, remove it temporarily"

### Bug Fixes

- add unit test, fix interceptor bug

# 1.0.17 (2020/09/28)

### Features

- goreportcard not working, remove it temporarily
- support return () syntax (#101)
- rename prommetric to prometheous, add unit tests
- update wechat and etcd yaml
- update example
- export AddOptions, AddStreamInterceptors, AddUnaryInterceptors
- query from cache first when do cache.Take
- rename (#98)
- add more clear error when rpc service is not started

### Bug Fixes

- fix typo of prometheus
- fix bug: module parse error (#97)

# 1.0.16 (2020/09/22)

### Features

- add test (#95)
- goctl support import api file (#94)
- add tracing logs in server side and client side
- remove unnecessary tag
- use options instead of opts in error message
- 修改不能编辑代码注释 (#92)
- feature: goctl jwt  (#91)
- update doc (#90)
- remove no need (#87)
- add trace/span in http logs
- use package level defined contextKey as context key
- printing context key friendly
- use contextType as string type
- rename ngin to rest in goctl
- optimize AtomicError (#82)
- update doc
- update rpc example

### Bug Fixes

- fix bug: release empty struct limit (#96)
- fix redis error (#88)
- chore: fix typos (#85)
- fix: golint: context.WithValue should should not use basic type as key (#83)
- fix rpc client examle (#81)

# 1.0.15 (2020/09/18)

### Features

- rename rpcx to zrpc
- update wechat qrcode

### Bug Fixes

- fix example tracing edge config (#76)

# 1.0.14 (2020/09/16)

### Features

- add more tests
- remove markdown linter temporarily
- simplify mapreduce code
- rename file and function name (#74)
- add mapping readme (#75)
- print message when starting api server
- rename function
- optimize: api generating for idea plugin (#68)
- api support empty request or empty response (#72)
- optimize route parse (#70)
- Sharedcalls.md (#69)
- optimized api new with absolute path like: goctl api new $PWD/xxxx (#67)
- update doc, add metric link
- drain pipe if reducer not drained
- Metric (#65)
- Markdown lint (#58)
- update goctl makefile

### Bug Fixes

- fix goctl api (#71)

# 1.0.13 (2020/09/11)

### Features

- add model&rpc doc (#62)
- update doc (#64)
- update quick start (#63)
- add in-process cache doc
- add fast create api demo service (#59)
- quickly generating rpc demo service (#60)
- update wechat image

### Bug Fixes

- fix goctl model (#61)
- fix GOMOD env fetch bug (#55)

# 1.0.12 (2020/09/08)

### Features

- add (#54)
- make chinese readme as default
- update readme to add mapreduce link
- add mr tool doc (#50)
- refactor (#49)
- refactor (#48)
- refactor gomod logic (#47)
- add unit test for mapreduce
- add language link
- add bookstore english tutorial
- add language link
- add shorturl english tutorial
- add english readme
- add goctl description
- update readme
- update example doc
- add wechat qrcode
- add bookstore example
- add bookstore example

### Bug Fixes

- fix goctl model path (#53)
- fix command run path bug (#52)
- fix typoin doc
- fix mapreduce problem when reducer doesn't write
- fix readme typo
- fix typo (#38)
- fix LF (#37)
- fix bookstore example

# 1.0.11 (2020/09/03)

### Features

- replace clickhouse driver to the official one

# 1.0.10 (2020/09/03)

### Features

- refactor code
- refactor

### Bug Fixes

- fix bug: miss time import (#36)
- fix shorturl example code (#35)

# 1.0.9 (2020/09/02)

### Features

- add shorturl example code
- support go 1.13
- update shorturl doc
- update shorturl doc
- trim space  (#31)

### Bug Fixes

- fix: root path on windows bug (#34)

# 1.0.8 (2020/09/01)

### Features

- remove no need empty line (#29)
- update readme
- update docs
- update shorturl doc
- make svcCtx as a member for better code generation
- remove files
- remove makefile generation
- update readme
- update goctl makefile
- make goctl work on linux
- update shorturl doc

### Bug Fixes

- fix doc errors
- fix dockerfile generation
- fix typo in doc
- fix doc error

# 1.0.7 (2020/08/29)

### Features

- update shorturl doc
- reorg imports
- rpc generation support windows (#28)
- update shorturl doc
- update shorturl doc
- update handler generation (#27)
- refine rpc generator
- refine goctl rpc generator
- rpc service generation (#26)
- update shorturl doc
- update shorturl doc
- better image rendering
- add quick example
- sort imports on api generation
- return zero value instead of nil on generated logic
- disable cpu stat in wsl linux
- use yaml, and detect go.mod in current dir
- move test code into internal package
- update ci configuration
- add more tests
- add more tests
- add more tests
- add more tests

### Bug Fixes

- fix config yaml gen (#25)
- fix ci script
- fix format (#23)

# 1.0.6 (2020/08/25)

### Features

- use predefined endpoint separator
- add fatal to stderr
- add etcd deploy yaml

# 1.0.5 (2020/08/24)

### Features

- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- accelerate tests
- make tests parallel
- add more tests
- update readme
- update doc, add architecture picture
- make test stable
- update readme
- add more tests
- add more tests
- add more tests
- reorg imports
- update keywords.md
- gocctl model v20200819 (#18)
- update doc

### Bug Fixes

- fix generate api demo (#19)

# 1.0.4 (2020/08/19)

### Features

- update stringx doc
- update stringx doc
- update image alt in doc
- update image scale in doc
- add keywords utility example
- add release badge
- add go report badge
- support customized mask char on trie
- add benchmark
- goctl model  reactor (#15)
- goctl生成Kotlin代码优化 (#16)
- rename rest files
- support direct scheme on rpc resolver
- remove utils.Report
- rename-Api
- multi-http-method-support
- remove-logx
- Add goctl kotlin support
- refactor names
- add license badge
- update codecov settings
- add codecov report
- add codecov badge
- use default decay value from finagle
- format readme
- update doc

### Bug Fixes

- fix render problem in doc
- fix golint warnings
- fix-log-fatal
- fix FileNotFoundException when response code is 4xx or 5xx
- fix-lang-must-not-found
- fix-break-line
- chore: fix typo

# 1.0.3 (2020/08/14)

### Features

- move lang.Must into logx.Must to make sure output fatal message as json
- confirm addition after add called in periodical executor
- add more tests

# 1.0.2 (2020/08/14)

### Bug Fixes

- fix data race

# 1.0.1 (2020/08/13)

### Features

- add queue package
- remove pdf
- update workflow
- update readme
- rename files
- update readme
- rename files
- auto generate go mod if need
- remove bodyless check
- parse body only if content length > 0
- return ErrBodylessRequest on get method etc.
- remove bodyless check
- parse body only if content length > 0
- return ErrBodylessRequest on get method etc.
- export httpx.GetRemoteAddr
- export router
- export token parser for refresh token service
- remove unused method
- export httpx.GetRemoteAddr
- export router
- export token parser for refresh token service
- update readme.md
- remove unused method
- update readme.md
- use strings.Contains instead of strings.Index
- refactor rpcx, export WithDialOption and WithTimeout
- move auth interceptor into serverinterceptors
- use fmt.Println instead of println
- remove unused method

### Bug Fixes

- fix windows slash
- fix windows bug
- fix windows slash
- fix windows bug
- fix windows slash

# 1.0.0 (2020/08/11)

# 0.1.1 (2022/12/08)

### Features

- Merge latest code
- feat: vben code generation via api file
- chore: tidy go.sum (#2675)
- chore: update deps (#2674)
- wip: vben code generation
- chore: upgrade dependencies (#2658)
- Fixes #2603 bump goctl cobra version to macos completion help bug (#2656)
- feature : responses whit context (#2637)
- feat: add trace.SpanIDFromContext and trace.TraceIDFromContext (#2654)
- replace strings.Title to cases.Title (#2650)
- The default port is used when there is no port number for k8s (#2598)
- feat: add stringx.ToCamelCase (#2622)
- chore: update deps (#2621)
- feat: validate value in options for mapping (#2616)
- chore: update dependencies (#2594)
- feat: support bool for env tag (#2593)
- feat: support env tag in config (#2577)

### Bug Fixes

- fix: update ErrorCtx logic
- fix: fix client side in zeromicro#2109 (zeromicro#2116) (#2659)
- fix: log currentSize should not be 0 when file exists and size is not 0 (#2639)
- fix(rest): fix issues#2628 (#2629)
- fix: fix conflict with the import package name (#2610)

# 0.1.0 (2022/12/03)

### Features

- feat: service port parameter
- feat: service port parameter
- feat: api crud generation by proto
- feat: auto migrate for rpc generation
- wip: api code generation
- wip: api code generation

### Bug Fixes

- fix: optimize api url

# 0.1.0-beta-1 (2022/11/27)

# 0.1.0-beta (2022/12/03)

### Features

- feat: service port parameter
- feat: api crud generation by proto
- feat: auto migrate for rpc generation
- wip: api code generation
- wip: api code generation
- feat: generate docker file
- feat: proto file generation and logic code generation with ent
- feat: proto file generation and logic code generation with ent
- wip: ent logic generating

### Bug Fixes

- fix: optimize api url
- fix: delete ent in tools

# 0.0.9 (2022/11/12)

### Features

- Update change log

### Bug Fixes

- fix: modify error code

# 0.0.8 (2022/11/11)

### Features

- merge latest code
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.3 to 1.11.0 (#2588)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.23.0 to 2.23.1 (#2587)
- chore(deps): bump github.com/jhump/protoreflect from 1.13.0 to 1.14.0 (#2579)
- feat: error translation
- feat: conf inherit (#2568)
- Modify comment syntax error (#2572)
- feat: gorm logger

### Bug Fixes

- fix: optimize go gen types
- fix: inherit issue when parent after inherits (#2586)
- fix(change model template file type): All model template variables ar… (#2573)
- fix(goctl): Fix #2561 (#2562)
- fix: potential slice append issue (#2560)

# 0.0.7.6 (2022/10/28)

### Features

- delete github file
- perf: optimize route swagger generation
- update `genstruct` with swagger annotion
- chore: update "DO NOT EDIT" format (#2559)
- add generate swagger annotation option in gogen
- update `Accept-Language` parser
- chore(action): enable cache dependency (#2549)
- add swagger support add gorm support add validate support
- chore: add more tests (#2547)
- feat: add logger.WithFields (#2546)
- chore: refactor (#2545)

### Bug Fixes

- fix: rocketmq config add optional tag
- fix: change default file name into snake format
- fix: producer and consumer pointer error
- fix(goctl): fix redundant import (#2551)
- typo(mapping): fix typo for key (#2548)

# 0.0.7.3-beta (2022/10/26)

### Features

- feat: rocket mq plugin
- feat(trace): support for disabling tracing of specified `spanName` (#2363)
- chore: adjust rpc comment format (#2501)
- feat: remove info log when disable log (#2525)
- chore(action): upgrade action (#2521)
- feat: support uuid.UUID in mapping (#2537)
- chore: add more tests (#2536)
- Fix typo (#2531)
- chore(deps): bump google.golang.org/grpc from 1.50.0 to 1.50.1 (#2527)
- Fix the wrong key about FindOne in mongo of goctl. (#2523)
- chore: add golangci-lint config file (#2519)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2514)
- Fix mongo insert tpl (#2512)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2511)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2510)
- chore: sqlx's metric name is different from redis (#2505)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.10.0 to 1.11.0 (#2504)

### Bug Fixes

- fix: redis's pipeline logs are not printed completely (#2538)
- fix(goctl): Fix issues (#2543)
- chore: fix lint errors (#2520)
- chore: fix naming problem (#2500)

# 0.0.7.2-beta (2022/10/22)

### Bug Fixes

- fix: bugs in accept language parsing

# 0.0.7.1-beta (2022/10/22)

### Bug Fixes

- fix: bugs in accept language parsing
- fix: delete log message reference from simple-admin-core

# 0.0.7 (2022/10/12)

### Features

- chore: remove unnecessary code (#2499)
- token limit support context (#2335)
- feat(goctl): better generate the api code of typescript (#2483)
- chore(deps): bump google.golang.org/grpc from 1.49.0 to 1.50.0 (#2487)
- chore: remove init if possible (#2485)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.2 to 1.10.3 (#2484)

### Bug Fixes

- fix: etc template
- fix: replace Infof() with Errorf() in DurationInterceptor (#2495) (#2497)
- fix a few function names on comments (#2496)
- fix: update deployment in k8s

# 0.0.7-beta-3 (2022/10/09)

### Bug Fixes

- fix: StackCoolDownMillis name

# 0.0.7-beta-2 (2022/10/09)

### Features

- revert: remove consul yaml config

# 0.0.7-beta-1 (2022/10/09)

### Bug Fixes

- fix: yaml key name

# 0.0.7-beta (2022/10/09)

### Features

- revert: cancel the consul and use k8s in generation
- chore: refactor to reduce duplicated code (#2477)
- chore: better shedding algorithm, make sure recover from shedding (#2476)
- feat(redis):add timeout method to extend blpop (#2472)
- chore: sort methods (#2470)
- feat: add logc package, support AddGlobalFields for both logc and logx. (#2463)
- feat: add string to map in httpx parse method (#2459)
- Readme Tweak (#2436)
- refactor: adjust http request slow log format (#2440)
- chore: gofumpt (#2439)
- feat: add color to debug (#2433)

### Bug Fixes

- fix: JSON tag in config files
- fix: system info in swagger
- fix(mongo): fix file name generation errors (#2479)
- fix: etcd reconnecting problem (#2478)
- fix #2343 (#2349)
- fix: add more tests (#2473)
- fix(goctl): fix the unit test bug of goctl (#2458)
- fix #2435 (#2442)
- fix: bugs when run goctls new

# 0.0.6 (2022/09/23)

### Features

- feat: gen consul code

# 0.0.6-beta-7 (2022/09/23)

### Bug Fixes

- fix: rest inline bug

# 0.0.6-beta-6 (2022/09/23)

### Bug Fixes

- fix: inline bug
- fix: restore field for consul config
- fix: json field for consul conf

# 0.0.6-beta-5 (2022/09/23)

### Bug Fixes

- fix: add yaml tag for all configuration
- fix: bug in load
- fix: change interface into pointer

# 0.0.6-beta-4 (2022/09/23)

### Bug Fixes

- fix: load function circle implement

# 0.0.6-beta-3 (2022/09/23)

### Features

- feat: consul kv store configuration

# 0.0.6-beta-2 (2022/09/22)

# 0.0.6-beta-1 (2022/09/22)

### Features

- feat: consul support

### Bug Fixes

- fix: update change log

# 0.0.5 (2022/09/21)

### Features

- chore: replace fmt.Fprint (#2425)
- cleanup: deprecated field and func (#2416)
- refactor: redis error for prometheus metric label (#2412)

### Bug Fixes

- fix: fix log output (#2424)

# 0.0.5-beta-2 (2022/09/20)

### Bug Fixes

- fix: package access

# 0.0.5-beta-1 (2022/09/20)

### Features

- feat: add log debug level (#2411)
- chore: add more tests (#2410)
- feat(goctl):Add ignore-columns flag (#2407)
- chore: add more tests (#2409)
- chore: update go-zero to v1.4.1
- feat: support caller skip in logx (#2401)
- chore: refactor the imports (#2406)
- feat: mysql and redis metric support (#2355)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2402)
- chore(deps): bump go.etcd.io/etcd/client/v3 from 3.5.4 to 3.5.5 (#2395)
- chore(deps): bump go.etcd.io/etcd/api/v3 from 3.5.4 to 3.5.5 (#2394)
- chore(deps): bump github.com/jhump/protoreflect from 1.12.0 to 1.13.0 (#2393)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2389)
- chore: refactor (#2388)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2385)
- feat: add grpc export (#2379)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 (#2383)

### Bug Fixes

- fix: add validator test
- fix goctl help message (#2414)

# 0.0.5-beta (2022/09/19)

### Features

- refactor: change api error pkg

### Bug Fixes

- fix: add validator

# 0.0.4 (2022/09/13)

### Features

- feat: support targetPort option in goctl kube (#2378)
- Update readme-cn.md
- chore(deps): bump github.com/lib/pq from 1.10.6 to 1.10.7 (#2373)
- feat: support baggage propagation in httpc (#2375)
- chore(deps): bump go.uber.org/goleak from 1.1.12 to 1.2.0 (#2371)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.1 to 1.10.2 (#2370)
- chore: refactor (#2365)

### Bug Fixes

- fix #2364 (#2377)
- fix: issue #2359 (#2368)
- fix:trace graceful stop,pre loss trace (#2358)

# 0.0.3 (2022/09/06)

### Bug Fixes

- fix: recover go gen type

# 0.0.2 (2022/09/06)

### Features

- correct test case (#2340)
- Hidden java (#2333)
- make logx#getWriter concurrency-safe (#2233)
- generates nested types in doc (#2201)
- Add strict flag (#2248)
- improve: number range compare left and righ value (#2315)
- refactor: sequential range over safemap (#2316)
- safemap add Range method (#2314)
- chore: remove unused packages (#2312)
- Fix/del server interceptor duplicate copy md 20220827 (#2309)
- chore: refactor gateway (#2303)
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.3 to 2.0.5 (#2305)

### Bug Fixes

- fix:etcd get&watch not atomic (#2321)
- fix: thread-safe in getWriter of logx (#2319)
- fix: range validation on mapping (#2317)
- fix: handle the scenarios that content-length is invalid (#2313)
- fix: more accurate panic message on mapreduce (#2311)
- fix #2301,package conflict generated by ddl (#2307)
- fix: logx disable not working in some cases (#2306)
- fix:duplicate copy MD (#2304)

# 0.0.1 (2022/08/26)

### Features

- feat: add go swagger support
- feat: casbin util
- feat: error message
- feat: gorm conf

### Bug Fixes

- fix: gen system info
- fix: gen types swagger doc
- fix: bug in rest response
- fix: error msg
- fix: package name

# tools/goctl/v1.4.3 (2022/12/13)

### Features

- Add more test (#2692)
- feat: add dev server and health (#2665)
- feat: accept camelcase for config keys (#2651)
- chore: tidy go.sum (#2675)
- chore: update deps (#2674)
- chore: upgrade dependencies (#2658)
- Fixes #2603 bump goctl cobra version to macos completion help bug (#2656)
- feature : responses whit context (#2637)
- feat: add trace.SpanIDFromContext and trace.TraceIDFromContext (#2654)
- replace strings.Title to cases.Title (#2650)
- The default port is used when there is no port number for k8s (#2598)
- feat: add stringx.ToCamelCase (#2622)
- chore: update deps (#2621)
- feat: validate value in options for mapping (#2616)
- chore: update dependencies (#2594)
- feat: support bool for env tag (#2593)
- feat: support env tag in config (#2577)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.3 to 1.11.0 (#2588)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.23.0 to 2.23.1 (#2587)
- chore(deps): bump github.com/jhump/protoreflect from 1.13.0 to 1.14.0 (#2579)
- feat: conf inherit (#2568)
- Modify comment syntax error (#2572)
- chore: update "DO NOT EDIT" format (#2559)
- chore(action): enable cache dependency (#2549)
- chore: add more tests (#2547)

### Bug Fixes

- fix: #2684 (#2693)
- fix: Fix string.title (#2687)
- fix: #2672 (#2681)
- fix:Remove duplicate code (#2686)
- fix: fix client side in zeromicro#2109 (zeromicro#2116) (#2659)
- fix: log currentSize should not be 0 when file exists and size is not 0 (#2639)
- fix(rest): fix issues#2628 (#2629)
- fix: fix conflict with the import package name (#2610)
- fix: inherit issue when parent after inherits (#2586)
- fix(change model template file type): All model template variables ar… (#2573)
- fix(goctl): Fix #2561 (#2562)
- fix: potential slice append issue (#2560)
- fix(goctl): fix redundant import (#2551)
- typo(mapping): fix typo for key (#2548)

# tools/goctl/v1.4.2 (2022/10/22)

### Features

- feat: add logger.WithFields (#2546)
- chore: refactor (#2545)
- feat(trace): support for disabling tracing of specified `spanName` (#2363)
- chore: adjust rpc comment format (#2501)
- feat: remove info log when disable log (#2525)
- chore(action): upgrade action (#2521)
- feat: support uuid.UUID in mapping (#2537)
- chore: add more tests (#2536)
- Fix typo (#2531)
- chore(deps): bump google.golang.org/grpc from 1.50.0 to 1.50.1 (#2527)
- Fix the wrong key about FindOne in mongo of goctl. (#2523)
- chore: add golangci-lint config file (#2519)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2514)
- Fix mongo insert tpl (#2512)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2511)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2510)
- chore: sqlx's metric name is different from redis (#2505)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.10.0 to 1.11.0 (#2504)
- chore: remove unnecessary code (#2499)
- token limit support context (#2335)
- feat(goctl): better generate the api code of typescript (#2483)
- chore(deps): bump google.golang.org/grpc from 1.49.0 to 1.50.0 (#2487)
- chore: remove init if possible (#2485)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.2 to 1.10.3 (#2484)
- chore: refactor to reduce duplicated code (#2477)
- chore: better shedding algorithm, make sure recover from shedding (#2476)
- feat(redis):add timeout method to extend blpop (#2472)
- chore: sort methods (#2470)
- feat: add logc package, support AddGlobalFields for both logc and logx. (#2463)
- feat: add string to map in httpx parse method (#2459)
- Readme Tweak (#2436)
- refactor: adjust http request slow log format (#2440)
- chore: gofumpt (#2439)
- feat: add color to debug (#2433)
- chore: replace fmt.Fprint (#2425)
- cleanup: deprecated field and func (#2416)
- refactor: redis error for prometheus metric label (#2412)
- feat: add log debug level (#2411)
- chore: add more tests (#2410)
- feat(goctl):Add ignore-columns flag (#2407)
- chore: add more tests (#2409)

### Bug Fixes

- fix: redis's pipeline logs are not printed completely (#2538)
- fix(goctl): Fix issues (#2543)
- chore: fix lint errors (#2520)
- chore: fix naming problem (#2500)
- fix: replace Infof() with Errorf() in DurationInterceptor (#2495) (#2497)
- fix a few function names on comments (#2496)
- fix(mongo): fix file name generation errors (#2479)
- fix: etcd reconnecting problem (#2478)
- fix #2343 (#2349)
- fix: add more tests (#2473)
- fix(goctl): fix the unit test bug of goctl (#2458)
- fix #2435 (#2442)
- fix: fix log output (#2424)
- fix goctl help message (#2414)

# tools/goctl/v1.4.1 (2022/09/17)

### Features

- chore: update go-zero to v1.4.1
- feat: support caller skip in logx (#2401)
- chore: refactor the imports (#2406)
- feat: mysql and redis metric support (#2355)
- chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#2402)
- chore(deps): bump go.etcd.io/etcd/client/v3 from 3.5.4 to 3.5.5 (#2395)
- chore(deps): bump go.etcd.io/etcd/api/v3 from 3.5.4 to 3.5.5 (#2394)
- chore(deps): bump github.com/jhump/protoreflect from 1.12.0 to 1.13.0 (#2393)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2389)
- chore: refactor (#2388)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2385)
- feat: add grpc export (#2379)
- chore(deps): bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 (#2383)
- feat: support targetPort option in goctl kube (#2378)
- Update readme-cn.md
- chore(deps): bump github.com/lib/pq from 1.10.6 to 1.10.7 (#2373)
- feat: support baggage propagation in httpc (#2375)
- chore(deps): bump go.uber.org/goleak from 1.1.12 to 1.2.0 (#2371)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.1 to 1.10.2 (#2370)
- chore: refactor (#2365)
- correct test case (#2340)
- Hidden java (#2333)
- make logx#getWriter concurrency-safe (#2233)
- generates nested types in doc (#2201)
- Add strict flag (#2248)
- improve: number range compare left and righ value (#2315)
- refactor: sequential range over safemap (#2316)
- safemap add Range method (#2314)
- chore: remove unused packages (#2312)
- Fix/del server interceptor duplicate copy md 20220827 (#2309)
- chore: refactor gateway (#2303)
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.3 to 2.0.5 (#2305)
- chore: refactor stat (#2299)
- Initialize CPU stat code only if used (#2020)
- chore(deps): bump google.golang.org/grpc from 1.48.0 to 1.49.0 (#2297)
- feat(redis): add ZaddFloat & ZaddFloatCtx (#2291)
- doc(readme): add star history (#2275)
- feat: rpc add health check function configuration optional (#2288)
- Update readme-cn.md
- Update issues.yml
- chore: Update readme (#2280)
- Update readme-cn.md
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.2 to 2.0.3 (#2267)
- chore: refactor logx (#2262)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.22.0 to 2.23.0 (#2260)
- test: add more tests (#2261)
- chore(deps): bump github.com/prometheus/client_golang (#2244)
- chore(deps): bump github.com/fullstorydev/grpcurl from 1.8.6 to 1.8.7 (#2245)

### Bug Fixes

- fix #2364 (#2377)
- fix: issue #2359 (#2368)
- fix:trace graceful stop,pre loss trace (#2358)
- fix:etcd get&watch not atomic (#2321)
- fix: thread-safe in getWriter of logx (#2319)
- fix: range validation on mapping (#2317)
- fix: handle the scenarios that content-length is invalid (#2313)
- fix: more accurate panic message on mapreduce (#2311)
- fix #2301,package conflict generated by ddl (#2307)
- fix: logx disable not working in some cases (#2306)
- fix:duplicate copy MD (#2304)
- fix resource manager dead lock (#2302)
- fix #2163 (#2283)
- fix(logx): display garbled characters in windows(DOS, Powershell) (#2232)
- fix #2240 (#2271)
- fix #2240 (#2263)
- fix: unsignedTypeMap type error (#2246)
- fix: test failure, due to go 1.19 compatibility (#2256)
- fix: time repr wrapper (#2255)

# tools/goctl/v1.4.0 (2022/08/07)

### Features

- chore: release action for goctl (#2239)
- Update readme-cn.md
- Update readme.md
- feat: more meaningful error messages, close body on httpc requests (#2238)
- Update readme.md
- Update readme.md
- docs: update docs for gateway (#2236)
- chore: renaming configs (#2234)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.10.0 to 1.10.1 (#2225)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2222)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2223)
- chore: refactor redislock (#2210)
- feat(redislock): support set context (#2208)
- chore(deps): bump google.golang.org/protobuf from 1.28.0 to 1.28.1 (#2205)
- support mulitple protoset files (#2190)
- chore: let logx.SetWriter can be called anytime (#2186)
- chore: refactoring (#2182)
- chore: refactoring logx (#2181)
- feat: logx support logs rotation based on size limitation. (#1652) (#2167)
- Update goctl version (#2178)
- feat: Support for multiple rpc service generation and rpc grouping (#1972)
- Update readme-cn.md
- Update api template (#2172)
- chore: refactoring mapping name (#2168)
- feat: support customized header to metadata processor (#2162)
- feat: support google.api.http in gateway (#2161)
- feat: set content-type to application/json (#2160)
- feat: verify RpcPath on startup (#2159)
- feat: support form values in gateway (#2158)
- feat: export gateway.Server to let users add middlewares (#2157)
- Update readme-cn.md
- Update readme.md
- feat: restful -> grpc gateway (#2155)
- chore(deps): bump google.golang.org/grpc from 1.47.0 to 1.48.0 (#2147)
- chore(deps): bump go.mongodb.org/mongo-driver from 1.9.1 to 1.10.0 (#2150)
- chore: remove unimplemented gateway (#2139)
- docs: update goctl readme (#2136)
- chore: refactor (#2130)
- chore: add more tests (#2129)
- feat:Add `Routes` method for server (#2125)
- feat: support logx.WithFields (#2128)
- feat: add Wrap and Wrapf in errorx (#2126)
- chore: coding style (#2120)

### Bug Fixes

- fix: #2216 (#2235)
- fix(logx): need to wait for the first caller to complete the execution. (#2213)
- fix: fix comment typo (#2220)
- fix: handling rpc error on gateway (#2212)
- fix: only setup logx once (#2188)
- fix: logx test foo (#2144)
- fix(httpc): fix typo errors (#2189)
- fix: remove invalid log fields in notLoggingContentMethods (#2187)
- fix:duplicate route check (#2154)
- fix: fix switch doesn't work bug (#2183)
- fix: fix #2102, #2108 (#2131)
- fix: goctl genhandler duplicate rest/httpx & goctl genhandler template support custom import httpx package (#2152)
- fix: generated sql query fields do not match template (#2004)
- fix goctl rpc protoc strings.EqualFold Service.Name GoPackage (#2046)

# tools/goctl/v1.3.9 (2022/07/09)

### Features

- chore: remove blank lines (#2117)
- feat:`goctl model mongo ` add `easy` flag  for easy declare. (#2073)
- refactor:remove duplicate codes (#2101)
- chore(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#2115)
- feat: add method to jsonx (#2049)
- chore(deps): bump go.opentelemetry.io/otel/exporters/zipkin (#2112)
- chore: update goctl version to 1.3.9 (#2111)
- chore: refactor (#2087)
- remove legacy code (#2086)
- chore: refactor (#2085)
- Update readme-cn.md
- remove legacy code (#2084)
- Update readme.md
- feat: CompareAndSwapInt32 may be better than AddInt32 (#2077)
- chore(deps): bump github.com/pelletier/go-toml/v2 from 2.0.1 to 2.0.2 (#2072)
- chore(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.4.2 (#2066)
- chore(deps): bump github.com/stretchr/testify from 1.7.2 to 1.8.0 (#2068)
- chore(deps): bump github.com/alicebob/miniredis/v2 from 2.21.0 to 2.22.0 (#2067)
- chore(deps): bump github.com/ClickHouse/clickhouse-go/v2 (#2064)
- Create dependabot.yml
- [ci skip] Fix dead doc link (#2047)
- chore: remove lifecycle preStop because sh not exist in scratch (#2042)
- Update readme-cn.md
- Update readme-cn.md
- chore: Add command desc & color commands (#2013)
- Update readme.md
- Update readme.md
- feat: rest.WithChain to replace builtin middlewares (#2033)
- Update readme-cn.md
- chore: refactor to simplify disabling builtin middlewares (#2031)
- add user middleware chain function (#1913)
- Add fig (#2008)
- feat: support build Dockerfile from current dir (#2021)
- chore: upgrade action version (#2027)
- feat: add trace in httpc (#2011)
- Update readme-cn.md
- chore: coding style (#2012)
- Update readme.md
- Update readme.md
- feat: Replace mongo package with monc & mon (#2002)
- feat: convert grpc errors to http status codes (#1997)
- chore: rename methods (#1998)
- periodlimit new function TakeWithContext (#1983)
- typo: add type keyword (#1992)
- feat: add 'imagePullPolicy' parameter for 'goctl kube deploy' (#1996)
- chore: update dependencies (#1985)

### Bug Fixes

- fix #2109 (#2116)
- fix: type matching supports string to int (#2038)
- fix concurrent map writes (#2079)
- fix 当表有唯一键时，update()的形参和实参不一致 (#2010)
- fix: `\u003cnil\u003e` log output when http server shutdown. (#2055)
- fix:typo in readme.md (#2061)
- fix: quickstart wrong package when go.mod exists in parent dir (#2048)
- fix #1977 (#2034)
- fix: 修复 clientinterceptors/tracinginterceptor.go 显示接受消息字节为0 (#2003)
- fix goctl api clone template fail (#1990)

# tools/goctl/v1.3.8 (2022/06/06)

### Features

- chore: update goctl version to 1.3.8 (#1981)
- Fix pg subcommand level error (#1979)

### Bug Fixes

- fix: generate bad Dockerfile on given dir (#1980)

# tools/goctl/v1.3.7 (2022/06/05)

### Features

- chore: make methods consistent in signatures (#1971)
- test: make tests stable (#1968)
- chore: make print pretty (#1967)
- chore: better mongo logs (#1965)
- feat: print routes (#1964)
- chore: update dependencies (#1963)
- Update readme-cn.md
- Update readme.md
- Chore/goctl version (#1962)

### Bug Fixes

- fix: The validation of tag "options" is not working with int/uint type (#1969)
- test: fix fails (#1970)
- 🐞 fix: fixed typo (#1916)

# tools/goctl/v1.3.6 (2022/06/03)

### Features

- chore: update version (#1961)
- Update readme-cn.md
- Update readme.md
- Update readme.md
- docs: add docs for logx (#1960)
- chore: refactoring mapping string to slice (#1959)
- chore: update roadmap (#1948)
- Delete duplicated crash recover logic. (#1950)
- Update readme-cn.md
- Update readme-cn.md
- chore: refine docker for better compatible with package main (#1944)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Update readme-cn.md
- Update readme.md
- Update readme-cn.md
- 优化代码 (#1939)
- core/mr:a little optimization for collector initialization in ForEach function (#1937)
- chore(action): simplified release configuration (#1935)
- chore: add release action to auto build binaries (#1884)
- feat: update docker alpine package mirror (#1924)
- Support built-in shorthand flags (#1925)
- feat: set default connection idle time for grpc servers (#1922)
- chore: update k8s.io/client-go for security reason, go is upgrade to 1.16 (#1912)
- Update readme-cn.md
- Update readme-cn.md
- Fix process blocking problem during check (#1911)
- feat: support WithStreamClientInterceptor for zrpc clients (#1907)
- Update FUNDING.yml
- chore: use get for quickstart, plain logs for easy understanding (#1905)
- use goproxy properly, remove files (#1903)
- feat: add toml config (#1899)
- chore: coding style for quickstart (#1902)
- refactor: refactor trace in redis & sql & mongo (#1865)
- feat: Add goctl quickstart (#1889)
- Fix code generation (#1897)
- Update readme-cn.md
- Update readme-cn.md
- chore: improve codecov (#1878)
- chore: update some logs (#1875)
- feat: logx with color (#1872)
- feat: Replace cli to cobra (#1855)
- Update readme.md
- Update readme-cn.md
- Update readme.md
- add conf documents (#1869)
- chore: refine tests (#1864)
- test: add codecov (#1863)
- test: add codecov (#1861)
- chore: use time.Now() instead of timex.Time() because go optimized it (#1860)
- feat: add fields with logx methods, support using third party logging libs. (#1847)
- test: add more tests (#1856)
- docs: update readme (#1849)

### Bug Fixes

- fix: panic on convert to string on fillSliceFromString() (#1951)
- fix: Useless delete cache logic in update (#1923)
- fix ts tpl (#1879)
- fix:tools/goctl/rpc/generator/template_test.go file has wrong parameters (#1882)
- chore: fix deprecated usages (#1871)
- fix time, duration, slice types on logx.Field (#1868)
- fix typo (#1857)

# tools/goctl/v1.3.5 (2022/04/28)

### Features

- refactor: move json related header vars to internal (#1840)
- feat: Support model code generation for multi tables (#1836)
- Update readme-cn.md
- refactor: simplify the code (#1835)
- feat: support sub domain for cors (#1827)
- feat: upgrade grpc to 1.46, and remove the deprecated grpc.WithBalancerName (#1820)
- chore: optimize code (#1818)
- chore: remove gofumpt -s flag, default to be enabled (#1816)
- chore: refactor (#1814)
- feat: add trace in redis & mon & sql (#1799)
- chore: Embed unit test data (#1812)
- add go-grpc_opt and go_opt for grpc new command (#1769)
- feat: use mongodb official driver instead of mgo (#1782)
- update rpc generate sample proto file (#1709)
- feat(goctl): go work multi-module support (#1800)
- docs(goctl): goctl 1.3.4 migration note (#1780)
- chore: use grpc.WithTransportCredentials and insecure.NewCredentials() instead of grpc.WithInsecure (#1798)
- goctl api new should given a service_name explictly (#1688)
- show help when running goctl api without any flags (#1678)
- show help when running goctl docker without any flags (#1679)
- show help when running goctl rpc protoc without any flags (#1683)
- improve goctl rpc new (#1687)
- revert postgres package refactor (#1796)
- refactor: move postgres to pg package (#1781)
- feat: add httpc.Do & httpc.Service.Do (#1775)
- feat(goctl): supports api  multi-level importing (#1747)
- show help when running goctl rpc template without any flags (#1685)
- chore: avoid deadlock after stopping TimingWheel (#1768)
- Fix #1765 (#1767)
- chore: remove legacy code (#1766)
- chore: add doc (#1764)
- add more tests (#1763)
- feat: add goctl docker build scripts (#1760)
- Update readme-cn.md
- Update readme.md
- Update readme.md
- Update readme-cn.md
- feat: support ctx in kv methods (#1759)
- feat: use go:embed to embed templates (#1756)

### Bug Fixes

- fix #1838 (#1839)
- fix: remove deprecated dependencies (#1837)
- fix #1806 (#1833)
- fix: rest: WriteJson get 200 when Marshal failed. (#1803)
- fix: Fix issue #1810 (#1811)
- fix: ignore timeout on websocket (#1802)
- fix: Hdel check result & Pfadd check result (#1801)
- chhore: fix usage typo (#1797)
- fix marshal ptr in httpc (#1789)
- fix(goctl): api/new/api.tpl (#1788)
- fix #1729 (#1783)
- fix bug: crash when generate model with goctl. (#1777)
- fix nil pointer if group not exists (#1773)
- fix: model unique keys generated differently in each re-generation (#1771)
- fix #1754 (#1757)

# tools/goctl/v1.3.4 (2022/04/03)

### Features

- Support `goctl env install` (#1752)
- chore: update go-zero to v1.3.2 in goctl (#1750)
- feat: simplify httpc (#1748)
- feat: return original value of setbit in redis (#1746)
- chore: update goctl version to 1.3.4 (#1742)
- feat: let model customizable (#1738)
- Fix zrpc code generation error with --remote (#1739)
- chore: refactor to use const instead of var (#1731)
- feat(goctl): supports model code 'DO NOT EDIT' (#1728)
- Fix unit test (#1730)
- refactor: guard timeout on API files (#1726)
- Added support for setting the parameter size accepted by the interface and custom timeout and maxbytes in API file (#1713)
- chore: refactor code (#1708)
- feat: remove reentrance in redislock, timeout bug (#1704)
- chore: refactor code (#1700)
- chore: refactor code (#1699)
- feat: add getset command in redis and kv (#1693)
- feat: add httpc.Parse (#1698)
- Add verbose flag (#1696)
- Update LICENSE
- refactor: simplify the code (#1670)
- Remove debug log (#1669)
- feat: support -base to specify base image for goctl docker (#1668)
- Remove unused code (#1667)
- feat: add Dockerfile for goctl (#1666)
- feat: Remove  command `goctl rpc proto`  (#1665)

### Bug Fixes

- fix: model generation bug on with cache (#1743)
- fix(goctl): api format with reader input (#1722)
- fix: empty slice are set to nil (#1702)
- fix -cache=true insert no clean cache (#1672)
- chore: fix lint issue (#1694)
- fix: the new  RawFieldNames considers the tag with options. (#1663)

# tools/goctl/v1.3.3 (2022/03/17)

### Features

- Mkdir if not exists (#1659)
- chore: remove unnecessary env (#1654)
- refactor: httpc package for easy to use (#1645)
- refactor: httpc package for easy to use (#1643)
- FindOneBy 漏 Context (#1642)
- feat: add httpc/Service for convinience (#1641)
- feat: add httpc/Get httpc/Post (#1640)
- feat: add rest/httpc to make http requests governacible (#1638)
- Update ROADMAP.md
- Update ROADMAP.md
- feat: support cpu stat on cgroups v2 (#1636)
- feat: support oracle :N dynamic parameters (#1552)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Support for referencing types in different API files using format (#1630)
- feat: support scratch as the base docker image (#1634)
- chore: reduce the docker image size (#1633)
- Fix #1585 #1547 (#1624)
- chore: update goctl version to 1.3.3, change docker build temp dir (#1621)
- Fix #1609 (#1617)
- Fix #1614 (#1616)
- chore: refactor code (#1613)
- chore: add unit tests (#1615)
- model中db标签增加'-'符号以支持数据库查询时忽略对应字段. (#1612)
- feat(goctl): api dart support flutter v2 (#1603)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- test: add more tests (#1604)
- Update readme-cn.md
- chore: update go-zero to v1.3.1 in goctl (#1599)
- chore: upgrade etcd (#1597)
- Update readme.md
- build: update goctl dependency ddl-parser to v1.0.3 (#1586)
- test: add testcase for FIFO Queue in collection module (#1589)
- Update readme-cn.md
- Update readme-cn.md
- Update readme.md
- Update readme.md
- Fix bug int overflow while build goctl on arch 386 (#1582)
- chore: add goctl command help (#1578)
- Update readme.md
- Update readme.md
- feat: supports `importValue` for more path formats (#1569)
- update goctl to go 1.16 for io/fs usage (#1571)
- feat: support pg serial type for auto_increment (#1563)
- Feature: Add goctl env (#1557)
- feat: log 404 requests with traceid (#1554)
- feat: support ctx in sql model generation (#1551)
- feat: support ctx in sqlx/sqlc, listed in ROADMAP (#1535)
- docs: add go-zero users (#1546)
- ignore context.Canceled for redis breaker (#1545)
- chore: update help message (#1544)
- add the serviceAccount of deployment (#1543)
- chore:use struct pointer (#1538)
- docs: update roadmap (#1537)

### Bug Fixes

- fix(goctl): model method FindOneCtx should be FindOne (#1656)
- typo (#1655)
- fix: typo (#1646)
- fix: Update unix-like path regex (#1637)
- fix(goctl): kotlin code generation (#1632)
- fix(goctl): dart gen user defined struct array (#1620)
- fix: HitQuota should be returned instead of Allowed when limit is equal to 1. (#1581)
- fix: fix(gctl): apiparser_parser auto format (#1607)
- Revert "🐞 fix(gen): pg gen of insert (#1591)" (#1598)
- 🐞 fix(gen): pg gen of insert (#1591)
- fix: goctl api dart support `form` tag (#1596)
- chore: fix data race (#1593)
- fix #1541 (#1542)

# tools/goctl/v1.3.2 (2022/02/14)

### Features

- Update readme-cn.md
- chore: refactor cache (#1532)
- feat: support ctx in `Cache` (#1518)
- chore: goctl format issue (#1531)
- upgrade grpc version (#1530)
- chore: update goctl version to 1.3.2 (#1524)
- refactor: refactor yaml unmarshaler (#1517)
- chore: optimize yaml unmarshaler (#1513)
- chore: make error clearer (#1514)
- feat: update go-redis to v8, support ctx in redis methods (#1507)
- Update readme-cn.md
- feature: Add `goctl completion` (#1505)
- test: change fuzz tests (#1504)
- ci: add test for win (#1503)
- chore: update command comment (#1501)

### Bug Fixes

- fix issue of default migrate version (#1536)
- fix #1525 (#1527)
- fix: fix a typo (#1522)
- fixes typo (#1511)
- fix typo: goctl protoc usage (#1502)

# tools/goctl/v1.3.1-alpha (2022/02/01)

# tools/goctl/v1.3.1 (2022/02/01)

### Features

- docs: update tal-tech to zeromico in docs (#1498)
- chore: update goctl version (#1497)
- feat: add runtime stats monitor (#1496)
- feat: handling panic in mapreduce, panic in calling goroutine, not inside goroutines (#1490)
- Update readme-cn.md
- chore: improve migrate confirmation (#1488)

### Bug Fixes

- fix: goctl not compile on windows (#1500)
- fix: goroutine stuck on edge case (#1495)

# tools/goctl/v1.3.0-beta1 (2022/01/26)

### Features

- chore: update warning message (#1487)
- patch: goctl migrate (#1485)
- chore: update go version for goctl (#1484)

# tools/goctl/v1.3.0-alpha (2022/01/25)

# tools/goctl/v1.3.0 (2022/02/09)

### Features

- refactor: refactor yaml unmarshaler (#1517)
- chore: optimize yaml unmarshaler (#1513)
- chore: make error clearer (#1514)
- feat: update go-redis to v8, support ctx in redis methods (#1507)
- Update readme-cn.md
- feature: Add `goctl completion` (#1505)
- test: change fuzz tests (#1504)
- ci: add test for win (#1503)
- chore: update command comment (#1501)
- docs: update tal-tech to zeromico in docs (#1498)
- chore: update goctl version (#1497)
- feat: add runtime stats monitor (#1496)
- feat: handling panic in mapreduce, panic in calling goroutine, not inside goroutines (#1490)
- Update readme-cn.md
- chore: improve migrate confirmation (#1488)
- chore: update warning message (#1487)
- patch: goctl migrate (#1485)
- chore: update go version for goctl (#1484)
- refactor: rename from tal-tech to zeromicro for goctl (#1481)
- Feature/trie ac automation (#1479)
- chore: optimize string search with Aho–Corasick algorithm (#1476)
- Polish the words in readme.md (#1475)
- docs: add go-zero users (#1473)
- chore: update unauthorized callback calling order (#1469)
- Fix/issue#1289 (#1460)
- patch: save missing templates to disk (#1463)
- Fix/issue#1447 (#1458)
- feat: implement console plain output for debug logs (#1456)
- chore: check interface satisfaction w/o allocating new variable (#1454)
- chore: remove jwt deprecated (#1452)
- feat: 支持redis的LTrim方法 (#1443)
- chore: upgrade dependencies (#1444)
- ci: add translator action (#1441)
- Feature rpc protoc (#1251)
- chore: refactor periodlimit (#1428)
- docs: add go-zero users (#1425)
- docs: add go-zero users (#1424)
- update docs (#1421)
- Fix pg model generation without tag (#1407)
- feat: Add migrate (#1419)
- docs: update install readme (#1417)
- chore: refactor rest/timeouthandler (#1415)
- feat: rename module from tal-tech to zeromicro (#1413)

### Bug Fixes

- fixes typo (#1511)
- fix typo: goctl protoc usage (#1502)
- fix: goctl not compile on windows (#1500)
- fix: goroutine stuck on edge case (#1495)
- fix #1468 (#1478)
- chore: fix typo (#1437)
- remove unnecessary drain, fix data race (#1435)
- fix: mr goroutine leak on context deadline (#1433)
- fix: golint issue (#1423)

# tools/goctl/v1.2.5 (2022/01/03)

### Features

- chore: update go-zero to v1.2.5 (#1410)
- refactor file|path (#1409)
- docs: update roadmap (#1405)
- feat: support tls for etcd client (#1390)
- refactor: optimize fx (#1404)
- docs: update goctl installation command (#1403)
- feat: implement fx.NoneMatch, fx.First, fx.Last (#1402)
- chore: simplify mapreduce (#1401)
- chore: update goctl version (#1394)
- ci: remove 386 binaries (#1393)

### Bug Fixes

- fix #1330 (#1382)
- chore: fix golint issues (#1396)

# tools/goctl/v1.2.4 (2021/12/30)

### Features

- ci: remove windows 386 binary (#1392)
- test: add more tests (#1391)
- feat: Add --remote (#1387)
- feat: support array in default and options tags (#1386)
- docs: add go-zero users (#1381)
- docs: update slack invitation link (#1378)
- chore: update goctl version to 1.2.4 for release tools/goctl/v1.2.4 (#1372)
- Updated MySQL生成表结构体遇到关键字db部分保持原字段名定义 (#1369)
- ci: add release action to auto build binaries (#1371)
- docs: update goctl markdown (#1370)
- feat: support context in MapReduce (#1368)
- chore: add 1s for tolerance in redislock (#1367)
- chore: coding style and comments (#1361)
- chore: optimize `ParseJsonBody` (#1353)
- chose: cancel the assignment and judge later (#1359)
- chore: put error message in error.log for verbose mode (#1355)
- test: add more tests (#1352)
- Revert "排除客户端中断导致的503错误 (#1343)" (#1351)
- feat: treat client closed requests as code 499 (#1350)
- 排除客户端中断导致的503错误 (#1343)
- Update FUNDING.yml
- chore: add tests & refactor (#1346)
- Feature: support adding custom cache to mongoc and sqlc (#1313)
- chore: add comments (#1345)
- chore: update goctl version to 1.2.5 (#1337)
- Update template (#1335)
- Feat goctl bug (#1332)
- feat: tidy mod, update go-zero to latest (#1334)
- feat: tidy mod, update go-zero to latest (#1333)
- chore: refactor (#1331)
- Update types.go (#1314)
- feat: tidy mod, add go.mod for goctl (#1328)
- chore: format code (#1327)
- commit  missing method for redis (#1325)
- docs: add go-zero users (#1323)
- style: format code (#1322)
- chore: rename service context from ctx to svcCtx (#1299)
- docs: add go-zero users (#1294)
- Revert "feat: reduce dependencies of framework by add go.mod in goctl (#1290)" (#1291)
- feat: reduce dependencies of framework by add go.mod in goctl (#1290)
- chore: update cli version (#1287)
- feat: support third party orm to interact with go-zero (#1286)
- Feature api root path (#1261)
- chore: cleanup zRPC retry code (#1280)
- feature(retry): Delete retry mechanism (#1279)
- feat: support %w in logx.Errorf (#1278)
- chore: only allow cors middleware to change headers (#1276)
- chore: avoid superfluous WriteHeader call errors (#1275)
- feat: add rest.WithCustomCors to let caller customize the response (#1274)
- Cli (#1272)
- docs: update readme to use goctl@cli (#1255)
- chore: update goctl version (#1250)
- Revert "Revert "feat: enable retry for zrpc (#1237)"" (#1246)
- Revert "feat: enable retry for zrpc (#1237)" (#1245)
- Duplicate temporary variable (#1244)
- Update template (#1243)
- feat: enable retry for zrpc (#1237)
- chore: remove conf.CheckedDuration (#1235)
- reset link goctl (#1232)
- feat: disable grpc retry, enable it in v1.2.4 (#1233)
- chore: refactor, better goctl message (#1228)
- chore: remove unused const (#1224)
- Update FUNDING.yml
- chore: redislock use stringx.randn replace randomStr func (#1220)
- feat: exit with non-zero code on errors (#1218)
- feat: support CORS, better implementation (#1217)
- Create FUNDING.yml
- chore: refine code (#1215)
- docs: add go-zero users (#1214)
- Fix issue 1205 (#1211)
- feat: support CORS by using rest.WithCors(...) (#1212)
- update dependencies. (#1210)
- test: add more tests (#1209)
- goctl docker command add -version (#1206)
- feat: support customizing timeout for specific route (#1203)
- feat: add NewSessionFromTx to interact with other orm (#1202)
- feat: simplify the grpc tls authentication (#1199)
- feat: use WithBlock() by default, NonBlock can be set in config or WithNonBlock() (#1198)
- chore: remove semicolon for routes of services in api files (#1195)
- chore: update goctl version to 1.2.3, prepare for release (#1193)
- feat: slow threshold customizable in zrpc (#1191)
- feat: slow threshold customizable in rest (#1189)
- feat: slow threshold customizable in sqlx (#1188)
- feat: slow threshold customizable in redis (#1187)
- feat: slow threshold customizable in mongo (#1186)
- feat: slow threshold customizable in redis (#1185)
- docs: update roadmap (#1184)
- feat: support multiple trace agents (#1183)
- feat: let different services start prometheus on demand (#1182)
- refactor: simplify tls config in rest (#1181)
- [update] add plugin config (#1180)
- test: add more tests (#1179)
- docs: update roadmap (#1178)
- docs: add go-zero users (#1176)
- feat: support auth account for etcd (#1174)
- feat: support ssl on zrpc, simplify the config (#1175)
- support RpcClient Vertify With Unilateralism and Mutual (#647)
- docs: add go-zero users (#1172)
- Feature add template version (#1152)
- test: add more tests (#1166)
- chore: reorg imports, format code, make MaxRetires default to 0 (#1165)
- Add grpc retry (#1160)
- test: add more tests (#1163)
- chore: reverse the order of stopping services (#1159)
- docs: update qr code (#1158)
- test: add more tests (#1154)
- test: add more tests (#1150)
- Mark deprecated syntax (#1148)
- test: add more tests (#1149)
- test: add more tests (#1147)
- docs: add go-zero users (#1141)
- test: add more tests (#1138)
- test: add more tests (#1137)
- docs: add go-zero users (#1135)
- test: add more tests (#1134)
- Fix issue #1127 (#1131)
- docs: add go-zero users (#1130)
- chore: refine rpc template in goctl (#1129)
- go-zero/core/hash/hash_test.go  增加测试 TestMd5Hex (#1128)
- Add `opts ...grpc.CallOption` in grpc client (#1122)
- Add request method in http log (#1120)
- update goctl version to 1.2.2 (#1125)
- add cncf landscape (#1123)
- test: add more tests (#1119)
- add more tests (#1115)
- add more tests (#1114)
- add more tests (#1113)
- test: add more tests (#1112)
- feat: opentelemetry integration, removed self designed tracing (#1111)
- docs: update roadmap (#1110)
- feat: reflection grpc service (#1107)
- test: add more tests (#1106)
- Fix the `resources` variable not reset after the resource manager is closed (#1105)
- chore: replace redis.NewRedis with redis.New (#1103)
- chore: mark redis.NewRedis as Deprecated, use redis.New instead. (#1100)
- update grpc package (#1099)
- Update Makefile (#1098)
- ci: add reviewdog (#1096)
- docs: update roadmap (#1094)
- docs: update roadmap (#1093)
- ci: accurate error reporting on lint check (#1089)
- update zero-doc links in readme (#1088)
- docs: change organization from tal-tech to zeromicro in readme (#1087)
- ci: add Lint check on commits (#1086)
- Revert "chore: run unit test with go 1.14 (#1084)" (#1085)
- chore: run unit test with go 1.14 (#1084)
- update goctl api (#1052)
- chore: when run goctl-rpc, the order of proto message aliases should be (#1078)
- coding style (#1083)
- we can use otel.ErrorHandlerFunc instead of custom struct when we update OpenTelemetry to 1.0.0 (#1081)
- update go.mod (#1079)
- Create a symbol link file named protoc-gen-goctl from goctl (#1076)
- update OpenTelemetry to 1.0.0 (#1075)
- update issue templates (#1074)
- Update issue templates
- downgrade golang-jwt to support go 1.14 (#1073)
- Add MustTempDir (#1069)
- upgrade grpc version & replace github.com/golang/protobuf/protoc-gen-go with google.golang.org/protobuf (#1065)
- add repo moving notice (#1062)
- add go-zero users (#1061)
- chore: make comment accurate (#1055)
- mention cncf landscape (#1054)
- add go-zero users (#1051)
- update goctl version to 1.2.1 (#1042)
- remove goctl config command (#1035)
- update k8s.io/client-go etc to use go 1.15 (#1031)
- reorg imports, format code (#1024)
- revert changes
- add go-zero users (#1022)
- rename sharedcalls to singleflight (#1017)
- refactor for better error reporting on sql error (#1016)
- expose sql.DB to let orm operate on it (#1015)
- update codecov settings (#1010)
- refactoring tracing interceptors. (#1009)
- use sdktrace instead of trace for opentelemetry to avoid conflicts (#1005)
- api文件中使用group时生成的handler和logic的包名应该为group的名字 (#545)
- add api template file (#1003)
- add opentelemetry test (#1002)
- reorg imports, format code (#1000)
- 开启otel后，tracelog自动获取otle的traceId和spanId (#946)
- optimize unit test (#999)
- rest log with context (#998)
- refactor to shorter config name (#997)
- feat: change logger to traceLogger for getting traceId when recovering (#374)
- 修复使用 postgres 数据库时，位置参数重复，导致参数与值不对应的问题。 (#960)
- move opentelemetry into trace package, and refactoring (#996)
- Feature goctl error wrap (#995)
- disable codecov github checks (#993)
- implement k8s service discovery (#988)
- format code, and reorg imports (#991)
- Fix filepath (#990)
- remote handler blank line when .HasRequest is false (#986)
- use codecov action v1 (#985)
- format coding style (#983)
- httpx.Error response without body (#982)
- format code (#979)
- configurable for load and stat statistics logs (#980)
- add go-zero users (#978)
- refactor (#977)
- update go version to 1.14 for github workflow (#976)
- Update Codecov `action` (#974)
- format coding style (#970)
- Fix context error in grpc (#962)
- format coding style (#969)
- Fix issues (#965)
- Add a test case for database code generation `tool` (#961)
- 修复stream拦截器tracer名问题 (#944)
- rest otel support (#943)
- add the opentelemetry tracing (#908)
- make sure setting code happen before callback in rest (#936)
- Fix issues (#931)
- update slack invite url (#933)
- add go-zero users (#928)
- Optimize model naming (#910)
- add unit test (#921)
- refactor (#920)
- Add traceId to the response headers (#919)
- add stringx.FirstN with ellipsis (#916)
- redis.go，type StringCmd = red.StringCmd (#790)
- add stringx.FirstN (#914)
- export pathvar for user-defined routers (#911)
- add Errorv/Infov/Slowv (#909)
- optimize grpc generation env check (#900)
- add workflow for closing inactive issues (#906)
- Update readme for better description. (#904)
- format coding style (#905)
- 带下划线的项目,配置文件名字错误。 (#733)
- refactor goctl (#902)
- refactor (#878)
- remove unnecessary chars. (#898)
- simplify type definition in readme (#896)
- refactor rest code (#895)
- add logx.DisableStat() to disable stat logs (#893)
- better text rendering (#892)
- format code (#888)
- format code (#884)
- add go-zero users (#883)
- update goctl version to 1.1.10 (#874)
- add goctl rpc template home flag (#871)
- add correct example for pg's url (#857)
- optimize typo (#855)
- upgrade grpc package (#845)
- simplify timeoutinterceptor (#840)
- Fixed http listener error. (#843)
- Feature model postgresql (#842)
- remove faq for old versions (#828)
- add go-zero users, update faq (#827)
- Fix the error stream method name (#826)
- go format with extra rules (#821)
- Fix issues: #725, #740 (#813)
- optimized (#819)
- Fix rpc generator bug (#799)
- Add --go_opt flag to adapt to the version after 1.4.0 of protoc-gen-go (#767)
- [WIP]Add parse headers info (#805)
- refactor mapping (#782)
- Add Sinter,Sinterstore & Modify TestRedis_Set (#779)
- update readme images (#776)
- update image rendering in readme (#775)
- disable load & stat logs for goctl (#773)
- upgrade grpc & etcd dependencies (#771)
- Fix issue #747 (#765)
- remove useless annotation (#761)
- add roadmap (#764)
- add contributing guid (#762)
- add code of conduct (#760)
- refactor fx (#759)
- Add some stream features (#712)
- add go-zero users (#756)
- add go-zero users (#751)
- replace cache key with colon (#746)
- add go-zero users (#739)
- Fix a typo (#729)
- add go-zero users, update slack invite link (#728)
- add go-zero users (#726)
- add go-zero users. (#723)
- optimize nested conditional (#709)
- optimize nested conditional (#708)
- Add document & comment for spec (#703)
- print entire sql statements in logx if necessary (#704)
- chore(format): change by gofumpt tool (#697)
- update goctl version to 1.1.8 (#696)
- add go-zero users (#688)
- resolve #610 (#684)
- Optimize model nl (#686)
- update readme for documents links (#681)
- update readme (#673)
- replace antlr module (#672)
- modify the order of PrometheusHandler (#670)
- update wechat qrcode (#665)
- disable prometheus if not configured (#663)
- add go-zero users (#643)
- update readme (#640)
- update readme (#638)
- optimize code (#637)
- spelling mistakes (#634)
- chore: update code format. (#628)
- update go-zero users (#623)
- add syncx.Guard func (#620)
- update readme (#617)
- add code coverage (#615)
- add FAQs in readme (#612)
- update go-zero users (#611)
- update go-zero users (#609)
- add go-zero users registry notes (#608)
- add go-zero users (#607)
- simplify redis tls implementation (#606)
- redis增加tls支持 (#595)
- refactor - remove ShrinkDeadline, it's the same as context.WithTimeout (#599)
- Replace contextx.ShrinkDeadline with context.WithTimeout (#598)
- Simplify contextx.ShrinkDeadline (#596)
- update regression test comment (#590)
- update regression test comment (#589)
- remove rt mode log (#587)
- optimize code (#579)
- support postgresql (#583)
- avoid goroutine leak after timeout (#575)
- gofmt logs (#574)
- add timezone and timeformat (#572)
- zrpc timeout & unit tests (#573)
- make hijack more stable (#565)
- refactor, and add comments to describe graceful shutdown (#564)
- Feature mongo gen (#546)
- Hdel support for multiple key deletion (#542)
- add important notes in readme (#560)
- add http hijack methods (#555)
- update doc link (#552)
- rename (#543)
- 暴露redis EvalSha  以及ScriptLoad接口 (#538)
- patch 1.1.5 (#530)
- feature 1.1.5 (#411)
- golint core/discov (#525)
- 修正http转发头字段值错误 (#521)
- Code optimized (#523)
- Code optimized (#493)
- add redis bitmap command (#490)
- redis add bitcount (#483)
- prevent negative timeout settings (#482)
- zrpc client support block (#412)
- Code optimized (#474)
- add more tests for service (#463)
- add more tests for rest (#462)
- Update serviceconf.go (#460)
- move examples into zero-examples (#453)
- remove images, use zero-doc instead (#452)
- add api doc (#449)
- add discov tests (#448)
- remove etcd facade, added for testing purpose (#447)
- add more tests for stores (#446)
- add more tests for stores (#445)
- add more tests for mongoc (#443)
- add more tests for sqlx (#442)
- add more tests for zrpc (#441)
- add more tests for sqlx (#440)
- add more tests for proc (#439)
- ring struct add lock (#434)
- Update readme.md
- update readme for broken links (#432)
- Support redis command Rpop (#431)
- support hscan in redis (#428)
- use english readme as default, because of github ranking (#427)
- Modify the http content-length max range : 30MB --> 32MB (#424)
- modify the maximum content-length to 30MB (#413)
- optimize code (#417)
- support zunionstore in redis (#410)
- use env if necessary in loading config (#409)
- update goctl version to 1.1.3 (#402)
- optimized (#392)
- add more tests for codec (#391)
- update readme (#390)
- Update periodicalexecutor.go (#389)
- format code (#386)
- simplify cgroup controller separation (#384)
- make sure unlock safe even if listeners panic (#383)
- code optimized (#382)
- Java optimized (#376)
- feature: refactor api parse to g4 (#365)
- add more tests for conf (#371)
- update doc to use table to render plugins (#370)
- remove duplicated code in goctl (#369)
- update goctl version to 1.1.3 (#364)
- make sure offset less than size even it's checked inside (#354)
- add godoc for RollingWindow (#351)
- simple rolling windows code (#346)
- optimized goctl format (#336)
- close issue of #337 (#347)
- align bucket boundary to interval in rolling window (#345)
- simplify rolling window code, and make tests run faster (#343)
- set guarded to false only on quitting background flush (#342)
- simplify periodical executor background routine (#339)
- add discord chat group in readme
- modify the goctl gensvc template (#323)
- Java (#327)
- simplify http.Flusher implementation (#326)
- simplify code with http.Flusher type conversion (#325)
- The ResponseWriters defined in rest.handler add Flush interface. (#318)
- add more tests for prof (#322)
- add more tests for zrpc (#321)
- add more tests (#320)
- add more tests (#319)
- add go report card back (#313)
- Update codeql-analysis.yml
- format code (#312)
- add config load support env var (#309)
- Update readme.md
- add wechat micro practice qrcode image (#289)
- Update readme.md
- Update readme-en.md
- Update readme.md
- Create codeql-analysis.yml
- optimize dockerfile generation (#284)
- refactor (#283)
- format dockerfile on non-chinese mode (#282)
- Update readme-en.md
- add EXPOSE in dockerfile generation (#281)
- optimize test case of TestRpcGenerate (#279)
- add category docker & kube (#276)
- optimize dockerfile (#272)
- fmt code (#270)
- require go 1.14 (#263)
- feature plugin custom flag (#251)
- optimized parse tag (#256)
- refactor & format code (#255)
- Feature bookstore update (#253)
- improve data type conversion (#236)
- goctl add plugin support (#243)
- support k8s deployment yaml generation (#247)
- optimize docker file generation, make docker build faster (#244)
- optimization (#241)
- some optimize by kevwan and benying (#240)
- simplify code (#234)
- simplify code, format makefile (#233)
- optimization (#221)
- modify the service name from proto (#230)
- Improve Makefile robustness (#224)
- set default handler value (#228)
- update version (#226)
- check go.mod before build docker image (#225)
- feature model interface (#222)
- feature: file namestyle (#223)
- format import
- 1.use local variable i; 2.make sure limiter larger than timer period (#218)
- optimize api new (#216)
- patch model&rpc (#207)
- update readme
- modify image url (#213)
- type should not define nested (#212)
- add error handle tests
- support error customization
- support type def without struct token (#210)
- add redis geospatial (#209)
- optimize parser (#206)
- update goctl readme
- update example
- update example
- refactor parser and remove deprecated code (#204)
- 1. group support multi level folder 2. remove force flag (#203)
- api support for comment double slash // (#201)
- format code
- change grpc interceptor to chain interceptor (#200)
- update etcd yaml to avoid no such nost resolve problem
- no default metric (#199)
- add dockerfile into template
- format service and add test (#197)
- rename postgres
- default metric host (#196)
- rewrite (#194)
- add dockerfile generator
- simplify http server starter
- graceful shutdown refined
- update doc (#193)
- Close the process when shutdown is finished (#157)
- break generator when happen error (#192)
- update cli package
- add more test (#189)
- refine code style
- faster the tests
- reactor rpc (#179)
- add https listen and serve
- move redistest into redis package
- refine tests
- update doc
- route support no request and response (#178)
- add more tests
- add more tests
- Update sharedcalls.go (#174)
- add gitee url
- update doc
- update bookstore example for generation prototype
- update doc
- remove wechat image
- update wechat qrcode
- support https in rest
- update wechat qrcode
- update readme
- add images back because of gitee not showing
- add images back because of gitee not showing
- goctl add stdin flag (#170)
- update doc using raw images
- update deployment version (#165)
- refactor middleware generator (#159)
- gen api svc add middleware implement temp code (#151)
- add logo in readme
- api handler generate incompletely while has no request  (#158)
- update api template (#156)
- add vote link
- docs: format markdown and add go mod in demo (#155)
- ignore blank between bracket and service tag (#154)
- model support globbing patterns (#153)
- make tests faster
- update wechat info
- can only specify one origin in cors
- make tests faster
- gozero template (#147)
- support cors in rest server
- optimized generator formatted code (#148)
- export WithUnaryClientInterceptor
- let balancer to be customizable
- rename NewPatRouter to NewRouter
- use goctl template to generate all kinds of templates
- api add middleware support (#140)
- add more tests
- add redis Zrevrank (#137)
- update breaker doc
- to correct breaker interface annotation (#136)
- update doc
- update doc
- add logx.Alert
- add fx.Split
- add anonymous annotation (#134)
- update goctl rpc template log print url (#133)
- print more message when parse error (#131)
- delete goctl rpc main tpl no use import (#130)
- update doc
- support api templates
- assert len > 0
- fail fast when rolling window size is zero
- simplify code generation
- stop rpc server when main function exit (#120)
- faster the tests
- Gozero sqlgen patch (#119)
- update readme
- add qq qrcode
- update readme
- make tests faster
- Goctl rpc patch (#117)
- make tests race-free
- avoid bigint converted into float64 when unmarshaling
- add more tests
- add more tests
- parser ad test (#116)
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add fx.Count
- add more tests
- add more tests
- breaker: remover useless code (#114)
- update wechat qrcode
- update codecov settings
- add more tests
- add more tests
- add more tests
- add more tests
- remove markdown linter
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- GetBreaker need double-check (#112)
- add more tests
- perfect the bookstore and shorturl doc (#109)
- better lock practice in sharedcalls
- update shorturl doc
- export cache package, add client interceptor customization
- add zrpc client interceptor
- doc: update sharedcalls.md layout (#107)
- refactor and rename folder to group (#106)
- add api doc link
- use default mongo db (#103)
- unmarshal should be struct
- Revert "goreportcard not working, remove it temporarily"
- goreportcard not working, remove it temporarily
- support return () syntax (#101)
- rename prommetric to prometheous, add unit tests
- update wechat and etcd yaml
- update example
- export AddOptions, AddStreamInterceptors, AddUnaryInterceptors
- query from cache first when do cache.Take
- rename (#98)
- add more clear error when rpc service is not started
- add test (#95)
- goctl support import api file (#94)
- add tracing logs in server side and client side
- remove unnecessary tag
- use options instead of opts in error message
- 修改不能编辑代码注释 (#92)
- feature: goctl jwt  (#91)
- update doc (#90)
- remove no need (#87)
- add trace/span in http logs
- use package level defined contextKey as context key
- printing context key friendly
- use contextType as string type
- rename ngin to rest in goctl
- optimize AtomicError (#82)
- update doc
- update rpc example
- rename rpcx to zrpc
- update wechat qrcode
- add more tests
- remove markdown linter temporarily
- simplify mapreduce code
- rename file and function name (#74)
- add mapping readme (#75)
- print message when starting api server
- rename function
- optimize: api generating for idea plugin (#68)
- api support empty request or empty response (#72)
- optimize route parse (#70)
- Sharedcalls.md (#69)
- optimized api new with absolute path like: goctl api new $PWD/xxxx (#67)
- update doc, add metric link
- drain pipe if reducer not drained
- Metric (#65)
- Markdown lint (#58)
- update goctl makefile
- add model&rpc doc (#62)
- update doc (#64)
- update quick start (#63)
- add in-process cache doc
- add fast create api demo service (#59)
- quickly generating rpc demo service (#60)
- update wechat image
- add (#54)
- make chinese readme as default
- update readme to add mapreduce link
- add mr tool doc (#50)
- refactor (#49)
- refactor (#48)
- refactor gomod logic (#47)
- add unit test for mapreduce
- add language link
- add bookstore english tutorial
- add language link
- add shorturl english tutorial
- add english readme
- add goctl description
- update readme
- update example doc
- add wechat qrcode
- add bookstore example
- add bookstore example
- replace clickhouse driver to the official one
- refactor code
- refactor
- add shorturl example code
- support go 1.13
- update shorturl doc
- update shorturl doc
- trim space  (#31)
- remove no need empty line (#29)
- update readme
- update docs
- update shorturl doc
- make svcCtx as a member for better code generation
- remove files
- remove makefile generation
- update readme
- update goctl makefile
- make goctl work on linux
- update shorturl doc
- update shorturl doc
- reorg imports
- rpc generation support windows (#28)
- update shorturl doc
- update shorturl doc
- update handler generation (#27)
- refine rpc generator
- refine goctl rpc generator
- rpc service generation (#26)
- update shorturl doc
- update shorturl doc
- better image rendering
- add quick example
- sort imports on api generation
- return zero value instead of nil on generated logic
- disable cpu stat in wsl linux
- use yaml, and detect go.mod in current dir
- move test code into internal package
- update ci configuration
- add more tests
- add more tests
- add more tests
- add more tests
- use predefined endpoint separator
- add fatal to stderr
- add etcd deploy yaml
- add more tests
- add more tests
- add more tests
- add more tests
- add more tests
- accelerate tests
- make tests parallel
- add more tests
- update readme
- update doc, add architecture picture
- make test stable
- update readme
- add more tests
- add more tests
- add more tests
- reorg imports
- update keywords.md
- gocctl model v20200819 (#18)
- update doc
- update stringx doc
- update stringx doc
- update image alt in doc
- update image scale in doc
- add keywords utility example
- add release badge
- add go report badge
- support customized mask char on trie
- add benchmark
- goctl model  reactor (#15)
- goctl生成Kotlin代码优化 (#16)
- rename rest files
- support direct scheme on rpc resolver
- remove utils.Report
- rename-Api
- multi-http-method-support
- remove-logx
- Add goctl kotlin support
- refactor names
- add license badge
- update codecov settings
- add codecov report
- add codecov badge
- use default decay value from finagle
- format readme
- update doc
- move lang.Must into logx.Must to make sure output fatal message as json
- confirm addition after add called in periodical executor
- add more tests
- add queue package
- remove pdf
- update workflow
- update readme
- rename files
- update readme
- rename files
- auto generate go mod if need
- remove bodyless check
- parse body only if content length > 0
- return ErrBodylessRequest on get method etc.
- remove bodyless check
- parse body only if content length > 0
- return ErrBodylessRequest on get method etc.
- export httpx.GetRemoteAddr
- export router
- export token parser for refresh token service
- remove unused method
- export httpx.GetRemoteAddr
- export router
- export token parser for refresh token service
- update readme.md
- remove unused method
- update readme.md
- use strings.Contains instead of strings.Index
- refactor rpcx, export WithDialOption and WithTimeout
- move auth interceptor into serverinterceptors
- use fmt.Println instead of println
- remove unused method
- use strings.Contains instead of strings.Index
- refactor rpcx, export WithDialOption and WithTimeout
- move auth interceptor into serverinterceptors
- use fmt.Println instead of println
- remove unused method
- optimize reading http request body
- remove deprecated model generator
- remove deprecated
- check content length before reading
- disable logs in unit tests
- correct parent packet for gomod
- remove star trends, not rendering for reasons
- rename test method
- refactor compare versions
- 【rich function】VersionCompare replace
- 【rich function】VersionCompare replace
- 【rich function】benchmark once function
- 【rich function】CustomVersionCompare append
- add star trends
- update goctl go generator
- Update readme.md
- update ci
- update ci, remove build
- Create go.yml
- add community wechat invitation
- reorg imports
- remote auto import
- update package reference
- change module in go.mod
- merge master
- add runner
- add more test for subset
- update goctl doc
- update goctl doc
- update readme
- update config
- remove files
- update docs
- update docs
- add tests
- avoid multiply zero on calculation load
- remove files
- add tests
- add subset algorithm
- add p2c peak ewma load balancer
- backup
- try grpc lb interface
- add resolver
- remove rq
- add roundrobin
- update packages
- remove packages
- add runner
- refactor
- format
- format
- refactor ngin to rest
- move router to httpx
- rename ngin to rest
- rename ngin to rest
- add more tests
- add more tests
- add more tests
- adjust doc
- adjust doc
- add doc
- add line break
- add article references
- add git url
- Add LICENSE
- Delete LICENSE
- Add LICENSE
- refactor
- rename mapreduce to mr
- refactor
- refactor
- avoid race condition
- refactor
- refactor
- refactor
- goctl added
- update packages
- add doc
- add doc
- add doc
- add goctl doc
- initial import

### Bug Fixes

- fix readme-cn (#1388)
- fix #1070 (#1389)
- fix #1376 (#1380)
- fix:  command system info missing go version (#1377)
- fix redis try-lock bug (#1366)
- go-zero tools ,fix a func,api new can not choose style (#1356)
- fix: #1318 (#1321)
- fix #1309 (#1317)
- fix #1305 (#1307)
- fix: go issue 16206 (#1298)
- fix #1288 (#1292)
- fixes #987 (#1283)
- feat: add etcd resolver scheme, fix discov minor issue (#1281)
- fixes #1257 (#1271)
- fixes #1169 (#1229)
- fixes #1222 (#1223)
- feat: ignore rest.WithPrefix on empty prefix (#1208)
- Generate route with prefix (#1200)
- feat: add rest.WithPrefix to support route prefix (#1194)
- fix the package name of grpc client (#1170)
- fix(goctl): repeat creation protoc-gen-goctl symlink (#1162)
- fix: opentelemetry traceid not correct (#1108)
- fix bug: generating dart code error (#1090)
- fix AtomicError panic when Set nil (#1049) (#1050)
- fix jwt security issue by using golang-jwt package (#1066)
- fix #1058 (#1064)
- chore: fix comment issues (#1056)
- fix test error on ubuntu (#1048)
- fix typo parse.go error message (#1041)
- fix symlink issue on windows for goctl (#1034)
- fix golint issues (#1027)
- fix proc.Done not found in windows (#1026)
- fix #1014 (#1018)
- fix golint issues, update codecov settings. (#1011)
- fix #1006 (#1008)
- fix golint issues (#992)
- fix #971 (#972)
- fix #957 (#959)
- fix #556 (#938)
- fix lint errors (#937)
- fix #820 (#934)
- fix missing `updateMethodTemplateFile` (#924)
- fix #915 (#917)
- fix #889 (#912)
- refactor goctl, fix golint issues (#903)
- fix: 解决golint 部分警告 (#897)
- fix golint issues (#899)
- fix http header binding failure bug #885 (#887)
- optimize mongo generation without cache (fix #881) (#882)
- fix #792 (#873)
- fix context missing (#872)
- fix bug that proc.SetTimeToForceQuit not working in windows (#869)
- fix issue #861 (#862)
- fix issue #831 (#850)
- fix issue #836 (#849)
- fix #796 (#844)
- Added database prefix of cache key. (#835)
- To generate grpc stream, fix issue #616 (#815)
- fix bug that empty query in transaction (#801)
- fix: Fix problems with non support for multidimensional arrays and basic type pointer arrays (#778)
- fix bug that etcd stream cancelled without re-watch (#770)
- fix broken link (#763)
- fix #736 (#738)
- fix golint issues, and optimize code (#705)
- fix #683 (#690)
- fix invalid link (#689)
- fix #676 (#682)
- fix zh_cn document url (#678)
- fix issue: https://github.com/zeromicro/goctl-swagger/issues/6 (#680)
- fix some typo (#677)
- fix antlr mod (#669)
- fix some typo (#667)
- fix comment function names (#649)
- doc: fix spell mistake (#627)
- fix (#592)
- fix a simple typo (#588)
- fix typo (#586)
- fix typo (#585)
- fix golint issues (#584)
- fix golint issues (#561)
- fix spelling (#551)
- fix golint issues (#540)
- fix collection breaker (#537)
- fix golint issues (#535)
- fix golint issues (#533)
- fix golint issues (#532)
- fix golint issues in zrpc (#531)
- fix golint issues in rest (#529)
- fix broken build (#528)
- fix golint issues in core/stores (#527)
- fix golint issues in core/syncx (#526)
- fix golint issues in core/threading (#524)
- fix golint issues in core/utils (#520)
- fix golint issues in core/timex (#517)
- fix golint issues in core/stringx (#516)
- fix golint issues in core/stat (#515)
- fix misspelling (#513)
- fix golint issues in core/service (#512)
- fix golint issues in core/search (#509)
- fix golint issues in core/rescue (#508)
- fix golint issues in core/queue (#507)
- fix golint issues in core/prometheus (#506)
- fix broken links in readme (#505)
- fix golint issues in core/prof (#503)
- fix golint issues in core/proc (#502)
- fix golint issues in core/netx (#501)
- fix golint issues in core/mr (#500)
- fix golint issues in core/metric (#499)
- fix golint issues in core/mathx (#498)
- fix golint issues in core/mapping (#497)
- fix golint issues in core/logx (#496)
- fix golint issues in core/load (#495)
- fix golint issues in core/limit (#494)
- fix golint issues in core/lang (#492)
- fix golint issues in core/jsonx (#491)
- fix golint issues in core/jsontype (#489)
- fix golint issues in core/iox (#488)
- fix golint issues in core/hash (#487)
- fix golint issues in core/fx (#486)
- fix golint issues in core/filex (#485)
- fix golint issues in core/executors (#484)
- fix golint issues in core/errorx (#480)
- fix golint issues in core/discov (#479)
- fix golint issues in core/contextx (#477)
- fix golint issues in core/conf (#476)
- fix golint issues in core/collection, refine cache interface (#475)
- fix golint issues in core/codec (#473)
- fix issue #469 (#471)
- fix gocyclo warnings (#468)
- fix golint issues in core/cmdline (#467)
- fix golint issues in core/breaker (#466)
- fix golint issues in core/bloom (#465)
- fix golint issues (#459)
- fix golint issues (#458)
- fix golint issues, else blocks (#457)
- fix golint issues, redis methods (#455)
- fix golint issues, package comments (#454)
- fix golint issues, exported doc (#451)
- fix golint issues (#450)
- fixes issue #425 (#438)
- fix readme.md error (#429)
- fix type convert error (#395)
- fix cgroup bug (#380)
- fix server.start return nil points (#379)
- f-fix spell (#381)
-  fix inner type generate error (#377)
- fix return in for (#367)
- Feature model fix (#362)
- fix potential data race in PeriodicalExecutor (#344)
- fix rolling window bug (#340)
- optimize code that fixes issue #317 (#338)
- fix bug #317 (#335)
- fix issue #317 (#331)
- fix broken link.
- fix broken doc link
- fixes #286 (#315)
- feature model fix (#296)
- fix gocyclo warnings (#278)
- fix dockerfile generation bug (#277)
- fix issue #266 (#275)
- fix tracelogger_test TestTraceLog (#271)
- fix lint errors (#249)
- fix doc errors
- fix issue #205
- fix issue #186
- rpc generation fix (#184)
- fix duplicate alias (#183)
- fix url 404 (#180)
- model template fix (#169)
- spell fix (#167)
- fix bug: generate incomplete model code  in case findOneByField (#160)
- fix zrpc client interceptor calling problem
- fix name typo and format with newline (#143)
- fix golint issues
- fix golint issues
- fix: fx/fn.Head func will forever block when n is less than 1 (#128)
- fix syncx/barrier test case (#123)
- fix: template cache key (#121)
- fix data race in tests
- fix data race in tests
- fix data race
- fix data race
- fix int64 primary key problem
- update: fix wrong word (#110)
- add unit test, fix interceptor bug
- fix typo of prometheus
- fix bug: module parse error (#97)
- fix bug: release empty struct limit (#96)
- fix redis error (#88)
- chore: fix typos (#85)
- fix: golint: context.WithValue should should not use basic type as key (#83)
- fix rpc client examle (#81)
- fix example tracing edge config (#76)
- fix goctl api (#71)
- fix goctl model (#61)
- fix GOMOD env fetch bug (#55)
- fix goctl model path (#53)
- fix command run path bug (#52)
- fix typoin doc
- fix mapreduce problem when reducer doesn't write
- fix readme typo
- fix typo (#38)
- fix LF (#37)
- fix bookstore example
- fix bug: miss time import (#36)
- fix shorturl example code (#35)
- fix: root path on windows bug (#34)
- fix doc errors
- fix dockerfile generation
- fix typo in doc
- fix doc error
- fix config yaml gen (#25)
- fix ci script
- fix format (#23)
- fix generate api demo (#19)
- fix render problem in doc
- fix golint warnings
- fix-log-fatal
- fix FileNotFoundException when response code is 4xx or 5xx
- fix-lang-must-not-found
- fix-break-line
- chore: fix typo
- fix data race
- fix windows slash
- fix windows bug
- fix windows slash
- fix windows bug
- fix windows slash
- fix windows slash
- fix goctl path issue on windows
- fix panic on auth
- fix image broken
- fix compile errors
