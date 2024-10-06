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
```mermaid
flowchart TD
    A[开始] --> B[创建用户实例]
    B --> C[启动监听消息goroutine]
    C --> D{用户上线}
    D --> E[加锁]
    E --> F[添加用户到在线列表]
    F --> G[解锁]
    G --> H[广播用户上线消息]
    H --> I[用户下线]
    I --> J[加锁]
    J --> K[从在线列表删除用户]
    K --> L[解锁]
    L --> M[广播用户下线消息]
    M --> N[发送消息给用户]
    N --> O[接收消息]
    O --> P{消息类型}
    P -->|who| Q[查询在线用户]
    Q --> R[发送在线用户列表]
    P -->|rename| S[处理改名]
    S --> T{检查新用户名是否已存在}
    T -->|存在| U[发送用户名被使用消息]
    T -->|不存在| V[更新用户名]
    V --> W[发送用户名更新成功消息]
    P -->|to| X[处理私聊消息]
    X --> Y{检查目标用户是否存在}
    Y -->|不存在| Z[发送目标用户不存在消息]
    Y -->|存在| AA[发送私聊内容给目标用户]
    AA --> AB[广播消息]
    AB --> O
    O --> AC[监听消息Channel]
    AC --> AD[接收Channel消息]
    AD --> AE[发送消息到客户端]
    AE --> AC
```
client.go: mock user-side client
```mermaid
flowchart TD
    A[开始] --> B[解析命令行参数]
    B --> C[创建客户端实例]
    C --> D{连接服务器}
    D -->|成功| E[启动处理服务器回应的goroutine]
    E --> F[进入主循环]
    F --> G[显示菜单]
    G --> H{用户输入}
    H -->|合法| I[设置客户端模式]
    H -->|非法| J[提示错误并重新输入]
    I --> K{根据模式选择操作}
    K -->|1. 公聊模式| L[公聊]
    K -->|2. 私聊模式| M[私聊]
    K -->|3. 更新用户名| N[更新用户名]
    K -->|0. 退出| O[退出]
    L --> P[提示输入聊天内容]
    P --> Q{聊天内容}
    Q -->|exit| R[退出聊天]
    Q -->|非exit| S[发送消息到服务器]
    S --> P
    M --> T[选择私聊对象]
    T --> U[提示输入聊天内容]
    U --> V{聊天内容}
    V -->|exit| W[退出聊天]
    V -->|非exit| X[发送私聊消息到服务器]
    X --> T
    N --> Y[提示输入新用户名]
    Y --> Z{发送新用户名到服务器}
    Z -->|成功| AA[更新用户名]
    Z -->|失败| AB[提示错误]
    AB --> Y
    AA --> F
    R --> F
    W --> F
    O --> AC[结束]
```
main.go: server starter 
