<template>
  <view class="login">
    <!-- 顶部Logo -->
    <view class="login-header">
      <image src="/static/images/logo.png" class="logo"></image>
      <view class="app-name">二手电脑回收</view>
      <view class="app-slogan">让闲置设备重新焕发价值</view>
    </view>
    
    <!-- 登录表单 -->
    <view class="login-form card">
      <view class="form-title">欢迎登录</view>
      
      <view class="form-group">
        <view class="input-wrapper">
          <text class="input-icon">👤</text>
          <input 
            v-model="formData.username" 
            class="form-input" 
            placeholder="请输入用户名"
            :value="formData.username"
          />
        </view>
      </view>
      
      <view class="form-group">
        <view class="input-wrapper">
          <text class="input-icon">🔒</text>
          <input 
            v-model="formData.password" 
            class="form-input" 
            placeholder="请输入密码"
            :password="!showPassword"
            :value="formData.password"
          />
          <text class="password-toggle" @click="togglePassword">
            {{ showPassword ? '👁️' : '👁️‍🗨️' }}
          </text>
        </view>
      </view>
      
      <view class="form-actions">
        <button class="login-btn" @click="handleLogin" :disabled="isLoading">
          <text v-if="!isLoading">登录</text>
          <text v-else>登录中...</text>
        </button>
        
        <view class="action-links">
          <text class="link-text" @click="goToRegister">还没有账号？立即注册</text>
        </view>
      </view>
    </view>
    
    <!-- 快速登录 -->
    <view class="quick-login">
      <view class="divider">
        <text class="divider-text">或者</text>
      </view>
      
      <view class="quick-buttons">
        <!-- 微信登录（仅小程序） -->
        <button 
          class="quick-btn wechat" 
          v-if="isWeChat"
          @click="wechatLogin"
          open-type="getUserInfo"
          @getuserinfo="onWechatLogin"
        >
          <text class="quick-icon">微</text>
          <text class="quick-text">微信登录</text>
        </button>
        
        <!-- 游客登录 -->
        <button class="quick-btn guest" @click="guestLogin">
          <text class="quick-icon">👤</text>
          <text class="quick-text">游客体验</text>
        </button>
      </view>
    </view>
    
    <!-- 底部信息 -->
    <view class="login-footer">
      <view class="agreement">
        <text class="agreement-text">
          登录即表示您同意
          <text class="agreement-link" @click="showAgreement('privacy')">《隐私政策》</text>
          和
          <text class="agreement-link" @click="showAgreement('service')">《服务协议》</text>
        </text>
      </view>
    </view>
  </view>
</template>

<script>
import { useUserStore } from '@/store'

export default {
  name: 'Login',
  data() {
    return {
      formData: {
        username: '',
        password: ''
      },
      showPassword: false,
      isLoading: false,
      isWeChat: false
    }
  },
  
  onLoad() {
    this.checkPlatform()
  },
  
  methods: {
    // 检查平台
    checkPlatform() {
      // #ifdef MP-WEIXIN
      this.isWeChat = true
      // #endif
    },
    
    // 切换密码显示
    togglePassword() {
      this.showPassword = !this.showPassword
    },
    
    // 表单验证
    validateForm() {
      if (!this.formData.username.trim()) {
        uni.showToast({
          title: '请输入用户名',
          icon: 'none'
        })
        return false
      }
      
      if (!this.formData.password.trim()) {
        uni.showToast({
          title: '请输入密码',
          icon: 'none'
        })
        return false
      }
      
      if (this.formData.password.length < 6) {
        uni.showToast({
          title: '密码至少6位',
          icon: 'none'
        })
        return false
      }
      
      return true
    },
    
    // 处理登录
    async handleLogin() {
      if (!this.validateForm()) {
        return
      }
      
      if (this.isLoading) {
        return
      }
      
      this.isLoading = true
      
      try {
        const res = await this.$http.post('/api/v1/auth/login', {
          username: this.formData.username,
          password: this.formData.password
        })
        
        // 保存登录信息
        const userStore = useUserStore()
        userStore.login(res.token, res.user)
        
        uni.showToast({
          title: '登录成功',
          icon: 'success'
        })
        
        // 跳转回上一页或首页
        setTimeout(() => {
          const pages = getCurrentPages()
          if (pages.length > 1) {
            uni.navigateBack()
          } else {
            uni.reLaunch({
              url: '/pages/index/index'
            })
          }
        }, 1500)
        
      } catch (error) {
        console.error('登录失败:', error)
        uni.showToast({
          title: error.error || '登录失败',
          icon: 'none'
        })
      } finally {
        this.isLoading = false
      }
    },
    
    // 微信登录
    wechatLogin() {
      // 微信小程序登录逻辑
      uni.showToast({
        title: '微信登录功能开发中',
        icon: 'none'
      })
    },
    
    // 微信登录回调
    onWechatLogin(e) {
      console.log('微信用户信息:', e.detail)
      // 处理微信登录
    },
    
    // 游客登录
    async guestLogin() {
      try {
        // 使用预设的游客账号
        const res = await this.$http.post('/api/v1/auth/login', {
          username: 'guest',
          password: 'guest123'
        })
        
        const userStore = useUserStore()
        userStore.login(res.token, res.user)
        
        uni.showToast({
          title: '游客登录成功',
          icon: 'success'
        })
        
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/index/index'
          })
        }, 1500)
        
      } catch (error) {
        // 如果游客账号不存在，自动创建
        try {
          await this.$http.post('/api/v1/auth/register', {
            username: 'guest',
            password: 'guest123',
            phone: '13800138000',
            real_name: '游客用户'
          })
          
          // 注册成功后自动登录
          this.guestLogin()
          
        } catch (regError) {
          uni.showToast({
            title: '游客登录失败',
            icon: 'none'
          })
        }
      }
    },
    
    // 跳转到注册页
    goToRegister() {
      uni.navigateTo({
        url: '/pages/user/register'
      })
    },
    
    // 显示协议
    showAgreement(type) {
      const title = type === 'privacy' ? '隐私政策' : '服务协议'
      const content = type === 'privacy' ? 
        '我们非常重视您的隐私保护，您的个人信息将被严格保密...' :
        '欢迎使用我们的回收服务，请仔细阅读以下服务条款...'
      
      uni.showModal({
        title,
        content,
        showCancel: false
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.login {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
}

.login-header {
  text-align: center;
  margin-bottom: 80rpx;
  
  .logo {
    width: 120rpx;
    height: 120rpx;
    border-radius: 24rpx;
    margin-bottom: 30rpx;
  }
  
  .app-name {
    font-size: 48rpx;
    font-weight: bold;
    color: #fff;
    margin-bottom: 16rpx;
  }
  
  .app-slogan {
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.8);
  }
}

.login-form {
  background: #fff;
  border-radius: 24rpx;
  padding: 60rpx 40rpx;
  margin-bottom: 40rpx;
  
  .form-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #333;
    text-align: center;
    margin-bottom: 60rpx;
  }
  
  .form-group {
    margin-bottom: 40rpx;
    
    .input-wrapper {
      position: relative;
      display: flex;
      align-items: center;
      border: 2rpx solid #e0e0e0;
      border-radius: 50rpx;
      background: #f8f9fa;
      
      &:focus-within {
        border-color: #667eea;
        background: #fff;
      }
      
      .input-icon {
        width: 80rpx;
        text-align: center;
        font-size: 32rpx;
        color: #999;
      }
      
      .form-input {
        flex: 1;
        height: 88rpx;
        font-size: 32rpx;
        color: #333;
      }
      
      .password-toggle {
        width: 80rpx;
        text-align: center;
        font-size: 32rpx;
        color: #999;
      }
    }
  }
  
  .form-actions {
    .login-btn {
      width: 100%;
      height: 88rpx;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: #fff;
      border-radius: 44rpx;
      font-size: 32rpx;
      font-weight: bold;
      border: none;
      margin-bottom: 40rpx;
      
      &:disabled {
        opacity: 0.6;
      }
    }
    
    .action-links {
      text-align: center;
      
      .link-text {
        font-size: 28rpx;
        color: #667eea;
      }
    }
  }
}

.quick-login {
  .divider {
    position: relative;
    text-align: center;
    margin-bottom: 40rpx;
    
    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 0;
      right: 0;
      height: 1rpx;
      background: rgba(255, 255, 255, 0.3);
    }
    
    .divider-text {
      display: inline-block;
      padding: 0 20rpx;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: rgba(255, 255, 255, 0.8);
      font-size: 24rpx;
    }
  }
  
  .quick-buttons {
    display: flex;
    gap: 30rpx;
    
    .quick-btn {
      flex: 1;
      height: 88rpx;
      border-radius: 44rpx;
      border: 2rpx solid rgba(255, 255, 255, 0.3);
      background: rgba(255, 255, 255, 0.1);
      color: #fff;
      font-size: 28rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      
      .quick-icon {
        margin-right: 12rpx;
        font-size: 32rpx;
      }
      
      &.wechat {
        background: rgba(9, 187, 7, 0.2);
        border-color: rgba(9, 187, 7, 0.3);
      }
      
      &.guest {
        background: rgba(255, 255, 255, 0.2);
      }
    }
  }
}

.login-footer {
  margin-top: auto;
  
  .agreement {
    text-align: center;
    
    .agreement-text {
      font-size: 24rpx;
      color: rgba(255, 255, 255, 0.7);
      line-height: 1.5;
      
      .agreement-link {
        color: #fff;
        text-decoration: underline;
      }
    }
  }
}
</style>
