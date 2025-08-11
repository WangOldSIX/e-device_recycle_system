<template>
  <view class="device-detail" v-if="device">
    <!-- 设备图片 -->
    <view class="device-images">
      <swiper class="image-swiper" :indicator-dots="true" :circular="true">
        <swiper-item v-for="(image, index) in deviceImages" :key="index">
          <image :src="image" class="device-image" mode="aspectFit"></image>
        </swiper-item>
      </swiper>
    </view>
    
    <!-- 设备基本信息 -->
    <view class="device-info card">
      <view class="device-title">
        <view class="device-name">{{ device.name }}</view>
        <view class="device-brand">{{ device.brand }} {{ device.model }}</view>
      </view>
      
      <view class="device-specs">
        <view class="spec-row" v-if="device.cpu">
          <text class="spec-label">处理器</text>
          <text class="spec-value">{{ device.cpu }}</text>
        </view>
        <view class="spec-row" v-if="device.memory">
          <text class="spec-label">内存</text>
          <text class="spec-value">{{ device.memory }}</text>
        </view>
        <view class="spec-row" v-if="device.storage">
          <text class="spec-label">存储</text>
          <text class="spec-value">{{ device.storage }}</text>
        </view>
        <view class="spec-row" v-if="device.graphics">
          <text class="spec-label">显卡</text>
          <text class="spec-value">{{ device.graphics }}</text>
        </view>
        <view class="spec-row" v-if="device.screen">
          <text class="spec-label">屏幕</text>
          <text class="spec-value">{{ device.screen }}</text>
        </view>
        <view class="spec-row" v-if="device.year_bought">
          <text class="spec-label">购买年份</text>
          <text class="spec-value">{{ device.year_bought }}年</text>
        </view>
        <view class="spec-row">
          <text class="spec-label">成色</text>
          <text class="spec-value condition" :class="device.condition">
            {{ $utils.getConditionText(device.condition) }}
          </text>
        </view>
      </view>
    </view>
    
    <!-- 回收价格 -->
    <view class="price-info card">
      <view class="price-title">参考回收价格</view>
      <view class="price-content">
        <view class="base-price">
          <text class="price-label">基础价格</text>
          <text class="price-value">{{ $utils.formatPrice(device.base_price) }}</text>
        </view>
        <view class="estimated-price">
          <text class="price-label">预估价格</text>
          <text class="price-value highlight">{{ $utils.formatPrice(estimatedPrice) }}</text>
        </view>
      </view>
      <view class="price-note">
        <text class="note-text">* 最终价格以现场评估为准</text>
      </view>
    </view>
    
    <!-- 设备描述 -->
    <view class="device-description card" v-if="device.description">
      <view class="desc-title">设备描述</view>
      <view class="desc-content">{{ device.description }}</view>
    </view>
    
    <!-- 回收流程 -->
    <view class="recycle-process card">
      <view class="process-title">回收流程</view>
      <view class="process-steps">
        <view class="step-item">
          <view class="step-icon">📝</view>
          <view class="step-text">提交申请</view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-icon">📞</view>
          <view class="step-text">预约上门</view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-icon">🔍</view>
          <view class="step-text">现场评估</view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-icon">💰</view>
          <view class="step-text">即时付款</view>
        </view>
      </view>
    </view>
    
    <!-- 底部操作按钮 -->
    <view class="bottom-actions">
      <button class="action-btn secondary" @click="contactService">咨询客服</button>
      <button class="action-btn primary" @click="createOrder">立即回收</button>
    </view>
  </view>
  
  <!-- 加载状态 -->
  <view class="loading" v-else>
    <text class="loading-text">加载中...</text>
  </view>
</template>

<script>
import { useUserStore } from '@/store'

export default {
  name: 'DeviceDetail',
  data() {
    return {
      deviceId: null,
      device: null,
      estimatedPrice: 0
    }
  },
  
  computed: {
    deviceImages() {
      if (!this.device || !this.device.images) {
        return ['/static/images/device-placeholder.png']
      }
      try {
        const images = JSON.parse(this.device.images)
        return images.length > 0 ? images : ['/static/images/device-placeholder.png']
      } catch {
        return ['/static/images/device-placeholder.png']
      }
    }
  },
  
  onLoad(options) {
    this.deviceId = options.id
    if (this.deviceId) {
      this.loadDeviceDetail()
    }
  },
  
  methods: {
    // 加载设备详情
    async loadDeviceDetail() {
      try {
        const res = await this.$http.get(`/api/v1/devices/${this.deviceId}`)
        this.device = res.device
        this.calculateEstimatedPrice()
      } catch (error) {
        console.error('加载设备详情失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'none'
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 2000)
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
    
    // 联系客服
    contactService() {
      uni.showModal({
        title: '联系客服',
        content: '客服电话：400-123-4567\n工作时间：9:00-18:00',
        showCancel: false
      })
    },
    
    // 创建回收订单
    createOrder() {
      // 检查登录状态
      if (!this.$utils.checkLogin()) {
        return
      }
      
      // 跳转到订单创建页面
      uni.navigateTo({
        url: `/pages/order/create?deviceId=${this.deviceId}`
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.device-detail {
  padding-bottom: 140rpx; // 为底部按钮留出空间
}

.device-images {
  .image-swiper {
    width: 100%;
    height: 600rpx;
    background: #f8f9fa;
    
    .device-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}

.device-info {
  margin: 20rpx;
  
  .device-title {
    margin-bottom: 30rpx;
    padding-bottom: 20rpx;
    border-bottom: 1rpx solid #f0f0f0;
    
    .device-name {
      font-size: 36rpx;
      font-weight: bold;
      color: #333;
      margin-bottom: 8rpx;
    }
    
    .device-brand {
      font-size: 28rpx;
      color: #666;
    }
  }
  
  .device-specs {
    .spec-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16rpx 0;
      border-bottom: 1rpx solid #f8f9fa;
      
      &:last-child {
        border-bottom: none;
      }
      
      .spec-label {
        font-size: 28rpx;
        color: #666;
        flex-shrink: 0;
        width: 120rpx;
      }
      
      .spec-value {
        font-size: 28rpx;
        color: #333;
        text-align: right;
        flex: 1;
        
        &.condition {
          &.excellent {
            color: #27ae60;
          }
          
          &.good {
            color: #2196f3;
          }
          
          &.fair {
            color: #ff9800;
          }
          
          &.poor {
            color: #f44336;
          }
        }
      }
    }
  }
}

.price-info {
  margin: 20rpx;
  
  .price-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 20rpx;
  }
  
  .price-content {
    .base-price,
    .estimated-price {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16rpx 0;
      
      .price-label {
        font-size: 28rpx;
        color: #666;
      }
      
      .price-value {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        
        &.highlight {
          color: #667eea;
          font-size: 36rpx;
        }
      }
    }
    
    .estimated-price {
      padding-top: 20rpx;
      border-top: 1rpx solid #f0f0f0;
    }
  }
  
  .price-note {
    margin-top: 20rpx;
    
    .note-text {
      font-size: 24rpx;
      color: #999;
    }
  }
}

.device-description {
  margin: 20rpx;
  
  .desc-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 16rpx;
  }
  
  .desc-content {
    font-size: 28rpx;
    color: #666;
    line-height: 1.6;
  }
}

.recycle-process {
  margin: 20rpx;
  
  .process-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 30rpx;
  }
  
  .process-steps {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .step-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      flex: 1;
      
      .step-icon {
        font-size: 48rpx;
        margin-bottom: 12rpx;
      }
      
      .step-text {
        font-size: 24rpx;
        color: #666;
        text-align: center;
      }
    }
    
    .step-arrow {
      font-size: 24rpx;
      color: #ccc;
      margin: 0 10rpx;
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
  display: flex;
  gap: 20rpx;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.1);
  
  .action-btn {
    flex: 1;
    height: 80rpx;
    border-radius: 40rpx;
    font-size: 32rpx;
    font-weight: bold;
    border: none;
    
    &.secondary {
      background: #f8f9fa;
      color: #666;
    }
    
    &.primary {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: #fff;
    }
  }
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400rpx;
  
  .loading-text {
    font-size: 28rpx;
    color: #999;
  }
}
</style>
