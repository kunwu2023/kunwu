<template>
    <el-container>
      <el-container>
        <!-- 侧边栏 -->
        <el-aside class="sidebar">
          <el-menu
              default-active="1"
              class="el-menu-vertical-demo"
          >
            <el-menu-item index="1" @click="gotoPage('/add_task')">
              <span><el-icon><DocumentAdd /></el-icon>扫描任务</span>
            </el-menu-item>
            <el-menu-item index="2" @click="gotoPage('/see_task'); store.state.badge.taskList = 0">
              <span v-if="store.state.badge.taskList === 0">
                  <el-icon><Document /></el-icon>任务列表
              </span>
              <span v-if="store.state.badge.taskList !== 0">
                <el-badge :value="store.state.badge.taskList" :max="9" class="item">
                  <el-icon><Document /></el-icon>任务列表
                </el-badge>
              </span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <!-- 主要容器 -->
        <el-main class="main" height="250">
          <router-view></router-view>
        </el-main>
      </el-container>

      <!-- 底栏 -->
      <el-footer class="footer">© 2023 昆吾 V0.1.1</el-footer>
    </el-container>
</template>

<script>
import {useRouter} from "vue-router";
import {DocumentAdd, Document} from '@element-plus/icons-vue'
import {useStore} from "vuex";
export default {
  name: "index",
  components:{DocumentAdd, Document},
  setup(){
    const store = useStore()
    const router = useRouter()
    // 跳转
    const gotoPage = (routerName) => {
      router.push(routerName)
    }
    return{
      gotoPage,
      store,
    }
  }
}
</script>

<style scoped>
.sidebar {
  width: 150px;
}

.main{
  padding-bottom: 50px;
}

/* 手动配置全局样式 */
html,
body,
.app_container,
.el-container{
  padding: 0;
  margin: 0;
  height: 100vh;
}

.footer {
  position: absolute; /* 底栏设置为绝对定位 */
  bottom: 0;
  left: 0;
  height: 30px;
  right: 25px;
  background-color: white;
  z-index: 999;
}

.item {
  margin-top: 10px;
  margin-right: 40px;
}

</style>