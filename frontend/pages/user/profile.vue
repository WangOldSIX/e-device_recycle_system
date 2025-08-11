<template>
  <view class="profile">
    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-avatar" @click="chooseAvatar">
        <image :src="userAvatar" class="avatar-image"></image>
        <view class="avatar-edit">📷</view>
      </view>
      <view class="user-info">
        <view class="user-name">{{ userStore.realName }}</view>
        <view class="user-username">@{{ userStore.username }}</view>
        <view class="user-status">
          <text class="status-dot" :class="userStore.userInfo?.status"></text>
          <text class="status-text">{{ getStatusText(userStore.userInfo?.status) }}</text>
        </view>
      </view>
      <view class="user-edit" @click="editProfile">
        <text class="edit-icon">✏️</text>
      </view>
    </view>
    
    <!-- 统计信息 -->
    <view class="stats-card card">
      <view class="stats-grid">
        <view class="stat-item" @click="navigateToOrders()">
          <view class="stat-number">{{ userStats.totalOrders }}</view>
          <view class="stat-label">总订单</view>
        </view>
        <view class="stat-item" @click="navigateToOrders('completed')">
          <view class="stat-number">{{ userStats.completedOrders }}</view>
          <view class="stat-label">已完成</view>
        </view>
        <view class="stat-item">
          <view class="stat-number">{{ $utils.formatPrice(userStats.totalEarnings) }}</view>
          <view class="stat-label">总收益</view>
        </view>
      </view>
    </view>
    
    <!-- 功能菜单 -->
    <view class="menu-section">
      <view class="menu-group card">
        <view class="menu-title">订单管理</view>
        <view class="menu-item" @click="navigateToOrders()">
          <view class="menu-icon">📋</view>
          <view class="menu-text">我的订单</view>
          <view class="menu-badge" v-if="userStats.pendingOrders > 0">
            {{ userStats.pendingOrders }}
          </view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="navigateToOrders('pending')">
          <view class="menu-icon">⏳</view>
          <view class="menu-text">待处理订单</view>
          <view class="menu-badge" v-if="userStats.pendingOrders > 0">
            {{ userStats.pendingOrders }}
          </view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="navigateToOrders('completed')">
          <view class="menu-icon">✅</view>
          <view class="menu-text">已完成订单</view>
          <view class="menu-arrow">></view>
        </view>
      </view>
      
      <view class="menu-group card">
        <view class="menu-title">账户设置</view>
        <view class="menu-item" @click="editProfile">
          <view class="menu-icon">👤</view>
          <view class="menu-text">个人信息</view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="changePassword">
          <view class="menu-icon">🔒</view>
          <view class="menu-text">修改密码</view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="addressManagement">
          <view class="menu-icon">📍</view>
          <view class="menu-text">地址管理</view>
          <view class="menu-arrow">></view>
        </view>
      </view>
      
      <view class="menu-group card">
        <view class="menu-title">帮助与支持</view>
        <view class="menu-item" @click="contactService">
          <view class="menu-icon">📞</view>
          <view class="menu-text">联系客服</view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="showFeedback">
          <view class="menu-icon">💬</view>
          <view class="menu-text">意见反馈</view>
          <view class="menu-arrow">></view>
        </view>
        <view class="menu-item" @click="showAbout">
          <view class="menu-icon">ℹ️</view>
          <view class="menu-text">关于我们</view>
          <view class="menu-arrow">></view>
        </view>
      </view>
      
      <view class="menu-group card">
        <view class="menu-item logout" @click="handleLogout">
          <view class="menu-icon">🚪</view>
          <view class="menu-text">退出登录</view>
          <view class="menu-arrow">></view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { useUserStore } from '@/store'

export default {
  name: 'Profile',
  data() {
    return {
      userStats: {
        totalOrders: 0,
        completedOrders: 0,
        pendingOrders: 0,
        totalEarnings: 0
      }
    }
  },
  
  computed: {
    userStore() {
      return useUserStore()
    },
    
    userAvatar() {
      return this.userStore.avatar || '/static/images/default-avatar.png'
    }
  },
  
  onLoad() {
    this.loadUserStats()
  },
  
  onShow() {
    this.loadUserStats()
  },
  
  methods: {
    // 加载用户统计数据
    async loadUserStats() {
      try {
        // 获取用户订单统计
        const res = await this.$http.get('/api/v1/orders', {
          page: 1,
          page_size: 1000 // 获取所有订单用于统计
        })
        
        const orders = res.orders || []
        
        this.userStats.totalOrders = orders.length
        this.userStats.completedOrders = orders.filter(order => order.status === 'completed').length
        this.userStats.pendingOrders = orders.filter(order => 
          ['pending', 'confirmed', 'picked_up', 'evaluated'].includes(order.status)
        ).length
        
        // 计算总收益
        this.userStats.totalEarnings = orders
          .filter(order => order.status === 'completed' && order.final_price)
          .reduce((total, order) => total + (order.final_price || 0), 0)
          
      } catch (error) {
        console.error('加载用户统计失败:', error)
      }
    },
    
    // 获取状态文本
    getStatusText(status) {
      const statusMap = {
        'active': '正常',
        'banned': '已禁用'
      }
      return statusMap[status] || '未知'
    },
    
    // 选择头像
    chooseAvatar() {
      uni.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          // 这里应该上传头像到服务器
          // 暂时只更新本地显示
          this.userStore.updateUserInfo({
            avatar: res.tempFilePaths[0]
          })
          
          uni.showToast({
            title: '头像更新成功',
            icon: 'success'
          })
        }
      })
    },
    
    // 编辑个人信息
    editProfile() {
      uni.showModal({
        title: '功能开发中',
        content: '个人信息编辑功能正在开发中，敬请期待',
        showCancel: false
      })
    },
    
    // 修改密码
    changePassword() {
      uni.showModal({
        title: '功能开发中',
        content: '修改密码功能正在开发中，敬请期待',
        showCancel: false
      })
    },
    
    // 地址管理
    addressManagement() {
      uni.showModal({
        title: '功能开发中',
        content: '地址管理功能正在开发中，敬请期待',
        showCancel: false
      })
    },
    
    // 跳转到订单页面
    navigateToOrders(status = '') {
      let url = '/pages/order/list'
      if (status) {
        url += '?status=' + status
      }
      uni.navigateTo({
        url
      })
    },
    
    // 联系客服
    contactService() {
      uni.showActionSheet({
        itemList: ['电话客服', '在线客服', '邮箱联系'],
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
            case 2:
              uni.setClipboardData({
                data: 'service@example.com',
                success: () => {
                  uni.showToast({
                    title: '邮箱已复制',
                    icon: 'success'
                  })
                }
              })
              break
          }
        }
      })
    },
    
    // 意见反馈
    showFeedback() {
      uni.showModal({
        title: '意见反馈',
        content: '您可以通过客服电话或邮箱向我们反馈意见和建议',
        showCancel: false
      })
    },
    
    // 关于我们
    showAbout() {
      uni.showModal({
        title: '关于我们',
        content: '二手电脑回收平台 v1.0.0\n\n致力于为用户提供专业、便捷、安全的二手设备回收服务，让闲置设备重新焕发价值。',
        showCancel: false
      })
    },
    
    // 退出登录
    handleLogout() {
      uni.showModal({
        title: '确认退出',
        content: '确定要退出当前账号吗？',
        success: (res) => {
          if (res.confirm) {
            this.userStore.logout()
          }
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.profile {
  background: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 40rpx;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
  display: flex;
  align-items: center;
  color: #fff;
  
  .user-avatar {
    position: relative;
    margin-right: 30rpx;
    
    .avatar-image {
      width: 120rpx;
      height: 120rpx;
      border-radius: 60rpx;
      border: 4rpx solid rgba(255, 255, 255, 0.3);
    }
    
    .avatar-edit {
      position: absolute;
      bottom: 0;
      right: 0;
      width: 40rpx;
      height: 40rpx;
      background: rgba(0, 0, 0, 0.5);
      border-radius: 20rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20rpx;
    }
  }
  
  .user-info {
    flex: 1;
    
    .user-name {
      font-size: 36rpx;
      font-weight: bold;
      margin-bottom: 8rpx;
    }
    
    .user-username {
      font-size: 26rpx;
      opacity: 0.8;
      margin-bottom: 12rpx;
    }
    
    .user-status {
      display: flex;
      align-items: center;
      
      .status-dot {
        width: 12rpx;
        height: 12rpx;
        border-radius: 6rpx;
        margin-right: 8rpx;
        
        &.active {
          background: #27ae60;
        }
        
        &.banned {
          background: #f44336;
        }
      }
      
      .status-text {
        font-size: 24rpx;
        opacity: 0.8;
      }
    }
  }
  
  .user-edit {
    padding: 16rpx;
    
    .edit-icon {
      font-size: 32rpx;
    }
  }
}

.stats-card {
  margin: 20rpx;
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rpx;
    background: #f0f0f0;
    border-radius: 12rpx;
    overflow: hidden;
    
    .stat-item {
      background: #fff;
      padding: 40rpx 20rpx;
      text-align: center;
      
      .stat-number {
        font-size: 36rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 8rpx;
      }
      
      .stat-label {
        font-size: 26rpx;
        color: #666;
      }
    }
  }
}

.menu-section {
  margin: 0 20rpx;
  
  .menu-group {
    margin-bottom: 20rpx;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .menu-title {
      font-size: 28rpx;
      font-weight: bold;
      color: #333;
      padding: 20rpx 0 16rpx;
      border-bottom: 1rpx solid #f0f0f0;
      margin-bottom: 8rpx;
    }
    
    .menu-item {
      display: flex;
      align-items: center;
      padding: 24rpx 0;
      border-bottom: 1rpx solid #f8f9fa;
      
      &:last-child {
        border-bottom: none;
      }
      
      &.logout {
        color: #f44336;
        
        .menu-icon,
        .menu-text {
          color: #f44336;
        }
      }
      
      .menu-icon {
        font-size: 32rpx;
        margin-right: 20rpx;
        width: 40rpx;
        text-align: center;
      }
      
      .menu-text {
        flex: 1;
        font-size: 30rpx;
        color: #333;
      }
      
      .menu-badge {
        background: #f44336;
        color: #fff;
        font-size: 20rpx;
        padding: 4rpx 12rpx;
        border-radius: 12rpx;
        margin-right: 16rpx;
        min-width: 32rpx;
        text-align: center;
      }
      
      .menu-arrow {
        font-size: 24rpx;
        color: #999;
      }
    }
  }
}
</style>
