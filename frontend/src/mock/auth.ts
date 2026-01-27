/**
 * 认证模块 Mock 数据
 */

// Mock 用户数据
export const mockUsers = [
  {
    id: 1,
    username: 'admin',
    password: '123456',
    email: 'admin@example.com',
    name: '管理员',
    role: 'admin',
  },
  {
    id: 2,
    username: 'user',
    password: '123456',
    email: 'user@example.com',
    name: '普通用户',
    role: 'user',
  },
]

// 登录接口 Mock
export const mockLogin = (data: any) => {
  const { username, password } = data
  const user = mockUsers.find(
    (u) => u.username === username && u.password === password
  )

  if (user) {
    const { password: _, ...userInfo } = user
    return {
      code: 200,
      message: '登录成功',
      data: {
        token: 'mock-token-' + Date.now(),
        user: userInfo,
      },
    }
  }

  return {
    code: 401,
    message: '用户名或密码错误',
    data: null,
  }
}

// 获取当前用户信息 Mock
export const mockGetCurrentUser = () => {
  return {
    code: 200,
    message: '成功',
    data: {
      id: 1,
      username: 'admin',
      email: 'admin@example.com',
      name: '管理员',
      role: 'admin',
    },
  }
}

// 登出接口 Mock
export const mockLogout = () => {
  return {
    code: 200,
    message: '登出成功',
    data: null,
  }
}
