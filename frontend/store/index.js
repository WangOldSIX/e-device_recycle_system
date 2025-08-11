import { createPinia, defineStore } from 'pinia'

// 创建pinia实例
const pinia = createPinia()

// 用户状态管理
export const useUserStore = defineStore('user', {
  state: () => ({
    token: uni.getStorageSync('token') || '',
    userInfo: uni.getStorageSync('userInfo') || null,
    isLogin: false
  }),
  
  getters: {
    // 获取用户头像
    avatar: (state) => {
      return state.userInfo?.avatar || '/static/images/default-avatar.png'
    },
    
    // 获取用户名
    username: (state) => {
      return state.userInfo?.username || '未登录'
    },
    
    // 获取真实姓名
    realName: (state) => {
      return state.userInfo?.real_name || state.userInfo?.username || '未设置'
    }
  },
  
  actions: {
    // 登录
    login(token, userInfo) {
      this.token = token
      this.userInfo = userInfo
      this.isLogin = true
      
      // 保存到本地存储
      uni.setStorageSync('token', token)
      uni.setStorageSync('userInfo', userInfo)
    },
    
    // 退出登录
    logout() {
      this.token = ''
      this.userInfo = null
      this.isLogin = false
      
      // 清除本地存储
      uni.removeStorageSync('token')
      uni.removeStorageSync('userInfo')
      
      // 跳转到登录页
      uni.reLaunch({
        url: '/pages/user/login'
      })
    },
    
    // 更新用户信息
    updateUserInfo(userInfo) {
      this.userInfo = { ...this.userInfo, ...userInfo }
      uni.setStorageSync('userInfo', this.userInfo)
    },
    
    // 检查登录状态
    checkLogin() {
      if (!this.token) {
        uni.showModal({
          title: '提示',
          content: '请先登录',
          success: (res) => {
            if (res.confirm) {
              uni.navigateTo({
                url: '/pages/user/login'
              })
            }
          }
        })
        return false
      }
      return true
    }
  }
})

// 设备状态管理
export const useDeviceStore = defineStore('device', {
  state: () => ({
    devices: [],
    currentDevice: null,
    categories: [
      { value: 'laptop', label: '笔记本电脑' },
      { value: 'desktop', label: '台式电脑' },
      { value: 'tablet', label: '平板电脑' },
      { value: 'phone', label: '手机' }
    ],
    conditions: [
      { value: 'excellent', label: '95新以上' },
      { value: 'good', label: '85-95新' },
      { value: 'fair', label: '70-85新' },
      { value: 'poor', label: '70新以下' }
    ]
  }),
  
  actions: {
    // 设置设备列表
    setDevices(devices) {
      this.devices = devices
    },
    
    // 设置当前设备
    setCurrentDevice(device) {
      this.currentDevice = device
    },
    
    // 添加设备
    addDevice(device) {
      this.devices.unshift(device)
    }
  }
})

// 订单状态管理
export const useOrderStore = defineStore('order', {
  state: () => ({
    orders: [],
    currentOrder: null,
    orderStatuses: [
      { value: 'pending', label: '待处理', color: '#f39c12' },
      { value: 'confirmed', label: '已确认', color: '#3498db' },
      { value: 'picked_up', label: '已上门', color: '#9b59b6' },
      { value: 'evaluated', label: '已评估', color: '#e67e22' },
      { value: 'completed', label: '已完成', color: '#27ae60' },
      { value: 'cancelled', label: '已取消', color: '#95a5a6' }
    ]
  }),
  
  getters: {
    // 获取状态文本
    getStatusText: (state) => {
      return (status) => {
        const statusObj = state.orderStatuses.find(s => s.value === status)
        return statusObj ? statusObj.label : status
      }
    },
    
    // 获取状态颜色
    getStatusColor: (state) => {
      return (status) => {
        const statusObj = state.orderStatuses.find(s => s.value === status)
        return statusObj ? statusObj.color : '#999'
      }
    }
  },
  
  actions: {
    // 设置订单列表
    setOrders(orders) {
      this.orders = orders
    },
    
    // 设置当前订单
    setCurrentOrder(order) {
      this.currentOrder = order
    },
    
    // 添加订单
    addOrder(order) {
      this.orders.unshift(order)
    },
    
    // 更新订单状态
    updateOrderStatus(orderId, status) {
      const order = this.orders.find(o => o.id === orderId)
      if (order) {
        order.status = status
      }
      
      if (this.currentOrder && this.currentOrder.id === orderId) {
        this.currentOrder.status = status
      }
    }
  }
})

export default pinia
