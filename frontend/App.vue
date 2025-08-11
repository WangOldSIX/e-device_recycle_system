<template>
	<!-- 应用根组件，所有页面的共同外壳 -->
</template>

<script>
	export default {
		onLaunch: function() {
			console.log('App Launch')
			// 应用启动时的逻辑
			this.checkLogin()
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		},
		methods: {
			// 检查登录状态
			checkLogin() {
				const token = uni.getStorageSync('token')
				if (token) {
					// 验证token有效性
					uni.request({
						url: this.$baseUrl + '/api/v1/user/profile',
						method: 'GET',
						header: {
							'Authorization': 'Bearer ' + token
						},
						success: (res) => {
							if (res.statusCode === 200) {
								// token有效，保存用户信息
								uni.setStorageSync('userInfo', res.data.user)
							} else {
								// token无效，清除本地存储
								uni.removeStorageSync('token')
								uni.removeStorageSync('userInfo')
							}
						},
						fail: () => {
							// 网络错误，清除本地存储
							uni.removeStorageSync('token')
							uni.removeStorageSync('userInfo')
						}
					})
				}
			}
		},
		globalData: {
			baseUrl: 'http://localhost:8080'
		}
	}
</script>

<style lang="scss">
	/* 注意要写在第一行，同时给style标签加入lang="scss"属性 */
	@import "uview-ui/index.scss";
	
	/* 全局样式 */
	page {
		background-color: #f5f5f5;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
	}
	
	/* 通用样式 */
	.container {
		padding: 20rpx;
	}
	
	.card {
		background: #fff;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
	}
	
	.btn-primary {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: #fff;
		border-radius: 50rpx;
		padding: 24rpx 48rpx;
		font-size: 32rpx;
		border: none;
		box-shadow: 0 8rpx 20rpx rgba(102, 126, 234, 0.3);
	}
	
	.btn-secondary {
		background: #f8f9fa;
		color: #6c757d;
		border-radius: 50rpx;
		padding: 24rpx 48rpx;
		font-size: 32rpx;
		border: 1rpx solid #dee2e6;
	}
	
	.text-primary {
		color: #667eea;
	}
	
	.text-secondary {
		color: #6c757d;
	}
	
	.text-muted {
		color: #999;
		font-size: 26rpx;
	}
	
	.flex {
		display: flex;
	}
	
	.flex-center {
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.flex-between {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	
	.flex-column {
		display: flex;
		flex-direction: column;
	}
	
	.mt-20 {
		margin-top: 20rpx;
	}
	
	.mt-40 {
		margin-top: 40rpx;
	}
	
	.mb-20 {
		margin-bottom: 20rpx;
	}
	
	.mb-40 {
		margin-bottom: 40rpx;
	}
	
	.ml-20 {
		margin-left: 20rpx;
	}
	
	.mr-20 {
		margin-right: 20rpx;
	}
</style>
