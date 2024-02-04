<template>
  <el-row>
    <el-col :span="24" class="grid-cell">
      <el-text>云端引擎：</el-text>
      <el-switch
          v-model="store.state.scan_configuration.advanced.cloud_scan"
          inline-prompt
          active-text="开"
          inactive-text="关"
      />
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-text>检测模式：</el-text>
      <el-switch
          v-model="store.state.scan_configuration.advanced.detection_mode"
          inline-prompt
          style="--el-switch-on-color: #E6A23C"
          active-text="检出优先"
          inactive-text="减少干扰"
      />
      <el-text v-if="store.state.scan_configuration.advanced.detection_mode">
        <el-icon>
          <Warning />
        </el-icon>检出优先模式会大大增加误报的情况，请酌情选择。</el-text>
    </el-col>
  </el-row>
</template>

<script>
import { useStore } from "vuex";
import { Warning } from "@element-plus/icons-vue";


export default {
  name: "advanced",
  components: {
    Warning,
  },
  setup(){
    const store = useStore()
    return{
      store
    }
  }
}
</script>

<style scoped>

</style>