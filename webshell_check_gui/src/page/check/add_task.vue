<template>
  <el-tabs v-model="activeName" class="demo-tabs" tabPosition="top" @tab-click="gotoPage">
    <el-tab-pane label="快速扫描" name="/quick_check">
    </el-tab-pane>
    <el-tab-pane label="本地扫描" name="/local_scan">
    </el-tab-pane>
    <el-tab-pane label="远程扫描(ssh)" name="/remote_scan">
    </el-tab-pane>
    <el-tab-pane label="高级选项" name="/advanced">
    </el-tab-pane>
  </el-tabs>
  <router-view></router-view>
<!--  <el-button @click="testIPC">testIPC</el-button>-->
</template>

<script>
import {ref} from "vue";
import { ipcRenderer } from 'electron'
import {useRouter} from "vue-router";
export default {
  name: "add_task",
  setup(){
    const activeName = ref('/quick_check')
    const router = useRouter()
    const gotoPage = (TabsPaneContext) => {
      router.push(TabsPaneContext.props.name)
    }
    const testIPC = async () => {
      await ipcRenderer.invoke('go-request', {
        function: 'hello',
        data: 'caozhe'
      });
    }
    return{
      testIPC,gotoPage,
      activeName
    }
  }
}
</script>

<style scoped>

.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

</style>