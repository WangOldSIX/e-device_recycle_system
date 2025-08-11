<template>
  <view class="devices">
    <!-- 搜索和筛选 -->
    <view class="search-filter card">
      <view class="search-box">
        <input 
          v-model="searchKeyword" 
          class="search-input" 
          placeholder="搜索设备型号、品牌" 
          @input="onSearch"
        />
        <text class="search-icon">🔍</text>
      </view>
      
      <view class="filter-tabs">
        <view 
          class="filter-tab" 
          :class="{ active: selectedCategory === '' }" 
          @click="selectCategory('')"
        >
          全部
        </view>
        <view 
          class="filter-tab" 
          :class="{ active: selectedCategory === category.value }" 
          v-for="category in categories" 
          :key="category.value"
          @click="selectCategory(category.value)"
        >
          {{ category.label }}
        </view>
      </view>
    </view>
    
    <!-- 设备列表 -->
    <view class="device-list">
      <view 
        class="device-card" 
        v-for="device in devices" 
        :key="device.id"
        @click="navigateToDetail(device.id)"
      >
        <image 
          :src="device.images || '/static/images/device-placeholder.png'" 
          class="device-image"
        ></image>
        <view class="device-info">
          <view class="device-name">{{ device.name }}</view>
          <view class="device-specs">
            <text class="spec-item">{{ device.brand }}</text>
            <text class="spec-item" v-if="device.cpu">{{ device.cpu }}</text>
            <text class="spec-item" v-if="device.memory">{{ device.memory }}</text>
          </view>
          <view class="device-condition">
            <text class="condition-tag" :class="device.condition">
              {{ $utils.getConditionText(device.condition) }}
            </text>
            <text class="year-info" v-if="device.year_bought">
              {{ device.year_bought }}年
            </text>
          </view>
          <view class="device-price">
            <text class="price-label">参考回收价</text>
            <text class="price-value">{{ $utils.formatPrice(device.base_price) }}</text>
          </view>
        </view>
        <view class="action-btn">
          <text class="btn-text">立即回收</text>
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
    <view class="empty-state" v-if="devices.length === 0 && !isLoading">
      <image src="/static/images/empty.png" class="empty-image"></image>
      <text class="empty-text">暂无设备数据</text>
    </view>
  </view>
</template>

<script>
import { useDeviceStore } from '@/store'

export default {
  name: 'Devices',
  data() {
    return {
      searchKeyword: '',
      selectedCategory: '',
      devices: [],
      categories: [
        { value: 'laptop', label: '笔记本' },
        { value: 'desktop', label: '台式机' },
        { value: 'tablet', label: '平板' },
        { value: 'phone', label: '手机' }
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
    // 如果有传入分类参数，则设置选中的分类
    if (options.category) {
      this.selectedCategory = options.category
    }
    this.loadDevices(true)
  },
  
  methods: {
    // 加载设备列表
    async loadDevices(reset = false) {
      if (this.isLoading) return
      
      this.isLoading = true
      
      if (reset) {
        this.pagination.page = 1
        this.devices = []
        this.hasMore = true
      }
      
      try {
        const params = {
          page: this.pagination.page,
          page_size: this.pagination.pageSize
        }
        
        if (this.selectedCategory) {
          params.category = this.selectedCategory
        }
        
        if (this.searchKeyword) {
          params.brand = this.searchKeyword
        }
        
        const res = await this.$http.get('/api/v1/devices', params)
        
        if (reset) {
          this.devices = res.devices || []
        } else {
          this.devices = [...this.devices, ...(res.devices || [])]
        }
        
        this.pagination.total = res.pagination?.total || 0
        this.hasMore = this.devices.length < this.pagination.total
        
        if (!reset) {
          this.pagination.page++
        }
        
      } catch (error) {
        console.error('加载设备列表失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'none'
        })
      } finally {
        this.isLoading = false
      }
    },
    
    // 搜索
    onSearch() {
      // 防抖处理
      clearTimeout(this.searchTimer)
      this.searchTimer = setTimeout(() => {
        this.loadDevices(true)
      }, 500)
    },
    
    // 选择分类
    selectCategory(category) {
      this.selectedCategory = category
      this.loadDevices(true)
    },
    
    // 加载更多
    loadMore() {
      if (this.hasMore && !this.isLoading) {
        this.pagination.page++
        this.loadDevices()
      }
    },
    
    // 跳转到设备详情
    navigateToDetail(deviceId) {
      uni.navigateTo({
        url: `/pages/devices/detail?id=${deviceId}`
      })
    }
  },
  
  onReachBottom() {
    this.loadMore()
  }
}
</script>

<style lang="scss" scoped>
.devices {
  padding-bottom: 20rpx;
}

.search-filter {
  margin: 20rpx;
  
  .search-box {
    position: relative;
    margin-bottom: 30rpx;
    
    .search-input {
      width: 100%;
      height: 80rpx;
      padding: 0 50rpx 0 20rpx;
      border: 1rpx solid #e0e0e0;
      border-radius: 40rpx;
      font-size: 28rpx;
      background: #f8f9fa;
    }
    
    .search-icon {
      position: absolute;
      right: 20rpx;
      top: 50%;
      transform: translateY(-50%);
      font-size: 28rpx;
      color: #999;
    }
  }
  
  .filter-tabs {
    display: flex;
    gap: 20rpx;
    
    .filter-tab {
      padding: 16rpx 32rpx;
      border-radius: 40rpx;
      background: #f8f9fa;
      color: #666;
      font-size: 28rpx;
      border: 1rpx solid #e0e0e0;
      
      &.active {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #fff;
        border-color: transparent;
      }
    }
  }
}

.device-list {
  margin: 0 20rpx;
  
  .device-card {
    display: flex;
    background: #fff;
    border-radius: 16rpx;
    padding: 30rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
    
    .device-image {
      width: 160rpx;
      height: 160rpx;
      border-radius: 12rpx;
      margin-right: 24rpx;
      object-fit: cover;
    }
    
    .device-info {
      flex: 1;
      
      .device-name {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 12rpx;
        line-height: 1.4;
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
        display: flex;
        align-items: center;
        gap: 12rpx;
        margin-bottom: 12rpx;
        
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
        
        .year-info {
          font-size: 24rpx;
          color: #999;
        }
      }
      
      .device-price {
        display: flex;
        align-items: baseline;
        gap: 8rpx;
        
        .price-label {
          font-size: 24rpx;
          color: #999;
        }
        
        .price-value {
          font-size: 32rpx;
          font-weight: bold;
          color: #667eea;
        }
      }
    }
    
    .action-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 120rpx;
      height: 60rpx;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 30rpx;
      margin-top: 50rpx;
      
      .btn-text {
        font-size: 26rpx;
        color: #fff;
        font-weight: bold;
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
  padding: 120rpx 0;
  
  .empty-image {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 30rpx;
    opacity: 0.5;
  }
  
  .empty-text {
    font-size: 28rpx;
    color: #999;
  }
}
</style>
