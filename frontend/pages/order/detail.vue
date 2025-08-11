<template>
  <view class="order-detail" v-if="order">
    <!-- 订单状态 -->
    <view class="status-card">
      <view class="status-header">
        <view class="status-icon" :style="{ background: $utils.getOrderStatusColor(order.status) }">
          {{ getStatusIcon(order.status) }}
        </view>
        <view class="status-info">
          <view class="status-text">{{ $utils.getOrderStatusText(order.status) }}</view>
          <view class="status-desc">{{ getStatusDescription(order.status) }}</view>
        </view>
      </view>
      
      <!-- 进度条 -->
      <view class="progress-bar">
        <view class="progress-steps">
          <view 
            class="progress-step" 
            :class="{ active: isStepActive(step.status), completed: isStepCompleted(step.status) }"
            v-for="step in progressSteps" 
            :key="step.status"
          >
            <view class="step-dot"></view>
            <view class="step-label">{{ step.label }}</view>
          </view>
        </view>
        <view class="progress-line" :style="{ width: progressWidth }"></view>
      </view>
    </view>
    
    <!-- 订单信息 -->
    <view class="order-info card">
      <view class="info-title">订单信息</view>
      <view class="info-list">
        <view class="info-item">
          <text class="info-label">订单号</text>
          <text class="info-value">{{ order.order_no }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">下单时间</text>
          <text class="info-value">{{ $utils.formatTime(order.created_at) }}</text>
        </view>
        <view class="info-item" v-if="order.pickup_time">
          <text class="info-label">预约时间</text>
          <text class="info-value">{{ $utils.formatTime(order.pickup_time) }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">联系人</text>
          <text class="info-value">{{ order.contact_name }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">联系电话</text>
          <text class="info-value">{{ order.contact_phone }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">上门地址</text>
          <text class="info-value">{{ order.pickup_address }}</text>
        </view>
      </view>
    </view>
    
    <!-- 设备信息 -->
    <view class="device-info card" v-if="order.device">
      <view class="info-title">设备信息</view>
      <view class="device-content">
        <image 
          :src="getDeviceImage(order.device.images)" 
          class="device-image"
        ></image>
        <view class="device-details">
          <view class="device-name">{{ order.device.name }}</view>
          <view class="device-brand">{{ order.device.brand }} {{ order.device.model }}</view>
          <view class="device-specs">
            <text class="spec-item" v-if="order.device.cpu">{{ order.device.cpu }}</text>
            <text class="spec-item" v-if="order.device.memory">{{ order.device.memory }}</text>
            <text class="spec-item" v-if="order.device.storage">{{ order.device.storage }}</text>
          </view>
          <view class="device-condition">
            <text class="condition-tag" :class="order.device.condition">
              {{ $utils.getConditionText(order.device.condition) }}
            </text>
          </view>
        </view>
      </view>
      
      <!-- 设备详细描述 -->
      <view class="device-description" v-if="order.device_info">
        <view class="desc-title">设备描述</view>
        <view class="desc-content">{{ order.device_info }}</view>
      </view>
    </view>
    
    <!-- 价格信息 -->
    <view class="price-info card">
      <view class="info-title">价格信息</view>
      <view class="price-list">
        <view class="price-item">
          <text class="price-label">预估价格</text>
          <text class="price-value">{{ $utils.formatPrice(order.estimated_price) }}</text>
        </view>
        <view class="price-item" v-if="order.final_price">
          <text class="price-label">最终价格</text>
          <text class="price-value final">{{ $utils.formatPrice(order.final_price) }}</text>
        </view>
      </view>
    </view>
    
    <!-- 评估信息 -->
    <view class="evaluation-info card" v-if="order.evaluation">
      <view class="info-title">评估信息</view>
      <view class="evaluation-content">
        <view class="score-grid">
          <view class="score-item">
            <view class="score-value">{{ order.evaluation.appearance_score }}</view>
            <view class="score-label">外观评分</view>
          </view>
          <view class="score-item">
            <view class="score-value">{{ order.evaluation.function_score }}</view>
            <view class="score-label">功能评分</view>
          </view>
          <view class="score-item">
            <view class="score-value">{{ order.evaluation.performance_score }}</view>
            <view class="score-label">性能评分</view>
          </view>
          <view class="score-item">
            <view class="score-value">{{ order.evaluation.overall_score.toFixed(1) }}</view>
            <view class="score-label">综合评分</view>
          </view>
        </view>
        
        <view class="evaluation-report" v-if="order.evaluation.evaluation_report">
          <view class="report-title">评估报告</view>
          <view class="report-content">{{ order.evaluation.evaluation_report }}</view>
        </view>
      </view>
    </view>
    
    <!-- 备注信息 -->
    <view class="remark-info card" v-if="order.remark">
      <view class="info-title">备注信息</view>
      <view class="remark-content">{{ order.remark }}</view>
    </view>
    
    <!-- 底部操作 -->
    <view class="bottom-actions" v-if="showActions">
      <button 
        class="action-btn cancel" 
        v-if="order.status === 'pending'"
        @click="cancelOrder"
      >
        取消订单
      </button>
      <button 
        class="action-btn contact" 
        @click="contactService"
      >
        联系客服
      </button>
      <button 
        class="action-btn reorder" 
        v-if="order.status === 'completed'"
        @click="reorder"
      >
        再次回收
      </button>
    </view>
  </view>
  
  <!-- 加载状态 -->
  <view class="loading" v-else>
    <text class="loading-text">加载中...</text>
  </view>
</template>

<script>
export default {
  name: 'OrderDetail',
  data() {
    return {
      orderId: null,
      order: null,
      progressSteps: [
        { status: 'pending', label: '待处理' },
        { status: 'confirmed', label: '已确认' },
        { status: 'picked_up', label: '已上门' },
        { status: 'evaluated', label: '已评估' },
        { status: 'completed', label: '已完成' }
      ]
    }
  },
  
  computed: {
    showActions() {
      return this.order && ['pending', 'confirmed', 'picked_up', 'evaluated', 'completed'].includes(this.order.status)
    },
    
    progressWidth() {
      if (!this.order) return '0%'
      
      const currentIndex = this.progressSteps.findIndex(step => step.status === this.order.status)
      if (currentIndex === -1) return '0%'
      
      const percentage = (currentIndex / (this.progressSteps.length - 1)) * 100
      return percentage + '%'
    }
  },
  
  onLoad(options) {
    this.orderId = options.id
    if (this.orderId) {
      this.loadOrderDetail()
    }
  },
  
  methods: {
    // 加载订单详情
    async loadOrderDetail() {
      try {
        const res = await this.$http.get(`/api/v1/orders/${this.orderId}`)
        this.order = res.order
      } catch (error) {
        console.error('加载订单详情失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'none'
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 2000)
      }
    },
    
    // 获取状态图标
    getStatusIcon(status) {
      const iconMap = {
        'pending': '⏳',
        'confirmed': '✅',
        'picked_up': '🚚',
        'evaluated': '📊',
        'completed': '💰',
        'cancelled': '❌'
      }
      return iconMap[status] || '❓'
    },
    
    // 获取状态描述
    getStatusDescription(status) {
      const descMap = {
        'pending': '我们已收到您的回收申请，将尽快处理',
        'confirmed': '订单已确认，我们将安排工作人员联系您',
        'picked_up': '工作人员已上门，正在进行设备检测',
        'evaluated': '设备评估完成，请确认回收价格',
        'completed': '回收完成，感谢您选择我们的服务',
        'cancelled': '订单已取消'
      }
      return descMap[status] || ''
    },
    
    // 判断步骤是否激活
    isStepActive(stepStatus) {
      return this.order && this.order.status === stepStatus
    },
    
    // 判断步骤是否完成
    isStepCompleted(stepStatus) {
      if (!this.order) return false
      
      const currentIndex = this.progressSteps.findIndex(step => step.status === this.order.status)
      const stepIndex = this.progressSteps.findIndex(step => step.status === stepStatus)
      
      return stepIndex < currentIndex
    },
    
    // 获取设备图片
    getDeviceImage(images) {
      if (!images) return '/static/images/device-placeholder.png'
      try {
        const imageList = JSON.parse(images)
        return imageList.length > 0 ? imageList[0] : '/static/images/device-placeholder.png'
      } catch {
        return '/static/images/device-placeholder.png'
      }
    },
    
    // 取消订单
    async cancelOrder() {
      uni.showModal({
        title: '确认取消',
        content: '确定要取消这个回收订单吗？',
        success: async (res) => {
          if (res.confirm) {
            try {
              await this.$http.put(`/api/v1/orders/${this.orderId}/cancel`)
              
              uni.showToast({
                title: '订单已取消',
                icon: 'success'
              })
              
              // 重新加载订单详情
              this.loadOrderDetail()
              
            } catch (error) {
              console.error('取消订单失败:', error)
              uni.showToast({
                title: error.error || '取消失败',
                icon: 'none'
              })
            }
          }
        }
      })
    },
    
    // 联系客服
    contactService() {
      uni.showActionSheet({
        itemList: ['电话客服', '在线客服'],
        success: (res) => {
          switch (res.tapIndex) {
            case 0:
              uni.makePhoneCall({
                phoneNumber: '400-123-4567'
              })
              break
            case 1:
              uni.showToast({
                title: '在线客服功能开发中',
                icon: 'none'
              })
              break
          }
        }
      })
    },
    
    // 再次回收
    reorder() {
      uni.navigateTo({
        url: `/pages/order/create?deviceId=${this.order.device_id}`
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.order-detail {
  background: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 140rpx;
}

.status-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40rpx;
  color: #fff;
  
  .status-header {
    display: flex;
    align-items: center;
    margin-bottom: 40rpx;
    
    .status-icon {
      width: 80rpx;
      height: 80rpx;
      border-radius: 40rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 40rpx;
      margin-right: 24rpx;
      background: rgba(255, 255, 255, 0.2);
    }
    
    .status-info {
      flex: 1;
      
      .status-text {
        font-size: 32rpx;
        font-weight: bold;
        margin-bottom: 8rpx;
      }
      
      .status-desc {
        font-size: 26rpx;
        opacity: 0.9;
        line-height: 1.4;
      }
    }
  }
  
  .progress-bar {
    position: relative;
    
    .progress-steps {
      display: flex;
      justify-content: space-between;
      position: relative;
      z-index: 2;
      
      .progress-step {
        display: flex;
        flex-direction: column;
        align-items: center;
        
        .step-dot {
          width: 24rpx;
          height: 24rpx;
          border-radius: 12rpx;
          background: rgba(255, 255, 255, 0.3);
          margin-bottom: 12rpx;
          
          .active &,
          .completed & {
            background: #fff;
          }
        }
        
        .step-label {
          font-size: 22rpx;
          opacity: 0.7;
          
          .active &,
          .completed & {
            opacity: 1;
            font-weight: bold;
          }
        }
      }
    }
    
    .progress-line {
      position: absolute;
      top: 12rpx;
      left: 0;
      height: 4rpx;
      background: #fff;
      border-radius: 2rpx;
      transition: width 0.3s ease;
      z-index: 1;
    }
  }
}

.order-info,
.device-info,
.price-info,
.evaluation-info,
.remark-info {
  margin: 20rpx;
  
  .info-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 24rpx;
    padding-bottom: 16rpx;
    border-bottom: 1rpx solid #f0f0f0;
  }
  
  .info-list {
    .info-item {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding: 16rpx 0;
      border-bottom: 1rpx solid #f8f9fa;
      
      &:last-child {
        border-bottom: none;
      }
      
      .info-label {
        font-size: 28rpx;
        color: #666;
        width: 140rpx;
        flex-shrink: 0;
      }
      
      .info-value {
        font-size: 28rpx;
        color: #333;
        text-align: right;
        flex: 1;
        word-break: break-all;
      }
    }
  }
}

.device-info {
  .device-content {
    display: flex;
    margin-bottom: 24rpx;
    
    .device-image {
      width: 160rpx;
      height: 160rpx;
      border-radius: 12rpx;
      margin-right: 24rpx;
      object-fit: cover;
    }
    
    .device-details {
      flex: 1;
      
      .device-name {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 12rpx;
      }
      
      .device-brand {
        font-size: 28rpx;
        color: #666;
        margin-bottom: 12rpx;
      }
      
      .device-specs {
        display: flex;
        flex-wrap: wrap;
        gap: 12rpx;
        margin-bottom: 12rpx;
        
        .spec-item {
          font-size: 24rpx;
          color: #666;
          background: #f8f9fa;
          padding: 4rpx 12rpx;
          border-radius: 8rpx;
        }
      }
      
      .device-condition {
        .condition-tag {
          font-size: 24rpx;
          padding: 6rpx 12rpx;
          border-radius: 8rpx;
          
          &.excellent {
            background: #e8f5e8;
            color: #27ae60;
          }
          
          &.good {
            background: #e3f2fd;
            color: #2196f3;
          }
          
          &.fair {
            background: #fff3e0;
            color: #ff9800;
          }
          
          &.poor {
            background: #ffebee;
            color: #f44336;
          }
        }
      }
    }
  }
  
  .device-description {
    .desc-title {
      font-size: 28rpx;
      font-weight: bold;
      color: #333;
      margin-bottom: 12rpx;
    }
    
    .desc-content {
      font-size: 26rpx;
      color: #666;
      line-height: 1.5;
      background: #f8f9fa;
      padding: 20rpx;
      border-radius: 12rpx;
    }
  }
}

.price-info {
  .price-list {
    display: flex;
    justify-content: space-around;
    
    .price-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      
      .price-label {
        font-size: 26rpx;
        color: #999;
        margin-bottom: 12rpx;
      }
      
      .price-value {
        font-size: 36rpx;
        font-weight: bold;
        color: #333;
        
        &.final {
          color: #667eea;
          font-size: 40rpx;
        }
      }
    }
  }
}

.evaluation-info {
  .evaluation-content {
    .score-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 20rpx;
      margin-bottom: 30rpx;
      
      .score-item {
        text-align: center;
        background: #f8f9fa;
        padding: 24rpx 12rpx;
        border-radius: 12rpx;
        
        .score-value {
          font-size: 36rpx;
          font-weight: bold;
          color: #667eea;
          margin-bottom: 8rpx;
        }
        
        .score-label {
          font-size: 24rpx;
          color: #666;
        }
      }
    }
    
    .evaluation-report {
      .report-title {
        font-size: 28rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 12rpx;
      }
      
      .report-content {
        font-size: 26rpx;
        color: #666;
        line-height: 1.5;
        background: #f8f9fa;
        padding: 20rpx;
        border-radius: 12rpx;
      }
    }
  }
}

.remark-info {
  .remark-content {
    font-size: 26rpx;
    color: #666;
    line-height: 1.5;
    background: #f8f9fa;
    padding: 20rpx;
    border-radius: 12rpx;
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
    font-size: 28rpx;
    font-weight: bold;
    border: none;
    
    &.cancel {
      background: #ffebee;
      color: #f44336;
    }
    
    &.contact {
      background: #e3f2fd;
      color: #2196f3;
    }
    
    &.reorder {
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
