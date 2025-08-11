<template>
  <view class="register">
    <!-- 顶部标题 -->
    <view class="register-header">
      <view class="header-title">创建账号</view>
      <view class="header-subtitle">加入我们，开启环保回收之旅</view>
    </view>
    
    <!-- 注册表单 -->
    <view class="register-form card">
      <view class="form-group">
        <text class="form-label">用户名 *</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.username" 
            class="form-input" 
            placeholder="请输入用户名（3-20位）"
            :value="formData.username"
            @blur="checkUsername"
          />
          <text class="input-status" :class="usernameStatus.class" v-if="usernameStatus.text">
            {{ usernameStatus.text }}
          </text>
        </view>
      </view>
      
      <view class="form-group">
        <text class="form-label">密码 *</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.password" 
            class="form-input" 
            placeholder="请输入密码（至少6位）"
            :password="!showPassword"
            :value="formData.password"
            @input="checkPasswordStrength"
          />
          <text class="password-toggle" @click="togglePassword">
            {{ showPassword ? '👁️' : '👁️‍🗨️' }}
          </text>
        </view>
        <view class="password-strength" v-if="passwordStrength.show">
          <text class="strength-text" :class="passwordStrength.class">
            密码强度：{{ passwordStrength.text }}
          </text>
          <view class="strength-bar">
            <view 
              class="strength-level" 
              :class="passwordStrength.class"
              :style="{ width: passwordStrength.width }"
            ></view>
          </view>
        </view>
      </view>
      
      <view class="form-group">
        <text class="form-label">确认密码 *</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.confirmPassword" 
            class="form-input" 
            placeholder="请再次输入密码"
            :password="!showConfirmPassword"
            :value="formData.confirmPassword"
            @blur="checkPasswordMatch"
          />
          <text class="password-toggle" @click="toggleConfirmPassword">
            {{ showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
          </text>
        </view>
        <text class="input-error" v-if="passwordMatchError">{{ passwordMatchError }}</text>
      </view>
      
      <view class="form-group">
        <text class="form-label">手机号 *</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.phone" 
            class="form-input" 
            placeholder="请输入手机号"
            type="number"
            :value="formData.phone"
            @blur="checkPhone"
          />
          <text class="input-status" :class="phoneStatus.class" v-if="phoneStatus.text">
            {{ phoneStatus.text }}
          </text>
        </view>
      </view>
      
      <view class="form-group">
        <text class="form-label">邮箱</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.email" 
            class="form-input" 
            placeholder="请输入邮箱（可选）"
            :value="formData.email"
            @blur="checkEmail"
          />
          <text class="input-status" :class="emailStatus.class" v-if="emailStatus.text">
            {{ emailStatus.text }}
          </text>
        </view>
      </view>
      
      <view class="form-group">
        <text class="form-label">真实姓名</text>
        <view class="input-wrapper">
          <input 
            v-model="formData.realName" 
            class="form-input" 
            placeholder="请输入真实姓名（可选）"
            :value="formData.realName"
          />
        </view>
      </view>
      
      <!-- 协议同意 -->
      <view class="agreement-section">
        <view class="agreement-check" @click="toggleAgreement">
          <text class="check-icon" :class="{ checked: agreeToTerms }">
            {{ agreeToTerms ? '✓' : '○' }}
          </text>
          <text class="agreement-text">
            我已阅读并同意
            <text class="agreement-link" @click.stop="showAgreement('service')">《服务协议》</text>
            和
            <text class="agreement-link" @click.stop="showAgreement('privacy')">《隐私政策》</text>
          </text>
        </view>
      </view>
      
      <!-- 注册按钮 -->
      <view class="form-actions">
        <button class="register-btn" @click="handleRegister" :disabled="!canRegister || isLoading">
          <text v-if="!isLoading">立即注册</text>
          <text v-else>注册中...</text>
        </button>
        
        <view class="action-links">
          <text class="link-text" @click="goToLogin">已有账号？立即登录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { useUserStore } from '@/store'

export default {
  name: 'Register',
  data() {
    return {
      formData: {
        username: '',
        password: '',
        confirmPassword: '',
        phone: '',
        email: '',
        realName: ''
      },
      showPassword: false,
      showConfirmPassword: false,
      agreeToTerms: false,
      isLoading: false,
      
      // 验证状态
      usernameStatus: { text: '', class: '' },
      phoneStatus: { text: '', class: '' },
      emailStatus: { text: '', class: '' },
      passwordMatchError: '',
      
      // 密码强度
      passwordStrength: {
        show: false,
        text: '',
        class: '',
        width: '0%'
      }
    }
  },
  
  computed: {
    canRegister() {
      return this.formData.username.length >= 3 &&
             this.formData.password.length >= 6 &&
             this.formData.confirmPassword === this.formData.password &&
             this.formData.phone.length === 11 &&
             this.agreeToTerms &&
             this.usernameStatus.class === 'success' &&
             this.phoneStatus.class === 'success'
    }
  },
  
  methods: {
    // 切换密码显示
    togglePassword() {
      this.showPassword = !this.showPassword
    },
    
    toggleConfirmPassword() {
      this.showConfirmPassword = !this.showConfirmPassword
    },
    
    // 切换协议同意
    toggleAgreement() {
      this.agreeToTerms = !this.agreeToTerms
    },
    
    // 检查用户名
    async checkUsername() {
      const username = this.formData.username.trim()
      
      if (!username) {
        this.usernameStatus = { text: '', class: '' }
        return
      }
      
      if (username.length < 3 || username.length > 20) {
        this.usernameStatus = { text: '用户名长度为3-20位', class: 'error' }
        return
      }
      
      // 检查用户名格式
      const usernameReg = /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/
      if (!usernameReg.test(username)) {
        this.usernameStatus = { text: '用户名只能包含字母、数字、下划线和中文', class: 'error' }
        return
      }
      
      // 这里可以调用API检查用户名是否已存在
      // 为了演示，我们简单检查一些常见的用户名
      const commonUsernames = ['admin', 'test', 'user', 'guest']
      if (commonUsernames.includes(username.toLowerCase())) {
        this.usernameStatus = { text: '该用户名已被使用', class: 'error' }
        return
      }
      
      this.usernameStatus = { text: '用户名可用', class: 'success' }
    },
    
    // 检查密码强度
    checkPasswordStrength() {
      const password = this.formData.password
      
      if (!password) {
        this.passwordStrength.show = false
        return
      }
      
      this.passwordStrength.show = true
      
      let score = 0
      let text = ''
      let className = ''
      let width = '0%'
      
      // 长度检查
      if (password.length >= 6) score += 1
      if (password.length >= 8) score += 1
      
      // 复杂度检查
      if (/[a-z]/.test(password)) score += 1
      if (/[A-Z]/.test(password)) score += 1
      if (/[0-9]/.test(password)) score += 1
      if (/[^a-zA-Z0-9]/.test(password)) score += 1
      
      if (score <= 2) {
        text = '弱'
        className = 'weak'
        width = '33%'
      } else if (score <= 4) {
        text = '中'
        className = 'medium'
        width = '66%'
      } else {
        text = '强'
        className = 'strong'
        width = '100%'
      }
      
      this.passwordStrength = { show: true, text, class: className, width }
    },
    
    // 检查密码匹配
    checkPasswordMatch() {
      if (!this.formData.confirmPassword) {
        this.passwordMatchError = ''
        return
      }
      
      if (this.formData.password !== this.formData.confirmPassword) {
        this.passwordMatchError = '两次输入的密码不一致'
      } else {
        this.passwordMatchError = ''
      }
    },
    
    // 检查手机号
    checkPhone() {
      const phone = this.formData.phone.trim()
      
      if (!phone) {
        this.phoneStatus = { text: '', class: '' }
        return
      }
      
      const phoneReg = /^1[3-9]\d{9}$/
      if (!phoneReg.test(phone)) {
        this.phoneStatus = { text: '请输入正确的手机号', class: 'error' }
        return
      }
      
      this.phoneStatus = { text: '手机号格式正确', class: 'success' }
    },
    
    // 检查邮箱
    checkEmail() {
      const email = this.formData.email.trim()
      
      if (!email) {
        this.emailStatus = { text: '', class: '' }
        return
      }
      
      const emailReg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailReg.test(email)) {
        this.emailStatus = { text: '请输入正确的邮箱地址', class: 'error' }
        return
      }
      
      this.emailStatus = { text: '邮箱格式正确', class: 'success' }
    },
    
    // 处理注册
    async handleRegister() {
      if (!this.canRegister || this.isLoading) {
        return
      }
      
      this.isLoading = true
      
      try {
        const registerData = {
          username: this.formData.username.trim(),
          password: this.formData.password,
          phone: this.formData.phone.trim(),
          email: this.formData.email.trim(),
          real_name: this.formData.realName.trim()
        }
        
        const res = await this.$http.post('/api/v1/auth/register', registerData)
        
        // 注册成功，自动登录
        const userStore = useUserStore()
        userStore.login(res.token, res.user)
        
        uni.showToast({
          title: '注册成功',
          icon: 'success'
        })
        
        // 跳转到首页
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/index/index'
          })
        }, 1500)
        
      } catch (error) {
        console.error('注册失败:', error)
        uni.showToast({
          title: error.error || '注册失败',
          icon: 'none'
        })
      } finally {
        this.isLoading = false
      }
    },
    
    // 跳转到登录页
    goToLogin() {
      uni.navigateBack()
    },
    
    // 显示协议
    showAgreement(type) {
      const title = type === 'privacy' ? '隐私政策' : '服务协议'
      const content = type === 'privacy' ? 
        '我们非常重视您的隐私保护，您的个人信息将被严格保密，仅用于提供更好的回收服务...' :
        '欢迎使用我们的回收服务，请仔细阅读以下服务条款：\n1. 我们承诺提供专业的设备回收服务\n2. 价格评估公正透明\n3. 保护用户数据安全...'
      
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
.register {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40rpx;
}

.register-header {
  text-align: center;
  margin-bottom: 60rpx;
  
  .header-title {
    font-size: 48rpx;
    font-weight: bold;
    color: #fff;
    margin-bottom: 16rpx;
  }
  
  .header-subtitle {
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.8);
  }
}

.register-form {
  background: #fff;
  border-radius: 24rpx;
  padding: 60rpx 40rpx;
  
  .form-group {
    margin-bottom: 40rpx;
    
    .form-label {
      display: block;
      font-size: 28rpx;
      font-weight: bold;
      color: #333;
      margin-bottom: 16rpx;
    }
    
    .input-wrapper {
      position: relative;
      display: flex;
      align-items: center;
      border: 2rpx solid #e0e0e0;
      border-radius: 12rpx;
      background: #f8f9fa;
      
      &:focus-within {
        border-color: #667eea;
        background: #fff;
      }
      
      .form-input {
        flex: 1;
        height: 80rpx;
        padding: 0 20rpx;
        font-size: 28rpx;
        color: #333;
      }
      
      .password-toggle {
        width: 60rpx;
        text-align: center;
        font-size: 28rpx;
        color: #999;
      }
      
      .input-status {
        position: absolute;
        right: 20rpx;
        font-size: 24rpx;
        
        &.success {
          color: #27ae60;
        }
        
        &.error {
          color: #f44336;
        }
      }
    }
    
    .input-error {
      display: block;
      font-size: 24rpx;
      color: #f44336;
      margin-top: 8rpx;
    }
    
    .password-strength {
      margin-top: 16rpx;
      
      .strength-text {
        font-size: 24rpx;
        margin-bottom: 8rpx;
        
        &.weak {
          color: #f44336;
        }
        
        &.medium {
          color: #ff9800;
        }
        
        &.strong {
          color: #27ae60;
        }
      }
      
      .strength-bar {
        height: 6rpx;
        background: #e0e0e0;
        border-radius: 3rpx;
        overflow: hidden;
        
        .strength-level {
          height: 100%;
          border-radius: 3rpx;
          transition: width 0.3s ease;
          
          &.weak {
            background: #f44336;
          }
          
          &.medium {
            background: #ff9800;
          }
          
          &.strong {
            background: #27ae60;
          }
        }
      }
    }
  }
  
  .agreement-section {
    margin-bottom: 40rpx;
    
    .agreement-check {
      display: flex;
      align-items: flex-start;
      
      .check-icon {
        width: 40rpx;
        height: 40rpx;
        border: 2rpx solid #e0e0e0;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24rpx;
        margin-right: 16rpx;
        margin-top: 4rpx;
        
        &.checked {
          background: #667eea;
          border-color: #667eea;
          color: #fff;
        }
      }
      
      .agreement-text {
        flex: 1;
        font-size: 26rpx;
        color: #666;
        line-height: 1.5;
        
        .agreement-link {
          color: #667eea;
          text-decoration: underline;
        }
      }
    }
  }
  
  .form-actions {
    .register-btn {
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
</style>
