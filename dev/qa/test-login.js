/**
 * Playwright 登录测试脚本
 */
const { chromium } = require('playwright');

(async () => {
  console.log('🚀 启动浏览器...');
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    // 访问登录页面
    console.log('📄 访问登录页面: http://localhost:7777');
    await page.goto('http://localhost:7777', { waitUntil: 'networkidle' });

    // 等待页面加载
    await page.waitForTimeout(2000);

    // 截图：登录页面
    await page.screenshot({ path: '/home/luo/code/remotegpu/dev/qa/screenshot-1-login-page.png' });
    console.log('📸 截图已保存: screenshot-1-login-page.png');

    // 查找用户名输入框
    console.log('🔍 查找登录表单元素...');
    const usernameInput = await page.locator('input[type="text"], input[placeholder*="用户名"], input[placeholder*="username"]').first();
    const passwordInput = await page.locator('input[type="password"]').first();

    if (await usernameInput.count() === 0) {
      console.error('❌ 未找到用户名输入框');
      await page.screenshot({ path: '/home/luo/code/remotegpu/dev/qa/screenshot-error.png' });
      return;
    }

    // 输入用户名和密码
    console.log('⌨️  输入用户名: testuser');
    await usernameInput.fill('testuser');

    console.log('⌨️  输入密码: Test123456');
    await passwordInput.fill('Test123456');

    // 截图：填写表单后
    await page.screenshot({ path: '/home/luo/code/remotegpu/dev/qa/screenshot-2-form-filled.png' });
    console.log('📸 截图已保存: screenshot-2-form-filled.png');

    // 查找并点击登录按钮
    console.log('🔍 查找登录按钮...');
    const loginButton = await page.locator('button:has-text("登录"), button:has-text("Login"), button[type="submit"]').first();

    if (await loginButton.count() === 0) {
      console.error('❌ 未找到登录按钮');
      return;
    }

    // 监听网络请求
    page.on('response', async (response) => {
      if (response.url().includes('/api/v1/user/login')) {
        console.log(`📡 登录请求响应: ${response.status()}`);
        const body = await response.text();
        console.log(`📦 响应内容: ${body}`);
      }
    });

    console.log('🖱️  点击登录按钮...');
    await loginButton.click();

    // 等待响应
    await page.waitForTimeout(3000);

    // 截图：登录后
    await page.screenshot({ path: '/home/luo/code/remotegpu/dev/qa/screenshot-3-after-login.png' });
    console.log('📸 截图已保存: screenshot-3-after-login.png');

    // 检查当前URL
    const currentUrl = page.url();
    console.log(`🌐 当前URL: ${currentUrl}`);

    // 检查是否有错误提示
    const errorMessage = await page.locator('.el-message--error, .error-message, [class*="error"]').first().textContent().catch(() => null);
    if (errorMessage) {
      console.log(`⚠️  错误提示: ${errorMessage}`);
    }

    console.log('✅ 测试完成');

  } catch (error) {
    console.error('❌ 测试失败:', error.message);
    await page.screenshot({ path: '/home/luo/code/remotegpu/dev/qa/screenshot-error.png' });
  } finally {
    await browser.close();
  }
})();
