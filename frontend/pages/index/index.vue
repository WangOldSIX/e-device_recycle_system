<template>
  <view class="index">
    <!-- 顶部轮播图 -->
    <swiper class="banner" :indicator-dots="true" :autoplay="true" :interval="3000" :duration="500">
      <swiper-item v-for="(item, index) in banners" :key="index">
        <image :src="item.image" class="banner-image"></image>
      </swiper-item>
    </swiper>
    
    <!-- 快捷功能 -->
    <view class="quick-actions card">
      <view class="action-title">快速回收</view>
      <view class="actions-grid">
        <view class="action-item" @click="navigateToDevices('laptop')">
          <image src="/static/icons/laptop.png" class="action-icon"></image>
          <text class="action-text">笔记本</text>
        </view>
        <view class="action-item" @click="navigateToDevices('desktop')">
          <image src="/static/icons/desktop.png" class="action-icon"></image>
          <text class="action-text">台式机</text>
        </view>
        <view class="action-item" @click="navigateToDevices('tablet')">
          <image src="/static/icons/tablet.png" class="action-icon"></image>
          <text class="action-text">平板</text>
        </view>
        <view class="action-item" @click="navigateToDevices('phone')">
          <image src="/static/icons/phone.png" class="action-icon"></image>
          <text class="action-text">手机</text>
        </view>
      </view>
    </view>
    
    <!-- 服务优势 -->
    <view class="advantages card">
      <view class="action-title">服务优势</view>
      <view class="advantage-list">
        <view class="advantage-item">
          <image src="/static/icons/price.png" class="advantage-icon"></image>
          <view class="advantage-content">
            <view class="advantage-title">价格公道</view>
            <view class="advantage-desc">专业评估，价格透明</view>
          </view>
        </view>
        <view class="advantage-item">
          <image src="/static/icons/pickup.png" class="advantage-icon"></image>
          <view class="advantage-content">
            <view class="advantage-title">上门回收</view>
            <view class="advantage-desc">免费上门，安全便捷</view>
          </view>
        </view>
        <view class="advantage-item">
          <image src="/static/icons/fast.png" class="advantage-icon"></image>
          <view class="advantage-content">
            <view class="advantage-title">快速到账</view>
            <view class="advantage-desc">当面验机，即时付款</view>
          </view>
        </view>
        <view class="advantage-item">
          <image src="/static/icons/security.png" class="advantage-icon"></image>
          <view class="advantage-content">
            <view class="advantage-title">数据安全</view>
            <view class="advantage-desc">专业清除，保护隐私</view>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 热门设备 -->
    <view class="hot-devices card">
      <view class="flex-between mb-20">
        <view class="action-title">热门设备</view>
        <view class="more-btn" @click="navigateToDevices()">
          <text class="more-text">查看更多</text>
          <text class="more-arrow">></text>
        </view>
      </view>
      <view class="device-list">
        <view class="device-item" v-for="device in hotDevices" :key="device.id" @click="navigateToDeviceDetail(device.id)">
          <image :src="device.images || '/static/images/device-placeholder.png'" class="device-image"></image>
          <view class="device-info">
            <view class="device-name">{{ device.name }}</view>
            <view class="device-brand">{{ device.brand }} {{ device.model }}</view>
            <view class="device-price">参考回收价：{{ $utils.formatPrice(device.base_price) }}</view>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 回收流程 -->
    <view class="process card">
      <view class="action-title">回收流程</view>
      <view class="process-steps">
        <view class="step-item">
          <view class="step-number">1</view>
          <view class="step-content">
            <view class="step-title">提交申请</view>
            <view class="step-desc">选择设备型号</view>
          </view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-number">2</view>
          <view class="step-content">
            <view class="step-title">预约上门</view>
            <view class="step-desc">选择合适时间</view>
          </view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-number">3</view>
          <view class="step-content">
            <view class="step-title">现场评估</view>
            <view class="step-desc">专业检测定价</view>
          </view>
        </view>
        <view class="step-arrow">→</view>
        <view class="step-item">
          <view class="step-number">4</view>
          <view class="step-content">
            <view class="step-title">即时付款</view>
            <view class="step-desc">当面交易结算</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { useDeviceStore } from '@/store'

export default {
  name: 'Index',
  data() {
    return {
      banners: [
        { image: '/static/images/banner1.jpg' },
        { image: '/static/images/banner2.jpg' },
        { image: '/static/images/banner3.jpg' }
      ],
      hotDevices: []
    }
  },
  
  onLoad() {
    this.loadHotDevices()
  },
  
  onShow() {
    // 每次显示页面时刷新数据
    this.loadHotDevices()
  },
  
  methods: {
    // 加载热门设备
    async loadHotDevices() {
      try {
        const res = await this.$http.get('/api/v1/devices', {
          page: 1,
          page_size: 4
        })
        this.hotDevices = res.devices || []
      } catch (error) {
        console.error('加载热门设备失败:', error)
      }
    },
    
    // 跳转到设备列表
    navigateToDevices(category = '') {
      let url = '/pages/devices/devices'
      if (category) {
        url += '?category=' + category
      }
      uni.navigateTo({
        url
      })
    },
    
    // 跳转到设备详情
    navigateToDeviceDetail(deviceId) {
      uni.navigateTo({
        url: `/pages/devices/detail?id=${deviceId}`
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.index {
  padding-bottom: 20rpx;
}

.banner {
  width: 100%;
  height: 360rpx;
  border-radius: 0 0 20rpx 20rpx;
  overflow: hidden;
  
  .banner-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.quick-actions {
  margin: 20rpx;
  
  .action-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 30rpx;
  }
  
  .actions-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 30rpx;
    
    .action-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      
      .action-icon {
        width: 80rpx;
        height: 80rpx;
        margin-bottom: 16rpx;
      }
      
      .action-text {
        font-size: 28rpx;
        color: #666;
      }
    }
  }
}

.advantages {
  margin: 20rpx;
  
  .advantage-list {
    .advantage-item {
      display: flex;
      align-items: center;
      padding: 20rpx 0;
      border-bottom: 1rpx solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .advantage-icon {
        width: 60rpx;
        height: 60rpx;
        margin-right: 20rpx;
      }
      
      .advantage-content {
        flex: 1;
        
        .advantage-title {
          font-size: 32rpx;
          font-weight: bold;
          color: #333;
          margin-bottom: 8rpx;
        }
        
        .advantage-desc {
          font-size: 26rpx;
          color: #999;
        }
      }
    }
  }
}

.hot-devices {
  margin: 20rpx;
  
  .more-btn {
    display: flex;
    align-items: center;
    
    .more-text {
      font-size: 28rpx;
      color: #667eea;
      margin-right: 8rpx;
    }
    
    .more-arrow {
      font-size: 28rpx;
      color: #667eea;
    }
  }
  
  .device-list {
    .device-item {
      display: flex;
      padding: 20rpx 0;
      border-bottom: 1rpx solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
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
          font-size: 32rpx;
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
}

.process {
  margin: 20rpx;
  
  .process-steps {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .step-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      flex: 1;
      
      .step-number {
        width: 60rpx;
        height: 60rpx;
        border-radius: 50%;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #fff;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28rpx;
        font-weight: bold;
        margin-bottom: 16rpx;
      }
      
      .step-content {
        text-align: center;
        
        .step-title {
          font-size: 28rpx;
          font-weight: bold;
          color: #333;
          margin-bottom: 8rpx;
        }
        
        .step-desc {
          font-size: 24rpx;
          color: #999;
        }
      }
    }
    
    .step-arrow {
      font-size: 32rpx;
      color: #ccc;
      margin: 0 10rpx;
    }
  }
}
</style>
