<template>
  <el-row>
    <el-col :span="3" class="grid-cell"><el-text>任务名：</el-text></el-col>
    <el-col :span="21" class="grid-cell">
      <el-input size="small" v-model="store.state.scan_configuration.mission_name" placeholder="请输入创建的任务名"></el-input>
    </el-col>
    <el-col :span="24" class="grid-cell">
      <el-upload
          @dragover.prevent
          @drop="handleDropUpFile"
          @click="upFile"
          ref="uploadRef"
          class="upload-demo"
          drag
          :auto-upload="false"
          :disabled="true"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将待检测的项目拖拽到此 或者<em>点击此处选择文件夹</em>
        </div>
      </el-upload>

      <div v-for="(scanMap, index) in store.state.scan_configuration.localScan.scanPath" :key="index">
        <el-row>
          <el-col :span="22">
            <el-input
                disabled
                size="small"
                v-model="scanMap.path"
                placeholder="当前路径未选择"
            />
          </el-col>
          <el-col :span="2">
            <el-button size="small" @click="delScanList(scanMap)" text><el-icon><Delete /></el-icon></el-button>
          </el-col>
        </el-row>
      </div>
      <div class="el-upload__tip">
        创建任务后检出的内容会在任务列表中出现。
      </div>
    </el-col>
    <el-col :offset="19" :span="4" class="grid-cell">
      <el-button type="primary" @click="addLocalScan" plain :loading="loadingFlag">
        创建任务<el-icon class="el-icon--right"><Pointer /></el-icon>
      </el-button>
    </el-col>
    <transition name="el-zoom-in-top">
      <el-col :span="24" v-show="check_result.flag" class="grid-cell">
        <el-table :data="check_result.tableData" max-height="calc(100vh - 150px)" :empty-text="'未发现可疑文件'" style="width: 100%">
          <el-table-column prop="path" label="文件名" show-overflow-tooltip sortable/>
          <el-table-column prop="results" label="结果" width="100" sortable>
            <template #default="scope">
              <el-button v-if="scope.row.results === '待检测'" size="small" type="info" :loading-icon="UploadFilled" loading plain>待检测</el-button>
              <el-button v-if="scope.row.results === '正常'" size="small" type="success" :icon="CircleCheck" plain>正常</el-button>
              <el-button v-if="scope.row.results === '恶意'" size="small" type="warning" :icon="Warning" plain>恶意</el-button>
              <el-button v-if="scope.row.results === '不支持'" size="small" type="danger" :icon="CircleClose" plain>不支持</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="size" label="大小" width="100" sortable>
          <template #default="scope">
            {{formatBytes(scope.row.size)}}
          </template>
          </el-table-column>
          <el-table-column prop="modificationTime" label="修改时间" sortable>
            <template #default="scope">
              {{ formatTimestamp(scope.row.modificationTime) }}
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </transition>
  </el-row>
</template>

<script>
import {UploadFilled, Pointer, Delete, CircleCheck, Warning, CircleClose} from "@element-plus/icons-vue";
import {ipcRenderer} from "electron";
import {useStore} from "vuex/dist/vuex.mjs";
import {computed, reactive, onUnmounted, ref} from "vue";
import {ElNotification} from "element-plus";
import {formatBytes, formatTimestamp, isNewData} from "@/utils/utils";
export default {
  name: "localScan",
  methods: {formatBytes},
  computed: {
    CircleClose() {
      return CircleClose
    },
    Warning() {
      return Warning
    },
    CircleCheck() {
      return CircleCheck
    },
    UploadFilled() {
      return UploadFilled
    }
  },
  components: {
    UploadFilled,
    Pointer,
    Delete
  },
  setup(){
    const store = useStore()
    const loadingFlag = ref(false)
    const scanPath = computed({
      get: () => store.state.scan_configuration.localScan.scanPath,
      set: (value) => store.commit('setScanPath', value),
    })
    const check_result = reactive({
      tableData: [],
      flag: false,
    })
    let stopLoop = false;
    async function handleDropUpFile(event) {
      // TODO 拖拽上传
      event.preventDefault(); // 阻止默认的拖拽行为
      const files = event.dataTransfer.files; // 获取拖拽的文件列表
      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        const folderInfo = {
          path: file.path,
          name: file.name,
          size: file.size,
          lastModified: Math.round(file.lastModified),
        };

        const exists = store.state.scan_configuration.localScan.scanPath.some(info => info.path === folderInfo.path) // 检测是否有重复的内容
        if (!exists) {
          store.state.scan_configuration.localScan.scanPath.push(folderInfo)
          loadingFlag.value = false
        }
      }
    }
    const upFile = async () => {
      // TODO 点击上传
      const file = await ipcRenderer.invoke('dialog:openFile', "openDirectory")
      const folderInfo = {
        path: file.path,
        name: file.name,
        size: file.size,
        lastModified: Math.round(file.lastModified),
      };
      const exists = store.state.scan_configuration.localScan.scanPath.some(info => info.path === folderInfo.path) // 检测是否有重复的内容
      if (!exists) {
        store.state.scan_configuration.localScan.scanPath.push(folderInfo)
        loadingFlag.value = false
      }
    }
    const delScanList = (folderInfo) => {
      const index = store.state.scan_configuration.localScan.scanPath.findIndex(
          (item) => item.path === folderInfo.path
      );
      if (index !== -1) {
        store.state.scan_configuration.localScan.scanPath.splice(index, 1);
      }
    };

    const addLocalScan = async () => {
      if (store.state.scan_configuration.localScan.scanPath.length === 0) {
        ElNotification({
          title: "注意",
          message: `没有选择文件`,
          type: "warning",
        })
        return
      }
      loadingFlag.value = true
      // TODO 创建本地扫描任务
      let cloudFlag
      if (store.state.scan_configuration.advanced.cloud_scan){
        cloudFlag = 1
      }else {
        cloudFlag = 2
      }
      const scanConfig = {
        scanPath: store.state.scan_configuration.localScan.scanPath,
        model: store.state.scan_configuration.localScan.model,
        cloud_scan: store.state.scan_configuration.advanced.cloud_scan,
        detection_mode: store.state.scan_configuration.advanced.detection_mode,
        mission_name: store.state.scan_configuration.mission_name,
      }
      const data = JSON.stringify(scanConfig) // 将数组转换成 JSON 字符串
      store.state.badge.taskList = store.state.badge.taskList + 1
      const response = await ipcRenderer.invoke('go-request', {
        function: 'addTask',
        data: data,
        cloudFlag:cloudFlag,
      });
      store.state.scan_configuration.localScan.scanPath = []
      // TODO 展示当前检测的任务
      let counter = 0; // 添加一个计数器变量
      while (!stopLoop) {
        // 在这里添加需要执行的代码
        if (response["data"].slice(0, 3) === "ID:" && response["data"].startsWith("ID:")) {
          const id = response["data"].slice(3); // 从第3个位置开始截取到字符串末尾
          const quickCheckList = await ipcRenderer.invoke('go-request', {
            function: 'getTaskList',
            data: parseInt(id),
          });
          if (isNewData(quickCheckList["taskBaseList"],check_result.tableData)){
            check_result.tableData = quickCheckList["taskBaseList"]
          }else if (counter >= 5){
            loadingFlag.value = false
          }
          check_result.flag = true
          // stopLoop = true;  // 记得要删
        }
        await new Promise(resolve => setTimeout(resolve, 1000)); // 延迟1秒
        counter = counter + 1
      }
    }
    onUnmounted(() => {
      stopLoop = true;
    })
    return{
      addLocalScan,upFile,handleDropUpFile,delScanList,formatTimestamp,
      store,scanPath,check_result,loadingFlag
    }
  },
}
</script>

<style scoped>

</style>