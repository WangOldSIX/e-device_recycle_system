<template>
  <view class="order-list">
    <!-- 状态筛选 -->
    <view class="status-filter card">
      <scroll-view class="filter-scroll" scroll-x="true">
        <view class="filter-tabs">
          <view 
            class="filter-tab" 
            :class="{ active: selectedStatus === '' }" 
            @click="selectStatus('')"
          >
            全部
          </view>
          <view 
            class="filter-tab" 
            :class="{ active: selectedStatus === status.value }" 
            v-for="status in orderStatuses" 
            :key="status.value"
            @click="selectStatus(status.value)"
          >
            {{ status.label }}
          </view>
        </view>
      </scroll-view>
    </view>
    
    <!-- 订单列表 -->
    <view class="orders-container">
      <view 
        class="order-card" 
        v-for="order in orders" 
        :key="order.id"
        @click="navigateToDetail(order.id)"
      >
        <!-- 订单头部 -->
        <view class="order-header">
          <view class="order-info">
            <text class="order-no">订单号：{{ order.order_no }}</text>
            <text class="order-time">{{ $utils.formatTime(order.created_at) }}</text>
          </view>
          <view class="order-status" :style="{ color: $utils.getOrderStatusColor(order.status) }">
            {{ $utils.getOrderStatusText(order.status) }}
          </view>
        </view>
        
        <!-- 设备信息 -->
        <view class="device-info" v-if="order.device">
          <image 
            :src="getDeviceImage(order.device.images)" 
            class="device-image"
          ></image>
          <view class="device-details">
            <view class="device-name">{{ order.device.name }}</view>
            <view class="device-brand">{{ order.device.brand }} {{ order.device.model }}</view>
            <view class="device-condition">
              <text class="condition-tag" :class="order.device.condition">
                {{ $utils.getConditionText(order.device.condition) }}
              </text>
            </view>
          </view>
        </view>
        
        <!-- 价格信息 -->
        <view class="price-info">
          <view class="price-item">
            <text class="price-label">预估价格</text>
            <text class="price-value">{{ $utils.formatPrice(order.estimated_price) }}</text>
          </view>
          <view class="price-item" v-if="order.final_price">
            <text class="price-label">最终价格</text>
            <text class="price-value final">{{ $utils.formatPrice(order.final_price) }}</text>
          </view>
        </view>
        
        <!-- 联系信息 -->
        <view class="contact-info">
          <view class="contact-item">
            <text class="contact-label">联系人</text>
            <text class="contact-value">{{ order.contact_name }}</text>
          </view>
          <view class="contact-item">
            <text class="contact-label">联系电话</text>
            <text class="contact-value">{{ order.contact_phone }}</text>
          </view>
        </view>
        
        <!-- 上门信息 -->
        <view class="pickup-info" v-if="order.pickup_time">
          <text class="pickup-label">预约时间</text>
          <text class="pickup-time">{{ $utils.formatTime(order.pickup_time) }}</text>
        </view>
        
        <!-- 操作按钮 -->
        <view class="order-actions" v-if="showActions(order.status)">
          <button 
            class="action-btn cancel" 
            v-if="order.status === 'pending'"
            @click.stop="cancelOrder(order.id)"
          >
            取消订单
          </button>
          <button 
            class="action-btn contact" 
            v-if="['confirmed', 'picked_up', 'evaluated'].includes(order.status)"
            @click.stop="contactService(order.order_no)"
          >
            联系客服
          </button>
          <button 
            class="action-btn reorder" 
            v-if="order.status === 'completed'"
            @click.stop="reorder(order.device_id)"
          >
            再次回收
          </button>
        </view>
      </view>
    </view>
    
    <!-- 加载更多 -->
    <view class="load-more" v-if="hasMore">
      <view class="load-btn" @click="loadMore" :class="{ loading: isLoading }">
        <text v-if="!isLoading">加载更多</text>
        <text v-else>加载中...</text>
      </view>
    </view>
    
    <!-- 空状态 -->
    <view class="empty-state" v-if="orders.length === 0 && !isLoading">
      <image src="/static/images/empty.png" class="empty-image"></image>
      <text class="empty-text">暂无订单记录</text>
      <button class="empty-btn" @click="navigateToDevices">
        立即回收
      </button>
    </view>
  </view>
</template>

<script>
import { useOrderStore } from '@/store'

export default {
  name: 'OrderList',
  data() {
    return {
      selectedStatus: '',
      orders: [],
      orderStatuses: [
        { value: 'pending', label: '待处理' },
        { value: 'confirmed', label: '已确认' },
        { value: 'picked_up', label: '已上门' },
        { value: 'evaluated', label: '已评估' },
        { value: 'completed', label: '已完成' },
        { value: 'cancelled', label: '已取消' }
      ],
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0
      },
      isLoading: false,
      hasMore: true
    }
  },
  
  onLoad(options) {
    if (options.status) {
      this.selectedStatus = options.status
    }
    this.loadOrders(true)
  },
  
  onShow() {
    // 从详情页返回时刷新列表
    this.loadOrders(true)
  },
  
  methods: {
    // 加载订单列表
    async loadOrders(reset = false) {
      if (this.isLoading) return
      
      this.isLoading = true
      
      if (reset) {
        this.pagination.page = 1
        this.orders = []
        this.hasMore = true
      }
      
      try {
        const params = {
          page: this.pagination.page,
          page_size: this.pagination.pageSize
        }
        
        if (this.selectedStatus) {
          params.status = this.selectedStatus
        }
        
        const res = await this.$http.get('/api/v1/orders', params)
        
        if (reset) {
          this.orders = res.orders || []
        } else {
          this.orders = [...this.orders, ...(res.orders || [])]
        }
        
        this.pagination.total = res.pagination?.total || 0
        this.hasMore = this.orders.length < this.pagination.total
        
        if (!reset) {
          this.pagination.page++
        }
        
      } catch (error) {
        console.error('加载订单列表失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'none'
        })
      } finally {
        this.isLoading = false
      }
    },
    
    // 选择状态
    selectStatus(status) {
      this.selectedStatus = status
      this.loadOrders(true)
    },
    
    // 加载更多
    loadMore() {
      if (this.hasMore && !this.isLoading) {
        this.pagination.page++
        this.loadOrders()
      }
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
    
    // 是否显示操作按钮
    showActions(status) {
      return ['pending', 'confirmed', 'picked_up', 'evaluated', 'completed'].includes(status)
    },
    
    // 取消订单
    async cancelOrder(orderId) {
      uni.showModal({
        title: '确认取消',
        content: '确定要取消这个回收订单吗？',
        success: async (res) => {
          if (res.confirm) {
            try {
              await this.$http.put(`/api/v1/orders/${orderId}/cancel`)
              
              uni.showToast({
                title: '订单已取消',
                icon: 'success'
              })
              
              // 刷新列表
              this.loadOrders(true)
              
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
    contactService(orderNo) {
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
    reorder(deviceId) {
      uni.navigateTo({
        url: `/pages/order/create?deviceId=${deviceId}`
      })
    },
    
    // 跳转到订单详情
    navigateToDetail(orderId) {
      uni.navigateTo({
        url: `/pages/order/detail?id=${orderId}`
      })
    },
    
    // 跳转到设备列表
    navigateToDevices() {
      uni.switchTab({
        url: '/pages/devices/devices'
      })
    }
  },
  
  onReachBottom() {
    this.loadMore()
  }
}
</script>

<style lang="scss" scoped>
.order-list {
  background: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 40rpx;
}

.status-filter {
  margin: 20rpx;
  padding: 0;
  
  .filter-scroll {
    white-space: nowrap;
    
    .filter-tabs {
      display: inline-flex;
      padding: 20rpx;
      gap: 20rpx;
      
      .filter-tab {
        padding: 16rpx 32rpx;
        border-radius: 40rpx;
        background: #f8f9fa;
        color: #666;
        font-size: 28rpx;
        white-space: nowrap;
        border: 1rpx solid #e0e0e0;
        
        &.active {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: #fff;
          border-color: transparent;
        }
      }
    }
  }
}

.orders-container {
  margin: 0 20rpx;
  
  .order-card {
    background: #fff;
    border-radius: 16rpx;
    padding: 30rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
    
    .order-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 24rpx;
      padding-bottom: 20rpx;
      border-bottom: 1rpx solid #f0f0f0;
      
      .order-info {
        flex: 1;
        
        .order-no {
          display: block;
          font-size: 28rpx;
          font-weight: bold;
          color: #333;
          margin-bottom: 8rpx;
        }
        
        .order-time {
          font-size: 24rpx;
          color: #999;
        }
      }
      
      .order-status {
        font-size: 26rpx;
        font-weight: bold;
        padding: 8rpx 16rpx;
        border-radius: 20rpx;
        background: rgba(102, 126, 234, 0.1);
      }
    }
    
    .device-info {
      display: flex;
      margin-bottom: 24rpx;
      
      .device-image {
        width: 100rpx;
        height: 100rpx;
        border-radius: 12rpx;
        margin-right: 20rpx;
        object-fit: cover;
      }
      
      .device-details {
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
        
        .device-condition {
          .condition-tag {
            font-size: 24rpx;
            padding: 4rpx 12rpx;
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
    
    .price-info {
      display: flex;
      justify-content: space-between;
      margin-bottom: 20rpx;
      
      .price-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        
        .price-label {
          font-size: 24rpx;
          color: #999;
          margin-bottom: 8rpx;
        }
        
        .price-value {
          font-size: 28rpx;
          font-weight: bold;
          color: #333;
          
          &.final {
            color: #667eea;
            font-size: 32rpx;
          }
        }
      }
    }
    
    .contact-info {
      display: flex;
      justify-content: space-between;
      margin-bottom: 20rpx;
      padding: 16rpx;
      background: #f8f9fa;
      border-radius: 12rpx;
      
      .contact-item {
        display: flex;
        flex-direction: column;
        
        .contact-label {
          font-size: 24rpx;
          color: #999;
          margin-bottom: 4rpx;
        }
        
        .contact-value {
          font-size: 26rpx;
          color: #333;
        }
      }
    }
    
    .pickup-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16rpx;
      background: #e3f2fd;
      border-radius: 12rpx;
      margin-bottom: 20rpx;
      
      .pickup-label {
        font-size: 26rpx;
        color: #2196f3;
        font-weight: bold;
      }
      
      .pickup-time {
        font-size: 26rpx;
        color: #2196f3;
      }
    }
    
    .order-actions {
      display: flex;
      gap: 20rpx;
      
      .action-btn {
        flex: 1;
        height: 60rpx;
        border-radius: 30rpx;
        font-size: 26rpx;
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
  }
}

.load-more {
  margin: 40rpx 20rpx;
  
  .load-btn {
    height: 80rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f8f9fa;
    border-radius: 40rpx;
    color: #666;
    font-size: 28rpx;
    
    &.loading {
      color: #999;
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 40rpx;
  
  .empty-image {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 30rpx;
    opacity: 0.5;
  }
  
  .empty-text {
    font-size: 28rpx;
    color: #999;
    margin-bottom: 40rpx;
  }
  
  .empty-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    border-radius: 40rpx;
    padding: 20rpx 40rpx;
    font-size: 28rpx;
    border: none;
  }
}
</style>
