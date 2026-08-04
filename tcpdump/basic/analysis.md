# 抓包分析：CollectScheduleRequestStat gRPC 接口

## 基本信息

- 抓包时间：2026-08-04 11:11:27
- 客户端（本机）：30.26.152.44:54901
- 服务端（gRPC 服务）：9.134.117.127:8086
- 传输协议：TCP + HTTP/2（gRPC 跑在 HTTP/2 之上，不是 HTTP/1.1）
- 总记录数：26 条
- 单次调用连接存活时长：约 38ms（.185462 → .223408）
- 说明：本文件只导出了 TCP 层「摘要行」（-nn 文本），未带 `-A` 的 ASCII 明文，
  因此下面的协议帧是通过 **报文长度 + 序列号(seq/ack) 推断** 的，属于典型 gRPC 帧结构。

## 报文标记含义（tcpdump 新手速查）

| 标记 | 含义 |
| --- | --- |
| `S` | SYN，发起连接（三次握手第 1 步） |
| `S.` | SYN+ACK，服务端回应（三次握手第 2 步） |
| `.` | 纯 ACK，确认收到数据 |
| `P.` | PSH+ACK，携带业务数据并尽快交给应用层 |
| `F.` | FIN+ACK，发起关闭连接 |

---

## 一、TCP 三次握手（记录 1-3）

> 作用：建立可靠 TCP 连接，协商初始序列号与窗口。

- **记录 1** `11:11:27.185462` 客户端 → 服务端 `Flags [S]`：客户端发 SYN，seq=2092668383，
  通告 mss 1460、wscale 6: window scale 6（即实际接收窗口可放大 2^6 倍）。
  - mss 指 TCP 单个报文段，承载的【应用层数据最大字节数】,Maximum Segment Size 最大分段大小
- **记录 2** `11:11:27.193759` 服务端 → 客户端 `Flags [S.]`：服务端回 SYN-ACK，ack=2092668384（=客户端 seq+1），
  通告 mss 1382、window scale 7。RTT ≈ 8ms，说明网络链路很好。
- **记录 3** `11:11:27.193906` 客户端 → 服务端 `Flags [.]`：客户端 ACK，三次握手完成。

> 提示：服务端 mss=1382 比客户端 1460 小，后续每个 TCP 段最大只能装 1382 字节业务数据（本例请求 175/175 字节远小于它，所以没分片）。

---

## 二、HTTP/2 连接建立（记录 4-10）

> gRPC 基于 HTTP/2，连接建立后要先交换连接前导(prefab)和 SETTINGS 帧。

- **记录 4** `11:11:27.194151` 客户端 → 服务端 `length 33`：发送 **HTTP/2 连接前导**。
  33 = 24 字节魔数(`PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n`) + 9 字节 SETTINGS 帧头，这是 gRPC 客户端的标配开场。
- **记录 5** `11:11:27.202783` 服务端 → 客户端 `length 15`：服务端 SETTINGS 帧（9 帧头 + 6 字节参数，如最大并发流）。
- **记录 6-7** 服务端的 `[.]` 与 `length 9`：ACK 确认客户端的 33 字节，并回一个纯 SETTINGS ACK（9 字节帧头、无载荷）。
- **记录 8-9** 客户端两个 `[.]`：确认收到服务端的前 16/25 字节 SETTINGS。
- **记录 10** `11:11:27.203209` 客户端 → 服务端 `length 9`：客户端对服务端 SETTINGS 的 ACK。

> 至此 HTTP/2 连接就绪，下面开始真正的 gRPC 调用。

---

## 三、gRPC 请求发送（记录 11-12）—— 重点是这两条

- **记录 11** `11:11:27.203272` 客户端 → 服务端 `length 111`：**HTTP/2 HEADERS 帧**，携带 gRPC 请求头，典型包含：
  `:method=POST`、`:path=/<service>/CollectScheduleRequestStat`、`:scheme=http`、
  `:authority=9.134.117.127:8086`、`content-type: application/grpc`、`te: trailers`。
  （这条就是你这次抓包最关心的「接口名」所在帧，:path 里能看到 CollectScheduleRequestStat。）
- **记录 12** `11:11:27.203805` 客户端 → 服务端 `length 175`：**HTTP/2 DATA 帧**，即序列化的 protobuf 请求体
  （对应代码 `schedule_client/main.go` 里 `CollectOption` 组装的 namespace=Test、service=lzb_test、lb=GLOBAL_P2C 等上报内容）。

> 客户端 HEADERS+DATA 一次性发完，seq 到 329 结束（43→154→329）。

---

## 四、服务端确认收到（记录 13-15）

- **记录 13-15** 服务端连续三个 `[.]`：分别 ACK 到 seq 43、154、329，确认已完整收下 175 字节请求体。

---

## 五、gRPC 响应返回（记录 16、19）—— 重点看这两条

- **记录 16** `11:11:27.212022` 服务端 → 客户端 `length 30`：**HTTP/2 HEADERS 帧**（响应头），
  含 `:status=200`、`content-type: application/grpc`，表示请求已被接受。
- **记录 19** `11:11:27.212645` 服务端 → 客户端 `length 139`：**HTTP/2 DATA 帧**，序列化的 protobuf 响应体
  （即代码里 `handle` 回调收到的 `StatResponse`，里面有 `code`、`info` 字段）。
  从「服务端收全请求(seq=329, .212021)」到「发出响应 DATA(.212645)」仅约 0.6ms，服务端处理很快。

> 记录 17-18、20-21 是双方为这些帧互发的 ACK 与少量 PING/流控制帧，属正常 HTTP/2 保活与流控，不影响主流程。

---

## 六、连接关闭（记录 22、25、26）

- **记录 22** `11:11:27.214208` 客户端 → 服务端 `Flags [F.]`：客户端发 FIN，单方面结束连接（gRPC 调用完即关）。
- **记录 23-24** 服务端两个 `[.]`：ACK 客户端的 FIN 与之前的数据。
- **记录 25** `11:11:27.223408` 服务端 → 客户端 `length 17`：**HTTP/2 GOAWAY 帧**（9 帧头 + 4 字节 last-stream-id + 4 字节错误码 = 17），
  优雅告知「不再接受新流」，正式关闭 HTTP/2 连接。
- **记录 26** `11:11:27.223408` 服务端 → 客户端 `Flags [F.]`：服务端回 FIN，四次挥手完成，连接关闭。

---

## 关键结论

1. **接口确实调通了**：三次握手 → HTTP/2 建连 → 请求 HEADERS/DATA → 响应 200 + DATA → 正常关闭，全链路无 RST、无重传、无超时。
2. **接口名位置**：`CollectScheduleRequestStat` 在记录 11（请求 HEADERS 帧的 `:path`）中，pcap 里可搜 `:path` 定位。
3. **为什么看不到明文 protobuf**：当前服务端是明文 HTTP/2（非 TLS），所以 protobuf 体本身是明文，只是本文件只导出了 TCP 摘要行；
   想直接看请求/响应体，用 `tshark -r schedule_capture.pcap -Y http2` 或 `tcpdump -A` 重新看即可。
4. **性能**：建连(RTT 8ms) + 调用(响应 0.6ms) 总耗时约 38ms，调用本身极快，耗时主要在 TCP/HTTP2 握手。

## 时序图（plantuml）

@startuml
title CollectScheduleRequestStat 抓包时序（记录 1-26 简化）
participant C as 客户端 30.26.152.44:54901
participant S as 服务端 9.134.117.127:8086

== TCP 三次握手 (记录1-3) ==
C -> S : [S] SYN (记录1)
S -> C : [S.] SYN-ACK (记录2)
C -> S : [.] ACK (记录3)

== HTTP/2 建连 (记录4-10) ==
C -> S : 连接前导+SETTINGS length33 (记录4)
S -> C : SETTINGS length15 + ACK (记录5-7)
C -> S : SETTINGS ACK length9 (记录10)

== gRPC 请求 (记录11-12) ==
C -> S : HEADERS :path=.../CollectScheduleRequestStat len111 (记录11)
C -> S : DATA protobuf 请求体 len175 (记录12)
S -> C : ACK seq329 (记录13-15)

== gRPC 响应 (记录16,19) ==
S -> C : HEADERS :status=200 len30 (记录16)
S -> C : DATA protobuf StatResponse len139 (记录19)
C -> S : ACK (记录17-18,20-21)

== 连接关闭 (记录22,25,26) ==
C -> S : [F.] FIN (记录22)
S -> C : GOAWAY len17 + ACK (记录23-25)
S -> C : [F.] FIN (记录26)
@enduml