# 登录失败添加验证码功能 - 实施计划

> 生成时间：2026-03-09
> 方案：A+A（后端 Redis + 自建验证码 / 前端并排 Flex 布局）

## 需求规格

| 项目 | 规格 |
|------|------|
| 触发阈值 | 连续失败 3 次后强制验证码 |
| 计数维度 | IP + 邮箱双重维度 |
| 验证码类型 | 图形验证码（6位字母数字） |
| 有效期 | 失败计数 15min / 验证码 5min |
| 成功行为 | 清除失败计数，正常登录 |

---

## 一、后端实施计划

### 1.1 新增文件

```
backend/internal/modules/auth/security/
├── guard.go          # Redis 失败计数器（IP/Email 维度）
├── captcha.go        # 验证码服务（生成、存储、校验）
└── image.go          # 图片渲染（6位字符 → PNG Base64）
```

### 1.2 修改文件

| 文件 | 变更内容 |
|------|----------|
| `backend/internal/config/loader.go` | 新增 `SecurityConfig`（阈值、TTL） |
| `backend/cmd/api/main.go` | 初始化 Redis 客户端，注入 auth 模块 |
| `backend/internal/modules/auth/core/entity.go` | `LoginRequest` 新增 `CaptchaID`, `CaptchaCode`, `ClientIP` |
| `backend/internal/modules/auth/core/service.go` | 注入 `LoginGuard` + `CaptchaService`，调整登录流程 |
| `backend/internal/modules/auth/transport/http/handler.go` | 新增 `/auth/captcha` 端点，处理错误码 |
| `backend/internal/modules/shared/errors/errors.go` | 新增 `ErrCaptchaRequired`, `ErrCaptchaInvalid` |
| `backend/internal/pkg/response/response.go` | 新增 `ErrorCode` 字段 |

### 1.3 API 契约

#### POST /api/v1/auth/captcha

**Request**: 无需参数（或 `{ "purpose": "login" }`）

**Response (201)**:
```json
{
  "code": 201,
  "data": {
    "captcha_id": "b565...",
    "captcha_image": "data:image/png;base64,iVBORw0...",
    "expires_in": 300
  }
}
```

#### POST /api/v1/auth/login（扩展）

**Request**:
```json
{
  "email": "alice@example.com",
  "password": "secret",
  "captcha_id": "b565...",    // 触发验证码后必填
  "captcha_code": "Q7X4BD"   // 触发验证码后必填
}
```

**Response - 需要验证码 (428)**:
```json
{
  "code": 428,
  "error_code": "AUTH_CAPTCHA_REQUIRED",
  "message": "captcha required",
  "data": {
    "failures": { "ip": 3, "email": 2 }
  }
}
```

**Response - 验证码错误 (400)**:
```json
{
  "code": 400,
  "error_code": "AUTH_CAPTCHA_INVALID",
  "message": "invalid or expired captcha"
}
```

### 1.4 Redis Key 设计

| Key Pattern | TTL | Value | Purpose |
|-------------|-----|-------|---------|
| `auth:login:fail:ip:{IP}` | 900s | int | IP 维度失败计数 |
| `auth:login:fail:email:{email}` | 900s | int | Email 维度失败计数 |
| `auth:captcha:{id}` | 300s | JSON `{hash, email, ip}` | 验证码答案（SHA-256 哈希） |

### 1.5 错误码定义

| 错误码 | HTTP Status | 说明 |
|--------|-------------|------|
| `AUTH_CAPTCHA_REQUIRED` | 428 | 需要验证码 |
| `AUTH_CAPTCHA_INVALID` | 400 | 验证码错误或过期 |
| `AUTH_INVALID_CREDENTIALS` | 401 | 账号密码错误 |

---

## 二、前端实施计划

### 2.1 修改文件

| 文件 | 变更内容 |
|------|----------|
| `frontend/src/types/index.ts` | 新增 `CaptchaResponse` 接口 |
| `frontend/src/api/auth.ts` | 扩展 `login` 参数，新增 `getCaptcha` 方法 |
| `frontend/src/views/admin/LoginView.vue` | 新增验证码状态、UI、交互逻辑 |

### 2.2 组件状态设计

```typescript
// 验证码相关状态
const MAX_FAILED_ATTEMPTS = 3
const failedAttempts = ref(0)      // 失败计数
const captchaId = ref('')          // 验证码 ID
const captchaImg = ref('')         // Base64 图片
const form.captcha = ''            // 用户输入
```

### 2.3 交互流程

```
用户登录 → 失败 → failedAttempts++
         ↓
  failedAttempts >= 3 ?
         ↓ 是
  自动获取验证码 → 显示验证码输入框
         ↓
  用户输入验证码 → 提交登录（携带 captcha_id, captcha_code）
         ↓
  成功 → 重置 failedAttempts，跳转
  失败 → 刷新验证码，清空输入
```

### 2.4 UI 布局（并排 Flex）

```html
<div class="flex items-center gap-3">
  <input class="input-base flex-1" maxlength="6" />
  <img :src="captchaImg" @click="loadCaptcha"
       class="h-10 w-28 cursor-pointer hover:opacity-80" />
</div>
```

---

## 三、分步实施（可并行）

### Phase 1: 基础设施（阻塞）

| Step | 任务 | 依赖 |
|------|------|------|
| 1.1 | 扩展 `SecurityConfig`，添加环境变量 | - |
| 1.2 | `main.go` 初始化 Redis，注入 auth 模块 | 1.1 |

### Phase 2: 后端核心（可并行）

| Step | 任务 | 依赖 |
|------|------|------|
| 2.1 | 实现 `security/guard.go`（Redis 计数器） | 1.2 |
| 2.2 | 实现 `security/captcha.go` + `image.go` | 1.2 |
| 2.3 | 扩展 `entity.go`（新增字段） | - |
| 2.4 | 新增错误码定义 | - |

### Phase 3: 集成与 API（依赖 Phase 2）

| Step | 任务 | 依赖 |
|------|------|------|
| 3.1 | 修改 `service.go`（注入 Guard + Captcha） | 2.1, 2.2 |
| 3.2 | 修改 `handler.go`（新增端点，错误码） | 3.1 |
| 3.3 | 扩展 `response.go`（ErrorCode 字段） | 2.4 |

### Phase 4: 前端（可与 Phase 2 并行）

| Step | 任务 | 依赖 |
|------|------|------|
| 4.1 | 新增 `CaptchaResponse` 类型 | - |
| 4.2 | 扩展 `authApi`（login + getCaptcha） | 4.1 |
| 4.3 | 修改 `LoginView.vue`（状态 + UI + 逻辑） | 4.2 |

### Phase 5: 测试与文档

| Step | 任务 | 依赖 |
|------|------|------|
| 5.1 | 后端单元测试 | 3.2 |
| 5.2 | 集成测试（手动验证 Redis） | 5.1 |
| 5.3 | 更新 API 文档 | 3.2 |

---

## 四、安全注意事项

1. **验证码答案哈希存储**：使用 SHA-256，不存明文
2. **IP 获取**：使用 `c.ClientIP()`，处理反向代理场景
3. **错误信息模糊**：避免泄露邮箱是否存在
4. **Redis 降级**：Redis 不可用时应记录日志但不阻断登录
5. **防 DoS**：验证码端点可加简单限流

---

## 五、验收标准

- [ ] 连续失败 3 次后，登录返回 `AUTH_CAPTCHA_REQUIRED`
- [ ] 前端自动显示验证码输入框并加载图片
- [ ] 验证码错误返回 `AUTH_CAPTCHA_INVALID`，前端自动刷新
- [ ] 登录成功后清除失败计数
- [ ] 15 分钟后失败计数自动过期
- [ ] 5 分钟后验证码自动过期
