import { createSSRApp } from 'vue'
import App from './App.vue'
import store from './store'

export function createApp() {
  const app = createSSRApp(App)
  
  // 配置全局属性
  app.config.globalProperties.$baseUrl = 'http://localhost:8080'
  
  // 全局API封装
  app.config.globalProperties.$http = {
    // GET请求
    get(url, data = {}) {
      return new Promise((resolve, reject) => {
        uni.request({
          url: app.config.globalProperties.$baseUrl + url,
          method: 'GET',
          data,
          header: this.getHeader(),
          success: (res) => {
            this.handleResponse(res, resolve, reject)
          },
          fail: reject
        })
      })
    },
    
    // POST请求
    post(url, data = {}) {
      return new Promise((resolve, reject) => {
        uni.request({
          url: app.config.globalProperties.$baseUrl + url,
          method: 'POST',
          data,
          header: this.getHeader(),
          success: (res) => {
            this.handleResponse(res, resolve, reject)
          },
          fail: reject
        })
      })
    },
    
    // PUT请求
    put(url, data = {}) {
      return new Promise((resolve, reject) => {
        uni.request({
          url: app.config.globalProperties.$baseUrl + url,
          method: 'PUT',
          data,
          header: this.getHeader(),
          success: (res) => {
            this.handleResponse(res, resolve, reject)
          },
          fail: reject
        })
      })
    },
    
    // DELETE请求
    delete(url, data = {}) {
      return new Promise((resolve, reject) => {
        uni.request({
          url: app.config.globalProperties.$baseUrl + url,
          method: 'DELETE',
          data,
          header: this.getHeader(),
          success: (res) => {
            this.handleResponse(res, resolve, reject)
          },
          fail: reject
        })
      })
    },
    
    // 获取请求头
    getHeader() {
      const token = uni.getStorageSync('token')
      return {
        'Content-Type': 'application/json',
        'Authorization': token ? 'Bearer ' + token : ''
      }
    },
    
    // 处理响应
    handleResponse(res, resolve, reject) {
      if (res.statusCode === 200) {
        resolve(res.data)
      } else if (res.statusCode === 401) {
        // 未授权，跳转到登录页
        uni.removeStorageSync('token')
        uni.removeStorageSync('userInfo')
        uni.reLaunch({
          url: '/pages/user/login'
        })
        reject(res.data)
      } else {
        // 显示错误信息
        uni.showToast({
          title: res.data.error || '请求失败',
          icon: 'none'
        })
        reject(res.data)
      }
    }
  }
  
  // 全局工具函数
  app.config.globalProperties.$utils = {
    // 格式化时间
    formatTime(time) {
      if (!time) return ''
      const date = new Date(time)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hour = String(date.getHours()).padStart(2, '0')
      const minute = String(date.getMinutes()).padStart(2, '0')
      return `${year}-${month}-${day} ${hour}:${minute}`
    },
    
    // 格式化价格
    formatPrice(price) {
      if (!price && price !== 0) return '面议'
      return '¥' + parseFloat(price).toFixed(2)
    },
    
    // 获取设备状态文本
    getConditionText(condition) {
      const map = {
        'excellent': '95新以上',
        'good': '85-95新',
        'fair': '70-85新',
        'poor': '70新以下'
      }
      return map[condition] || condition
    },
    
    // 获取订单状态文本
    getOrderStatusText(status) {
      const map = {
        'pending': '待处理',
        'confirmed': '已确认',
        'picked_up': '已上门',
        'evaluated': '已评估',
        'completed': '已完成',
        'cancelled': '已取消'
      }
      return map[status] || status
    },
    
    // 获取订单状态颜色
    getOrderStatusColor(status) {
      const map = {
        'pending': '#f39c12',
        'confirmed': '#3498db',
        'picked_up': '#9b59b6',
        'evaluated': '#e67e22',
        'completed': '#27ae60',
        'cancelled': '#95a5a6'
      }
      return map[status] || '#999'
    },
    
    // 检查登录状态
    checkLogin() {
      const token = uni.getStorageSync('token')
      if (!token) {
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
  
  app.use(store)
  return {
    app
  }
}
