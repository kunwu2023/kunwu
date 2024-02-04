import { createStore } from 'vuex'

export default createStore({
  // 创建数据仓库
  state:{
    badge:{
      taskList:0,
    },
    scan_configuration:{
      // TODO 扫描配置文件
      mission_name:"默认任务",
      localScan:{
        scanPath:[],
        model:1
      },
      remoteScan:{
        userName:"",
        passWord:"",
        serverIp:"",
        sshPort:22,
        scanPath:"",
        model:2
      },
      advanced:{
        cloud_scan: true,
        detection_mode: false,
      }
    }, // 扫描配置文件
  },
  // 使用以下方法调用数据仓库

  // 异步调用同步调用，然后用同步调用改值
  // 同步调用
  mutations:{
    setLocalScanPath(state, path) {
      state.scan_configuration.localScan.scanPath = path
    },
    setScanConfiguration(state, val){
      state.scan_configuration = val
    },
    setQuickCheckResult(state, val){
      state.quick_check_result = val
    },
    setScanPath(state, val) {
      state.scan_configuration.localScan.scanPath = val
    },
  },

  // 异步调用
  actions:{
    setLocalScanPath({ commit }, path) {
      commit('setLocalScanPath', path)
    },
    setScanConfiguration({ commit }, val){
      commit('setScanConfiguration', val)
    },
    setQuickCheckResult({ commit }, val){
      commit('setQuickCheckResult', val)
    },
  }
})