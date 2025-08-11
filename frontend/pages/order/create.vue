<template>
  <view class="order-create">
    <!-- 设备信息 -->
    <view class="device-summary card" v-if="device">
      <view class="summary-title">选择的设备</view>
      <view class="device-item">
        <image :src="deviceImage" class="device-image"></image>
        <view class="device-info">
          <view class="device-name">{{ device.name }}</view>
          <view class="device-brand">{{ device.brand }} {{ device.model }}</view>
          <view class="device-price">预估价格：{{ $utils.formatPrice(estimatedPrice) }}</view>
        </view>
      </view>
    </view>
    
    <!-- 联系信息 -->
    <view class="contact-info card">
      <view class="form-title">联系信息</view>
      <view class="form-group">
        <text class="form-label">联系人 *</text>
        <input 
          v-model="formData.contactName" 
          class="form-input" 
          placeholder="请输入联系人姓名"
          :value="formData.contactName"
        />
      </view>
      <view class="form-group">
        <text class="form-label">联系电话 *</text>
        <input 
          v-model="formData.contactPhone" 
          class="form-input" 
          placeholder="请输入联系电话"
          type="number"
          :value="formData.contactPhone"
        />
      </view>
    </view>
    
    <!-- 上门信息 -->
    <view class="pickup-info card">
      <view class="form-title">上门信息</view>
      <view class="form-group">
        <text class="form-label">上门地址 *</text>
        <textarea 
          v-model="formData.pickupAddress" 
          class="form-textarea" 
          placeholder="请输入详细的上门地址"
          :value="formData.pickupAddress"
        ></textarea>
      </view>
      <view class="form-group">
        <text class="form-label">预约时间</text>
        <view class="datetime-picker" @click="showDateTimePicker">
          <text class="datetime-text" :class="{ placeholder: !formData.pickupTime }">
            {{ formData.pickupTime ? $utils.formatTime(formData.pickupTime) : '请选择预约时间（可选）' }}
          </text>
          <text class="datetime-icon">📅</text>
        </view>
      </view>
    </view>
    
    <!-- 设备详情 -->
    <view class="device-details card">
      <view class="form-title">设备详情</view>
      <view class="form-group">
        <text class="form-label">设备详细描述</text>
        <textarea 
          v-model="formData.deviceInfo" 
          class="form-textarea" 
          placeholder="请描述设备的具体情况，如外观、功能等（可选）"
          :value="formData.deviceInfo"
        ></textarea>
      </view>
      <view class="form-group">
        <text class="form-label">设备图片</text>
        <view class="image-upload">
          <view class="image-list">
            <view class="image-item" v-for="(image, index) in uploadedImages" :key="index">
              <image :src="image" class="uploaded-image"></image>
              <view class="image-delete" @click="removeImage(index)">×</view>
            </view>
            <view class="upload-btn" @click="chooseImage" v-if="uploadedImages.length < 5">
              <text class="upload-icon">📷</text>
              <text class="upload-text">添加图片</text>
            </view>
          </view>
          <view class="upload-tip">最多可上传5张图片</view>
        </view>
      </view>
    </view>
    
    <!-- 备注信息 -->
    <view class="remark-info card">
      <view class="form-title">备注信息</view>
      <view class="form-group">
        <textarea 
          v-model="formData.remark" 
          class="form-textarea" 
          placeholder="其他需要说明的情况（可选）"
          :value="formData.remark"
        ></textarea>
      </view>
    </view>
    
    <!-- 服务说明 -->
    <view class="service-notice card">
      <view class="notice-title">服务说明</view>
      <view class="notice-list">
        <view class="notice-item">
          <text class="notice-icon">✓</text>
          <text class="notice-text">免费上门回收，专业评估</text>
        </view>
        <view class="notice-item">
          <text class="notice-icon">✓</text>
          <text class="notice-text">当面验机，即时付款</text>
        </view>
        <view class="notice-item">
          <text class="notice-icon">✓</text>
          <text class="notice-text">数据安全清除，保护隐私</text>
        </view>
        <view class="notice-item">
          <text class="notice-icon">✓</text>
          <text class="notice-text">价格透明，无隐性收费</text>
        </view>
      </view>
    </view>
    
    <!-- 底部提交按钮 -->
    <view class="bottom-actions">
      <button class="submit-btn" @click="submitOrder" :disabled="isSubmitting">
        <text v-if="!isSubmitting">提交申请</text>
        <text v-else>提交中...</text>
      </button>
    </view>
    
    <!-- 时间选择器 -->
    <picker 
      v-if="showPicker" 
      mode="datetime" 
      :value="pickerValue" 
      @change="onDateTimeChange"
      @cancel="showPicker = false"
    >
    </picker>
  </view>
</template>

<script>
import { useUserStore } from '@/store'

export default {
  name: 'OrderCreate',
  data() {
    return {
      deviceId: null,
      device: null,
      estimatedPrice: 0,
      formData: {
        contactName: '',
        contactPhone: '',
        pickupAddress: '',
        pickupTime: null,
        deviceInfo: '',
        remark: ''
      },
      uploadedImages: [],
      isSubmitting: false,
      showPicker: false,
      pickerValue: ''
    }
  },
  
  computed: {
    deviceImage() {
      if (!this.device || !this.device.images) {
        return '/static/images/device-placeholder.png'
      }
      try {
        const images = JSON.parse(this.device.images)
        return images.length > 0 ? images[0] : '/static/images/device-placeholder.png'
      } catch {
        return '/static/images/device-placeholder.png'
      }
    }
  },
  
  onLoad(options) {
    this.deviceId = options.deviceId
    if (this.deviceId) {
      this.loadDeviceInfo()
    }
    this.loadUserInfo()
  },
  
  methods: {
    // 加载设备信息
    async loadDeviceInfo() {
      try {
        const res = await this.$http.get(`/api/v1/devices/${this.deviceId}`)
        this.device = res.device
        this.calculateEstimatedPrice()
      } catch (error) {
        console.error('加载设备信息失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'none'
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 2000)
      }
    },
    
    // 加载用户信息
    loadUserInfo() {
      const userInfo = uni.getStorageSync('userInfo')
      if (userInfo) {
        this.formData.contactName = userInfo.real_name || userInfo.username
        this.formData.contactPhone = userInfo.phone || ''
      }
    },
    
    // 计算预估价格
    calculateEstimatedPrice() {
      if (!this.device) return
      
      const basePrice = this.device.base_price
      const condition = this.device.condition
      const yearBought = this.device.year_bought
      const currentYear = new Date().getFullYear()
      
      // 计算年份折旧
      const yearsSinceBoought = currentYear - yearBought
      let depreciationRate = 0.1 * yearsSinceBoought
      if (depreciationRate > 0.8) {
        depreciationRate = 0.8
      }
      
      // 成色系数
      const conditionMultiplier = {
        'excellent': 1.0,
        'good': 0.8,
        'fair': 0.6,
        'poor': 0.4
      }
      
      const multiplier = conditionMultiplier[condition] || 0.5
      
      // 计算最终价格
      let finalPrice = basePrice * (1 - depreciationRate) * multiplier
      
      // 确保价格不低于基础价格的10%
      const minPrice = basePrice * 0.1
      if (finalPrice < minPrice) {
        finalPrice = minPrice
      }
      
      this.estimatedPrice = finalPrice
    },
    
    // 显示时间选择器
    showDateTimePicker() {
      this.showPicker = true
      
      // 设置默认时间为明天9点
      const tomorrow = new Date()
      tomorrow.setDate(tomorrow.getDate() + 1)
      tomorrow.setHours(9, 0, 0, 0)
      this.pickerValue = tomorrow.toISOString().slice(0, 16)
    },
    
    // 时间选择改变
    onDateTimeChange(e) {
      this.formData.pickupTime = e.detail.value
      this.showPicker = false
    },
    
    // 选择图片
    chooseImage() {
      uni.chooseImage({
        count: 5 - this.uploadedImages.length,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          this.uploadedImages = [...this.uploadedImages, ...res.tempFilePaths]
        }
      })
    },
    
    // 删除图片
    removeImage(index) {
      this.uploadedImages.splice(index, 1)
    },
    
    // 表单验证
    validateForm() {
      if (!this.formData.contactName.trim()) {
        uni.showToast({
          title: '请输入联系人姓名',
          icon: 'none'
        })
        return false
      }
      
      if (!this.formData.contactPhone.trim()) {
        uni.showToast({
          title: '请输入联系电话',
          icon: 'none'
        })
        return false
      }
      
      // 验证手机号格式
      const phoneReg = /^1[3-9]\d{9}$/
      if (!phoneReg.test(this.formData.contactPhone)) {
        uni.showToast({
          title: '请输入正确的手机号',
          icon: 'none'
        })
        return false
      }
      
      if (!this.formData.pickupAddress.trim()) {
        uni.showToast({
          title: '请输入上门地址',
          icon: 'none'
        })
        return false
      }
      
      return true
    },
    
    // 提交订单
    async submitOrder() {
      if (!this.validateForm()) {
        return
      }
      
      if (this.isSubmitting) {
        return
      }
      
      this.isSubmitting = true
      
      try {
        const orderData = {
          device_id: this.deviceId,
          contact_name: this.formData.contactName,
          contact_phone: this.formData.contactPhone,
          pickup_address: this.formData.pickupAddress,
          pickup_time: this.formData.pickupTime,
          device_info: this.formData.deviceInfo,
          images: JSON.stringify(this.uploadedImages),
          remark: this.formData.remark
        }
        
        const res = await this.$http.post('/api/v1/orders', orderData)
        
        uni.showToast({
          title: '申请提交成功',
          icon: 'success'
        })
        
        // 跳转到订单详情页
        setTimeout(() => {
          uni.redirectTo({
            url: `/pages/order/detail?id=${res.order.id}`
          })
        }, 1500)
        
      } catch (error) {
        console.error('提交订单失败:', error)
        uni.showToast({
          title: error.error || '提交失败',
          icon: 'none'
        })
      } finally {
        this.isSubmitting = false
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.order-create {
  padding-bottom: 140rpx;
}

.device-summary {
  margin: 20rpx;
  
  .summary-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 20rpx;
  }
  
  .device-item {
    display: flex;
    align-items: center;
    
    .device-image {
      width: 120rpx;
      height: 120rpx;
      border-radius: 12rpx;
      margin-right: 20rpx;
      object-fit: cover;
    }
    
    .device-info {
      flex: 1;
      
      .device-name {
        font-size: 30rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 8rpx;
      }
      
      .device-brand {
        font-size: 26rpx;
        color: #666;
        margin-bottom: 8rpx;
      }
      
      .device-price {
        font-size: 28rpx;
        color: #667eea;
        font-weight: bold;
      }
    }
  }
}

.contact-info,
.pickup-info,
.device-details,
.remark-info {
  margin: 20rpx;
  
  .form-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 30rpx;
  }
  
  .form-group {
    margin-bottom: 30rpx;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .form-label {
      display: block;
      font-size: 28rpx;
      color: #333;
      margin-bottom: 12rpx;
    }
    
    .form-input {
      width: 100%;
      height: 80rpx;
      padding: 0 20rpx;
      border: 1rpx solid #e0e0e0;
      border-radius: 12rpx;
      font-size: 28rpx;
      background: #f8f9fa;
    }
    
    .form-textarea {
      width: 100%;
      min-height: 120rpx;
      padding: 20rpx;
      border: 1rpx solid #e0e0e0;
      border-radius: 12rpx;
      font-size: 28rpx;
      background: #f8f9fa;
      line-height: 1.5;
    }
    
    .datetime-picker {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 80rpx;
      padding: 0 20rpx;
      border: 1rpx solid #e0e0e0;
      border-radius: 12rpx;
      background: #f8f9fa;
      
      .datetime-text {
        font-size: 28rpx;
        color: #333;
        
        &.placeholder {
          color: #999;
        }
      }
      
      .datetime-icon {
        font-size: 28rpx;
        color: #666;
      }
    }
  }
  
  .image-upload {
    .image-list {
      display: flex;
      flex-wrap: wrap;
      gap: 20rpx;
      margin-bottom: 12rpx;
      
      .image-item {
        position: relative;
        width: 160rpx;
        height: 160rpx;
        
        .uploaded-image {
          width: 100%;
          height: 100%;
          border-radius: 12rpx;
          object-fit: cover;
        }
        
        .image-delete {
          position: absolute;
          top: -10rpx;
          right: -10rpx;
          width: 40rpx;
          height: 40rpx;
          border-radius: 50%;
          background: #f44336;
          color: #fff;
          font-size: 24rpx;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }
      
      .upload-btn {
        width: 160rpx;
        height: 160rpx;
        border: 2rpx dashed #ccc;
        border-radius: 12rpx;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        background: #f8f9fa;
        
        .upload-icon {
          font-size: 48rpx;
          margin-bottom: 8rpx;
        }
        
        .upload-text {
          font-size: 24rpx;
          color: #666;
        }
      }
    }
    
    .upload-tip {
      font-size: 24rpx;
      color: #999;
    }
  }
}

.service-notice {
  margin: 20rpx;
  
  .notice-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 20rpx;
  }
  
  .notice-list {
    .notice-item {
      display: flex;
      align-items: center;
      margin-bottom: 16rpx;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .notice-icon {
        width: 32rpx;
        height: 32rpx;
        border-radius: 50%;
        background: #27ae60;
        color: #fff;
        font-size: 20rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16rpx;
        flex-shrink: 0;
      }
      
      .notice-text {
        font-size: 28rpx;
        color: #666;
        line-height: 1.5;
      }
    }
  }
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 20rpx;
  border-top: 1rpx solid #f0f0f0;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.1);
  
  .submit-btn {
    width: 100%;
    height: 80rpx;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    border-radius: 40rpx;
    font-size: 32rpx;
    font-weight: bold;
    border: none;
    
    &:disabled {
      opacity: 0.6;
    }
  }
}
</style>
