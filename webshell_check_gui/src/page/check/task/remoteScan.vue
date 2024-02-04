<template>
  <el-row>
    <el-col :span="24" class="grid-cell">
      <el-text>ssh用户名：</el-text>
      <el-input
          v-model="store.state.scan_configuration.remoteScan.userName"
          placeholder="请输入用户名"
      />
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-text>ssh密码：</el-text>
      <el-input
          v-model="store.state.scan_configuration.remoteScan.passWord"
          type="password"
          placeholder="请输入密码"
          show-password
      />
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-text>远程IP：</el-text>
      <el-input
          v-model="store.state.scan_configuration.remoteScan.serverIp"
          placeholder="请输入服务器ip"
      />
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-text>ssh端口号：</el-text>
      <el-input
          v-model="store.state.scan_configuration.remoteScan.sshPort"
          placeholder="请输入ssh端口号"
      />
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-text>扫描路径：</el-text>
      <el-input
          v-model="store.state.scan_configuration.remoteScan.scanPath"
          placeholder="请输入需要扫描的绝对路径"
      />
    </el-col>
    <el-col :offset="19" :span="4" class="grid-cell">
      <el-button type="primary" @click="addRemoteScan" :loading="loadingFlag" plain>
        创建任务<el-icon class="el-icon--right"><Pointer /></el-icon>
      </el-button>
    </el-col>
  </el-row>
</template>

<script>
import { Pointer } from "@element-plus/icons-vue";
import { useStore } from "vuex";
import { ipcRenderer } from "electron";
import {ElMessage} from "element-plus";
import {ref} from "vue";

export default {
  name: "remoteScan",
  components: {
    Pointer
  },
  setup(){
    const store = useStore()
    const loadingFlag = ref(false)
    const addRemoteScan = async () => {
      loadingFlag.value = true
      // TODO 创建远程扫描任务
      const scanConfig = {
        cloudScanPath: store.state.scan_configuration.remoteScan.scanPath,
        model: store.state.scan_configuration.remoteScan.model,
        userName: store.state.scan_configuration.remoteScan.userName,
        passWord: store.state.scan_configuration.remoteScan.passWord,
        serverIp: store.state.scan_configuration.remoteScan.serverIp,
        sshPort: store.state.scan_configuration.remoteScan.sshPort,
      }
      let cloudFlag
      if (store.state.scan_configuration.advanced.cloud_scan){
        cloudFlag = 1
      }else {
        cloudFlag = 2
      }
      ElMessage({
        message: "远程扫描任务创建中",
        type: 'success',
      })
      store.state.badge.taskList = store.state.badge.taskList + 1

      const response = await ipcRenderer.invoke('go-request', {
        function: 'addCloudTask',
        data: scanConfig,
        cloudFlag:cloudFlag,
      });
      if ("error" in response){
        ElMessage({
          message: response["error"],
          type: 'error',
        })
        store.state.badge.taskList = store.state.badge.taskList - 1
      }else{
        ElMessage({
          message: "远程扫描任务完成",
          type: 'success',
        })
      }
      loadingFlag.value = false
    }
    return{
      store,loadingFlag,
      addRemoteScan,
    }
  }
}
</script>

<style scoped>

</style>