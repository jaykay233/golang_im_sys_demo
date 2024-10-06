# golang_im_sys_demo

## course link
https://www.bilibili.com/video/BV1gf4y1r79E/?spm_id_from=333.337.search-card.all.click

## structure
server.go: process server-side business
```mermaid
flowchart TD
    A[开始] --> B[创建服务器实例]
    B --> C[监听TCP端口]
    C --> D{等待连接}
    D --> E[接受连接]
    E --> F[创建用户实例]
    F --> G[用户上线]
    G --> H[启动消息监听goroutine]
    H --> I[启动用户消息处理goroutine]
    I --> J{读取消息}
    J -->|成功| K[处理消息并广播]
    K --> J
    J -->|失败| L[用户下线]
    L --> M[关闭资源]
    M --> N[关闭连接]
    N --> D
    H --> O{监听消息广播}
    O -->|收到消息| P[发送消息给所有在线用户]
    P --> O
```
user.go: process user-side business  
client.go: mock user-side client  
main.go: server starter 
